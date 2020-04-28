package main

import (
	"github/drakos74/oremi/covid"
	oremi "github/drakos74/oremi/internal"
	"github/drakos74/oremi/internal/data/source/web"
	"log"

	"gioui.org/layout"
)

func main() {

	// get the source data
	b, err := web.Html("https://covid.ourworldindata.org/data/owid-covid-data.csv")
	if err != nil {
		log.Fatalf("could not get html contents: %v", err)
	}

	// parse the input into our model
	infections := covid.Parse(b)

	// create a data collection out of the gathered data
	collections, err := infections.ToCollection()
	if err != nil {
		log.Fatalf("could not convert data to collection: %v", err)
	}

	// draw the data collection
	oremi.DrawGraph("covid-19", layout.Vertical, 1600, 800, collections)

}
