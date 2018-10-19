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

	inputDataSlice := strings.Split(string(inputDataBytes), "\n")

	s := selector.NewSelector("AWS Accounts")

	selection, err := s.SelectFromSlice(inputDataSlice)
	if err != nil {
		log.Fatalf("Error selecting from list: %s", err)
	}

	log.Infof("Selection: %s", selection)

}