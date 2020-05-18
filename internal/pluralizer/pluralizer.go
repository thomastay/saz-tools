package pluralizer

import "strconv"

// FormatOrdinal appends a proper suffix (st, nd, rd or the) to an ordinal
// number according to the English grammar.
func FormatOrdinal(number int) string {
	value := strconv.Itoa(number)
	switch number % 10 {
	case 1:
		return value + "st"
	case 2:
		return value + "nd"
	case 3:
		return value + "rd"
	default:
		return value + "th"
	}
}
