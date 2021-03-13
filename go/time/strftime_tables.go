// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time

import "time"

type std int

const (
	_                        = iota
	stdLongMonth             = iota + stdNeedDate  // "January"
	stdMonth                                       // "Jan"
	stdNumMonth                                    // "1"
	stdZeroMonth                                   // "01"
	stdLongWeekDay                                 // "Monday"
	stdWeekDay                                     // "Mon"
	stdDay                                         // "2"
	stdUnderDay                                    // "_2"
	stdZeroDay                                     // "02"
	stdUnderYearDay                                // "__2"
	stdZeroYearDay                                 // "002"
	stdHour                  = iota + stdNeedClock // "15"
	stdHour12                                      // "3"
	stdZeroHour12                                  // "03"
	stdMinute                                      // "4"
	stdZeroMinute                                  // "04"
	stdSecond                                      // "5"
	stdZeroSecond                                  // "05"
	stdLongYear              = iota + stdNeedDate  // "2006"
	stdYear                                        // "06"
	stdPM                    = iota + stdNeedClock // "PM"
	stdpm                                          // "pm"
	stdTZ                    = iota                // "MST"
	stdISO8601TZ                                   // "Z0700"  // prints Z for UTC
	stdISO8601SecondsTZ                            // "Z070000"
	stdISO8601ShortTZ                              // "Z07"
	stdISO8601ColonTZ                              // "Z07:00" // prints Z for UTC
	stdISO8601ColonSecondsTZ                       // "Z07:00:00"
	stdNumTZ                                       // "-0700"  // always numeric
	stdNumSecondsTz                                // "-070000"
	stdNumShortTZ                                  // "-07"    // always numeric
	stdNumColonTZ                                  // "-07:00" // always numeric
	stdNumColonSecondsTZ                           // "-07:00:00"
	stdFracSecond0                                 // ".0", ".00", ... , trailing zeros included
	stdFracSecond9                                 // ".9", ".99", ..., trailing zeros omitted

	// extend for strftime
	stdNop
	stdCharNewLine                                          // New-line character ('\n')
	stdCharHorizontalTab                                    // Horizontal-tab character ('\t')
	stdCharPercentSign                                      // A % sign	 ('%')
	stdZeroNumWeek                                          // numerical week representation (0 - Sunday ~ 6 - Saturday)
	stdNumWeekDay                                           // numerical week representation (1 - Monday ~ 7- Sunday)
	stdSundayFirstWeekOfYear                                // week of the year (Sunday first)
	stdMonFirstWeekOfYear                                   // week of the year (Monday first)
	stdFirstTwoDigitYear                                    // "20"
	stdSecondsSinceEpoch                                    // The number of seconds since the Epoch, 1970-01-01 00:00:00 +0000 (UTC). (TZ)
	stdDayOfYear                                            // day of the year (range [001,366])
	stdDateAndTime                                          // Date and time representation; Thu Aug 23 14:55:02 2001
	stdShortSlashDate                                       // Short MM/DD/YY date, equivalent to %m/%d/%y; 08/23/01
	stdShortDashDate                                        // Short YYYY-MM-DD date, equivalent to %Y-%m-%d; 2001-08-23
	stdHour12ClockTime                                      // 12-hour clock time; 02:55:02 pm
	stdHourClockTime                                        // Time representation; 14:55:02
	stdHourHourMinuteTime                                   // 24-hour HH:MM time, equivalent to %H:%M; 14:55
	stdISO8601WeekYear       = iota + stdNeedISOISO8601Week // last two digits of ISO 8601 week-based year
	stdISO8601LongWeekYear                                  // ISO 8601 week-based year
	stdISO8601Week                                          // ISO 8601 week
	stdISO8601NumWeek                                       // ISO 8601 week number (01-53); 34
	stdISO8601Time                                          // ISO 8601 time format (HH:MM:SS), equivalent to %H:%M:%S; 14:55:02

	stdNeedDate           = 1 << 8             // need month, day, year
	stdNeedClock          = 2 << 8             // need hour, minute, second
	stdNeedISOISO8601Week = 4 << 8             // need ISO8601 week and year
	stdNeedEModifier      = 8 << 8             // need to use alternative numeric symbols (say, roman numerals) %Ec, %EC, %Ex, %EX, %Ey, %EY
	stdNeedOModifier      = 16 << 8            // need to use a locale-dependent alternative representation %Od, %Oe, %OH, %OI, %Om, %OM, %OS, %Ou, %OU, %OV, %Ow, %OW, %Oy
	stdArgShift           = 16                 // extra argument in high bits, above low stdArgShift
	stdMask               = 1<<stdArgShift - 1 // mask out argument

)

