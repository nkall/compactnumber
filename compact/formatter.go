package compact

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gotnospirit/makeplural/plural"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
)

// Formatter is a struct containing a method to format an integer based on the specified language and compaction type.
type Formatter struct {
	lang        language.Tag
	compactType CompactType
}

// NewFormatter creates a new formatter based on the specified language and compaction type.
func NewFormatter(lang language.Tag, compactType CompactType) Formatter {
	return Formatter{
		lang:        lang,
		compactType: compactType,
	}
}

// Format takes in an integer and options and formats it according to the formatter's locale and compaction settings.
// Documented in CLDR spec: http://www.unicode.org/reports/tr35/tr35-numbers.html#Compact_Number_Formats
func (f *Formatter) Format(n int, numOptions ...number.Option) (string, error) {
	compactForms, ok := compactFormsByLanguage[f.lang.String()]
	if !ok {
		return "", errors.New(fmt.Sprintf("missing compact forms for language %s", f.lang.String()))
	}

	compactForm := compactForms[f.compactType]

	// To format a number N, the greatest type less than or equal to N is used, with the appropriate plural category.
	var rule CompactFormRule
	for _, compactFormRule := range compactForm {
		if n >= compactFormRule.Type {
			rule = compactFormRule
			break
		}
	}

	// N is divided by the type, after removing the number of zeros in the pattern, less 1.
	shortN := f.shortNum(n, rule)

	plurForm, err := f.pluralForm(shortN)
	if err != nil {
		return "", err
	}

	pattern, ok := rule.PatternsByPluralForm[plurForm]
	if !ok {
		// Attempt to fall back to catch-all "other" pattern first, then default to sentinel value
		pattern, ok = rule.PatternsByPluralForm["other"]
		if !ok {
			pattern = "0"
		}
	}

	// If the value is precisely “0”, either explicit or defaulted, then the normal number format pattern for that sort of object is supplied
	baseNumPrinter := message.NewPrinter(f.lang)
	if pattern == "0" {
		return baseNumPrinter.Sprintf("%v", number.Decimal(n, numOptions...)), nil
	}

	outPattern, err := formatPattern(pattern)
	if err != nil {
		return "", err
	}
	return baseNumPrinter.Sprintf(outPattern, number.Decimal(shortN, numOptions...)), nil
}

// Divides number to be used in compact display according to logic in CLDR spec: http://www.unicode.org/reports/tr35/tr35-numbers.html#Compact_Number_Formats
func (f *Formatter) shortNum(n int, rule CompactFormRule) float64 {
	typeDivisor := rule.Type
	for i := 0; i < rule.ZeroesInPattern-1; i++ {
		typeDivisor /= 10
	}
	shortN := float64(n)
	if typeDivisor != 0 {
		shortN /= float64(typeDivisor)
	}

	return shortN
}

// Gets the pluralized form of the number, as per CLDR spec: http://cldr.unicode.org/index/cldr-spec/plural-rules
// We use gotnospirit/makeplural for this as golang.org/x/text/plural does not expose a suitable PluralForm method.
func (f *Formatter) pluralForm(n float64) (string, error) {
	plurFunc, err := plural.GetFunc(f.lang.String())
	if err != nil {
		return "", err
	}

	return plurFunc(n, false), nil
}

// Process CLDR pattern to a format suitable for use in Printer.Sprintf in golang.org/x/text/message.
// Documentation for these special characters can be found in the CLDR spec: http://cldr.unicode.org/translation/number-patterns
func formatPattern(pattern string) (string, error) {
	// Remove special pattern symbols, as this formatting is already handled by golang.org/x/text/message
	pattern = strings.ReplaceAll(pattern, ".", "")
	pattern = strings.ReplaceAll(pattern, ",", "")
	pattern = strings.ReplaceAll(pattern, "#", "")
	pattern = strings.ReplaceAll(pattern, "'", "")

	// Replace all 0s with a single %v for number formatting
	zeroIndex := strings.IndexRune(pattern, '0')
	if zeroIndex == -1 {
		return "", errors.New(fmt.Sprintf("invalid pattern (no digit pattern characters): %s", pattern))
	}

	pattern = strings.ReplaceAll(pattern, "0", "")
	patternRunes := []rune(pattern)

	pattern = fmt.Sprintf("%s%s%s", string(patternRunes[:zeroIndex]), "%v", string(patternRunes[zeroIndex:]))
	return pattern, nil
}
