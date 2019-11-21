package signal_test

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	signal_ "github.com/searKing/golang/go/os/signal"
)

func ExampleDumpSignalTo() {
	signal_.DumpSignalTo(syscall.Stderr)
	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)

	// simulate to send a SIGINT to this example test
	go func() {
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()
	for{
		// Block until a signal is received.
		select {
		case s, ok := <-c:
			if !ok {
				return
			}
			fmt.Printf("Got signal: %s\n", s)
			_, _ = fmt.Fprintf(os.Stderr, "Previous run crashed:\n%s\n", signal_.PreviousStacktrace())
			signal.Stop(c)
			close(c)
		case <-time.After(time.Second):
			_, _ = fmt.Fprintf(os.Stderr, "time overseed:\n")
			return
		}
		// set os.Stderr tp pass test, for the stacktrace is random.
		//stderr prints something like:
		//Signal received(2).
		//Stacktrace dumped to file: stacktrace.dump.
		//Previous run crashed:
		// 0# searking::SignalHandler::operator()(int, __siginfo*, void*) in /private/var/folders/12/870qx8rd0_d96nt6g078wp080000gn/T/___ExampleSignalAction_in_github_com_searKing_golang_go_os_signal
		// 1# _sigtramp in /usr/lib/system/libsystem_platform.dylib

	}
	// Output:
	// Got signal: interrupt
}
