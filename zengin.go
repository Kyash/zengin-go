package zengin

import (
	zengin "github.com/Kyash/zengin-go/internal"
)

// ParseToString /*
// Parse 全銀 format file and return rows
// which consists of fields in this order:
// 振込名義人,振込日,金融機関コード,支店コード,科目,口座番号,口座名義人,金額
func ParseToString(reader zengin.Reader) ([][]string, error) {

	transfers, err := zengin.Parse(reader)
	if err != nil {
		return nil, err
	}

	return zengin.ConvertToTable(transfers), nil
}

// Parse /*
// Parse 全銀 format file and return rows with all fields
func Parse(reader zengin.Reader) {
	// TODO: implement
}
