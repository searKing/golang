// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mysql

// Error codes for server-side errors.
// Originally found in https://dev.mysql.com/doc/mysql-errors/8.0/en/server-error-reference.html
// See above reference for more information on each code.
const (
	// unknown
	ERUnknownError = 1105

	// internal
	ERInternalError = 1815

	// unimplemented
	ERNotSupportedYet = 1235
	ERUnsupportedPS   = 1295

	// resource exhausted
	ERDiskFull               = 1021
	EROutOfMemory            = 1037
	EROutOfSortMemory        = 1038
	ERConCount               = 1040
	EROutOfResources         = 1041
	ERRecordFileFull         = 1114
	ERHostIsBlocked          = 1129
	ERCantCreateThread       = 1135
	ERTooManyDelayedThreads  = 1151
	ERNetPacketTooLarge      = 1153
	ERTooManyUserConnections = 1203
	ERLockTableFull          = 1206
	ERUserLimitReached       = 1226

	// deadline exceeded
	ERLockWaitTimeout = 1205

	// unavailable
	ERServerShutdown = 1053

	// not found
	ERCantFindFile          = 1017
	ERFormNotFound          = 1029
	ERKeyNotFound           = 1032
	ERBadFieldError         = 1054
	ERNoSuchThread          = 1094
	ERUnknownTable          = 1109
	ERCantFindUDF           = 1122
	ERNonExistingGrant      = 1141
	ERNoSuchTable           = 1146
	ERNonExistingTableGrant = 1147
	ERKeyDoesNotExist       = 1176
	ERDbDropExists          = 1008

	// permissions
	ERDBAccessDenied            = 1044
	ERAccessDeniedError         = 1045
	ERKillDenied                = 1095
	ERNoPermissionToCreateUsers = 1211
	ERSpecifiedAccessDenied     = 1227

	// failed precondition
	ERNoDb                          = 1046
	ERNoSuchIndex                   = 1082
	ERCantDropFieldOrKey            = 1091
	ERTableNotLockedForWrite        = 1099
	ERTableNotLocked                = 1100
	ERTooBigSelect                  = 1104
	ERNotAllowedCommand             = 1148
	ERTooLongString                 = 1162
	ERDelayedInsertTableLocked      = 1165
	ERDupUnique                     = 1169
	ERRequiresPrimaryKey            = 1173
	ERCantDoThisDuringAnTransaction = 1179
	ERReadOnlyTransaction           = 1207
	ERCannotAddForeign              = 1215
	ERNoReferencedRow               = 1216
	ERRowIsReferenced               = 1217
	ERCantUpdateWithReadLock        = 1223
	ERNoDefault                     = 1230
	EROperandColumns                = 1241
	ERSubqueryNo1Row                = 1242
	ERWarnDataOutOfRange            = 1264
	ERNonUpdateableTable            = 1288
	ERFeatureDisabled               = 1289
	EROptionPreventsStatement       = 1290
	ERDuplicatedValueInType         = 1291
	ERSPDoesNotExist                = 1305
	ERRowIsReferenced2              = 1451
	ErNoReferencedRow2              = 1452
	ErSPNotVarArg                   = 1414
	ERInnodbReadOnly                = 1874

	// already exists
	ERTableExists    = 1050
	ERDupEntry       = 1062
	ERFileExists     = 1086
	ERUDFExists      = 1125
	ERDbCreateExists = 1007

	// aborted
	ERGotSignal          = 1078
	ERForcingClose       = 1080
	ERAbortingConnection = 1152
	ERLockDeadlock       = 1213

	// invalid arg
	ERUnknownComError              = 1047
	ERBadNullError                 = 1048
	ERBadDb                        = 1049
	ERBadTable                     = 1051
	ERNonUniq                      = 1052
	ERWrongFieldWithGroup          = 1055
	ERWrongGroupField              = 1056
	ERWrongSumSelect               = 1057
	ERWrongValueCount              = 1058
	ERTooLongIdent                 = 1059
	ERDupFieldName                 = 1060
	ERDupKeyName                   = 1061
	ERWrongFieldSpec               = 1063
	ERParseError                   = 1064
	EREmptyQuery                   = 1065
	ERNonUniqTable                 = 1066
	ERInvalidDefault               = 1067
	ERMultiplePriKey               = 1068
	ERTooManyKeys                  = 1069
	ERTooManyKeyParts              = 1070
	ERTooLongKey                   = 1071
	ERKeyColumnDoesNotExist        = 1072
	ERBlobUsedAsKey                = 1073
	ERTooBigFieldLength            = 1074
	ERWrongAutoKey                 = 1075
	ERWrongFieldTerminators        = 1083
	ERBlobsAndNoTerminated         = 1084
	ERTextFileNotReadable          = 1085
	ERWrongSubKey                  = 1089
	ERCantRemoveAllFields          = 1090
	ERUpdateTableUsed              = 1093
	ERNoTablesUsed                 = 1096
	ERTooBigSet                    = 1097
	ERBlobCantHaveDefault          = 1101
	ERWrongDbName                  = 1102
	ERWrongTableName               = 1103
	ERUnknownProcedure             = 1106
	ERWrongParamCountToProcedure   = 1107
	ERWrongParametersToProcedure   = 1108
	ERFieldSpecifiedTwice          = 1110
	ERInvalidGroupFuncUse          = 1111
	ERTableMustHaveColumns         = 1113
	ERUnknownCharacterSet          = 1115
	ERTooManyTables                = 1116
	ERTooManyFields                = 1117
	ERTooBigRowSize                = 1118
	ERWrongOuterJoin               = 1120
	ERNullColumnInIndex            = 1121
	ERFunctionNotDefined           = 1128
	ERWrongValueCountOnRow         = 1136
	ERInvalidUseOfNull             = 1138
	ERRegexpError                  = 1139
	ERMixOfGroupFuncAndFields      = 1140
	ERIllegalGrantForTable         = 1144
	ERSyntaxError                  = 1149
	ERWrongColumnName              = 1166
	ERWrongKeyColumn               = 1167
	ERBlobKeyWithoutLength         = 1170
	ERPrimaryCantHaveNull          = 1171
	ERTooManyRows                  = 1172
	ERLockOrActiveTransaction      = 1192
	ERUnknownSystemVariable        = 1193
	ERSetConstantsOnly             = 1204
	ERWrongArguments               = 1210
	ERWrongUsage                   = 1221
	ERWrongNumberOfColumnsInSelect = 1222
	ERDupArgument                  = 1225
	ERLocalVariable                = 1228
	ERGlobalVariable               = 1229
	ERWrongValueForVar             = 1231
	ERWrongTypeForVar              = 1232
	ERVarCantBeRead                = 1233
	ERCantUseOptionHere            = 1234
	ERIncorrectGlobalLocalVar      = 1238
	ERWrongFKDef                   = 1239
	ERKeyRefDoNotMatchTableRef     = 1240
	ERCyclicReference              = 1245
	ERCollationCharsetMismatch     = 1253
	ERCantAggregate2Collations     = 1267
	ERCantAggregate3Collations     = 1270
	ERCantAggregateNCollations     = 1271
	ERVariableIsNotStruct          = 1272
	ERUnknownCollation             = 1273
	ERWrongNameForIndex            = 1280
	ERWrongNameForCatalog          = 1281
	ERBadFTColumn                  = 1283
	ERTruncatedWrongValue          = 1292
	ERTooMuchAutoTimestampCols     = 1293
	ERInvalidOnUpdate              = 1294
	ERUnknownTimeZone              = 1298
	ERInvalidCharacterString       = 1300
	ERIllegalReference             = 1247
	ERDerivedMustHaveAlias         = 1248
	ERTableNameNotAllowedHere      = 1250
	ERQueryInterrupted             = 1317
	ERTruncatedWrongValueForField  = 1366
	ERDataTooLong                  = 1406
	ERForbidSchemaChange           = 1450
	ERDataOutOfRange               = 1690

	// server not available
	ERServerIsntAvailable = 3168
)
