package main

import (
	"encoding/json"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {

	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error when marshal object to json string"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// assuming items in numbersAlreadyGot is order by the number desc
func getMissedNumbers(startedNumber int64, pageSize int, numbersAlreadyGot []int64) []int64 {
	missedLen := pageSize - len(numbersAlreadyGot)
	if len(numbersAlreadyGot) == pageSize {
		return []int64{}
	}
	notChecked := numbersAlreadyGot[:]
	missed := make([]int64, 0, missedLen)
	for number := startedNumber; number > startedNumber-int64(pageSize); number-- {
		if len(notChecked) > 0 && number == notChecked[0] {
			notChecked = notChecked[1:]
			continue
		}
		missed = append(missed, number)
	}
	return missed
}

func cors(handleFunc func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		handleFunc(w, r)
	}
}
