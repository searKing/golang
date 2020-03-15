// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exception

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime/debug"

	"github.com/searKing/golang/go/util/object"
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

var (
	/**
	 * A shared value for an empty stack.
	 */
	UnassignedStack = make([]byte, 0)
	// Setting this static field introduces an acceptable
	// initialization dependency on a few java.util classes.
	SuppressedSentinel = make([]Throwable, 0)

	EmptyThrowableArray = make([]Throwable, 0)
)

type Throwable interface {
	GetMessage() string
	GetLocalizedMessage() string
	GetCause() Throwable
	InitCause(cause Throwable) Throwable
	String() string
	PrintStackTrace()
	PrintStackTrace1(writer io.Writer)
	PrintEnclosedStackTrace(writer io.Writer, enclosingTrace []byte,
		caption string, prefix string, dejaVu map[Throwable]struct{})
	fillInStackTrack() Throwable
	GetSuppressed() []Throwable
	GetStackTrace() []byte
	Error() string
}

/**
 * The {@code Throwable} class is the superclass of all errors and
 * exceptions in the Java language. Only objects that are instances of this
 * class (or one of its subclasses) are thrown by the Java Virtual Machine or
 * can be thrown by the Java {@code throw} statement. Similarly, only
 * this class or one of its subclasses can be the argument type in a
 * {@code catch} clause.
 *
 * For the purposes of compile-time checking of exceptions, {@code
 * Throwable} and any subclass of {@code Throwable} that is not also a
 * subclass of either {@link RuntimeException} or {@link Error} are
 * regarded as checked exceptions.
 *
 * <p>Instances of two subclasses, {@link java.lang.Error} and
 * {@link java.lang.Exception}, are conventionally used to indicate
 * that exceptional situations have occurred. Typically, these instances
 * are freshly created in the context of the exceptional situation so
 * as to include relevant information (such as stack trace data).
 *
 * <p>A throwable contains a snapshot of the execution stack of its
 * thread at the time it was created. It can also contain a message
 * string that gives more information about the error. Over time, a
 * throwable can {@linkplain Throwable#addSuppressed suppress} other
 * throwables from being propagated.  Finally, the throwable can also
 * contain a <i>cause</i>: another throwable that caused this
 * throwable to be constructed.  The recording of this causal information
 * is referred to as the <i>chained exception</i> facility, as the
 * cause can, itself, have a cause, and so on, leading to a "chain" of
 * exceptions, each caused by another.
 *
 * <p>One reason that a throwable may have a cause is that the class that
 * throws it is built atop a lower layered abstraction, and an operation on
 * the upper layer fails due to a failure in the lower layer.  It would be bad
 * design to let the throwable thrown by the lower layer propagate outward, as
 * it is generally unrelated to the abstraction provided by the upper layer.
 * Further, doing so would tie the API of the upper layer to the details of
 * its implementation, assuming the lower layer's exception was a checked
 * exception.  Throwing a "wrapped exception" (i.e., an exception containing a
 * cause) allows the upper layer to communicate the details of the failure to
 * its caller without incurring either of these shortcomings.  It preserves
 * the flexibility to change the implementation of the upper layer without
 * changing its API (in particular, the set of exceptions thrown by its
 * methods).
 *
 * <p>A second reason that a throwable may have a cause is that the method
 * that throws it must conform to a general-purpose interface that does not
 * permit the method to throw the cause directly.  For example, suppose
 * a persistent collection conforms to the {@link java.util.Collection
 * Collection} interface, and that its persistence is implemented atop
 * {@code java.io}.  Suppose the internals of the {@code add} method
 * can throw an {@link java.io.IOException IOException}.  The implementation
 * can communicate the details of the {@code IOException} to its caller
 * while conforming to the {@code Collection} interface by wrapping the
 * {@code IOException} in an appropriate unchecked exception.  (The
 * specification for the persistent collection should indicate that it is
 * capable of throwing such exceptions.)
 *
 * <p>A cause can be associated with a throwable in two ways: via a
 * constructor that takes the cause as an argument, or via the
 * {@link #initCause(Throwable)} method.  New throwable classes that
 * wish to allow causes to be associated with them should provide constructors
 * that take a cause and delegate (perhaps indirectly) to one of the
 * {@code Throwable} constructors that takes a cause.
 *
 * Because the {@code initCause} method is public, it allows a cause to be
 * associated with any throwable, even a "legacy throwable" whose
 * implementation predates the addition of the exception chaining mechanism to
 * {@code Throwable}.
 *
 * <p>By convention, class {@code Throwable} and its subclasses have two
 * constructors, one that takes no arguments and one that takes a
 * {@code String} argument that can be used to produce a detail message.
 * Further, those subclasses that might likely have a cause associated with
 * them should have two more constructors, one that takes a
 * {@code Throwable} (the cause), and one that takes a
 * {@code String} (the detail message) and a {@code Throwable} (the
 * cause).
 *
 * @author  unascribed
 * @author  Josh Bloch (Added exception chaining and programmatic access to
 *          stack trace in 1.4.)
 * @jls 11.2 Compile-Time Checking of Exceptions
 * @since 1.0
 */