func (s std) String() string {
	s_ := int(s & stdMask)
	label, has := stdChunkNames[s_]
	if has {
		return label
	}
	return strftimeChunkNames[s_]
}

func (s std) SimilarString() string {
	s_ := int(s)
	label, has := stdChunkNames[s_]
	if has {
		return label
	}
	label, has = stdSimilarChunkNames[s_]
	if has {
		return label
	}
	return strftimeChunkNames[s_]
}

func (s std) StrftimeString() string {
	s_ := int(s & stdMask)
	label, has := strftimeChunkNames[s_]
	if has {
		return label
	}
	return stdChunkNames[s_]
}

func (s std) SimilarStrftimeString() string {
	s_ := int(s)
	label, has := strftimeChunkNames[s_]
	if has {
		return label
	}
	label, has = strftimeSimilarChunkNames[s_]
	if has {
		return label
	}
	return stdChunkNames[s_]
}

// stdChunkNames maps from nextStdChunk results to the matched strings.
var stdChunkNames = map[int]string{
	0:                               "",
	stdLongMonth:                    "January",
	stdMonth:                        "Jan",
	stdNumMonth:                     "1",
	stdZeroMonth:                    "01",
	stdLongWeekDay:                  "Monday",
	stdWeekDay:                      "Mon",
	stdDay:                          "2",
	stdUnderDay:                     "_2",
	stdZeroDay:                      "02",
	stdUnderYearDay:                 "__2",
	stdZeroYearDay:                  "002",
	stdHour:                         "15",
	stdHour12:                       "3",
	stdZeroHour12:                   "03",
	stdMinute:                       "4",
	stdZeroMinute:                   "04",
	stdSecond:                       "5",
	stdZeroSecond:                   "05",
	stdLongYear:                     "2006",
	stdYear:                         "06",
	stdPM:                           "PM",
	stdpm:                           "pm",
	stdTZ:                           "MST",
	stdISO8601TZ:                    "Z0700",
	stdISO8601SecondsTZ:             "Z070000",
	stdISO8601ShortTZ:               "Z07",
	stdISO8601ColonTZ:               "Z07:00",
	stdISO8601ColonSecondsTZ:        "Z07:00:00",
	stdNumTZ:                        "-0700",
	stdNumSecondsTz:                 "-070000",
	stdNumShortTZ:                   "-07",
	stdNumColonTZ:                   "-07:00",
	stdNumColonSecondsTZ:            "-07:00:00",
	stdFracSecond0 | 1<<stdArgShift: ".0",
	stdFracSecond0 | 2<<stdArgShift: ".00",
	stdFracSecond0 | 3<<stdArgShift: ".000",
	stdFracSecond0 | 4<<stdArgShift: ".0000",
	stdFracSecond0 | 5<<stdArgShift: ".00000",
	stdFracSecond0 | 6<<stdArgShift: ".000000",
	stdFracSecond0 | 7<<stdArgShift: ".0000000",
	stdFracSecond0 | 8<<stdArgShift: ".00000000",
	stdFracSecond0 | 9<<stdArgShift: ".000000000",
	stdFracSecond9 | 1<<stdArgShift: ".9",
	stdFracSecond9 | 2<<stdArgShift: ".99",
	stdFracSecond9 | 3<<stdArgShift: ".999",
	stdFracSecond9 | 4<<stdArgShift: ".9999",
	stdFracSecond9 | 5<<stdArgShift: ".99999",
	stdFracSecond9 | 6<<stdArgShift: ".999999",
	stdFracSecond9 | 7<<stdArgShift: ".9999999",
	stdFracSecond9 | 8<<stdArgShift: ".99999999",
	stdFracSecond9 | 9<<stdArgShift: ".999999999",
	stdNop:                          "",
	stdCharNewLine:                  "\n",
	stdCharHorizontalTab:            "\t",
	stdCharPercentSign:              "%",
	stdDateAndTime:                  time.ANSIC,
	stdShortSlashDate:               "2006/01/02",
	stdShortDashDate:                "2006-01-02",
	stdHour12ClockTime:              "03:04:05 pm",
	stdHourClockTime:                "15:04:05",
	stdHourHourMinuteTime:           "15:04",
	stdISO8601Time:                  "15:04:05",
}

