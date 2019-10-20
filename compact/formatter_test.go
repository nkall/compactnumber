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

	out, err = formatter.Format(-5999)
	mustMatch(t, out, err, "-5K", nil)

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

	out, err = formatter.Format(-17999999)
	mustMatch(t, out, err, "-17M", nil)

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

	out, err = formatter.Format(-420000)
	mustMatch(t, out, err, "-420 thousand", nil)

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

	out, err = formatter.Format(-5999)
	mustMatch(t, out, err, "−5k", nil)

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

	out, err = formatter.Format(-9999999999)
	mustMatch(t, out, err, "-9 миллиардов", nil)

	out, err = formatter.Format(1999999999999)
	mustMatch(t, out, err, "1 триллион", nil)

	out, err = formatter.Format(2999999999999)
	mustMatch(t, out, err, "2 триллиона", nil)

	out, err = formatter.Format(9999999999999)
	mustMatch(t, out, err, "9 триллионов", nil)
}

func TestFormatterFormatVariousMajorLanguages(t *testing.T) {
	tests := []struct {
		localeStr   string
		expectedOut string
	}{
		{localeStr: "en-US", expectedOut: "-69M"},
		{localeStr: "ar-SA", expectedOut: "‎-69 مليون"},
		{localeStr: "bg-BG", expectedOut: "-69 млн."},
		{localeStr: "da-DK", expectedOut: "-69 mio."},
		{localeStr: "de-DE", expectedOut: "-69 Mio."},
		{localeStr: "el-GR", expectedOut: "-69 εκ."},
		{localeStr: "es-ES", expectedOut: "-69 M"},
		{localeStr: "es-MX", expectedOut: "-69 M"},
		{localeStr: "fi-FI", expectedOut: "−69 milj."},
		{localeStr: "fr-FR", expectedOut: "-69 M"},
		{localeStr: "hu-HU", expectedOut: "-69 M"},
		{localeStr: "it-IT", expectedOut: "-69 Mln"},
		{localeStr: "ja-JP", expectedOut: "-6,954万"},
		{localeStr: "nl-NL", expectedOut: "-69 mln."},
		{localeStr: "no-NO", expectedOut: "-69 mill."},
		{localeStr: "pl-PL", expectedOut: "-69 mln"},
		{localeStr: "pt-BR", expectedOut: "-69 mi"},
		{localeStr: "pt-PT", expectedOut: "-69 M"},
		{localeStr: "ro-RO", expectedOut: "-69 mil."},
		{localeStr: "ru-RU", expectedOut: "-69 млн"},
		{localeStr: "sk-SK", expectedOut: "-69 mil."},
		{localeStr: "sv-SE", expectedOut: "−69 mn"},
		{localeStr: "th-TH", expectedOut: "-69M"},
		{localeStr: "tr-TR", expectedOut: "-69 Mn"},
		{localeStr: "vi-VN", expectedOut: "-69 Tr"},
		{localeStr: "zh-CN", expectedOut: "-6,954萬"},
		{localeStr: "zh-TW", expectedOut: "-6,954萬"},
	}
	for _, tt := range tests {
		t.Run(tt.localeStr, func(t *testing.T) {
			formatter := compact.NewFormatter(language.Make(tt.localeStr), compact.Short)
			out, err := formatter.Format(-69540001)
			mustMatch(t, out, err, tt.expectedOut, nil)
		})
	}
}

func mustMatch(t *testing.T, out string, err error, expectedOut string, expectedErr error) {
	if err != expectedErr {
		t.Error(fmt.Sprintf("got unexpected error %v (wanted %v)", err, expectedErr))
	}

	if out != expectedOut {
		t.Error(fmt.Sprintf("got unexpected output %s (wanted %s)", out, expectedOut))
	}
}