type ThrowableObject struct {
	/**
	 * Specific details about the Throwable.  For example, for
	 * {@code FileNotFoundException}, this contains the name of
	 * the file that could not be found.
	 *
	 * @serial
	 */
	detailMessage string

	/**
	 * The throwable that caused this throwable to get thrown, or null if this
	 * throwable was not caused by another throwable, or if the causative
	 * throwable is unknown.  If this field is equal to this throwable itself,
	 * it indicates that the cause of this throwable has not yet been
	 * initialized.
	 *
	 * @serial
	 * @since 1.4
	 */
	cause Throwable

	/**
	 * The stack trace, as returned by {@link #getStackTrace()}.
	 *
	 * The field is initialized to a zero-length array.  A {@code
	 * null} value of this field indicates subsequent calls to {@link
	 * #setStackTrace(StackTraceElement[])} and {@link
	 * #fillInStackTrace()} will be no-ops.
	 *
	 * @serial
	 * @since 1.4
	 */
	stackTrace []byte

	/**
	 * The list of suppressed exceptions, as returned by {@link
	 * #getSuppressed()}.  The list is initialized to a zero-element
	 * unmodifiable sentinel list.  When a serialized Throwable is
	 * read in, if the {@code suppressedExceptions} field points to a
	 * zero-element list, the field is reset to the sentinel value.
	 *
	 * @serial
	 * @since 1.7
	 */
	suppressedExceptions []Throwable
}

/**
 * Constructs a new throwable with {@code null} as its detail message.
 * The cause is not initialized, and may subsequently be initialized by a
 * call to {@link #initCause}.
 *
 * <p>The {@link #fillInStackTrace()} method is called to initialize
 * the stack trace data in the newly created throwable.
 */
func NewThrowable() *ThrowableObject {
	return NewThrowable1("")
}

/**
 * Constructs a new throwable with the specified detail message.  The
 * cause is not initialized, and may subsequently be initialized by
 * a call to {@link #initCause}.
 *
 * <p>The {@link #fillInStackTrace()} method is called to initialize
 * the stack trace data in the newly created throwable.
 *
 * @param   message   the detail message. The detail message is saved for
 *          later retrieval by the {@link #getMessage()} method.
 */
func NewThrowable1(message string) *ThrowableObject {
	t := NewThrowable2(message, nil)
	t.cause = t
	return t
}

/**
 * Constructs a new throwable with the specified detail message and
 * cause.  <p>Note that the detail message associated with
 * {@code cause} is <i>not</i> automatically incorporated in
 * this throwable's detail message.
 *
 * <p>The {@link #fillInStackTrace()} method is called to initialize
 * the stack trace data in the newly created throwable.
 *
 * @param  message the detail message (which is saved for later retrieval
 *         by the {@link #getMessage()} method).
 * @param  cause the cause (which is saved for later retrieval by the
 *         {@link #getCause()} method).  (A {@code null} value is
 *         permitted, and indicates that the cause is nonexistent or
 *         unknown.)
 * @since  1.4
 */
func NewThrowable2(message string, cause Throwable) *ThrowableObject {
	return NewThrowable4(message, cause, true, true)
}

