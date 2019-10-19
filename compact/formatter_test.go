package compact_test

import (
	"fmt"
	"golang.org/x/text/number"
	"testing"

	"github.com/nkall/compactnumber/compact"
	"golang.org/x/text/language"
)

func TestFormatterFormatAmericanEnglish(t *testing.T) {
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

	out, err = formatter.Format(17999999, number.Scale(1))
	mustMatch(t, out, err, "17M", nil)

}

func mustMatch(t *testing.T, out string, err error, expectedOut string, expectedErr error) {
	if err != expectedErr {
		t.Error(fmt.Sprintf("got unexpected error %v (wanted %v)", err, expectedErr))
	}

	if out != expectedOut {
		t.Error(fmt.Sprintf("got unexpected output %s (wanted %s)", out, expectedOut))
	}
}
