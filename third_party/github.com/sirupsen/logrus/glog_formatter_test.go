// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logrus

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
	"testing"
	"time"

	time_ "github.com/searKing/golang/go/time"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ExampleNewGlogFormatter() {
	getPid = func() int {
		return 0
	}
	tf := NewGlogFormatter()
	entry := logrus.WithField("foo", "bar").
		WithError(fmt.Errorf("opps, an error occured"))
	entry.Message = "Hello World"
	b, _ := tf.Format(entry)

	fmt.Printf("%s", string(b))

	// Output:
	// P00010101 00:00:00.000000 0 run_example.go:64] Hello World, error="opps, an error occured", foo=bar
}

func ExampleNewGlogEnhancedFormatter() {
	getPid = func() int {
		return 0
	}
	tf := NewGlogEnhancedFormatter()
	entry := logrus.WithField("foo", "bar").
		WithError(fmt.Errorf("opps, an error occured"))
	entry.Message = "Hello World"
	b, _ := tf.Format(entry)
	fmt.Printf("%s", string(b))

	// Output:
	// [PANIC] [00010101 00:00:00.000000] [0] [run_example.go:64] Hello World, error=opps, an error occured, foo=bar
}

func TestFormatting(t *testing.T) {
	getPid = func() int {
		return 0
	}

	tf := &GlogFormatter{
		DisableColors: true,
	}

	testCases := []struct {
		value    string
		expected string
	}{
		{`foo`, "P00010101 00:00:00.000000 0 testing.go:1259] , test=foo\n"},
	}

	for i, tc := range testCases {
		b, _ := tf.Format(logrus.WithField("test", tc.value))

		if string(b) != tc.expected {
			t.Errorf("#%d: formatting expected for %q, got %q; expected %q", i, tc.value, string(b), tc.expected)
		}
	}
}

func TestQuoting(t *testing.T) {
	getPid = func() int {
		return 0
	}

	tf := &GlogFormatter{DisableColors: true}

	checkQuoting := func(q bool, value interface{}) {
		b, _ := tf.Format(logrus.WithField("test", value))
		idx := bytes.Index(b, ([]byte)("test="))
		cont := bytes.Contains(b[idx+5:], []byte("\""))
		if cont != q {
			if q {
				t.Errorf("quoting expected for: %#v", value)
			} else {
				t.Errorf("quoting not expected for: %#v", value)
			}
		}
	}

	checkQuoting(false, "")
	checkQuoting(false, "abcd")
	checkQuoting(false, "v1.0")
	checkQuoting(false, "1234567890")
	checkQuoting(false, "/foobar")
	checkQuoting(false, "foo_bar")
	checkQuoting(false, "foo@bar")
	checkQuoting(false, "foobar^")
	checkQuoting(false, "+/-_^@f.oobar")
	checkQuoting(true, "foo\n\rbar")
	checkQuoting(true, "foobar$")
	checkQuoting(true, "&foobar")
	checkQuoting(true, "x y")
	checkQuoting(true, "x,y")
	checkQuoting(false, errors.New("invalid"))
	checkQuoting(true, errors.New("invalid argument"))

	// Test for quoting empty fields.
	tf.QuoteEmptyFields = true
	checkQuoting(true, "")
	checkQuoting(false, "abcd")
	checkQuoting(true, "foo\n\rbar")
	checkQuoting(true, errors.New("invalid argument"))

	// Test forcing quotes.
	tf.ForceQuote = true
	checkQuoting(true, "")
	checkQuoting(true, "abcd")
	checkQuoting(true, "foo\n\rbar")
	checkQuoting(true, errors.New("invalid argument"))

	// Test forcing quotes when also disabling them.
	tf.DisableQuote = true
	checkQuoting(true, "")
	checkQuoting(true, "abcd")
	checkQuoting(true, "foo\n\rbar")
	checkQuoting(true, errors.New("invalid argument"))

	// Test disabling quotes
	tf.ForceQuote = false
	tf.QuoteEmptyFields = false
	checkQuoting(false, "")
	checkQuoting(false, "abcd")
	checkQuoting(false, "foo\n\rbar")
	checkQuoting(false, errors.New("invalid argument"))
}

