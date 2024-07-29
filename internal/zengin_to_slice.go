package internal

import (
	"github.com/Kyash/zengin-go/types"
	"strconv"
)

func ToTable(t []types.Transfer) [][]string {
	var result [][]string
	header := []string{
		"SenderName",
		"RecipientBankCode",
		"RecipientBranchCode",
		"RecipientAccountType",
		"RecipientAccountNumber",
		"RecipientName",
		"Amount",
	}
	result = append(result, header)
	for _, transfer := range t {
		result = append(result, transferToStrings(transfer))
	}
	return result
}

func ToTableJa(t []types.Transfer) [][]string {
	var result [][]string
	header := []string{
		"振込名義人",
		"金融機関コード",
		"支店コード",
		"科目",
		"口座番号",
		"口座名義人",
		"金額",
	}
	result = append(result, header)
	for _, transfer := range t {
		result = append(result, transferToStrings(transfer))
	}
	return result
}

func transferToStrings(t types.Transfer) []string {
	return []string{
		t.SenderName,
		t.RecipientBankCode,
		t.RecipientBranchCode,
		strconv.Itoa(int(t.RecipientAccountType)),
		t.RecipientAccountNumber,
		t.RecipientName,
		strconv.Itoa(int(t.Amount)),
	}
}
