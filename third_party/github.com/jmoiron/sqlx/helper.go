package sqlx

//go:generate go-enum -type SqlOperator -trimprefix=SqlOperator --transform=upper
type SqlOperator int

const (
	SqlOperatorAnd SqlOperator = 0
	SqlOperatorOr  SqlOperator = 1
	SqlOperatorNot SqlOperator = 2
)

//go:generate go-enum -type SqlCompare -trimprefix=SqlCompare --linecomment
type SqlCompare int

const (
	SqlCompareEqual            SqlCompare = iota //=
	SqlCompareNotEqual         SqlCompare = iota //<>
	SqlCompareGreaterThan      SqlCompare = iota //>
	SqlCompareLessThan         SqlCompare = iota //<
	SqlCompareGreatEqual       SqlCompare = iota //>=
	SqlCompareLessAndEqualThan SqlCompare = iota //<=
	SqlCompareLike             SqlCompare = iota //LIKE
)

// NamedColumns returns the []string{value1, value2 ...}
// query := NamedColumns("foo", "bar")
// // []string{"foo", "bar"}
func NamedColumns(cols ...string) []string {
	return NamedTableColumns("", cols...)
}

// NamedValues returns the []string{:value1, :value2 ...}
// query := NamedValues("foo", "bar")
// // []string{":foo", ":bar"}
func NamedValues(cols ...string) []string {
	return NamedTableValues(cols...)
}

// NamedColumnsValues returns the []string{value1=:value1, value2=:value2 ...}
// query := NamedColumnsValues("foo", "bar")
// // []string{"foo=:foo", bar=:bar"}
func NamedColumnsValues(cmp SqlCompare, cols ...string) []string {
	return NamedTableColumnsValues(cmp, "", cols...)
}

// JoinColumns concatenates the elements of cols to column1, column2, ...
// query := JoinColumns("foo", "bar")
// // "foo,bar"
func JoinColumns(cols ...string) string {
	return JoinTableColumns("", cols...)
}

// JoinNamedValues concatenates the elements of values to :value1, :value2, ...
// query := JoinNamedValues("foo", "bar")
// // ":foo,:bar"
// query := JoinNamedValues()
// // "DEFAULT"
// Deprecated: Use NamedInsertArguments instead.
func JoinNamedValues(cols ...string) string {
	return JoinNamedTableValues(cols...)
}

// JoinNamedColumnsAndValues concatenates the elements of values to value1=:value1, value2=:value2 ...
// Deprecated: Use NamedUpdateArguments instead.
func JoinNamedColumnsValues(cols ...string) string {
	return JoinNamedTableColumnsValues("", cols...)
}

// JoinNamedCondition concatenates the elements of values to value1=:value1 AND value2=:value2 ...
// query := JoinNamedCondition(SqlCompareEqual,SqlOperatorAnd,"foo", "bar")
// // "foo=:foo AND bar=:bar"
func JoinNamedCondition(cmp SqlCompare, operator SqlOperator, cols ...string) string {
	return JoinNamedTableCondition(cmp, operator, "", cols...)
}

// JoinNamedColumns concatenates the elements of cols to column1, column2, ...
// Deprecated: Use NamedInsertArguments instead.
func JoinNamedColumns(cols ...string) string {
	return JoinNamedTableColumns("", cols...)
}
