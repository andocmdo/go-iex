package iex

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// the URL preamble for the IEX API
const urlPre = "https://api.iextrading.com/1.0/stock/"
const urlAPIEndPoint = "/chart/2y"

// ChartRecord is the format of JSON data returned from IEX API
type ChartRecord struct {
	Date             string  `json:"date"`
	Open             float64 `json:"open"`
	High             float64 `json:"high"`
	Low              float64 `json:"low"`
	Close            float64 `json:"close"`
	Volume           float64 `json:"volume"`
	UnadjustedVolume float64 `json:"unadjustedVolume"`
	Change           float64 `json:"change"`
	ChangePercent    float64 `json:"changePercent"`
	Vwap             float64 `json:"vwap"`
	Label            string  `json:"label"`
	ChangeOverTime   float64 `json:"changeOverTime"`
}

// ChartRecordDateForm is used to set the date format for parsing from JSON reply
const ChartRecordDateForm = "2006-01-02"

// FetchRecordsByMillis gets data from IEX API and returns a slice of ChartRecords
func FetchRecordsByMillis(symbol string, startUnix int, endUnix int) ([]ChartRecord, error) {
	requestURL := urlPre + symbol + urlAPIEndPoint
	resp, err := http.Get(requestURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var fullRecords []ChartRecord
	records := make([]ChartRecord, 1)
	err = json.Unmarshal(body, &fullRecords)
	if err != nil {
		return nil, err
	}
	// TODO not sure if i want to make the dates 1 sec before midnight, or just use
	// 	text dates for the start and end times
	for i := 0; i < len(fullRecords); i++ {
		t, _ := time.Parse(ChartRecordDateForm, fullRecords[i].Date)
		millis := t.UnixNano() / 1000000
		eod := int(millis + 86399000)
		// using this scheme, start date will be inclusive, end date is non-inclusive
		if eod > startUnix && eod < endUnix {
			records = append(records, fullRecords[i])
		}
	}
	return records, nil
}
