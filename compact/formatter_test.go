package compact_test

import (
	"fmt"
	"testing"

	"github.com/nkall/compactnumber/compact"
	"golang.org/x/text/language"
)

func TestFormatterFormatAmericanEnglishShort(t *testing.T) {
	enLang := language.Make("en-US")
	formatter := compact.NewFormatter(enLang, compact.Short)

	out, err := formatter.Format(0)
	mustMatch(t, out, err, "0", nil)

	out, err = formatter.Format(999)
	mustMatch(t, out, err, "999", nil)

	out, err = formatter.Format(5000)
	mustMatch(t, out, err, "5K", nil)

	out, err = formatter.Format(5501)
	mustMatch(t, out, err, "5K", nil)

	out, err = formatter.Format(5999)
	mustMatch(t, out, err, "5K", nil)

	out, err = formatter.Format(19000)
	mustMatch(t, out, err, "19K", nil)

	out, err = formatter.Format(19999)
	mustMatch(t, out, err, "19K", nil)

	out, err = formatter.Format(420000)
	mustMatch(t, out, err, "420K", nil)

	out, err = formatter.Format(420999)
	mustMatch(t, out, err, "420K", nil)

	out, err = formatter.Format(17000000)
	mustMatch(t, out, err, "17M", nil)

	out, err = formatter.Format(17999999)
	mustMatch(t, out, err, "17M", nil)

	out, err = formatter.Format(999999999999)
	mustMatch(t, out, err, "999B", nil)

	out, err = formatter.Format(999999999999999)
	mustMatch(t, out, err, "999T", nil)
}

func TestFormatterFormatAmericanEnglishLong(t *testing.T) {
	enLang := language.Make("en-US")
	formatter := compact.NewFormatter(enLang, compact.Long)

	out, err := formatter.Format(0)
	mustMatch(t, out, err, "0", nil)

	out, err = formatter.Format(5501)
	mustMatch(t, out, err, "5 thousand", nil)

	out, err = formatter.Format(420000)
	mustMatch(t, out, err, "420 thousand", nil)

	out, err = formatter.Format(17000000)
	mustMatch(t, out, err, "17 million", nil)

	out, err = formatter.Format(999999999999)
	mustMatch(t, out, err, "999 billion", nil)

	out, err = formatter.Format(999999999999999)
	mustMatch(t, out, err, "999 trillion", nil)
}

func TestFormatterFormatNorskBokmaal(t *testing.T) {
	noLang := language.Make("nb-NO")
	formatter := compact.NewFormatter(noLang, compact.Short)

	out, err := formatter.Format(0)
	mustMatch(t, out, err, "0", nil)

	out, err = formatter.Format(999)
	mustMatch(t, out, err, "999", nil)

	out, err = formatter.Format(5999)
	mustMatch(t, out, err, "5k", nil)

	out, err = formatter.Format(17999999)
	mustMatch(t, out, err, "17 mill.", nil)

	out, err = formatter.Format(999999999999)
	mustMatch(t, out, err, "999 mrd.", nil)

	out, err = formatter.Format(5999999999999)
	mustMatch(t, out, err, "5 bill.", nil)
}

func TestFormatterFormatRussian(t *testing.T) {
	ruLang := language.Make("ru")
	formatter := compact.NewFormatter(ruLang, compact.Long)

	out, err := formatter.Format(1999999999)
	mustMatch(t, out, err, "1 миллиард", nil)

	out, err = formatter.Format(2999999999)
	mustMatch(t, out, err, "2 миллиарда", nil)

	out, err = formatter.Format(9999999999)
	mustMatch(t, out, err, "9 миллиардов", nil)

	out, err = formatter.Format(1999999999999)
	mustMatch(t, out, err, "1 триллион", nil)

	out, err = formatter.Format(2999999999999)
	mustMatch(t, out, err, "2 триллиона", nil)

	out, err = formatter.Format(9999999999999)
	mustMatch(t, out, err, "9 триллионов", nil)
}

func mustMatch(t *testing.T, out string, err error, expectedOut string, expectedErr error) {
	if err != expectedErr {
		t.Error(fmt.Sprintf("got unexpected error %v (wanted %v)", err, expectedErr))
	}

	if out != expectedOut {
		t.Error(fmt.Sprintf("got unexpected output %s (wanted %s)", out, expectedOut))
	}
}
