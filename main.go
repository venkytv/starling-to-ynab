package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gocarina/gocsv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// InRecord is used to load transaction records from Starling Bank statements
type InRecord struct {
	Date         string `csv:"Date"`
	CounterParty string `csv:"Counter Party"`
	Reference    string `csv:"Reference"`
	Type         string `csv:"Type"`
	Amount       string `csv:"Amount (GBP)"`
	Balance      string `csv:"Balance (GBP)"`
	Category     string `csv:"Spending Category"`
	Notes        string `csv:"Notes"`
}

// OutRecord is used to dump transaction records in YNAB format
type OutRecord struct {
	Date   string `csv:"Date"`
	Payee  string `csv:"Payee"`
	Memo   string `csv:"Memo"`
	Amount string `csv:"Amount"`
}

func loadCsv(file string) {
	f, err := os.Open(file)
	check(err)
	defer f.Close()

	records := []*InRecord{}
	err = gocsv.UnmarshalFile(f, &records)
	check(err)

	out := []*OutRecord{}

	for _, record := range records {
		out = append(out, &OutRecord{
			Date:   record.Date,
			Payee:  record.CounterParty,
			Memo:   record.Reference,
			Amount: record.Amount,
		})
	}

	csv, err := gocsv.MarshalString(&out)
	check(err)
	fmt.Println(csv)
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s <csv> [<csv>...]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}
	for _, arg := range flag.Args() {
		loadCsv(arg)
	}
}
