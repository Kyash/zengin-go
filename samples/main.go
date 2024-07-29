package main

import (
	"fmt"
	"github.com/Kyash/zengin-go"
	"log"
	"os"
)

// This is a sample program that reads a Zengin format file and prints its content.

func main() {
	// Get file name from command line arguments
	var fileName string
	if len(os.Args) > 2 {
		log.Fatal("too many arguments")
	}
	if len(os.Args) < 2 {
		fileName = "sample.txt"
	} else {
		fileName = os.Args[1]
	}

	// Open file and get CSV-like table from Zengin format file with Japanese field names
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	records, err := zengin.ToCSVJa(file)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Transfers (ja):")
	fmt.Println()
	for _, record := range records {
		for _, field := range record {
			print(field + ", ")
		}
		println()
	}

	// Reopen file and get CSV-like table from Zengin format file with English field names
	file, err = os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	records, err = zengin.ToCSV(file)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Transfers (en):")
	fmt.Println()
	for _, record := range records {
		for _, field := range record {
			print(field + ", ")
		}
		println()
	}

	// Reopen file and get rows as a pre-defined structure with all fields from Zengin format file
	file, err = os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	transfers, err := zengin.Parse(file)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Transfers (all fields):")
	fmt.Println()
	for i, transfer := range transfers {
		fmt.Println("Transfer", i+1)
		println("SenderName: ", transfer.SenderName)
		println("SenderCode: ", transfer.SenderCode)
		println("SenderAccountType: ", transfer.SenderAccountType)
		println("SenderBranchName: ", transfer.SenderBranchName)
		println("Transfer Category: ", transfer.TransferCategory)
		println("Date: ", transfer.TransferDate)
		println("Recipient Name: ", transfer.RecipientName)
		println("Recipient Account Type: ", transfer.RecipientAccountType)
		println("Recipient Bank Code", transfer.RecipientBankCode)
		println("Amount: ", transfer.Amount)
		println("Category Code: ", transfer.CategoryCode)
		println("Edi Present? -- ", transfer.EdiPresent)
		println("Total Amount: ", transfer.TotalAmount)
		println("Extra Field: ", transfer.Extra)
		println("Exchange Office Code: ", transfer.ExchangeOfficeCode)
		println("New code: ", transfer.NewCode)
		println("====================================")
		println()
	}
}
