package main

import (
	"flag"
	"fmt"
	"github.com/dhowden/tag"
	"log"
	"os"
)

func main() {
	var in, out string

	flag.StringVar(&in, "i", "", "give input file")
	flag.StringVar(&out, "o", "hh", "give output file")
	flag.Parse()

	if in == "" {
		log.Fatalln("please give a valid input file")
	}

	if out == "" {
		log.Fatalln("please give a valid output file")
	}

	if _, err := os.Stat(in); err != nil {
		log.Fatalln("couldn't find the input file")
	}

	f, err := os.Open(in)
	if err != nil {
		panic(err)
	}

	m, err := tag.ReadFrom(f)
	if err != nil {
		panic(err)
	}

	fmt.Println("Song's Title:", m.Title())
	fmt.Println("Song's Artist:", m.Artist())
	fmt.Println("Song's Genre:", m.Genre())

}
