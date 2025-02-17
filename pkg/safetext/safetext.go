// safetext provides methods for determining the use of steganographic
// characters in lines of text. This is an important function for
// remaining informed as a journalist worldwide.
package safetext

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Summary provides user information to the caller about the safetext
// results.
type Summary struct {
	Count       int      `json:"count"`
	Total       int      `json:"total_steganographic"`
	Percent     float32  `json:"percent_steganographic"`
	Positives   []string `json:"positives"`
	Line        string   `json:"original"`
	Appearances string   `json:"appearances"`
}

// SummaryReport provides a way to return results about a slice of
// Summary objects.
type SummaryReport struct {
	Count       int      `json:"count"`
	Total       int      `json:"total_steganographic"`
	Percent     float32  `json:"percent_steganographic"`
	Positives   []string `json:"positives"`
	Original    string   `json:"original"`
	Appearances string   `json:"appearances"`
}

// charSummary records what steganographic character we have encountered and
// the number.
type charSummary struct {
	char  string
	count int
}

// analysisResult provides a way to collect results in between
// composite functions.
type analysisResult struct {
	line     string
	chars    []charSummary
	markedUp string
}

var (
	// Replacement for an init function that ensures data is better
	// behaved when loaded.
	_ = initSafeTextJSON()
)

// initSafeTextJSON loads safetext.json into a structure on compilation.
// init method via:
//
//	https://medium.com/random-go-tips/init-without-init-ebf2f62e7c4a
func initSafeTextJSON() struct{} {
	json.Unmarshal([]byte(safetextJSON), &state)
	return struct{}{}
}

// SafetextAnalysis provides an interface for common processing of
// text data.
type SafetextAnalysis interface {
	Analyse(res *analysisResult) *analysisResult
}

// analysisBase provides a structure for our composite interface. The
// composite interface will run all configured analyses once it has
// been constructed.
type analysisBase struct{}

// Analyse is used to implement an analysis function, e.g. identify
// zero-width characters, count appearances, and so on. As required.
func (safetext analysisBase) Analyse(res *analysisResult) *analysisResult {
	return res
}

type zerowidth struct{ SafetextAnalysis }

// Analyse returns information about zero-width characters.
func (safetext zerowidth) Analyse(res *analysisResult) *analysisResult {
	processText(res, state.ZeroWidthChars)
	safetext.SafetextAnalysis.Analyse(res)
	return res
}

type nonStandardSpaces struct{ SafetextAnalysis }

// Analyse returns information about non-standard space characters.
func (safetext nonStandardSpaces) Analyse(res *analysisResult) *analysisResult {
	processText(res, state.NonStandardSpaces)
	safetext.SafetextAnalysis.Analyse(res)
	return res
}

type homoglyphsEN struct{ SafetextAnalysis }

// Analyse returns information about homoglyphs (EN) characters.
func (safetext homoglyphsEN) Analyse(res *analysisResult) *analysisResult {
	processText(res, state.HomoglyphsEN)
	safetext.SafetextAnalysis.Analyse(res)
	return res
}

// processText allows us to process text against a list of characters
// that we want to identify.
func processText(res *analysisResult, charList map[string]string) *analysisResult {
	if res.markedUp == "" {
		res.markedUp = res.line
	}
	for identifier, str := range charList {
		if !strings.Contains(res.line, str) {
			continue
		}
		cs := charSummary{}
		cs.char = identifier
		cs.count = strings.Count(res.line, str)
		res.chars = append(res.chars, cs)
		res.markedUp = showChars(str, res.markedUp)
	}
	return res
}

// showChars marks where steganographic markings were present in an
// original line of code.
func showChars(char string, line string) string {
	return strings.Replace(line, char, fmt.Sprintf("(%s)", char), -1)
}

func configDefaultAnalysis(analysis SafetextAnalysis) SafetextAnalysis {
	// run third.
	analysis = homoglyphsEN{analysis}
	// run second.
	analysis = nonStandardSpaces{analysis}
	// run first.
	analysis = zerowidth{analysis}
	return analysis
}

// Default config is a helper to return the default english config.
// other config options will be included for other languages or to
// remove english as default as required by users of the library.
func DefaultConfig() SafetextAnalysis {
	analysis := &analysisBase{}
	return configDefaultAnalysis(analysis)
}

// IdentifyNonSafeChara is the primary runner for SafeText. It calls
// configured analyses and then returns a summary structure.
func IdentifyNonSafeChars(analysis SafetextAnalysis, line string) (Summary, error) {
	res := analysisResult{}
	res.line = line
	analysis.Analyse(&res)
	var summary Summary
	summary.Count = len(line)
	total := 0
	positives := []string{}
	for _, value := range res.chars {
		positives = append(positives, value.char)
		total = total + value.count
	}
	summary.Total = total
	summary.Percent = (float32(summary.Total) / float32(summary.Count)) * 100
	summary.Positives = positives
	summary.Line = res.line
	summary.Appearances = res.markedUp
	return summary, nil
}

// SummarizeResults collates the information in a slice of Summary
// structs for a complete picutre of the text file.
func SummarizeResults(results []Summary) SummaryReport {
	report := SummaryReport{}
	for _, value := range results {
		report.Count = report.Count + value.Count
		report.Total = report.Total + value.Total
		report.Positives = append(report.Positives, value.Positives...)
		report.Original = fmt.Sprintf("%s%s\n", report.Original, value.Line)
		report.Appearances = fmt.Sprintf("%s%s\n", report.Appearances, value.Appearances)
	}
	report.Percent = (float32(report.Total) / float32(report.Count)) * 100
	return report
}

// Interface compliance assertions.
var _ SafetextAnalysis = analysisBase{}
var _ SafetextAnalysis = zerowidth{}
var _ SafetextAnalysis = nonStandardSpaces{}

// Custom language interfaces.
var _ SafetextAnalysis = homoglyphsEN{}
