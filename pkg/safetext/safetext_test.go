package safetext

import (
	"bufio"
	"os"
	"testing"
)

type testFilesData struct {
	filepath      string
	expectedCount int
	expectedSteg  int
}

var testFiles = []testFilesData{
	{
		"testdata/digipres",
		50,
		10,
	},
	{
		"testdata/email_example",
		634,
		200,
	},
	{
		"testdata/one_word",
		340,
		102,
	},
	{
		"testdata/safetext_source",
		156,
		6,
	},
	{
		"testdata/no_steganography",
		38,
		0,
	},
}

// TestIntegration runs some basic tests over the testdata in this
// repo..
func TestIntegration(t *testing.T) {
	for _, v := range testFiles {
		file, _ := os.Open(v.filepath)
		analysis := DefaultConfig()
		scanner := bufio.NewScanner(file)
		all := []Summary{}
		for scanner.Scan() {
			line := scanner.Text()
			summary, _ := IdentifyNonSafeChars(analysis, line)
			if summary.Count == 0 {
				continue
			}
			all = append(all, summary)
		}
		res := SummarizeResults(all)
		if res.Count != v.expectedCount {
			t.Errorf(
				"expected character count is wrong ('%s'): '%d', expected: '%d'",
				v.filepath, res.Count, v.expectedCount,
			)
		}
		if res.Total != v.expectedSteg {
			t.Errorf("expected steg count is wrong ('%s'): '%d', expected: '%d'",
				v.filepath, res.Total, v.expectedSteg,
			)
		}
	}
}