/**
 * Constructs a new throwable with the specified detail message,
 * cause, {@linkplain #addSuppressed suppression} enabled or
 * disabled, and writable stack trace enabled or disabled.  If
 * suppression is disabled, {@link #getSuppressed} for this object
 * will return a zero-length array and calls to {@link
 * #addSuppressed} that would otherwise append an exception to the
 * suppressed list will have no effect.  If the writable stack
 * trace is false, this constructor will not call {@link
 * #fillInStackTrace()}, a {@code null} will be written to the
 * {@code stackTrace} field, and subsequent calls to {@code
 * fillInStackTrace} and {@link
 * #setStackTrace(StackTraceElement[])} will not set the stack
 * trace.  If the writable stack trace is false, {@link
 * #getStackTrace} will return a zero length array.
 *
 * <p>Note that the other constructors of {@code Throwable} treat
 * suppression as being enabled and the stack trace as being
 * writable.  Subclasses of {@code Throwable} should document any
 * conditions under which suppression is disabled and document
 * conditions under which the stack trace is not writable.
 * Disabling of suppression should only occur in exceptional
 * circumstances where special requirements exist, such as a
 * virtual machine reusing exception objects under low-memory
 * situations.  Circumstances where a given exception object is
 * repeatedly caught and rethrown, such as to implement control
 * flow between two sub-systems, is another situation where
 * immutable throwable objects would be appropriate.
 *
 * @param  message the detail message.
 * @param cause the cause.  (A {@code null} value is permitted,
 * and indicates that the cause is nonexistent or unknown.)
 * @param enableSuppression whether or not suppression is enabled or disabled
 * @param writableStackTrace whether or not the stack trace should be
 *                           writable
 *
 * @see OutOfMemoryError
 * @see NullPointerException
 * @see ArithmeticException
 * @since 1.7
 */
func NewThrowable4(message string, cause Throwable, enableSuppression, writableStackTrace bool) *ThrowableObject {
	t := &ThrowableObject{}
	t.init()
	if writableStackTrace {
		t.fillInStackTrack()
	} else {
		t.stackTrace = nil
	}
	t.detailMessage = message
	t.cause = cause
	if !enableSuppression {
		t.suppressedExceptions = nil
	}
	return t
}

// member variable init as C++
func (t *ThrowableObject) init() {
	t.cause = t
	t.stackTrace = UnassignedStack
	t.suppressedExceptions = SuppressedSentinel
}

func (t *ThrowableObject) Error() string {
	return t.GetMessage()
}

/**
 * Returns the detail message string of this throwable.
 *
 * @return  the detail message string of this {@code Throwable} instance
 *          (which may be {@code null}).
 */
func (t *ThrowableObject) GetMessage() string {
	return t.detailMessage
}

/**
 * Creates a localized description of this throwable.
 * Subclasses may override this method in order to produce a
 * locale-specific message.  For subclasses that do not override this
 * method, the default implementation returns the same result as
 * {@code getMessage()}.
 *
 * @return  The localized description of this throwable.
 * @since   1.1
 */
func (t *ThrowableObject) GetLocalizedMessage() string {
	return t.GetMessage()
}

/**
 * Returns the cause of this throwable or {@code null} if the
 * cause is nonexistent or unknown.  (The cause is the throwable that
 * caused this throwable to get thrown.)
 *
 * <p>This implementation returns the cause that was supplied via one of
 * the constructors requiring a {@code Throwable}, or that was set after
 * creation with the {@link #initCause(Throwable)} method.  While it is
 * typically unnecessary to override this method, a subclass can override
 * it to return a cause set by some other means.  This is appropriate for
 * a "legacy chained throwable" that predates the addition of chained
 * exceptions to {@code Throwable}.  Note that it is <i>not</i>
 * necessary to override any of the {@code PrintStackTrace} methods,
 * all of which invoke the {@code getCause} method to determine the
 * cause of a throwable.
 *
 * @return  the cause of this throwable or {@code null} if the
 *          cause is nonexistent or unknown.
 * @since 1.4
 */
func (t *ThrowableObject) GetCause() Throwable {
	if t.cause == t {
		return nil
	}
	return t.cause
}

