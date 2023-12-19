package internal

func ConvertToTable(t []Transfer) [][]string {
	var result [][]string
	for _, transfer := range t {
		result = append(result, transferToStrings(transfer))
	}
	return nil
}

func transferToStrings(t Transfer) []string {
	return nil
}
