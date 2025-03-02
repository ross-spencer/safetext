package safetext

import (
	"encoding/json"
	"fmt"
)

var example string = "supercalifragilist\u2060icexpialidotious"

func ExampleIdentifyNonSafeChars() {
	analysis := DefaultConfig()
	res, err := IdentifyNonSafeChars(analysis, example)
	if err == nil {
		// handle err
	}
	resJSON, _ := json.MarshalIndent(res, "", "  ")
	fmt.Println(string(resJSON))
}
