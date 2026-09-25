package ui

import (
	"context"
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

func StartSpinner(ctx context.Context, message string) context.CancelFunc {
	_, cancel := StartProgressSpinner(ctx, message, 0)
	return cancel
}

func StartProgressSpinner(ctx context.Context, message string, total int) (func(), context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	var completed atomic.Int32
	done := make(chan struct{})

	increment := func() {
		completed.Add(1)
	}

	go func() {
		defer close(done)
		frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		for {
			select {
			case <-ctx.Done():
				fmt.Fprintf(os.Stdout, "\r\033[K")
				return
			default:
				comp := completed.Load()
				var prog string
				if total > 0 {
					prog = fmt.Sprintf(" (%d/%d)", comp, total)
				}
				fmt.Fprintf(os.Stdout, "\r\033[36m%s\033[0m %s%s", frames[i%len(frames)], message, prog)
				i++
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	stopFunc := func() {
		cancel()
		<-done
	}

	return increment, stopFunc
}
