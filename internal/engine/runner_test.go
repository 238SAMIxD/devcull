package engine_test

import (
	"context"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"github.com/238SAMIxD/devcull/internal/cleaner"
	"github.com/238SAMIxD/devcull/internal/engine"
)

type mockCleaner struct {
	name      string
	installed bool

	activeCount *int32
	maxActive   *int32
	blockCh     chan struct{}

	isInstalledCount atomic.Int32
	cleanCount       atomic.Int32
}

func (m *mockCleaner) Name() string               { return m.name }
func (m *mockCleaner) Category() cleaner.Category { return cleaner.Category("Test") }
func (m *mockCleaner) Aliases() []string          { return nil }

func (m *mockCleaner) IsInstalled(ctx context.Context) bool {
	m.isInstalledCount.Add(1)

	if m.activeCount != nil && m.maxActive != nil && m.blockCh != nil {
		current := atomic.AddInt32(m.activeCount, 1)

		for {
			max := atomic.LoadInt32(m.maxActive)
			if current > max {
				if atomic.CompareAndSwapInt32(m.maxActive, max, current) {
					break
				}
			} else {
				break
			}
		}

		select {
		case <-m.blockCh:
		case <-ctx.Done():
		}

		atomic.AddInt32(m.activeCount, -1)
	}

	return m.installed
}

func (m *mockCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	return 100, nil
}

func (m *mockCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	m.cleanCount.Add(1)
	return 100, nil
}

func TestRunner_ConcurrencyLimit(t *testing.T) {
	numCPU := runtime.NumCPU()
	var activeCount, maxActive int32
	blockCh := make(chan struct{})

	numCleaners := numCPU * 2
	if numCleaners < 4 {
		numCleaners = 4
	}

	cleaners := make([]cleaner.Cleaner, numCleaners)
	for i := 0; i < numCleaners; i++ {
		cleaners[i] = &mockCleaner{
			name:        "mock",
			installed:   true,
			activeCount: &activeCount,
			maxActive:   &maxActive,
			blockCh:     blockCh,
		}
	}

	done := make(chan struct{})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		engine.Run(ctx, cleaners, false, nil)
		close(done)
	}()

	time.Sleep(200 * time.Millisecond)

	close(blockCh)

	<-done

	observedMax := atomic.LoadInt32(&maxActive)
	if int(observedMax) > numCPU {
		t.Errorf("Expected max concurrency <= %d, got %d", numCPU, observedMax)
	}
	if int(observedMax) == 0 {
		t.Errorf("Expected some concurrency, got 0")
	}
}

func TestRunner_ContextCancellation(t *testing.T) {
	numCleaners := runtime.NumCPU() + 5
	cleaners := make([]cleaner.Cleaner, numCleaners)
	mocks := make([]*mockCleaner, numCleaners)
	for i := 0; i < numCleaners; i++ {
		mocks[i] = &mockCleaner{
			name:      "mock",
			installed: true,
		}
		cleaners[i] = mocks[i]
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	results := engine.Run(ctx, cleaners, false, nil)

	if len(results) != numCleaners {
		t.Errorf("Expected %d results, got %d", numCleaners, len(results))
	}

	for _, res := range results {
		if res.Err != context.Canceled {
			t.Errorf("Expected context.Canceled error, got %v", res.Err)
		}
	}

	for _, m := range mocks {
		if m.isInstalledCount.Load() > 0 {
			t.Errorf("IsInstalled was called despite pre-canceled context")
		}
		if m.cleanCount.Load() > 0 {
			t.Errorf("Clean was called despite pre-canceled context")
		}
	}
}
