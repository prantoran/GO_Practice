package main

import (
	"fmt"
	"io/ioutil"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"gopkg.in/Iwark/spreadsheet.v2"
)

const (
	SpreadSheetID = "1mYiA2T4_QTFUkAXk0BE3u7snN2o5FgSRqxmRrn_Dzh4"
)

func main() {
	data, err := ioutil.ReadFile("client_secret.json")
	checkError(err)
	conf, err := google.JWTConfigFromJSON(data, spreadsheet.Scope)
	checkError(err)
	client := conf.Client(context.TODO())

	service := spreadsheet.NewServiceWithClient(client)
	spreadsheet, err := service.FetchSpreadsheet("1mYiA2T4_QTFUkAXk0BE3u7snN2o5FgSRqxmRrn_Dzh4")
	checkError(err)
	sheet, err := spreadsheet.SheetByIndex(0)
	checkError(err)
	for _, row := range sheet.Rows {
		for _, cell := range row {
			fmt.Println(cell.Value)
		}
	}

	// Update cell content
	sheet.Update(0, 0, "hogehoge")

	// Make sure call Synchronize to reflect the changes
	err = sheet.Synchronize()
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
