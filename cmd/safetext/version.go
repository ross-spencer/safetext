package main

import "fmt"

var appname string
var version string
var commit string = "0000000000000000000000000000000000000000"
var date string

func getVersion() string {
	if version == "" {
		return "0.0.0-dev"
	}
	return fmt.Sprintf("%s/%s (commit: '%s' (%s))", appname, version, commit[:6], date)
}
