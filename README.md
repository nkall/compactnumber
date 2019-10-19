# Compact Number Formats
A wrapper over `golang.org/x/text` with support for compact number formats in various languages.  For example, in English,
19,000,000 has the compact formats of 19 million (long) or 19M (short), but this varies widely across languages. You can
see Unicode's documentation of Compact Number Formats [here](http://www.unicode.org/reports/tr35/tr35-numbers.html#Compact_Number_Formats).

This implementation should be considered best-effort. I hope to eventually add this functionality to the core
`golang.org/x/text` repository (see [my proposal here](https://github.com/golang/go/issues/34989)).

## Usage
```
import (
	"fmt"

	"github.com/nkall/compactnumber/compact"
	"golang.org/x/text/language"
)

func main() {
	enLang := language.Make("en-US")
	formatter := compact.NewFormatter(enLang, compact.Short)
	out, err := formatter.Format(17999999)
	if err != nil {
		panic(err)
	}

	fmt.Println(out) // 17M
}
```

## Generating Compact Forms
Compact forms can be regenerated with the latest CLDR data by following these steps:

1. Download the latest JSON CLDR distribution from https://github.com/unicode-cldr/cldr-numbers-modern
1. Extract the contents of the main directory to `compactnumber/cldr`.
1. Run `make generate` and check in the updated file `forms.gen.go`.