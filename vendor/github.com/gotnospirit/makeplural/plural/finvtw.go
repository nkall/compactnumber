package plural

import (
	"math"
	"strconv"
	"strings"
)

func finvtw(value interface{}) (int64, int64, float64, int, int64, int) {
	// @see http://unicode.org/reports/tr35/tr35-numbers.html#Operands
	//
	// Symbol	Value
	// n	    absolute value of the source number (integer and decimals).
	// i	    integer digits of n.
	// v	    number of visible fraction digits in n, with trailing zeros.
	// w	    number of visible fraction digits in n, without trailing zeros.
	// f	    visible fractional digits in n, with trailing zeros.
	// t	    visible fractional digits in n, without trailing zeros.
	var strval string
	var pos int

	var f int64
	var i int64
	var n float64
	var v int
	var t int64

	switch value.(type) {
	case int:
		i = int64(value.(int))
		return 0, i, math.Abs(float64(i)), 0, 0, 0

	case int64:
		i = value.(int64)
		return 0, i, math.Abs(float64(i)), 0, 0, 0

	case float64:
		floatval := value.(float64)
		strval = strconv.FormatFloat(floatval, 'f', -1, 64)
		pos = strings.Index(strval, ".")
		if -1 == pos {
			return 0, int64(floatval), math.Abs(floatval), 0, 0, 0
		}
		n = math.Abs(floatval)

	case string:
		strval = value.(string)
		pos = strings.Index(strval, ".")
		if -1 == pos {
			intvalue, err := strconv.ParseInt(strval, 10, 64)
			if nil != err {
				return 0, 0, 0, 0, 0, 0
			}
			return 0, intvalue, math.Abs(float64(intvalue)), 0, 0, 0
		}
		floatval, err := strconv.ParseFloat(strval, 64)
		if nil != err {
			return 0, 0, 0, 0, 0, 0
		}

		n = math.Abs(floatval)
	}

	strf := strval[pos+1:]

	f, err := strconv.ParseInt(strf, 10, 64)
	if nil != err {
		return 0, 0, 0, 0, 0, 0
	}

	i, err = strconv.ParseInt(strval[:pos], 10, 64)
	if nil != err {
		return 0, 0, 0, 0, 0, 0
	}

	v = len(strf)

	var strt string

	offset := len(strf)
loop:
	for i := offset - 1; i >= 0; i-- {
		switch strf[i] {
		case '0':
			offset--

		default:
			break loop
		}
	}
	if offset >= 1 {
		strt = strf[:offset]
	}

	if "" == strt {
		return f, i, n, v, 0, 0
	} else if strf != strt {
		t, err = strconv.ParseInt(strt, 10, 64)
		if nil != err {
			return 0, 0, 0, 0, 0, 0
		}
		return f, i, n, v, t, len(strt)
	}
	return f, i, n, v, f, v
}
