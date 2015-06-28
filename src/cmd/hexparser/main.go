package main

import (
	"csvfiles"
	"flag"
)

func main() {
	var (
		filepath, dbconn string
	)

	flag.StringVar(&filepath, "path", "temp.zip", "path to zip file")
	flag.StringVar(&dbconn, "db", "prices.db", "path to sqlite database file")
	flag.Parse()

	r := csvfiles.NewResultData(dbconn)
	r.SetParsedFiles(filepath)
	r.ProcessData()
}
