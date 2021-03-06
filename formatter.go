// Package compactnumber allows for localized CLDR compact number formatting.
package compactnumber

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gotnospirit/makeplural/plural"
	"github.com/nkall/compactnumber/internal/models"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
)

// Formatter is a struct containing a method to format an integer based on the specified language and compaction type.
type Formatter struct {
	lang        language.Tag
	compactType CompactType
}

// FormatterAPI is an interface implemented by Formatter that can be used for mocking purposes
type FormatterAPI interface {
	Format(n int, numOptions ...number.Option) (string, error)
}

// NewFormatter creates a new formatter based on the specified language and compaction type.
func NewFormatter(lang string, compactType CompactType) Formatter {
	return Formatter{
		lang:        language.Make(lang),
		compactType: compactType,
	}
}

// Format takes in an integer and options and formats it according to the formatter's locale and compaction settings.
// Note: this method truncates numbers and does not support fractions (e.g. 11.5M).
//
// Documented in CLDR spec: http://www.unicode.org/reports/tr35/tr35-numbers.html#Compact_Number_Formats
func (f *Formatter) Format(n int, numOptions ...number.Option) (string, error) {
	numOptions = append(numOptions, number.Scale(0))

	compactForms, ok := compactFormsByLanguage[f.lang.String()]
	if !ok {
		// Fall back to base language
		base, confidence := f.lang.Base()
		if confidence == language.No {
			return "", errors.New(fmt.Sprintf("no compact forms or fallback for language %s", f.lang.String()))
		}

		compactForms, ok = compactFormsByLanguage[base.String()]
		if !ok {
			return "", errors.New(fmt.Sprintf("missing compact forms for language %s and fallback %s", f.lang.String(), base))
		}
	}

	compactForm := compactForms[f.compactType]

	// Apply negative modifier at the end if dealing with negative number
	negativeModifier := 1
	if n < 0 {
		negativeModifier = -1
		n *= -1
	}

	// To format a number N, the greatest type less than or equal to N is used, with the appropriate plural category.
	var rule models.CompactFormRule
	for _, compactFormRule := range compactForm {
		if int64(n) >= compactFormRule.Type {
			rule = compactFormRule
		} else {
			break
		}
	}

	// N is divided by the type, after removing the number of zeros in the pattern, less 1.
	shortN := f.shortNum(n, rule)

	// Best effort fetching plural form
	plurForm := f.pluralForm(shortN)

	pattern, ok := rule.PatternsByPluralForm[plurForm]
	if !ok {
		// Attempt to fall back to catch-all "other" pattern if none for current plural form found
		pattern = rule.PatternsByPluralForm["other"]
	}

	var err error
	pattern, err = formatPattern(pattern)
	if err != nil {
		return "", err
	}

	// If the value is precisely “0”, either explicit or defaulted, then the normal number format pattern for that sort of object is supplied
	baseNumPrinter := message.NewPrinter(f.lang)
	if pattern == "0" {
		return baseNumPrinter.Sprintf("%v", number.Decimal(n*negativeModifier, numOptions...)), nil
	}

	return baseNumPrinter.Sprintf(pattern, number.Decimal(shortN*int64(negativeModifier), numOptions...)), nil
}

// Divides number to be used in compact display according to logic in CLDR spec: http://www.unicode.org/reports/tr35/tr35-numbers.html#Compact_Number_Formats
func (f *Formatter) shortNum(n int, rule models.CompactFormRule) int64 {
	typeDivisor := rule.Type
	for i := 0; i < rule.ZeroesInPattern-1; i++ {
		typeDivisor /= 10
	}

	outNum := int64(n)
	if typeDivisor != 0 {
		outNum /= typeDivisor
	}

	return outNum
}

// Gets the pluralized form of the number, as per CLDR spec: http://cldr.unicode.org/index/cldr-spec/plural-rules
// We use gotnospirit/makeplural for this as golang.org/x/text/plural does not expose a suitable PluralForm method.
// This is a best effort function since the languages might not match up perfectly between packages.
func (f *Formatter) pluralForm(n interface{}) string {
	base, confidence := f.lang.Base()
	if confidence == language.No {
		return "other"
	}

	plurFunc, err := plural.GetFunc(base.String())
	if err != nil {
		return "other"
	}

	return plurFunc(n, false)
}

// Process CLDR pattern to a format suitable for use in Printer.Sprintf in golang.org/x/text/message.
// Documentation for these special characters can be found in the CLDR spec: http://cldr.unicode.org/translation/number-patterns
func formatPattern(pattern string) (string, error) {
	// Default to 0 (sentinel value for no pattern)
	if pattern == "" {
		return "0", nil
	}

	// Default to the first form if there's multiple
	pattern = strings.Split(pattern, ";")[0]

	// Remove special pattern symbols, as this formatting is already handled by golang.org/x/text/message
	pattern = strings.Replace(pattern, "'", "", -1)

	// Replace all 0s with a single %v for number formatting
	zeroIndex := strings.IndexRune(pattern, '0')
	if zeroIndex == -1 {
		return "", errors.New(fmt.Sprintf("invalid pattern (no digit pattern characters): %s", pattern))
	}

	pattern = strings.Replace(pattern, "0", "", -1)
	patternRunes := []rune(pattern)

	endStr := ""
	if zeroIndex < len(patternRunes) {
		endStr = string(patternRunes[zeroIndex:])
	}

	pattern = fmt.Sprintf("%s%s%s", string(patternRunes[:zeroIndex]), "%v", endStr)
	return pattern, nil
}
