package safetext

import _ "embed"

//go:embed safetext.json
var safetextJSON string

// SafetextState provides a way to persist the Safetext character list
// in memory to process data against.
type SafetextState struct {
	ZeroWidthChars    map[string]string `json:"ZERO_WIDTH_CHARS"`
	NonStandardSpaces map[string]string `json:"NON_STANDARD_SPACES"`
	HomoglyphsEN      map[string]string `json:"HOMOGLYPHS_EN"`
}
