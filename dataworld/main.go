package main

import (
	"fmt"
	"os"

	"github.com/datadotworld/dwapi-go/dwapi"
)

func main() {
	// new client
	token := os.Getenv("DW_AUTH_TOKEN")
	dw := dwapi.NewClient(token)

	// sql query
	owner := "jonloyens"
	datasetid := "intermediate-data-world"
	acceptType := "text/csv"
	savePath := "jonloyens-intermediate-data-world.csv"
	requestBody := &dwapi.SQLQueryRequest{
		Query:              "SELECT * FROM fatal_police_shootings_data",
		IncludeTableSchema: false,
	}
	response, err := dw.Query.ExecuteSQLAndSave(owner, datasetid,
		acceptType, savePath, requestBody)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Query.ExecuteSQL() returned an error:", err)
	}
	fmt.Println(response.Message)
}