var stdSimilarChunkNames = map[int]string{
	0:                        "",
	stdNop:                   "",
	stdZeroNumWeek:           stdChunkNames[stdWeekDay],
	stdNumWeekDay:            stdChunkNames[stdWeekDay],
	stdSundayFirstWeekOfYear: stdChunkNames[stdWeekDay],
	stdMonFirstWeekOfYear:    stdChunkNames[stdWeekDay],
	stdFirstTwoDigitYear:     stdChunkNames[stdLongYear],
	stdSecondsSinceEpoch:     strftimeChunkNames[stdSecondsSinceEpoch],
	stdDayOfYear:             strftimeChunkNames[stdDayOfYear],
	stdISO8601WeekYear:       stdChunkNames[stdYear],
	stdISO8601LongWeekYear:   stdChunkNames[stdLongYear],
	stdISO8601Week:           strftimeChunkNames[stdISO8601Week],
	stdISO8601NumWeek:        strftimeChunkNames[stdISO8601NumWeek],

	stdDateAndTime | stdNeedEModifier:           stdChunkNames[stdDateAndTime],
	stdFirstTwoDigitYear | stdNeedEModifier:     stdChunkNames[stdLongYear],
	stdShortSlashDate | stdNeedEModifier:        stdChunkNames[stdShortSlashDate],
	stdHourClockTime | stdNeedEModifier:         stdChunkNames[stdHourClockTime],
	stdYear | stdNeedEModifier:                  stdChunkNames[stdYear],
	stdLongYear | stdNeedEModifier:              stdChunkNames[stdLongYear],
	stdZeroDay | stdNeedOModifier:               stdChunkNames[stdZeroDay],
	stdUnderDay | stdNeedOModifier:              stdChunkNames[stdUnderDay],
	stdHour | stdNeedOModifier:                  stdChunkNames[stdHour],
	stdZeroHour12 | stdNeedOModifier:            stdChunkNames[stdZeroHour12],
	stdZeroMonth | stdNeedOModifier:             stdChunkNames[stdZeroMonth],
	stdZeroMinute | stdNeedOModifier:            stdChunkNames[stdZeroMinute],
	stdZeroSecond | stdNeedOModifier:            stdChunkNames[stdZeroSecond],
	stdNumWeekDay | stdNeedOModifier:            stdChunkNames[stdWeekDay],
	stdSundayFirstWeekOfYear | stdNeedOModifier: stdChunkNames[stdWeekDay],
	stdISO8601Week | stdNeedOModifier:           strftimeChunkNames[stdISO8601Week],
	stdZeroNumWeek | stdNeedOModifier:           stdChunkNames[stdWeekDay],
	stdMonFirstWeekOfYear | stdNeedOModifier:    stdChunkNames[stdWeekDay],
	stdYear | stdNeedOModifier:                  stdChunkNames[stdYear],
}

