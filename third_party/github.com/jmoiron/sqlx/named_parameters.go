package sqlx

import "strings"

// JoinNamedColumns concatenates the elements of cols to column1, column2, ...
func JoinNamedColumns(cols []string) string {
	return strings.Join(cols, ",")
}

// JoinNamedValues concatenates the elements of values to :value1, :value2, ...
func JoinNamedValues(cols []string) string {
	if len(cols) == 1 {
		return ":" + cols[0]
	}
	return strings.Join(cols, ", :")
}

// JoinNamedCondition concatenates the elements of values to value1=:value1 AND value2=:value2 ...
func JoinNamedCondition(cols []string) string {
	for i, col := range cols {
		cols[i] = col + "=:" + col
	}
	return strings.Join(cols, " AND ")
}
