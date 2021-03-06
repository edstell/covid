package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/edstell/covid/covid"
)

func main() {
	client := covid.NewClient()
	result, err := client.GetData(
		context.Background(),
		&structure{},
		covid.FormatJSON,
		covid.AreaTypeUTLA(),
		covid.AreaName("lewisham"),
		covid.Date(time.Now().Add(-1*24*time.Hour)),
	)
	if err != nil {
		log.Fatal(err)
	}
	bytes, err := ioutil.ReadAll(result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(string(bytes))
}

type structure struct{}

func (s *structure) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"date":     "date",
		"areaName": "areaName",
		"newCases": "newCasesByPublishDate",
	})
}
