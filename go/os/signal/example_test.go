package signal_test

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

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
	// Block until a signal is received.
	s := <-c
	fmt.Printf("Got signal: %s\n", s)
	_, _ = fmt.Fprintf(os.Stderr, "Previous run crashed:\n%s\n", signal_.PreviousStacktrace())
	// set os.Stderr tp pass test, for the stacktrace is random.
	//stderr prints something like:
	//Signal received(2).
	//Backtrace dumped to file: stacktrace.dump.
	//Previous run crashed:
	// 0# searking::SignalHandler::operator()(int, __siginfo*, void*) in /private/var/folders/12/870qx8rd0_d96nt6g078wp080000gn/T/___ExampleSignalAction_in_github_com_searKing_golang_go_os_signal
	// 1# _sigtramp in /usr/lib/system/libsystem_platform.dylib
	_ = s
	// Output:
	// Got signal: interrupt
}
