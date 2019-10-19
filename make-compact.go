package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/nkall/compactnumber/compact"
)

type generationParams struct {
	CompactFormsByLanguage map[string]compact.CompactForms
	Timestamp              time.Time
}

func main() {
	flag.Parse()

	dirs, err := ioutil.ReadDir("./cldr")
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range dirs {
		if !d.IsDir() {
			log.Fatal("not a directory: %s", d.Name())
		}

		b, err := ioutil.ReadFile(fmt.Sprintf("./cldr/%s/numbers.json", d.Name()))
		if err != nil {
			log.Fatal(err)
		}

	}
}