// strftimeChunkNames maps from nextStrftimeChunk results to the matched strings.
// see http://www.cplusplus.com/reference/ctime/strftime/
// see https://man7.org/linux/man-pages/man3/strftime.3.html
var strftimeChunkNames = map[int]string{
	0:                        "",
	stdWeekDay:               "%a", // Thu; Abbreviated weekday name
	stdLongWeekDay:           "%A", // Thursday; Full weekday name
	stdMonth:                 "%b", // Aug; Abbreviated month name, (same as %h)
	stdLongMonth:             "%B", // August; Full month name
	stdDateAndTime:           "%c", // Thu Aug 23 14:55:02 2001; Date and time representation， %a %b %e %H:%M:%S %Y
	stdFirstTwoDigitYear:     "%C", // 20; Year divided by 100 and truncated to integer (00-99)
	stdZeroDay:               "%d", // 23; Day of the month, zero-padded (01-31)
	stdShortSlashDate:        "%D", // 08/23/01; Short MM/DD/YY date, equivalent to %m/%d/%y, (same as %x)
	stdUnderDay:              "%e", // 23; Day of the month, space-padded ( 1-31)
	stdShortDashDate:         "%F", // 2001-08-23; Short YYYY-MM-DD date, equivalent to %Y-%m-%d
	stdISO8601WeekYear:       "%g", // 01; Week-based year, last two digits (00-99)
	stdISO8601LongWeekYear:   "%G", // 2001; Week-based year
	stdHour:                  "%H", // 14; Hour in 24h format (00-23), (same as %k)
	stdZeroHour12:            "%I", // 02; Hour in 12h format (01-12), (same as %l)
	stdDayOfYear:             "%j", // 235; Day of the year (001-366)
	stdZeroMonth:             "%m", // 08; Month as a decimal number (01-12)
	stdZeroMinute:            "%M", // 55; Minute (00-59)
	stdCharNewLine:           "%n", // '\n'; New-line character ('\n')
	stdPM:                    "%p", // PM; AM or PM designation
	stdpm:                    "%P", // om; am or pm designation
	stdHour12ClockTime:       "%r", // 02:55:02 pm; 12-hour clock time
	stdHourHourMinuteTime:    "%R", // 14:55; 24-hour HH:MM time, equivalent to %H:%M
	stdSecondsSinceEpoch:     "%s", // ; The number of seconds since the Epoch, 1970-01-01 00:00:00 +0000 (UTC). (TZ)
	stdZeroSecond:            "%S", // 02; Second (00-61)
	stdCharHorizontalTab:     "%t", // '\t'; Horizontal-tab character ('\t')
	stdISO8601Time:           "%T", // 14:55:02; ISO 8601 time format (HH:MM:SS), equivalent to %H:%M:%S
	stdNumWeekDay:            "%u", // 4; ISO 8601 weekday as number with Monday as 1 (1-7)
	stdSundayFirstWeekOfYear: "%U", // 33; Week number with the first Sunday as the first day of week one (00-53)
	stdISO8601NumWeek:        "%v", // 34; ISO 8601 week number (01-53)
	stdISO8601Week:           "%V", // 34; ISO 8601 week number (01-53)
	stdZeroNumWeek:           "%w", // 4; Weekday as a decimal number with Sunday as 0 (0-6)
	stdMonFirstWeekOfYear:    "%W", // 34; Week number with the first Monday as the first day of week one (00-53)
	stdHourClockTime:         "%X", // 14:55:02; Time representation
	stdYear:                  "%y", // 01; Year, last two digits (00-99)
	stdLongYear:              "%Y", // 2001; Year
	stdNumTZ:                 "%z", // +100; ISO 8601 offset from UTC in timezone (1 minute=1, 1 hour=100), If timezone cannot be determined, no characters
	stdTZ:                    "%Z", // CDT; Timezone name or abbreviation, If timezone cannot be determined, no characters
	stdCharPercentSign:       "%%", // '%'; A % sign
	stdISO8601TZ:             "%z",
	stdNop:                   "",
}

