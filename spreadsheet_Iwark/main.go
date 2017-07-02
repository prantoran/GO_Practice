package main

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"gopkg.in/Iwark/spreadsheet.v2"
)

const (
	SpreadSheetID = "1N4DJ8RiUV00E3EwZP4q_A7qR4F9BD5mhtapLJn7s6gI"
)

func main() {
	data, err := ioutil.ReadFile("client_secret.json")
	checkError(err)
	conf, err := google.JWTConfigFromJSON(data, spreadsheet.Scope)
	checkError(err)
	fmt.Printf("conf type: %T \n\n", conf)
	fmt.Printf("email: %v\nexpires: %v\nprivatekey: %v\nprivatekeyid: %v\n", conf.Email, conf.Expires, conf.PrivateKey, conf.PrivateKeyID)
	fmt.Printf("scope: %v\nsubject: %v\n,tokenurl: %v\n", conf.Scopes, conf.Subject, conf.TokenURL)

	n := bytes.Index(conf.PrivateKey, []byte{0})
	fmt.Printf("n: %v\n", n)
	s := string(conf.PrivateKey[:1690])
	fmt.Printf("s: %v \n", s)

	client := conf.Client(context.TODO())

	service := spreadsheet.NewServiceWithClient(client)
	spreadsheet, err := service.FetchSpreadsheet(SpreadSheetID)
	checkError(err)
	sheet, err := spreadsheet.SheetByIndex(0)
	checkError(err)
	for _, row := range sheet.Rows {
		for _, cell := range row {
			fmt.Printf(cell.Value)
		}
		fmt.Printf("\n")
	}

	//Update cell content
	sheet.Update(0, 4, "sakura")

	//Make sure call Synchronize to reflect the changes
	err = sheet.Synchronize()
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