/**
 * Initializes the <i>cause</i> of this throwable to the specified value.
 * (The cause is the throwable that caused this throwable to get thrown.)
 *
 * <p>This method can be called at most once.  It is generally called from
 * within the constructor, or immediately after creating the
 * throwable.  If this throwable was created
 * with {@link #Throwable(Throwable)} or
 * {@link #Throwable(String,Throwable)}, this method cannot be called
 * even once.
 *
 * <p>An example of using this method on a legacy throwable type
 * without other support for setting the cause is:
 *
 * <pre>
 * try {
 *     lowLevelOp();
 * } catch (LowLevelException le) {
 *     throw (HighLevelException)
 *           new HighLevelException().initCause(le); // Legacy constructor
 * }
 * </pre>
 *
 * @param  cause the cause (which is saved for later retrieval by the
 *         {@link #getCause()} method).  (A {@code null} value is
 *         permitted, and indicates that the cause is nonexistent or
 *         unknown.)
 * @return  a reference to this {@code Throwable} instance.
 * @throws IllegalArgumentException if {@code cause} is this
 *         throwable.  (A throwable cannot be its own cause.)
 * @throws IllegalStateException if this throwable was
 *         created with {@link #Throwable(Throwable)} or
 *         {@link #Throwable(String,Throwable)}, or this method has already
 *         been called on this throwable.
 * @since  1.4
 */
func (t *ThrowableObject) InitCause(cause Throwable) Throwable {
	if t.cause != t {
		return NewIllegalStateException2("Can't overwrite cause with "+object.ToString(cause, "a nil"), t)
	}
	if cause == t {
		return NewIllegalArgumentException2("Self-causation not permitted", t)
	}
	t.cause = cause
	return t
}

/**
 * Returns a short description of this throwable.
 * The result is the concatenation of:
 * <ul>
 * <li> the {@linkplain Class#getName() name} of the class of this object
 * <li> ": " (a colon and a space)
 * <li> the result of invoking this object's {@link #getLocalizedMessage}
 *      method
 * </ul>
 * If {@code getLocalizedMessage} returns {@code null}, then just
 * the class name is returned.
 *
 * @return a string representation of this throwable.
 */
func (t *ThrowableObject) String() string {
	s := object.GetFunc().Name()
	message := t.GetLocalizedMessage()
	if len(message) != 0 {
		return s + ":" + message
	}
	return s
}

