package safetext_test

import (
	"encoding/json"
	"fmt"

	"github.com/ross-spencer/safetext/pkg/safetext"
)

var example string = "supercalifragilist\u2060icexpialidotious"

func ExampleIdentifyNonSafeChars() {
	analysis := safetext.DefaultConfig()
	res, err := safetext.IdentifyNonSafeChars(analysis, example)
	if err == nil {
		// handle err
	}
	resJSON, _ := json.MarshalIndent(res, "", "  ")
	fmt.Println(string(resJSON))
}
