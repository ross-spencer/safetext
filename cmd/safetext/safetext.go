package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/ross-spencer/safetext/pkg/logformatter"
	"github.com/ross-spencer/safetext/pkg/safetext"
)

var (
	inputFile   string
	versionFlag bool
)

func setLogger(utc bool) {
	lw := new(logformatter.LogWriter)
	lw.Appname = "SafeText"
	lw.UTC = utc
	log.SetOutput(lw)
}

func init() {
	setLogger(true)
	flag.StringVar(&inputFile, "f", "", "file to scan")
	flag.BoolVar(&versionFlag, "version", false, "return version")
}

func handlePipedInput() string {
	reader := bufio.NewReader(os.Stdin)
	var output []rune
	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		output = append(output, input)
	}
	return string(output)
}

func isPipeInput() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		log.Println("problem determining input type")
		os.Exit(1)
	}
	if (info.Mode() & os.ModeNamedPipe) != 0 {
		return true
	}
	return false
}

func safetext_pipe_runner(data string) {
	sc := bufio.NewScanner(strings.NewReader(data))
	analysis := safetext.DefaultConfig()
	all := []safetext.Summary{}
	for sc.Scan() {
		line := sc.Text()
		summary, _ := safetext.IdentifyNonSafeChars(analysis, line)
		all = append(all, summary)
	}
	res := safetext.SummarizeResults(all)
	resJSON, _ := json.MarshalIndent(res, "", "  ")
	fmt.Println(string(resJSON))
}

func safetext_file_runner() {
	file, err := os.Open(inputFile)
	if err != nil {
		log.Printf("cannot process input: %s", err)
	}
	defer file.Close()
	analysis := safetext.DefaultConfig()
	scanner := bufio.NewScanner(file)
	all := []safetext.Summary{}
	for scanner.Scan() {
		line := scanner.Text()
		summary, _ := safetext.IdentifyNonSafeChars(analysis, line)
		if summary.Count == 0 {
			continue
		}
		all = append(all, summary)
	}
	res := safetext.SummarizeResults(all)
	resJSON, _ := json.MarshalIndent(res, "", "  ")
	fmt.Println(string(resJSON))
}

func main() {
	if isPipeInput() {
		data := handlePipedInput()
		safetext_pipe_runner(data)
		os.Exit(0)
	}
	flag.Parse()
	if versionFlag {
		fmt.Fprintf(os.Stderr, "%s \n", getVersion())
		os.Exit(0)
	} else if flag.NFlag() == 0 {
		fmt.Fprintln(os.Stderr, "safetext: identify steganographic characters in text content")
		fmt.Fprintln(os.Stderr, "usage: safetext <input file>              ")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "output: [JSON]   {safetext summary}")
		fmt.Fprintf(os.Stderr, "output: [STRING] '%s'\n\n", version)
		flag.Usage()
		os.Exit(0)
	}
	safetext_file_runner()
}
