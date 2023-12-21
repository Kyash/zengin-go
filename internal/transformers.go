package internal

import (
	"fmt"
	"strconv"
	"strings"
)

func createTransfers(header Header, data []Data, trailer Trailer) ([]Transfer, error) {
	if header == (Header{}) {
		return nil, fmt.Errorf("header is empty")
	}
	if trailer == (Trailer{}) {
		return nil, fmt.Errorf("trailer is empty")
	}
	if len(data) != trailer.TotalCount {
		return nil, fmt.Errorf("total count mismatch: %d != %d", len(data), trailer.TotalCount)
	}
	if trailer.TotalAmount != sumAmount(data) {
		return nil, fmt.Errorf("total amount mismatch: %d != %d", trailer.TotalAmount, sumAmount(data))
	}

	var transfers []Transfer
	var transfer Transfer
	for _, block := range data {
		transfer.SenderName = strings.TrimSpace(header.SenderName)
		transfer.TransferDate = header.TransactionDate
		transfer.BankCode = block.RecipientBankCode
		transfer.BranchCode = block.RecipientBranchCode
		transfer.AccountType = strconv.Itoa(int(block.RecipientAccountType))
		transfer.AccountNumber = block.RecipientAccountNumber
		transfer.AccountName = strings.TrimSpace(block.RecipientName)
		transfer.Amount = block.TransferAmount
		transfers = append(transfers, transfer)
	}
	return transfers, nil
}

func sumAmount(data []Data) uint64 {
	var sum uint64
	for _, block := range data {
		sum += block.TransferAmount
	}
	return sum
}