func TestEscaping(t *testing.T) {
	tf := &GlogFormatter{DisableColors: true}

	testCases := []struct {
		value    string
		expected string
	}{
		{`ba"r`, `ba\"r`},
		{`ba'r`, `ba'r`},
	}

	for _, tc := range testCases {
		b, _ := tf.Format(logrus.WithField("test", tc.value))
		if !bytes.Contains(b, []byte(tc.expected)) {
			t.Errorf("escaping expected for %q (result was %q instead of %q)", tc.value, string(b), tc.expected)
		}
	}
}

func TestEscaping_Interface(t *testing.T) {
	tf := &GlogFormatter{DisableColors: true}

	ts := time.Now()

	testCases := []struct {
		value    interface{}
		expected string
	}{
		{ts, fmt.Sprintf("\"%s\"", ts.String())},
		{errors.New("error: something went wrong"), "\"error: something went wrong\""},
	}

	for _, tc := range testCases {
		b, _ := tf.Format(logrus.WithField("test", tc.value))
		if !bytes.Contains(b, []byte(tc.expected)) {
			t.Errorf("escaping expected for %q (result was %q instead of %q)", tc.value, string(b), tc.expected)
		}
	}
}

func TestTimestampFormat(t *testing.T) {
	getPid = func() int {
		return 0
	}
	checkTimeStr := func(format string) {
		customFormatter := &GlogFormatter{DisableColors: true, TimestampFormat: format}
		customStr, _ := customFormatter.Format(logrus.WithField("test", "test"))
		timeStart := bytes.Index(customStr, ([]byte)("P"))
		timeEnd := bytes.Index(customStr, ([]byte)(" 0 glog_formatter_test.go"))
		timeStr := customStr[timeStart+1 : timeEnd]
		if format == "" {
			format = time_.GLogDate
		}
		_, e := time.Parse(format, strings.TrimSpace(string(timeStr)))
		if e != nil {
			t.Errorf("time string %q did not match provided time format %q: %s", timeStr, format, e)
		}
	}

	checkTimeStr("2006-01-02T15:04:05.000000000Z07:00")
	checkTimeStr("Mon Jan _2 15:04:05 2006")
	checkTimeStr("")
}

func TestDisableLevelTruncation(t *testing.T) {
	entry := &logrus.Entry{
		Time:    time.Now(),
		Message: "testing",
	}
	timestampFormat := "Mon Jan 2 15:04:05 -0700 MST 2006"
	checkDisableTruncation := func(disabled bool, level logrus.Level) {
		tf := &GlogFormatter{LevelTruncationLimit: 5}
		tf.TimestampFormat = timestampFormat
		var b bytes.Buffer
		entry.Level = level
		buf, _ := tf.Format(entry)
		b.Write(buf)
		logLine := (&b).String()
		if disabled {
			expected := strings.ToUpper(level.String())
			if !strings.Contains(logLine, expected) {
				t.Errorf("level string expected to be %s when truncation disabled", expected)
			}
		} else {
			expected := strings.ToUpper(level.String())
			if len(level.String()) > 5 {
				if strings.Contains(logLine, expected) {
					t.Errorf("level string %s expected to be truncated to %s when truncation is enabled", expected, expected[0:4])
				}
			} else {
				if !strings.Contains(logLine, expected) {
					t.Errorf("level string expected to be %s when truncation is enabled and level string is below truncation threshold", expected)
				}
			}
		}
	}

	checkDisableTruncation(true, logrus.DebugLevel)
	checkDisableTruncation(true, logrus.InfoLevel)
	checkDisableTruncation(false, logrus.ErrorLevel)
	checkDisableTruncation(false, logrus.InfoLevel)
}

