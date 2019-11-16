package signal_test

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	signal_ "github.com/searKing/golang/go/os/signal"
)

func ExampleSignalAction() {
	signal_.DumpBacktrace(true)
	signal_.SignalDumpTo(syscall.Stdout)
	signal_.SignalAction(syscall.SIGINT)

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)

	go func() {
		time.Sleep(time.Second)
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()
	// Block until a signal is received.
	s := <-c
	fmt.Println("Got signal:", s)

	// Output:
	// Got signal: interrupt
}
