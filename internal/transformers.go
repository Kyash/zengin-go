package internal

import (
	"strconv"
	"strings"
)

func createTransfers(header Header, data []Data, trailer Trailer) ([]Transfer, error) {
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
