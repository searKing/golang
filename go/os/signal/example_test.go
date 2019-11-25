package signal_test

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"syscall"
	"time"

	signal_ "github.com/searKing/golang/go/os/signal"
)

func ExampleDumpSignalTo() {
	signal_.DumpSignalTo(syscall.Stderr)
	signal_.SetSigInvokeChain(syscall.SIGUSR1, syscall.SIGUSR2, 0, syscall.SIGSEGV)

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	c := make(chan os.Signal, 1)
	//signal.Notify(c, syscall.SIGINT, syscall.SIGSEGV)
	signal_.Notify(c, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGSEGV)

	// simulate to send a SIGINT to this example test
	go func() {
		//_ = syscall.Kill(syscall.Getpid(), syscall.SIGQUIT)
		//_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
		//_ = syscall.Kill(syscall.Getpid(), syscall.SIGSEGV)
		//_ = syscall.Kill(syscall.Getpid(), syscall.SIGFPE)
		//_ = syscall.Kill(syscall.Getpid(), syscall.SIGSEGV)
		//_ = syscall.Kill(syscall.Getpid(), syscall.SIGSEGV)
		//internal.Raise(syscall.SIGSEGV)

		time.Sleep(time.Second)
		var err error
		fmt.Println(err.Error())
	}()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		defer wg.Done()
		for {
			// Block until a signal is received.
			select {
			case s, ok := <-c:
				if !ok {
					return
				}
				//signal.Stop(c)
				//close(c)
				if s == syscall.SIGUSR1 {
					_, _ = fmt.Fprintf(os.Stderr, "Previous run crashed:\n%s\n", signal_.PreviousStacktrace())
					_ = syscall.Kill(syscall.Getpid(), syscall.SIGUSR2)
				} else {
					fmt.Printf("Got signal: %s\n", s)
				}

			case <-time.After(time.Minute):
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
	}()
	wg.Wait()

	// Output:
	// Got signal: interrupt
}