func TestPadLevelText(t *testing.T) {
	// A note for future maintainers / committers:
	//
	// This test denormalizes the level text as a part of its assertions.
	// Because of that, its not really a "unit test" of the PadLevelText functionality.
	// So! Many apologies to the potential future person who has to rewrite this test
	// when they are changing some completely unrelated functionality.
	params := []struct {
		name            string
		level           logrus.Level
		paddedLevelText string
	}{
		{
			name:            "PanicLevel",
			level:           logrus.PanicLevel,
			paddedLevelText: "PANIC", // 2 extra spaces
		},
		{
			name:            "FatalLevel",
			level:           logrus.FatalLevel,
			paddedLevelText: "FATAL", // 2 extra spaces
		},
		{
			name:            "ErrorLevel",
			level:           logrus.ErrorLevel,
			paddedLevelText: "ERROR", // 2 extra spaces
		},
		{
			name:  "WarnLevel",
			level: logrus.WarnLevel,
			// WARNING is already the max length, so we don't need to assert a paddedLevelText
		},
		{
			name:            "DebugLevel",
			level:           logrus.DebugLevel,
			paddedLevelText: "DEBUG", // 2 extra spaces
		},
		{
			name:            "TraceLevel",
			level:           logrus.TraceLevel,
			paddedLevelText: "TRACE", // 2 extra spaces
		},
		{
			name:            "InfoLevel",
			level:           logrus.InfoLevel,
			paddedLevelText: "INFO", // 3 extra spaces
		},
	}

	// We create a "default" GlogFormatter to do a control test.
	// We also create a GlogFormatter with PadLevelText, which is the parameter we want to do our most relevant assertions against.
	tfDefault := GlogFormatter{}
	tfWithPadding := GlogFormatter{HumanReadable: true, PadLevelText: true, LevelTruncationLimit: -1}

	for _, val := range params {
		t.Run(val.name, func(t *testing.T) {
			// GlogFormatter writes into these bytes.Buffers, and we make assertions about their contents later
			var bytesDefault bytes.Buffer
			var bytesWithPadding bytes.Buffer

			// The GlogFormatter instance and the bytes.Buffer instance are different here
			// all the other arguments are the same. We also initialize them so that they
			// fill in the value of levelTextMaxLength.
			tfDefault.init(&logrus.Entry{})
			b1, _ := tfDefault.Format(&logrus.Entry{Level: val.level})
			bytesDefault.Write(b1)
			tfWithPadding.init(&logrus.Entry{})
			b2, _ := tfWithPadding.Format(&logrus.Entry{Level: val.level})
			bytesWithPadding.Write(b2)
			// turn the bytes back into a string so that we can actually work with the data
			logLineDefault := (&bytesDefault).String()
			logLineWithPadding := (&bytesWithPadding).String()

			// Control: the level text should not be padded by default
			if val.paddedLevelText != "" && strings.Contains(logLineDefault, val.paddedLevelText) {
				t.Errorf("log line %q should not contain the padded level text %q by default", logLineDefault, val.paddedLevelText)
			}

			// Assertion: the level text should still contain the string representation of the level
			if !strings.Contains(strings.ToLower(logLineWithPadding), val.level.String()) {
				t.Errorf("log line %q should contain the level text %q when padding is enabled", logLineWithPadding, val.level.String())
			}

			// Assertion: the level text should be in its padded form now
			if val.paddedLevelText != "" && !strings.Contains(logLineWithPadding, val.paddedLevelText) {
				t.Errorf("log line %q should contain the padded level text %q when padding is enabled", logLineWithPadding, val.paddedLevelText)
			}

		})
	}
}

func TestDisableTimestampWithColoredOutput(t *testing.T) {
	tf := &GlogFormatter{DisableTimestamp: true, ForceColors: true}

	b, _ := tf.Format(logrus.WithField("test", "test"))
	if strings.Contains(string(b), "[0000]") {
		t.Error("timestamp not expected when DisableTimestamp is true")
	}
}

