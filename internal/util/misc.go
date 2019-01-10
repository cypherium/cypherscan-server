package util

import (
	"encoding/json"
	"fmt"
)

// PrintStructInJSON is to help debugging by print the struct in json format
func PrintStructInJSON(obj interface{}) {
	b, err := json.Marshal(obj)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}
