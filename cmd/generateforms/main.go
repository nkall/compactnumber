package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/nkall/compactnumber/compact"
)

type generationParams struct {
	CompactFormsByLanguage map[string]compact.CompactForms
	Timestamp              time.Time
	CLDRVersion            string
}

type FileExtractResult struct {
	Language     string
	CLDRVersion  string
	CompactForms compact.CompactForms
}

func extractFromFile(b []byte) (FileExtractResult, error) {
	var file fileJson
	err := json.Unmarshal(b, &file)
	if err != nil {
		return FileExtractResult{}, err
	}

	// Each file has just one language body, but the key isn't consistent
	if len(file.Main) > 1 {
		return FileExtractResult{}, errors.New("got multiple language bodies")
	}

	var body bodyJson
	for _, fileBody := range file.Main {
		body = fileBody
		break
	}

	cldrVersion := body.Identity.Version.CLDRVersion
	language := body.Identity.Language
	if body.Identity.Territory != "" {
		language = fmt.Sprintf("%s-%s", language, body.Identity.Territory)
	}

	defaultNumberingSystemBytes, ok := body.Numbers["defaultNumberingSystem"]
	if !ok {
		return FileExtractResult{}, errors.New("missing default numbering system")
	}

	var defaultNumberingSystem string
	err = json.Unmarshal(defaultNumberingSystemBytes, &defaultNumberingSystem)
	if err != nil {
		return FileExtractResult{}, err
	}

	decimalFormatsBytes, ok := body.Numbers[fmt.Sprintf("decimalFormats-numberSystem-%s", defaultNumberingSystem)]
	if !ok {
		return FileExtractResult{}, errors.New(fmt.Sprintf("missing decimal formats for default numbering system %s", defaultNumberingSystem))
	}

	var decimalFormats decimalFormatsJson
	err = json.Unmarshal(decimalFormatsBytes, &decimalFormats)
	if err != nil {
		return FileExtractResult{}, err
	}

	forms, err := extractCompactForms(decimalFormats)
	if err != nil {
		return FileExtractResult{}, err
	}

	return FileExtractResult{
		Language:     language,
		CLDRVersion:  cldrVersion,
		CompactForms: forms,
	}, nil
}

func extractCompactForms(formats decimalFormatsJson) (compact.CompactForms, error) {
	if len(formats.Long.DecimalFormat) == 0 {
		return nil, errors.New("missing long formats")
	} else if len(formats.Short.DecimalFormat) == 0 {
		return nil, errors.New("missing short formats")
	}

	longForms, err := extractCompactFormRules(formats.Long.DecimalFormat)
	if err != nil {
		return nil, err
	}

	shortForms, err := extractCompactFormRules(formats.Short.DecimalFormat)
	if err != nil {
		return nil, err
	}

	return compact.CompactForms{
		compact.Short: shortForms,
		compact.Long:  longForms,
	}, nil
}

func extractCompactFormRules(formatRules map[string]string) ([]compact.CompactFormRule, error) {
	countString := "-count-"
	rules := make([]compact.CompactFormRule, 0, len(formatRules))

	// We expect to iterate in sorted order
	formatNames := make([]string, 0, len(formatRules))
	for formatName := range formatRules {
		formatNames = append(formatNames, formatName)
	}
	sort.Strings(formatNames)

	currType := int64(-1)
	var currRule *compact.CompactFormRule
	for _, formatName := range formatNames {
		countIndex := strings.Index(formatName, countString)
		if countIndex == -1 {
			return nil, errors.New(fmt.Sprintf("missing count from %s", formatName))
		}
		typeStr := formatName[:countIndex]
		typeNum, err := strconv.ParseInt(typeStr, 10, 0)

		if err != nil {
			return nil, err
		}

		formatPattern := formatRules[formatName]
		zeroesCount := strings.Count(formatPattern, "0")
		if zeroesCount == 0 {
			return nil, errors.New(fmt.Sprintf("missing zeroes from pattern: %s", formatPattern))
		}

		if typeNum != currType {
			if currRule != nil {
				rules = append(rules, *currRule)
			}
			currType = typeNum
			currRule = &compact.CompactFormRule{
				Type:                 typeNum,
				ZeroesInPattern:      zeroesCount,
				PatternsByPluralForm: make(map[string]string),
			}
		}

		afterCountIndex := countIndex + len(countString)
		pluralForm := formatName[afterCountIndex:]

		if currRule != nil {
			currRule.PatternsByPluralForm[pluralForm] = formatPattern
		} else {
			return nil, errors.New("current rule is nil (bad programmer)")
		}
	}

	return rules, nil
}

func writeFormsFile(params generationParams) error {
	templateFile, err := ioutil.ReadFile("./cmd/generateforms/forms.tmpl")
	if err != nil {
		return errors.New(fmt.Sprintf("error reading template file: %s", err.Error()))
	}

	templ, err := template.New("").Parse(string(templateFile))
	if err != nil {
		return errors.New(fmt.Sprintf("error parsing template file: %s", err.Error()))
	}

	f, err := os.Create("./compact/forms.gen.go")
	if err != nil {
		return errors.New(fmt.Sprintf("error parsing template file: %s", err.Error()))
	}

	defer f.Close()

	w := bufio.NewWriter(f)
	err = templ.Execute(w, params)
	if err != nil {
		return err
	}

	w.Flush()

	return nil
}

func main() {
	cldrPath := "./cldr"
	dirs, err := ioutil.ReadDir(cldrPath)
	if err != nil {
		log.Fatal(err)
	}

	var genParams generationParams
	genParams.CompactFormsByLanguage = make(map[string]compact.CompactForms)

	for _, d := range dirs {
		if !d.IsDir() {
			log.Printf("Skipping non-directory entry %s\n", d.Name())
			continue
		}
		log.Printf("Processing directory %s...\n", d.Name())

		b, err := ioutil.ReadFile(fmt.Sprintf("%s/%s/numbers.json", cldrPath, d.Name()))
		if err != nil {
			log.Fatalf("error reading file at %s: %s", d.Name(), err.Error())
		}

		forms, err := extractFromFile(b)
		if err != nil {
			log.Fatalf("error extracting forms from %s: %s", d.Name(), err.Error())
		}
		genParams.CLDRVersion = forms.CLDRVersion
		genParams.CompactFormsByLanguage[forms.Language] = forms.CompactForms
	}

	genParams.Timestamp = time.Now()

	err = writeFormsFile(genParams)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully generated compact formats file.")
}

// JSON unmarshalling structs
type fileJson struct {
	Main map[string]bodyJson `json:"main"`
}

type bodyJson struct {
	Identity identityJson               `json:"identity"`
	Numbers  map[string]json.RawMessage `json:"numbers"`
}

type identityJson struct {
	Version struct {
		CLDRVersion string `json:"_cldrVersion"`
	} `json:"version"`
	Language  string `json:"language"`
	Territory string `json:"territory"`
	Variant   string `json:"variant"`
}

type decimalFormatsJson struct {
	Long  decimalFormatJson `json:"long"`
	Short decimalFormatJson `json:"short"`
}

type decimalFormatJson struct {
	DecimalFormat map[string]string `json:"decimalFormat"`
}
