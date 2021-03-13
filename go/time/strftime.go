// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time

import (
	"strings"
)

func LayoutTimeToStrftime(layout string) string {
	return layoutTimeToStrftime(layout, false)
}

func LayoutTimeToSimilarStrftime(layout string) string {
	return layoutTimeToStrftime(layout, true)
}

func LayoutStrftimeToTime(layout string) string {
	return layoutStrftimeToTime(layout, false)
}

func LayoutStrftimeToSimilarTime(layout string) string {
	return layoutStrftimeToTime(layout, true)
}

func layoutTimeToStrftime(layout string, similar bool) string {
	var buf strings.Builder
	// Each iteration generates one std value.
	for layout != "" {
		prefix, std_, suffix := nextStdChunk(layout)
		if prefix != "" {
			buf.WriteString(prefix)
		}
		if std_ == 0 {
			break
		}

		if similar {
			buf.WriteString(std(std_).SimilarStrftimeString())
		} else {
			buf.WriteString(std(std_).StrftimeString())
		}
		layout = suffix
	}
	return buf.String()
}

func layoutStrftimeToTime(layout string, similar bool) string {
	var buf strings.Builder
	// Each iteration generates one std value.
	for layout != "" {
		prefix, std_, suffix := nextStrftimeChunk(layout)
		if prefix != "" {
			buf.WriteString(prefix)
		}
		if std_ == 0 {
			break
		}

		if similar {
			buf.WriteString(std(std_).SimilarString())
		} else {
			buf.WriteString(std(std_).String())
		}
		layout = suffix
	}
	return buf.String()
}