var strftimeSimilarChunkNames = map[int]string{
	stdNumMonth:                                 strftimeChunkNames[stdMonth],
	stdDay:                                      strftimeChunkNames[stdZeroDay],
	stdUnderYearDay:                             strftimeChunkNames[stdZeroDay],
	stdZeroYearDay:                              strftimeChunkNames[stdDayOfYear],
	stdHour12:                                   strftimeChunkNames[stdZeroHour12],
	stdMinute:                                   strftimeChunkNames[stdZeroMinute],
	stdSecond:                                   strftimeChunkNames[stdZeroSecond],
	stdISO8601SecondsTZ:                         strftimeChunkNames[stdNumTZ],
	stdISO8601ShortTZ:                           strftimeChunkNames[stdNumTZ],
	stdISO8601ColonTZ:                           strftimeChunkNames[stdNumTZ],
	stdISO8601ColonSecondsTZ:                    strftimeChunkNames[stdNumTZ],
	stdNumSecondsTz:                             strftimeChunkNames[stdNumTZ],
	stdNumShortTZ:                               strftimeChunkNames[stdNumTZ],
	stdNumColonTZ:                               strftimeChunkNames[stdNumTZ],
	stdNumColonSecondsTZ:                        strftimeChunkNames[stdNumTZ],
	stdDateAndTime | stdNeedEModifier:           "%Ec",
	stdFirstTwoDigitYear | stdNeedEModifier:     "%EC",
	stdShortSlashDate | stdNeedEModifier:        "%Ex",
	stdHourClockTime | stdNeedEModifier:         "%EX",
	stdYear | stdNeedEModifier:                  "%Ey",
	stdLongYear | stdNeedEModifier:              "%EY",
	stdZeroDay | stdNeedOModifier:               "%Od",
	stdUnderDay | stdNeedOModifier:              "%Oe",
	stdHour | stdNeedOModifier:                  "%OH",
	stdZeroHour12 | stdNeedOModifier:            "%OI",
	stdZeroMonth | stdNeedOModifier:             "%Om",
	stdZeroMinute | stdNeedOModifier:            "%OM",
	stdZeroSecond | stdNeedOModifier:            "%OS",
	stdNumWeekDay | stdNeedOModifier:            "%Ou",
	stdSundayFirstWeekOfYear | stdNeedOModifier: "%OU",
	stdISO8601Week | stdNeedOModifier:           "%OV",
	stdZeroNumWeek | stdNeedOModifier:           "%Ow",
	stdMonFirstWeekOfYear | stdNeedOModifier:    "%OW",
	stdYear | stdNeedOModifier:                  "%Oy",
	stdFracSecond0 | 1<<stdArgShift:             "",
	stdFracSecond0 | 2<<stdArgShift:             "",
	stdFracSecond0 | 3<<stdArgShift:             "",
	stdFracSecond0 | 4<<stdArgShift:             "",
	stdFracSecond0 | 5<<stdArgShift:             "",
	stdFracSecond0 | 6<<stdArgShift:             "",
	stdFracSecond0 | 7<<stdArgShift:             "",
	stdFracSecond0 | 8<<stdArgShift:             "",
	stdFracSecond0 | 9<<stdArgShift:             "",
	stdFracSecond9 | 1<<stdArgShift:             "",
	stdFracSecond9 | 2<<stdArgShift:             "",
	stdFracSecond9 | 3<<stdArgShift:             "",
	stdFracSecond9 | 4<<stdArgShift:             "",
	stdFracSecond9 | 5<<stdArgShift:             "",
	stdFracSecond9 | 6<<stdArgShift:             "",
	stdFracSecond9 | 7<<stdArgShift:             "",
	stdFracSecond9 | 8<<stdArgShift:             "",
	stdFracSecond9 | 9<<stdArgShift:             "",
	stdNop:                                      "",
}