/**
 * Prints this throwable and its backtrace to the
 * standard error stream. This method prints a stack trace for this
 * {@code Throwable} object on the error output stream that is
 * the value of the field {@code System.err}. The first line of
 * output contains the result of the {@link #toString()} method for
 * this object.  Remaining lines represent data previously recorded by
 * the method {@link #fillInStackTrace()}. The format of this
 * information depends on the implementation, but the following
 * example may be regarded as typical:
 * <blockquote><pre>
 * java.lang.NullPointerException
 *         at MyClass.mash(MyClass.java:9)
 *         at MyClass.crunch(MyClass.java:6)
 *         at MyClass.main(MyClass.java:3)
 * </pre></blockquote>
 * This example was produced by running the program:
 * <pre>
 * class MyClass {
 *     public static void main(String[] args) {
 *         crunch(null);
 *     }
 *     static void crunch(int[] a) {
 *         mash(a);
 *     }
 *     static void mash(int[] b) {
 *         System.out.println(b[0]);
 *     }
 * }
 * </pre>
 * The backtrace for a throwable with an initialized, non-null cause
 * should generally include the backtrace for the cause.  The format
 * of this information depends on the implementation, but the following
 * example may be regarded as typical:
 * <pre>
 * HighLevelException: MidLevelException: LowLevelException
 *         at Junk.a(Junk.java:13)
 *         at Junk.main(Junk.java:4)
 * Caused by: MidLevelException: LowLevelException
 *         at Junk.c(Junk.java:23)
 *         at Junk.b(Junk.java:17)
 *         at Junk.a(Junk.java:11)
 *         ... 1 more
 * Caused by: LowLevelException
 *         at Junk.e(Junk.java:30)
 *         at Junk.d(Junk.java:27)
 *         at Junk.c(Junk.java:21)
 *         ... 3 more
 * </pre>
 * Note the presence of lines containing the characters {@code "..."}.
 * These lines indicate that the remainder of the stack trace for this
 * exception matches the indicated number of frames from the bottom of the
 * stack trace of the exception that was caused by this exception (the
 * "enclosing" exception).  This shorthand can greatly reduce the length
 * of the output in the common case where a wrapped exception is thrown
 * from same method as the "causative exception" is caught.  The above
 * example was produced by running the program:
 * <pre>
 * public class Junk {
 *     public static void main(String args[]) {
 *         try {
 *             a();
 *         } catch(HighLevelException e) {
 *             e.printStackTrace();
 *         }
 *     }
 *     static void a() throws HighLevelException {
 *         try {
 *             b();
 *         } catch(MidLevelException e) {
 *             throw new HighLevelException(e);
 *         }
 *     }
 *     static void b() throws MidLevelException {
 *         c();
 *     }
 *     static void c() throws MidLevelException {
 *         try {
 *             d();
 *         } catch(LowLevelException e) {
 *             throw new MidLevelException(e);
 *         }
 *     }
 *     static void d() throws LowLevelException {
 *        e();
 *     }
 *     static void e() throws LowLevelException {
 *         throw new LowLevelException();
 *     }
 * }
 *
 * class HighLevelException extends Exception {
 *     HighLevelException(Throwable cause) { super(cause); }
 * }
 *
 * class MidLevelException extends Exception {
 *     MidLevelException(Throwable cause)  { super(cause); }
 * }
 *
 * class LowLevelException extends Exception {
 * }
 * </pre>
 * As of release 7, the platform supports the notion of
 * <i>suppressed exceptions</i> (in conjunction with the {@code
 * try}-with-resources statement). Any exceptions that were
 * suppressed in order to deliver an exception are printed out
 * beneath the stack trace.  The format of this information
 * depends on the implementation, but the following example may be
 * regarded as typical:
 *
 * <pre>
 * Exception in thread "main" java.lang.Exception: Something happened
 *  at Foo.bar(Foo.java:10)
 *  at Foo.main(Foo.java:5)
 *  Suppressed: Resource$CloseFailException: Resource ID = 0
 *          at Resource.close(Resource.java:26)
 *          at Foo.bar(Foo.java:9)
 *          ... 1 more
 * </pre>
 * Note that the "... n more" notation is used on suppressed exceptions
 * just at it is used on causes. Unlike causes, suppressed exceptions are
 * indented beyond their "containing exceptions."
 *
 * <p>An exception can have both a cause and one or more suppressed
 * exceptions:
 * <pre>
 * Exception in thread "main" java.lang.Exception: Main block
 *  at Foo3.main(Foo3.java:7)
 *  Suppressed: Resource$CloseFailException: Resource ID = 2
 *          at Resource.close(Resource.java:26)
 *          at Foo3.main(Foo3.java:5)
 *  Suppressed: Resource$CloseFailException: Resource ID = 1
 *          at Resource.close(Resource.java:26)
 *          at Foo3.main(Foo3.java:5)
 * Caused by: java.lang.Exception: I did it
 *  at Foo3.main(Foo3.java:8)
 * </pre>
 * Likewise, a suppressed exception can have a cause:
 * <pre>
 * Exception in thread "main" java.lang.Exception: Main block
 *  at Foo4.main(Foo4.java:6)
 *  Suppressed: Resource2$CloseFailException: Resource ID = 1
 *          at Resource2.close(Resource2.java:20)
 *          at Foo4.main(Foo4.java:5)
 *  Caused by: java.lang.Exception: Rats, you caught me
 *          at Resource2$CloseFailException.&lt;init&gt;(Resource2.java:45)
 *          ... 2 more
 * </pre>
 */
func (t *ThrowableObject) PrintStackTrace() {
	t.PrintStackTrace1(os.Stderr)
}

/**
 * Prints this throwable and its backtrace to the specified print stream.
 *
 * @param s {@code PrintStream} to use for output
 */
func (t *ThrowableObject) PrintStackTrace1(writer io.Writer) {
	// Guard against malicious overrides of Throwable.equals by
	// using a Set with identity equality semantics.
	var dejaVu = map[Throwable]struct{}{t: {}}

	// Print our stack trace
	_, _ = writer.Write([]byte(fmt.Sprintln(t)))

	trace := t.GetOurStackTrace()
	_, _ = writer.Write(t.GetOurStackTrace())

	// Print suppressed exceptions, if any
	for _, se := range t.GetSuppressed() {
		se.PrintEnclosedStackTrace(writer, trace, SuppressedCaption, "\t", dejaVu)
	}

	// Print cause, if any
	ourCause := t.GetCause()
	if ourCause != nil {
		ourCause.PrintEnclosedStackTrace(writer, trace, CauseCaption, "", dejaVu)
	}
}

/**
 * Print our stack trace as an enclosed exception for the specified
 * stack trace.
 */
