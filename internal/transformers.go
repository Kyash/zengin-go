package internal

import (
	"fmt"
	"github.com/Kyash/zengin-go/types"
	"strings"
)

func createTransfers(header types.Header, data []types.Data, trailer types.Trailer) ([]types.Transfer, error) {
	if header == (types.Header{}) {
		return nil, fmt.Errorf("header is empty")
	}
	if trailer == (types.Trailer{}) {
		return nil, fmt.Errorf("trailer is empty")
	}
	if len(data) != trailer.TotalCount {
		return nil, fmt.Errorf("total count mismatch: %d != %d", len(data), trailer.TotalCount)
	}
	if trailer.TotalAmount != sumAmount(data) {
		return nil, fmt.Errorf("total amount mismatch: %d != %d", trailer.TotalAmount, sumAmount(data))
	}

	var transfers []types.Transfer
	var transfer types.Transfer
	// Transfer Kyashが必要としているデータしかないですが、後で全部入れるようにする。
	for _, block := range data {
		transfer.SenderName = strings.TrimSpace(header.SenderName)
		transfer.TransferDate = header.TransferDate
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

func sumAmount(data []types.Data) uint64 {
	var sum uint64
	for _, block := range data {
		sum += block.Amount
	}
	return sum
}
