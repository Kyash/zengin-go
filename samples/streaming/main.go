package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Kyash/zengin-go"
)

// This demonstrates the streaming API for processing large Zengin files
// without loading all transfers into memory at once.

func main() {
	// Get file name from command line arguments
	var path string
	if len(os.Args) > 2 {
		log.Fatal("too many arguments")
	}
	if len(os.Args) < 2 {
		path = "sample.txt"
	} else {
		path = os.Args[1]
	}

	// Open file and create streaming iterator
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	iterator, err := zengin.NewTransferIterator(file)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Transfers (streaming API - all fields):")
	fmt.Println()

	count := 0

	// Process transfers one by one
	for {
		transfer, err := iterator.Next()
		if err == io.EOF {
			break // Done processing
		}
		if err != nil {
			log.Fatal(err)
		}

		count++

		fmt.Println("Transfer", count)
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
