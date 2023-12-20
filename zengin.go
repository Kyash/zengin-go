package zengin

import (
	zengin "Kyash/zengin/internal"
)

func ParseToString(reader zengin.Reader) ([][]string, error) {

	transfers, err := zengin.Parse(reader)
	if err != nil {
		return nil, err
	}

	return zengin.ConvertToTable(transfers), nil
}
