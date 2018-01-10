package iex

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// the URL preamble for the IEX API
const urlPre = "https://api.iextrading.com/1.0/stock/"

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

// FetchRecords returns slice of ChartRecords for symbol between start and end times
func FetchRecords(symbol string, start time.Time, end time.Time) ([]ChartRecord, error) {
	fullRecords, err := fetch2yRecords(symbol)
	if err != nil {
		return nil, err
	}
	records := make([]ChartRecord, 1)
	// TODO not sure if i want to make the dates 1 sec before midnight, or just use
	// 	text dates for the start and end times
	for i := 0; i < len(fullRecords); i++ {
		t, _ := time.Parse(ChartRecordDateForm, fullRecords[i].Date)
		millis := t.UnixNano() / 1000000
		// using this scheme, start date will be inclusive, end date is non-inclusive
		eod := millis + 86399000
		if eod > (start.UnixNano()/1000000) && eod < (end.UnixNano()/1000000) {
			records = append(records, fullRecords[i])
		}
	}
	return records, nil
}

// FetchRecordsByMillis gets data from IEX API and returns a slice of ChartRecords
func FetchRecordsByMillis(symbol string, startUnix int, endUnix int) ([]ChartRecord, error) {
	fullRecords, err := fetch2yRecords(symbol)
	if err != nil {
		return nil, err
	}
	records := make([]ChartRecord, 1)
	// TODO not sure if i want to make the dates 1 sec before midnight, or just use
	// 	text dates for the start and end times
	for i := 0; i < len(fullRecords); i++ {
		t, _ := time.Parse(ChartRecordDateForm, fullRecords[i].Date)
		millis := t.UnixNano() / 1000000
		// using this scheme, start date will be inclusive, end date is non-inclusive
		eod := int(millis + 86399000)
		if eod > startUnix && eod < endUnix {
			records = append(records, fullRecords[i])
		}
	}
	return records, nil
}

// fetch2yrRecords uses the 2y endpoint
// TODO if requested date is further than 2y then use the 5y endpoint?
func fetch2yRecords(symbol string) ([]ChartRecord, error) {
	requestURL := urlPre + symbol + "/chart/2y"
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
	err = json.Unmarshal(body, &fullRecords)
	if err != nil {
		return nil, err
	}
	return fullRecords, nil
}
