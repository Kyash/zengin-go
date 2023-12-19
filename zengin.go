package zengin

import (
	zengin "Kyash/zengin/internal"
	"bufio"
	"errors"
	"io"
	"strings"
)

type Reader interface {
	io.Reader
}

func ParseFile(reader Reader) ([][]string, error) {

	transfers, err := parse(reader)
	if err != nil {
		return nil, err
	}

	return zengin.ConvertToTable(transfers), nil
}

func parse(reader Reader) ([]zengin.Transfer, error) {
	scanner := bufio.NewScanner(reader)

	var transfers []zengin.Transfer
	var endFound bool

	var data []zengin.Data
	var header zengin.Header
	var trailer zengin.Trailer

	for scanner.Scan() {
		line := scanner.Text()
		var record zengin.Data

		var err error
		switch {
		case zengin.IsHeader(line):
			header, err = parseHeader(line)
		case zengin.IsData(line):
			record, err = parseData(line)
			data = append(data, record)
		case zengin.IsTrailer(line):
			trailer, err = parseTrailer(line)
		case zengin.IsEndRecord(line):
			endFound = true
			transfers, err = createTransfers(header, data, trailer), nil
		default:
			return nil, errors.New("unknown record type: " + line)
		}

		if err != nil {
			return nil, err
		}

	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if !endFound {
		return nil, errors.New("end reached without end record")
	}
	if len(transfers) == 0 {
		return nil, errors.New("no transfers found")
	}

	return transfers, nil
}

func createTransfers(header zengin.Header, data []zengin.Data, trailer zengin.Trailer) []zengin.Transfer {
	return []zengin.Transfer{}
}

func parseHeader(line string) (zengin.Header, error) {
	if len(line) < 120 { // Ensure line has enough characters
		return zengin.Header{}, errors.New("header line too short")
	}

	header := zengin.Header{}

	recordType := line[0:1]
	if recordType != "1" {
		return zengin.Header{}, errors.New("header record type is not 1")
	}
	header.RecordType = recordType

	categoryCode := line[1:3]
	switch categoryCode {
	case "21":
		header.CategoryCode = zengin.CategoryCodeCombination
	case "11", "71":
		header.CategoryCode = zengin.CategoryCodePayment
	case "12", "72":
		header.CategoryCode = zengin.CategoryCodeBonus
	default:
		return zengin.Header{}, errors.New("unknown category code: " + categoryCode)
	}

	encodingType := line[3:4]
	if encodingType != "0" {
		return zengin.Header{}, errors.New("unsupported encoding type: " + encodingType)
	}
	header.EncodingType = encodingType

	senderCode, err := parseSenderCode(line[4:14])
	if err != nil {
		return zengin.Header{}, errors.New("invalid sender code: " + err)
	}
	header.SenderCode = senderCode

	senderName := line[14:54] // Trim spaces for names
	header.SenderName = strings.TrimSpace(senderName)

	transactionDate := line[54:58]
	header.TransactionDate = transactionDate

	bankCode := line[58:62]

	bankName := strings.trimspace(line[62:77])

	branchCode := line[77:80]

	branchName := strings.trimspace(line[80:95])

	accountType := line[95:96]

	accountNumber := line[96:103]

	dummy := line[103:120] // assuming 17 characters for the dummy field

	return zengin.Header{}, nil
}

func parseSenderCode(senderCode string) (string, error) {
	if len(senderCode) != 10 {
		return "", errors.New("sender code must be 10 digits")
	}
	return senderCode, nil
}

func parseData(line string) (zengin.Data, error) {
	return zengin.Data{}, nil
}

func parseTrailer(line string) (zengin.Trailer, error) {
	return zengin.Trailer{}, nil
}

func parseEndRecord(line string) (zengin.EndRecord, error) {
	return zengin.EndRecord{}, nil
}
