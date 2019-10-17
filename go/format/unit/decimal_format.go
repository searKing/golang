package unit

func DecimalFormatInt(number int, precision int) string {
	return DecimalFormatInt64(int64(number), precision)
}

func DecimalFormatUint(number uint, precision int) string {
	return DecimalFormatUint64(uint64(number), precision)
}

func DecimalFormatInt64(number int64, precision int) string {
	return DecimalMultiplePrefixTODO.Copy().SetInt64(number).FormatInt64(number, precision)
}

func DecimalFormatUint64(number uint64, precision int) string {
	return DecimalMultiplePrefixTODO.Copy().SetUint64(number).FormatUint64(number, precision)
}

func DecimalFormatFloat(number float64, precision int) string {
	return DecimalMultiplePrefixTODO.Copy().SetFloat64(number).FormatFloat(number, precision)
}
