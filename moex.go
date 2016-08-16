package main

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const URL = "https://www.moex.com/iss/engines/currency/markets/selt/boardgroups/13/securities.xml?iss.meta=off&iss.only=marketdata&security_collection=173&lang=RU"

var last MOEXRow

type MOEXRow struct {
	High  float32 `xml:"HIGH,attr"`
	Low   float32 `xml:"LOW,attr"`
	Last  float32 `xml:"LAST,attr"`
	SecID string  `xml:"SECID,attr"`
}

type MOEXResponse struct {
	RowList []MOEXRow `xml:"data>rows>row"`
}

func GetUpdateChannel(SecID string) (<-chan MOEXRow, error) {
	rowChan := make(chan MOEXRow)

	go func() {
		for {
			resp, err := http.Get(URL)

			if err != nil {
				log.Print("get error: ", err)
				continue
			}

			body, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				log.Print("read error: ", err)
				continue
			}

			v := MOEXResponse{}
			err = xml.Unmarshal(body, &v)

			if err != nil {
				log.Printf("unmarshal error: %v", err)
				continue
			}

			for _, row := range v.RowList {
				if SecID == row.SecID {
					last = row
					rowChan <- row
					break
				}
			}

			time.Sleep(3 * time.Minute)
		}
	}()

	return rowChan, nil
}

func getLastRow() MOEXRow {
	return last
}
