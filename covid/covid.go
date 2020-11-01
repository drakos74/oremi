package covid

import (
	"fmt"
	"time"

	"github.com/drakos74/oremi/label"

	"github.com/drakos74/oremi"

	"github.com/drakos74/oremi/internal/data/model"

	"github.com/gocarina/gocsv"
)

type Infection struct {
	Country          string `csv:"location"`
	Date             Date   `csv:"date"`
	Cases            uint   `csv:"total_cases"`
	Deaths           uint   `csv:"total_deaths"`
	Tests            uint   `csv:"total_tests"`
	TotalCasesPerMil uint   `csv:"total_cases_per_million"`
}

func (i Infection) ToVector() (model.Vector, error) {
	return model.NewVector([]string{i.Country, i.Date.String()}, toDay(i.Date.Time), float64(i.Cases), float64(i.Deaths), float64(i.Tests), float64(i.TotalCasesPerMil)), nil
}

type Infections []Infection

func (i Infections) ToCollection() (map[string]map[string]oremi.Collection, error) {
	series := make(map[string]*model.Series)
	for _, infection := range i {
		s, ok := series[infection.Country]
		if !ok {
			s = model.NewSeries(
				label.Day("date"),
				label.Num("cases"),
				label.Num("deaths"),
				label.Num("tests"),
				label.Num("cases/mil"),
			)
		}
		v, err := infection.ToVector()
		if err != nil {
			return nil, fmt.Errorf("could not convert `%v` to vector: %w", infection, err)
		}
		s.Add(v)
		series[infection.Country] = s
	}

	// transform to collections
	collections := make(map[string]oremi.Collection)

	stop := false
	for key, collection := range series {
		if stop {
			continue
		}
		collections[key] = *oremi.New(collection)
		//stop = true
	}
	return map[string]map[string]oremi.Collection{"covid-19": collections}, nil
}

type Date struct {
	time.Time
}

// Convert the internal date as CSV string
func (date *Date) MarshalCSV() (string, error) {
	return date.Time.Format("2006-01-02"), nil
}

func (date *Date) String() string {
	return date.Time.Format("2006-02-01")
}

// Convert the CSV string as internal date
func (date *Date) UnmarshalCSV(csv string) (err error) {
	date.Time, err = time.Parse("2006-01-02", csv)
	return err
}

func Parse(file []byte) Infections {
	cases := []*Infection{}
	if err := gocsv.UnmarshalBytes(file, &cases); err != nil {
		panic(err)
	}

	infections := make([]Infection, len(cases))
	for i, c := range cases {
		infections[i] = *c
	}

	return infections
}