// nextStdChunk finds the first occurrence of a std string in
// layout and returns the text before, the std string, and the text after.
func nextStrftimeChunk(layout string) (prefix string, std int, suffix string) {
	for i := 0; i < len(layout); i++ {
		switch c := int(layout[i]); c {
		case '%': // %
			j := i + 1
			switch c := int(layout[j]); c {
			case 'a': // Mon
				return layout[:i], stdWeekDay, layout[j+1:]
			case 'A': // Monday
				return layout[0:i], stdLongWeekDay, layout[j+1:]
			case 'b', 'h': // Jan
				return layout[0:i], stdMonth, layout[j+1:]
			case 'B': // January
				return layout[0:i], stdLongMonth, layout[j+1:]
			case 'c': // "Mon Jan _2 15:04:05 2006" (assumes "C" locale)
				return layout[0:i], stdDateAndTime, layout[j+1:]
			case 'C': // 20
				return layout[0:i], stdFirstTwoDigitYear, layout[j+1:]
			case 'd': // 02
				return layout[0:i], stdZeroDay, layout[j+1:]
			case 'D', 'x': // %m/%d/%y
				return layout[0:i], stdShortSlashDate, layout[j+1:]
			case 'e': // _2
				return layout[0:i], stdUnderDay, layout[j+1:]
			case 'E': // E modifier is to use a locale-dependent alternative representation
				// %Ec, %EC, %Ex, %EX, %Ey, %EY
				k := j + 1
				switch c := int(layout[k]); c {
				case 'c': // "Mon Jan _2 15:04:05 2006" (assumes "C" locale)
					return layout[0:i], stdDateAndTime & stdNeedEModifier, layout[k+1:]
				case 'C': // 20
					return layout[0:i], stdFirstTwoDigitYear & stdNeedEModifier, layout[k+1:]
				case 'x': // %m/%d/%y
					return layout[0:i], stdShortSlashDate & stdNeedEModifier, layout[k+1:]
				case 'X': // locale depended time representation (assumes "C" locale)
					return layout[0:i], stdHourClockTime & stdNeedEModifier, layout[k+1:]
				case 'y':
					return layout[0:i], stdYear & stdNeedEModifier, layout[k+1:]
				case 'Y':
					return layout[0:i], stdLongYear & stdNeedEModifier, layout[k+1:]
				}
			case 'f': // fraction seconds in microseconds (Python)
				std = stdFracSecond0
				std |= 6 << stdArgShift // microseconds precision
				return layout[0:i], std, layout[j+1:]
			case 'F': // %Y-%m-%d
				return layout[0:i], stdShortDashDate, layout[j+1:]
			case 'g':
				return layout[0:i], stdISO8601WeekYear, layout[j+1:]
			case 'G':
				return layout[0:i], stdISO8601LongWeekYear, layout[j+1:]
			case 'H', 'k':
				return layout[0:i], stdHour, layout[j+1:]
			case 'I', 'l':
				return layout[0:i], stdZeroHour12, layout[j+1:]
			case 'j':
				return layout[0:i], stdDayOfYear, layout[j+1:]
			case 'm':
				return layout[0:i], stdZeroMonth, layout[j+1:]
			case 'M':
				return layout[0:i], stdZeroMinute, layout[j+1:]
			case 'n':
				return layout[0:i], stdCharNewLine, layout[j+1:]
			case 'O': // O modifier is to use alternative numeric symbols (say, roman numerals)
				// %Od %Oe %OH %OI %Om %OM %OS %Ou %OU %OV %Ow %OW %Oy
				k := j + 1
				switch c := int(layout[k]); c {
				case 'd': // 02
					return layout[0:i], stdZeroDay & stdNeedOModifier, layout[k+1:]
				case 'e': // _2
					return layout[0:i], stdUnderDay & stdNeedOModifier, layout[k+1:]
				case 'H':
					return layout[0:i], stdHour & stdNeedOModifier, layout[k+1:]
				case 'I':
					return layout[0:i], stdZeroHour12 & stdNeedOModifier, layout[k+1:]
				case 'm':
					return layout[0:i], stdZeroMonth & stdNeedOModifier, layout[k+1:]
				case 'M':
					return layout[0:i], stdZeroMinute & stdNeedOModifier, layout[k+1:]
				case 'S':
					return layout[0:i], stdZeroSecond & stdNeedOModifier, layout[k+1:]
				case 'u': // weekday as a decimal number, where Monday is 1
					return layout[0:i], stdNumWeekDay & stdNeedOModifier, layout[k+1:]
				case 'U': // week of the year as a decimal number (Sunday is the first day of the week)
					return layout[0:i], stdSundayFirstWeekOfYear & stdNeedOModifier, layout[k+1:]
				case 'V':
					return layout[0:i], stdISO8601Week & stdNeedOModifier, layout[k+1:]
				case 'w':
					return layout[0:i], stdZeroNumWeek & stdNeedOModifier, layout[k+1:]
				case 'W': // week of the year as a decimal number (Monday is the first day of the week)
					return layout[0:i], stdMonFirstWeekOfYear & stdNeedOModifier, layout[k+1:]
				case 'y':
					return layout[0:i], stdYear & stdNeedOModifier, layout[k+1:]
				}
			case 'p':
				return layout[0:i], stdPM, layout[j+1:]
			case 'P':
				return layout[0:i], stdpm, layout[j+1:]
			case 'r': // "%I:%M:%S %p"
				return layout[0:i], stdHour12ClockTime, layout[j+1:]
			case 'R': // %H:%M"
				return layout[0:i], stdHourHourMinuteTime, layout[j+1:]
			case 's':
				return layout[0:i], stdSecondsSinceEpoch, layout[j+1:]
			case 'S':
				return layout[0:i], stdZeroSecond, layout[j+1:]
			case 't':
				return layout[0:i], stdCharHorizontalTab, layout[j+1:]
			case 'T': // %H:%M:%S
				return layout[0:i], stdISO8601Time, layout[j+1:]
			case 'u': // weekday as a decimal number, where Monday is 1
				return layout[0:i], stdNumWeekDay, layout[j+1:]
			case 'U': // week of the year as a decimal number (Sunday is the first day of the week)
				return layout[0:i], stdSundayFirstWeekOfYear, layout[j+1:]
			case 'v':
				return layout[0:i], stdISO8601NumWeek, layout[j+1:]
			case 'V':
				return layout[0:i], stdISO8601Week, layout[j+1:]
			case 'w':
				return layout[0:i], stdZeroNumWeek, layout[j+1:]
			case 'W': // week of the year as a decimal number (Monday is the first day of the week)
				return layout[0:i], stdMonFirstWeekOfYear, layout[j+1:]
			case 'X': // locale depended time representation (assumes "C" locale)
				return layout[0:i], stdHourClockTime, layout[j+1:]
			case 'y':
				return layout[0:i], stdYear, layout[j+1:]
			case 'Y':
				return layout[0:i], stdLongYear, layout[j+1:]
			case 'z':
				return layout[0:i], stdISO8601ColonTZ, layout[j+1:]
			case 'Z':
				return layout[0:i], stdTZ, layout[j+1:]
			case '%':
				return layout[0:i], stdCharPercentSign, layout[j+1:]
			}
		}
	}
	return layout, 0, ""
}
