package exception

import (
	object2 "github.com/searKing/golang/go/util/object"
	"io"
	"os"
	"runtime/debug"
)

const (
	/** Message for trying to suppress a null exception. */
	NullCauseMessage string = "Cannot suppress a null exception."
	/** Message for trying to suppress oneself. */
	SelfSuppressionMessage string = "Self-suppression not permitted"
	/** Caption  for labeling causative exception stack traces */
	CauseCaption string = "Caused by: "
	/** Caption for labeling suppressed exception stack traces */
	SuppressedCaption string = "Suppressed: "
)

type ThrowableInterface interface {
	GetMessage() string
	GetLocalizedMessage() string
	GetCause() ThrowableInterface
	InitCause(cause ThrowableInterface) ThrowableInterface
	ToString() string
	PrintStackTrace()
	PrintStackTrace1(writer io.Writer)
	fillInStackTrack()
	GetSuppressed() []ThrowableInterface
	GetStackTrace() []byte
	Error() string
}
type Throwable struct {
	detailMessage        string
	cause                ThrowableInterface
	stackTrace           []byte
	suppressedExceptions []ThrowableInterface
}

func NewThrowable() *Throwable {
	return NewThrowable1("")
}
func NewThrowable1(message string) *Throwable {
	return NewThrowable2(message, nil)
}
func NewThrowable2(message string, cause ThrowableInterface) *Throwable {
	t := &Throwable{}
	t.fillInStackTrack()
	t.detailMessage = message
	t.cause = cause
	return t
}
func NewThrowable4(message string, cause *Throwable, enableSuppression, writableStackTrace bool) *Throwable {
	t := &Throwable{}
	if writableStackTrace {
		t.fillInStackTrack()
	} else {
		t.stackTrace = debug.Stack()
	}
	t.detailMessage = message
	t.cause = cause
	if !enableSuppression {
		t.suppressedExceptions = nil
	}
	t.suppressedExceptions = []ThrowableInterface{}
	return t
}
func (thiz *Throwable) Error() string {
	return thiz.GetMessage()
}

func (thiz *Throwable) GetMessage() string {
	return thiz.detailMessage
}
func (thiz *Throwable) GetLocalizedMessage() string {
	return thiz.GetMessage()
}
func (thiz *Throwable) GetCause() ThrowableInterface {
	if thiz.cause == thiz {
		return nil
	}
	return thiz.cause
}
func (thiz *Throwable) InitCause(cause ThrowableInterface) ThrowableInterface {
	if thiz.cause != thiz {
		return NewIllegalStateException2("Can't overwrite cause with "+object2.ToString(cause, "a nil"), thiz)
	}
	if cause == thiz {
		return NewIllegalArgumentException2("Self-causation not permitted", thiz)
	}
	return thiz.cause
}
func (thiz *Throwable) ToString() string {
	s := object2.GetStruct().Name()
	message := thiz.GetLocalizedMessage()
	if len(message) != 0 {
		return s + ":" + message
	}
	return s
}
func (thiz *Throwable) PrintStackTrace() {
	thiz.PrintStackTrace1(os.Stderr)
}
func (thiz *Throwable) PrintStackTrace1(writer io.Writer) {
	writer.Write(thiz.GetStackTrace())
	for _, se := range thiz.GetSuppressed() {
		se.PrintStackTrace1(writer)
	}
	ourCause := thiz.GetCause()
	if ourCause != nil {
		ourCause.PrintStackTrace1(writer)
	}
}
func (thiz *Throwable) fillInStackTrack() {

}
func (thiz *Throwable) GetSuppressed() []ThrowableInterface {
	return thiz.suppressedExceptions
}

func (thiz *Throwable) GetStackTrace() []byte {
	return object2.DeepClone(thiz.GetOurStackTrace()).([]byte)
}
func (thiz *Throwable) GetOurStackTrace() []byte {
	return thiz.stackTrace
}
func (thiz *Throwable) SetStackTrace(trace []byte) {
	thiz.stackTrace = object2.DeepClone(trace).([]byte)
}
func (thiz *Throwable) AddSuppressed(exception ThrowableInterface) {
	if exception == thiz {
		panic(NewIllegalArgumentException2(SelfSuppressionMessage, exception))
	}
	if exception == nil {
		panic(NewNullPointerException(NullCauseMessage))
	}
	if thiz.suppressedExceptions == nil {
		return
	}
	thiz.suppressedExceptions = append(thiz.suppressedExceptions, exception)
}
