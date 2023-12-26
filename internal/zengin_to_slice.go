package internal

import (
	"strconv"
)

func ConvertToTable(t []Transfer) [][]string {
	var result [][]string
	header := []string{
		"振込名義人",
		"振込日",
		"金融機関コード",
		"支店コード",
		"科目",
		"口座番号",
		"口座名義人",
		"金額",
		"\n",
	}
	result = append(result, header)
	for _, transfer := range t {
		result = append(result, transferToStrings(transfer))
	}
	return result
}

func transferToStrings(t Transfer) []string {
	// 振込名義人,振込日,金融機関コード,支店コード,科目,口座番号,口座名義人,金額
	return []string{
		t.SenderName,
		t.TransactionDate,
		t.RecipientBankCode,
		t.RecipientBranchCode,
		strconv.Itoa(int(t.RecipientAccountType)),
		t.RecipientAccountNumber,
		t.RecipientName,
		strconv.Itoa(int(t.Amount)),
		"\n",
	}
}
