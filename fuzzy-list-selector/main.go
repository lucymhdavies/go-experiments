package main

import (
	"io/ioutil"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/lucymhdavies/go-experiments/fuzzy-list-selector/selector"
)

func main() {
	log.SetLevel(log.DebugLevel)

	inputDataBytes, err := ioutil.ReadFile("data.txt")
	if err != nil {
		log.Fatalf("Error reading data.txt: %s", err)
	}

	inputDataSlice := strings.Split(strings.TrimSpace(string(inputDataBytes)), "\n")
	for _, data := range inputDataSlice {
		log.Infof(" - %s", data)
	}

	s := selector.NewSelector("AWS Accounts")

	selection, err := s.SelectFromSliceWithFilter(inputDataSlice, "")
	if err != nil {
		log.Fatalf("Error selecting from list: %s", err)
	}

	log.Infof("Selection: %s", selection)

}
