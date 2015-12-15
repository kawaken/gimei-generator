package main

import (
	"encoding/csv"
	"os"
	"strings"

	"github.com/mattn/go-gimei"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Toker struct {
	name    *gimei.Name
	address *gimei.Address
}
type Generator func(*Toker) string

func genParrot(s string) Generator {
	return func(*Toker) string {
		return s
	}
}

func guessSeparator(format string) string {
	if strings.Contains(format, "	") {
		return "	"
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
			g = func(t *Toker) string { return t.name.Kanji() }
		case "ふりがな", "せいめい", "なまえ":
			g = func(t *Toker) string { return t.name.Hiragana() }
		case "フリガナ", "セイメイ", "ナマエ":
			g = func(t *Toker) string { return t.name.Katakana() }
		case "姓", "氏":
			g = func(t *Toker) string { return t.name.Last.Kanji() }
		case "せい":
			g = func(t *Toker) string { return t.name.Last.Hiragana() }
		case "セイ":
			g = func(t *Toker) string { return t.name.Last.Katakana() }
		case "名":
			g = func(t *Toker) string { return t.name.First.Kanji() }
		case "めい":
			g = func(t *Toker) string { return t.name.First.Hiragana() }
		case "メイ":
			g = func(t *Toker) string { return t.name.First.Katakana() }
		case "住所":
			g = func(t *Toker) string { return t.address.Kanji() }
		case "じゅうしょ":
			g = func(t *Toker) string { return t.address.Hiragana() }
		case "ジュウショ":
			g = func(t *Toker) string { return t.address.Katakana() }
		default:
			g = genParrot(column)
		}
		gens = append(gens, g)
	}
	return columns, gens
}

func newRow(generators []Generator) []string {
	t := &Toker{
		name:    gimei.NewName(),
		address: gimei.NewAddress(),
	}

	row := make([]string, len(generators))
	for i, g := range generators {
		row[i] = g(t)
	}
	return row
}

func output(number int, header []string, generators []Generator) {
	w := csv.NewWriter(os.Stdout)
	w.Write(header)

	for i := 0; i < number; i++ {
		w.Write(newRow(generators))
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
	output(number, header, generators)
}