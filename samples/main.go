package main

import (
	"Kyash/zengin"
	"log"
	"os"
)

func main() {
	var fileName string

	if len(os.Args) > 2 {
		log.Fatal("too many arguments")
	}

	if len(os.Args) < 2 {
		fileName = "sample.txt"
	} else {
		fileName = os.Args[1]
	}

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	records, err := zengin.ParseToString(file)
	if err != nil {
		log.Fatal(err)
	}

	for _, record := range records {
		for _, field := range record {
			print(field + ", ")
		}
		println()
	}

}
