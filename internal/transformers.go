package internal

import (
	"fmt"
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
	// Transfer Kyashが必要としているデータしかないですが、後で全部入れるようにする。
	for _, block := range data {
		transfer.SenderName = strings.TrimSpace(header.SenderName)
		transfer.TransactionDate = header.TransactionDate
		transfer.RecipientBankCode = block.RecipientBankCode
		transfer.RecipientBranchCode = block.RecipientBranchCode
		transfer.RecipientAccountType = block.RecipientAccountType
		transfer.RecipientAccountNumber = block.RecipientAccountNumber
		transfer.RecipientName = strings.TrimSpace(block.RecipientName)
		transfer.Amount = block.Amount
		transfers = append(transfers, transfer)
	}
	return transfers, nil
}

func sumAmount(data []Data) uint64 {
	var sum uint64
	for _, block := range data {
		sum += block.Amount
	}
	return sum
}