func TestNewlineBehavior(t *testing.T) {
	tf := &GlogFormatter{ForceColors: true}

	// Ensure a single new line is removed as per stdlib log
	e := logrus.NewEntry(logrus.StandardLogger())
	e.Message = "test message\n"
	b, _ := tf.Format(e)
	if bytes.Contains(b, []byte("test message\n")) {
		t.Error("first newline at end of Entry.Message resulted in unexpected 2 newlines in output. Expected newline to be removed.")
	}

	// Ensure a double new line is reduced to a single new line
	e = logrus.NewEntry(logrus.StandardLogger())
	e.Message = "test message\n\n"
	b, _ = tf.Format(e)
	if bytes.Contains(b, []byte("test message\n\n")) {
		t.Error("Double newline at end of Entry.Message resulted in unexpected 2 newlines in output. Expected single newline")
	}
	if !bytes.Contains(b, []byte("test message\n")) {
		t.Error("Double newline at end of Entry.Message did not result in a single newline after formatting")
	}
}

func TestGlogFormatterFieldMap(t *testing.T) {
	formatter := &GlogFormatter{
		DisableColors: true,
		FieldMap: FieldMap{
			FieldKeyMsg:   "message",
			FieldKeyLevel: "somelevel",
			FieldKeyTime:  "timeywimey",
		},
	}

	entry := &logrus.Entry{
		Message: "oh hi",
		Level:   logrus.WarnLevel,
		Time:    time.Date(1981, time.February, 24, 4, 28, 3, 100, time.UTC),
		Data: logrus.Fields{
			"field1":     "f1",
			"message":    "messagefield",
			"somelevel":  "levelfield",
			"timeywimey": "timeywimeyfield",
		},
	}

	b, err := formatter.Format(entry)
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	assert.Equal(t,
		`W19810224 04:28:03.000000 0 testing.go:1259] `+
			`oh hi, `+
			`field1=f1, `+
			`fields.message=messagefield, `+
			`fields.somelevel=levelfield, `+
			`fields.timeywimey=timeywimeyfield`+"\n",
		string(b),
		"Formatted output doesn't respect FieldMap")
}