func (t *ThrowableObject) PrintEnclosedStackTrace(writer io.Writer, enclosingTrace []byte,
	caption string, prefix string, dejaVu map[Throwable]struct{}) {
	if _, has := dejaVu[t]; has {
		_, _ = writer.Write([]byte(fmt.Sprintln("\t[CIRCULAR REFERENCE:", t, "]")))
		return
	}

	dejaVu[t] = struct{}{}
	// Compute number of frames in common between this and enclosing trace
	trace := t.GetOurStackTrace()
	traces := bytes.Split(trace, []byte{'\n'})
	enclosingTraces := bytes.Split(enclosingTrace, []byte{'\n'})
	m := len(traces) - 1
	n := len(enclosingTraces) - 1
	for m >= 0 && n >= 0 && bytes.EqualFold(traces[m], enclosingTraces[n]) {
		m--
		n--
	}

	framesInCommon := len(traces) - 1 - m

	// Print our stack trace
	_, _ = writer.Write([]byte(fmt.Sprintln(prefix+caption, t)))
	for i := 0; i <= m; i++ {
		_, _ = writer.Write([]byte(fmt.Sprintln(traces[i])))
	}
	if framesInCommon != 0 {
		_, _ = writer.Write([]byte(fmt.Sprintln(prefix, "\t...", framesInCommon, " more")))
	}

	// Print suppressed exceptions, if any
	for _, se := range t.GetSuppressed() {
		se.PrintEnclosedStackTrace(writer, trace, SuppressedCaption, prefix+"\t", dejaVu)
	}

	// Print cause, if any
	ourCause := t.GetCause()
	if ourCause != nil {
		ourCause.PrintEnclosedStackTrace(writer, trace, CauseCaption, prefix, dejaVu)
	}
}

/**
 * Fills in the execution stack trace. This method records within this
 * {@code Throwable} object information about the current state of
 * the stack frames for the current thread.
 *
 * <p>If the stack trace of this {@code Throwable} {@linkplain
 * Throwable#Throwable(String, Throwable, boolean, boolean) is not
 * writable}, calling this method has no effect.
 *
 * @return  a reference to this {@code Throwable} instance.
 * @see     java.lang.Throwable#printStackTrace()
 */
func (t *ThrowableObject) fillInStackTrack() Throwable {
	if t.stackTrace != nil {
		t.stackTrace = UnassignedStack
	}
	return t
}

/**
 * Provides programmatic access to the stack trace information printed by
 * {@link #printStackTrace()}.  Returns an array of stack trace elements,
 * each representing one stack frame.  The zeroth element of the array
 * (assuming the array's length is non-zero) represents the top of the
 * stack, which is the last method invocation in the sequence.  Typically,
 * this is the point at which this throwable was created and thrown.
 * The last element of the array (assuming the array's length is non-zero)
 * represents the bottom of the stack, which is the first method invocation
 * in the sequence.
 *
 * <p>Some virtual machines may, under some circumstances, omit one
 * or more stack frames from the stack trace.  In the extreme case,
 * a virtual machine that has no stack trace information concerning
 * this throwable is permitted to return a zero-length array from this
 * method.  Generally speaking, the array returned by this method will
 * contain one element for every frame that would be printed by
 * {@code printStackTrace}.  Writes to the returned array do not
 * affect future calls to this method.
 *
 * @return an array of stack trace elements representing the stack trace
 *         pertaining to this throwable.
 * @since  1.4
 */
func (t *ThrowableObject) GetStackTrace() []byte {
	return object.DeepClone(t.GetOurStackTrace()).([]byte)
}

func (t *ThrowableObject) GetOurStackTrace() []byte {
	if t.stackTrace == nil {
		return UnassignedStack
	}

	if bytes.Equal(t.stackTrace, UnassignedStack) {
		t.stackTrace = debug.Stack()
	}
	return t.stackTrace
}

