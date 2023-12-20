package internal

import (
	"errors"
	"fmt"
	"strconv"
)

func Parse(file Reader) ([]Transfer, error) {

	scanner, encoding, err := guessEncoding(file)
	if err != nil {
		return nil, err
	}

	var transfers []Transfer
	var header Header

	for scanner.Scan() {
		line := []rune(scanner.Text())

		var err error
		switch {
		case IsHeader(line):
			var data []Data
			var trailer Trailer

			header, err = parseHeader(line, encoding)
			if err != nil {
				return nil, err
			}

			for scanner.Scan() {
				line = []rune(scanner.Text())

				if !IsData(line) {
					if !IsTrailer(line) {
						return nil, errors.New("unexpected record type: " + string(line))
					}
					trailer, err = parseTrailer(line)
					newTransfers, err := createTransfers(header, data, trailer)
					if err != nil {
						return nil, err
					}
					transfers = append(transfers, newTransfers...)
					break
				}

				record, err := parseData(line)
				if err != nil {
					return nil, err
				}

				data = append(data, record)
			}
		case IsEndRecord(line):
			continue
		default:
			return nil, errors.New("unknown record type: " + string(line))
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(transfers) == 0 {
		return nil, errors.New("no transfers found")
	}

	return transfers, nil
}

func parseHeader(line []rune, encoding Encoding) (Header, error) {
	if len(line) < MinHeaderLength { // Ensure line has enough characters
		return Header{}, errors.New("header line too short")
	}

	header := Header{}

	recordType := string(line[0:1])
	if recordType != "1" {
		return Header{}, errors.New("header record type is not 1")
	}
	header.RecordType = recordType

	categoryCode, err := parseCategoryCode(string(line[1:3]))
	if err != nil {
		return Header{}, err
	}
	header.CategoryCode = categoryCode

	encodingType := string(line[3:4])
	if encoding != EncodingShiftJIS && encodingType == "0" {
		return Header{}, errors.New("unsupported encoding type: " + encodingType)
	}
	header.EncodingType = encodingType

	senderCode, err := parseSenderCode(string(line[4:14]))
	if err != nil {
		return Header{}, errors.New("invalid sender code: " + err.Error())
	}
	header.SenderCode = senderCode

	header.SenderName = string(line[14:54])

	date, err := parseDate(string(line[54:58]))
	if err != nil {
		return Header{}, err
	}
	header.TransactionDate = date

	bankCode, err := parseBankCode(string(line[58:62]))
	if err != nil {
		return Header{}, err
	}
	header.SenderBankCode = bankCode

	header.SenderBankName = string(line[62:77]) // optional

	branchCode, err := parseBranchCode(string(line[77:80]))
	if err != nil {
		return Header{}, err
	}
	header.SenderBankCode = branchCode

	// Fields below are optional

	if len(line) >= 95 {
		header.SenderBranchName = string(line[80:95])
	}

	if len(line) >= 96 {
		accountType, err := parseAccountType(string(line[95:96]))
		if err != nil {
			return Header{}, err
		}
		header.SenderAccountType = accountType
	}

	if len(line) >= 103 {
		header.SenderAccountNumber = string(line[96:103])
	}

	return header, nil
}

func parseData(line []rune) (Data, error) {
	if len(line) < MinDataLength { // Ensure the line is of expected length
		return Data{}, errors.New("data line too short")
	}

	data := Data{}

	recordType := string(line[0:1])
	if recordType != "2" {
		return Data{}, errors.New("data record type is not 2")
	}
	data.RecordType = recordType

	bankCode, err := parseBankCode(string(line[1:5]))
	if err != nil {
		return Data{}, err
	}
	data.RecipientBankCode = bankCode

	data.RecipientBankName = string(line[5:20]) // optional

	branchCode, err := parseBranchCode(string(line[20:23]))
	if err != nil {
		return Data{}, err
	}
	data.RecipientBranchCode = branchCode

	data.RecipientBranchName = string(line[23:38]) // optional

	exchangeOfficeCode := string(line[38:42]) // optional
	if exchangeOfficeCode != "    " {
		if _, err := strconv.Atoi(exchangeOfficeCode); err != nil {
			return Data{}, errors.New("invalid exchange office code: contains non-numeric characters")
		}
	}
	data.ExchangeOfficeCode = exchangeOfficeCode

	accountType, err := parseAccountType(string(line[42:43]))
	if err != nil {
		return Data{}, err
	}
	data.RecipientAccountType = accountType

	accountNumber := string(line[43:50])
	if _, err := strconv.Atoi(accountNumber); err != nil {
		return Data{}, errors.New("invalid account number: contains non-numeric characters")
	}
	data.RecipientAccountNumber = accountNumber

	data.RecipientName = string(line[50:80])

	// Parse transfer amount as integer
	amount, err := strconv.ParseUint(string(line[80:90]), 10, 64)
	if err != nil {
		return Data{}, fmt.Errorf("invalid transfer amount: %v", err)
	}
	data.TransferAmount = amount

	newCode, err := parseNewCode(string(line[90:91])) // unused
	if err != nil {
		return Data{}, err
	}
	data.NewCode = newCode

	// Fields below are optional

	if len(line) >= 111 {
		data.Extra = string(line[91:111])
	}

	if len(line) >= 112 {
		data.TransferCategory = string(line[111:112]) // unused
		if _, err := strconv.Atoi(data.TransferCategory); err != nil {
			return Data{}, errors.New("invalid transfer category: contains non-numeric characters")
		}
	}

	if len(line) >= 113 {
		ediPresent := string(line[112:113])
		if ediPresent == "Y" {
			data.EdiPresent = true
		} else {
			data.EdiPresent = false
		}
	}

	return data, nil
}

func parseTrailer(line []rune) (Trailer, error) {
	if len(line) < MinTrailerLength { // Ensure the line is of expected length
		return Trailer{}, errors.New("trailer line too short")
	}

	trailer := Trailer{
		RecordType: string(line[0:1]),
	}
	if trailer.RecordType != "8" {
		return Trailer{}, errors.New("invalid record type for trailer")
	}

	// Parse TotalCount as integer
	totalCount, err := strconv.Atoi(string(line[1:7]))
	if err != nil {
		return Trailer{}, fmt.Errorf("invalid total count: %v", err)
	}
	trailer.TotalCount = totalCount

	// Parse TotalAmount as integer
	totalAmount, err := strconv.ParseUint(string(line[7:19]), 10, 64)
	if err != nil {
		return Trailer{}, fmt.Errorf("invalid total amount: %v", err)
	}
	trailer.TotalAmount = totalAmount

	return trailer, nil
}
