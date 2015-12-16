package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/mattn/go-gimei"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Fake struct {
	name    *gimei.Name
	address *gimei.Address
}
type Generator func(*Fake) string

func genParrot(s string) Generator {
	return func(*Fake) string {
		return s
	}
}

func guessSeparator(format string) string {
	if strings.Contains(format, "\t") {
		return "\t"
	}
	return ","
}

func parseFormat(format string, separator string) ([]string, []Generator) {

	columns := strings.Split(format, separator)
	gens := make([]Generator, 0, len(columns))

	for _, column := range columns {
		var g Generator
		switch column {
		case "姓名", "氏名", "名前":
			g = func(f *Fake) string { return f.name.Kanji() }
		case "ふりがな", "せいめい", "なまえ":
			g = func(f *Fake) string { return f.name.Hiragana() }
		case "フリガナ", "セイメイ", "ナマエ":
			g = func(f *Fake) string { return f.name.Katakana() }
		case "姓", "氏":
			g = func(f *Fake) string { return f.name.Last.Kanji() }
		case "せい":
			g = func(f *Fake) string { return f.name.Last.Hiragana() }
		case "セイ":
			g = func(f *Fake) string { return f.name.Last.Katakana() }
		case "名":
			g = func(f *Fake) string { return f.name.First.Kanji() }
		case "めい":
			g = func(f *Fake) string { return f.name.First.Hiragana() }
		case "メイ":
			g = func(f *Fake) string { return f.name.First.Katakana() }
		case "住所":
			g = func(f *Fake) string { return f.address.Kanji() }
		case "じゅうしょ":
			g = func(f *Fake) string { return f.address.Hiragana() }
		case "ジュウショ":
			g = func(f *Fake) string { return f.address.Katakana() }
		default:
			g = genParrot(column)
		}
		gens = append(gens, g)
	}
	return columns, gens
}

func newRow(generators []Generator) []string {
	f := &Fake{
		name:    gimei.NewName(),
		address: gimei.NewAddress(),
	}

	row := make([]string, len(generators))
	for i, g := range generators {
		row[i] = g(f)
	}
	return row
}

func output(number int, header []string, generators []Generator, separator string) {
	w := csv.NewWriter(os.Stdout)
	for _, r := range separator {
		w.Comma = r
	}

	err := w.Write(header)
	if err != nil {
		fmt.Printf("Cant write header. %s\n", err)
		return
	}

	for i := 0; i < number; i++ {
		w.Write(newRow(generators))
		if err != nil {
			w.Flush()
			fmt.Printf("Cant write gimei. %d/%d %s\n", i, number, err)
			return
		}

		if i%10000 == 0 {
			w.Flush()
		}
	}
	w.Flush()
}

func main() {
	flagN := kingpin.Flag("number", "Number of lines.").Short('n').Int()
	flagS := kingpin.Flag("separator", "Column separator.").Short('s').String()
	arg := kingpin.Arg("header", "CSV Header as format.").Required().String()
	kingpin.Parse()

	number := *flagN
	separator := *flagS
	format := *arg
	if separator == "" {
		separator = guessSeparator(format)
	}

	header, generators := parseFormat(format, separator)
	output(number, header, generators, separator)
}