/**
 * Sets the stack trace elements that will be returned by
 * {@link #getStackTrace()} and printed by {@link #printStackTrace()}
 * and related methods.
 *
 * This method, which is designed for use by RPC frameworks and other
 * advanced systems, allows the client to override the default
 * stack trace that is either generated by {@link #fillInStackTrace()}
 * when a throwable is constructed or deserialized when a throwable is
 * read from a serialization stream.
 *
 * <p>If the stack trace of this {@code Throwable} {@linkplain
 * Throwable#Throwable(String, Throwable, boolean, boolean) is not
 * writable}, calling this method has no effect other than
 * validating its argument.
 *
 * @param   stackTrace the stack trace elements to be associated with
 * this {@code Throwable}.  The specified array is copied by this
 * call; changes in the specified array after the method invocation
 * returns will have no affect on this {@code Throwable}'s stack
 * trace.
 *
 * @throws NullPointerException if {@code stackTrace} is
 *         {@code null} or if any of the elements of
 *         {@code stackTrace} are {@code null}
 *
 * @since  1.4
 */
func (t *ThrowableObject) SetStackTrace(trace []byte) {
	// Validate argument
	if len(trace) == 0 {
		panic(NewNullPointerException1("stackTrace[]"))
		return
	}
	if t.stackTrace == nil { // Immutable stack
		return
	}
	defensiveCopy := object.DeepClone(trace).([]byte)
	t.stackTrace = defensiveCopy
}

/**
 * Appends the specified exception to the exceptions that were
 * suppressed in order to deliver this exception. This method is
 * thread-safe and typically called (automatically and implicitly)
 * by the {@code try}-with-resources statement.
 *
 * <p>The suppression behavior is enabled <em>unless</em> disabled
 * {@linkplain #Throwable(String, Throwable, boolean, boolean) via
 * a constructor}.  When suppression is disabled, this method does
 * nothing other than to validate its argument.
 *
 * <p>Note that when one exception {@linkplain
 * #initCause(Throwable) causes} another exception, the first
 * exception is usually caught and then the second exception is
 * thrown in response.  In other words, there is a causal
 * connection between the two exceptions.
 *
 * In contrast, there are situations where two independent
 * exceptions can be thrown in sibling code blocks, in particular
 * in the {@code try} block of a {@code try}-with-resources
 * statement and the compiler-generated {@code finally} block
 * which closes the resource.
 *
 * In these situations, only one of the thrown exceptions can be
 * propagated.  In the {@code try}-with-resources statement, when
 * there are two such exceptions, the exception originating from
 * the {@code try} block is propagated and the exception from the
 * {@code finally} block is added to the list of exceptions
 * suppressed by the exception from the {@code try} block.  As an
 * exception unwinds the stack, it can accumulate multiple
 * suppressed exceptions.
 *
 * <p>An exception may have suppressed exceptions while also being
 * caused by another exception.  Whether or not an exception has a
 * cause is semantically known at the time of its creation, unlike
 * whether or not an exception will suppress other exceptions
 * which is typically only determined after an exception is
 * thrown.
 *
 * <p>Note that programmer written code is also able to take
 * advantage of calling this method in situations where there are
 * multiple sibling exceptions and only one can be propagated.
 *
 * @param exception the exception to be added to the list of
 *        suppressed exceptions
 * @throws IllegalArgumentException if {@code exception} is this
 *         throwable; a throwable cannot suppress itself.
 * @throws NullPointerException if {@code exception} is {@code null}
 * @since 1.7
 */
func (t *ThrowableObject) AddSuppressed(exception Throwable) {
	if exception == t {
		panic(NewIllegalArgumentException2(SelfSuppressionMessage, exception))
	}
	if exception == nil {
		panic(NewNullPointerException1(NullCauseMessage))
	}
	if t.suppressedExceptions == nil {
		return
	}
	t.suppressedExceptions = append(t.suppressedExceptions, exception)
}

/**
 * Returns an array containing all of the exceptions that were
 * suppressed, typically by the {@code try}-with-resources
 * statement, in order to deliver this exception.
 *
 * If no exceptions were suppressed or {@linkplain
 * #Throwable(String, Throwable, boolean, boolean) suppression is
 * disabled}, an empty array is returned.  This method is
 * thread-safe.  Writes to the returned array do not affect future
 * calls to this method.
 *
 * @return an array containing all of the exceptions that were
 *         suppressed to deliver this exception.
 * @since 1.7
 */
func (t *ThrowableObject) GetSuppressed() []Throwable {
	if reflect.DeepEqual(t.suppressedExceptions, SuppressedSentinel) || t.suppressedExceptions == nil {
		return EmptyThrowableArray
	}
	return t.suppressedExceptions
}
