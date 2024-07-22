package zengin

import (
	zengin "github.com/Kyash/zengin-go/internal"
	"github.com/Kyash/zengin-go/types"
)

// Parse Zengin format file and return rows with all fields
func Parse(reader zengin.Reader) ([]types.Transfer, error) {
	transfers, err := zengin.Parse(reader)
	if err != nil {
		return nil, err
	}

	return transfers, nil
}

// ToCSV
// Parse Zengin format file and return a csv like table with field names as below
// SenderName,TransferDate,BankCode,BranchCode,AccountType,AccountNumber,AccountName,Amount
func ToCSV(reader zengin.Reader) ([][]string, error) {

	transfers, err := zengin.Parse(reader)
	if err != nil {
		return nil, err
	}

	return zengin.ToTable(transfers), nil
}

// ToCSVJa
// Parse Zengin format file and return a csv like table with field names as below
// 振込名義人,振込日,金融機関コード,支店コード,科目,口座番号,口座名義人,金額
func ToCSVJa(reader zengin.Reader) ([][]string, error) {

	transfers, err := zengin.Parse(reader)
	if err != nil {
		return nil, err
	}

	return zengin.ToTableJa(transfers), nil
}