func TestGlogFormatterIsColored(t *testing.T) {
	params := []struct {
		name               string
		expectedResult     bool
		isTerminal         bool
		disableColor       bool
		forceColor         bool
		envColor           bool
		clicolorIsSet      bool
		clicolorForceIsSet bool
		clicolorVal        string
		clicolorForceVal   string
	}{
		// Default values
		{
			name:               "testcase1",
			expectedResult:     false,
			isTerminal:         false,
			disableColor:       false,
			forceColor:         false,
			envColor:           false,
			clicolorIsSet:      false,
			clicolorForceIsSet: false,
		},
		// Output on terminal
		{
			name:               "testcase2",
			expectedResult:     true,
			isTerminal:         true,
			disableColor:       false,
			forceColor:         false,
			envColor:           false,
			clicolorIsSet:      false,
			clicolorForceIsSet: false,
		},
		// Output on terminal with color disabled
		{
			name:               "testcase3",
			expectedResult:     false,
			isTerminal:         true,
			disableColor:       true,
			forceColor:         false,
			envColor:           false,
			clicolorIsSet:      false,
			clicolorForceIsSet: false,
		},
		// Output not on terminal with color disabled
		{
			name:               "testcase4",
			expectedResult:     false,
			isTerminal:         false,
			disableColor:       true,
			forceColor:         false,
			envColor:           false,
			clicolorIsSet:      false,
			clicolorForceIsSet: false,
		},
		// Output not on terminal with color forced
		{
			name:               "testcase5",
			expectedResult:     true,
			isTerminal:         false,
			disableColor:       false,
			forceColor:         true,
			envColor:           false,
			clicolorIsSet:      false,
			clicolorForceIsSet: false,
		},
		// Output on terminal with clicolor set to "0"
		{
			name:               "testcase6",
			expectedResult:     false,
			isTerminal:         true,
			disableColor:       false,
			forceColor:         false,
			envColor:           true,
			clicolorIsSet:      true,
			clicolorForceIsSet: false,
			clicolorVal:        "0",
		},
		// Output on terminal with clicolor set to "1"
		{
			name:               "testcase7",
			expectedResult:     true,
			isTerminal:         true,
			disableColor:       false,
			forceColor:         false,
			envColor:           true,
			clicolorIsSet:      true,
			clicolorForceIsSet: false,
			clicolorVal:        "1",
		},
		// Output not on terminal with clicolor set to "0"
		{
			name:               "testcase8",
			expectedResult:     false,
			isTerminal:         false,
			disableColor:       false,
			forceColor:         false,
			envColor:           true,
			clicolorIsSet:      true,
			clicolorForceIsSet: false,
			clicolorVal:        "0",
		},
		// Output not on terminal with clicolor set to "1"
		{
			name:               "testcase9",
			expectedResult:     false,
			isTerminal:         false,
			disableColor:       false,
			forceColor:         false,
			envColor:           true,
			clicolorIsSet:      true,
			clicolorForceIsSet: false,
			clicolorVal:        "1",
		},
		// Output not on terminal with clicolor set to "1" and force color
		{
			name:               "testcase10",
			expectedResult:     true,
			isTerminal:         false,
			disableColor:       false,
			forceColor:         true,
			envColor:           true,
			clicolorIsSet:      true,
			clicolorForceIsSet: false,
			clicolorVal:        "1",
		},
		// Output not on terminal with clicolor set to "0" and force color
		{
			name:               "testcase11",
			expectedResult:     false,
			isTerminal:         false,
			disableColor:       false,
			forceColor:         true,
			envColor:           true,
			clicolorIsSet:      true,
			clicolorForceIsSet: false,
			clicolorVal:        "0",
		},
		// Output not on terminal with clicolor_force set to "1"
		{
			name:               "testcase12",
			expectedResult:     true,
			isTerminal:         false,
			disableColor:       false,
			forceColor:         false,
			envColor:           true,
			clicolorIsSet:      false,
			clicolorForceIsSet: true,
			clicolorForceVal:   "1",
		},
		// Output not on terminal with clicolor_force set to "0"
		{
			name:               "testcase13",
			expectedResult:     false,
			isTerminal:         false,
			disableColor:       false,
			forceColor:         false,
			envColor:           true,
			clicolorIsSet:      false,
			clicolorForceIsSet: true,
			clicolorForceVal:   "0",
		},
		// Output on terminal with clicolor_force set to "0"
		{
			name:               "testcase14",
			expectedResult:     false,
			isTerminal:         true,
			disableColor:       false,
			forceColor:         false,
			envColor:           true,
			clicolorIsSet:      false,
			clicolorForceIsSet: true,
			clicolorForceVal:   "0",
		},
	}

	cleanenv := func() {
		os.Unsetenv("CLICOLOR")
		os.Unsetenv("CLICOLOR_FORCE")
	}

	defer cleanenv()

	for _, val := range params {
		t.Run("textformatter_"+val.name, func(subT *testing.T) {
			tf := GlogFormatter{
				isTerminal:                val.isTerminal,
				DisableColors:             val.disableColor,
				ForceColors:               val.forceColor,
				EnvironmentOverrideColors: val.envColor,
			}
			cleanenv()
			if val.clicolorIsSet {
				os.Setenv("CLICOLOR", val.clicolorVal)
			}
			if val.clicolorForceIsSet {
				os.Setenv("CLICOLOR_FORCE", val.clicolorForceVal)
			}
			res := tf.isColored()
			if runtime.GOOS == "windows" && !tf.ForceColors && !val.clicolorForceIsSet {
				assert.Equal(subT, false, res)
			} else {
				assert.Equal(subT, val.expectedResult, res)
			}
		})
	}
}

func TestCustomSorting(t *testing.T) {
	formatter := &GlogFormatter{
		DisableColors: true,
		SortingFunc: func(keys []string) {
			sort.Slice(keys, func(i, j int) bool {
				if keys[j] == "prefix" {
					return false
				}
				if keys[i] == "prefix" {
					return true
				}
				return strings.Compare(keys[i], keys[j]) == -1
			})
		},
	}

	entry := &logrus.Entry{
		Message: "Testing custom sort function",
		Time:    time.Now(),
		Level:   logrus.InfoLevel,
		Data: logrus.Fields{
			"test":      "testvalue",
			"prefix":    "the application prefix",
			"blablabla": "blablabla",
		},
	}
	b, err := formatter.Format(entry)
	require.NoError(t, err)
	require.True(t, strings.Contains(string(b), "prefix="), "format output is %q", string(b))
}
