package zengin

import (
	zengin "Kyash/zengin/internal"
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"io"
	"log"
	"strconv"
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

func guessEncoding(file Reader) (*bufio.Scanner, zengin.Encoding, error) {
	var encoding zengin.Encoding
	reader := bufio.NewReader(file)
	peekBytes, err := reader.Peek(1024)
	if err != nil && err != io.EOF {
		log.Fatal("couldn't read from file", "error", err)
		return nil, zengin.EncodingUndefined, err
	}

	// Ignore "certain" (3rd value), as during testing it was always false, even though it correctly detects utf-8.
	_, name, _ := charset.DetermineEncoding(peekBytes, "")

	var scanner *bufio.Scanner
	switch name {
	case "utf-8":
		encoding = zengin.EncodingUTF8
		scanner = bufio.NewScanner(reader)
	// Shift-JIS can't be reliably detected, so we'll assume it's Shift-JIS if it's not UTF-8.
	default:
		encoding = zengin.EncodingShiftJIS
		scanner = bufio.NewScanner(transform.NewReader(reader, japanese.ShiftJIS.NewDecoder()))
	}

	return scanner, encoding, nil
}

func parse(file Reader) ([]zengin.Transfer, error) {

	scanner, encoding, err := guessEncoding(file)
	if err != nil {
		return nil, err
	}

	var transfers []zengin.Transfer
	var header zengin.Header

	for scanner.Scan() {
		line := []rune(scanner.Text())

		var err error
		switch {
		case zengin.IsHeader(line):
			var data []zengin.Data
			var trailer zengin.Trailer

			header, err = parseHeader(line, encoding)
			if err != nil {
				return nil, err
			}

			for scanner.Scan() {
				line = []rune(scanner.Text())

				if !zengin.IsData(line) {
					if !zengin.IsTrailer(line) {
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
		case zengin.IsEndRecord(line):
			continue
		default:
			return nil, errors.New("unknown record type: " + string(line))
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if !zengin.IsEndRecord([]rune(scanner.Text())) {
		return nil, errors.New("end reached without end record")
	}
	if len(transfers) == 0 {
		return nil, errors.New("no transfers found")
	}

	return transfers, nil
}

func createTransfers(header zengin.Header, data []zengin.Data, trailer zengin.Trailer) ([]zengin.Transfer, error) {
	return []zengin.Transfer{}, nil
}

func parseHeader(line []rune, encoding zengin.Encoding) (zengin.Header, error) {
	if len(line) < zengin.MinHeaderLength { // Ensure line has enough characters
		return zengin.Header{}, errors.New("header line too short")
	}

	header := zengin.Header{}

	recordType := string(line[0:1])
	if recordType != "1" {
		return zengin.Header{}, errors.New("header record type is not 1")
	}
	header.RecordType = recordType

	categoryCode, err := parseCategoryCode(string(line[1:3]))
	fmt.Printf("categoryCode: %v\n", categoryCode)
	if err != nil {
		return zengin.Header{}, err
	}
	header.CategoryCode = categoryCode

	encodingType := string(line[3:4])
	fmt.Printf("encodingType: %v\n", encodingType)
	if encoding != zengin.EncodingShiftJIS && encodingType == "0" {
		return zengin.Header{}, errors.New("unsupported encoding type: " + encodingType)
	}
	header.EncodingType = encodingType

	senderCode, err := parseSenderCode(string(line[4:14]))
	fmt.Printf("senderCode: %v\n", senderCode)
	if err != nil {
		return zengin.Header{}, errors.New("invalid sender code: " + err.Error())
	}
	header.SenderCode = senderCode

	header.SenderName = string(line[14:54])
	fmt.Printf("senderName: %v\n", header.SenderName)

	header.TransactionDate = string(line[54:58])
	fmt.Printf("transactionDate: %v\n", header.TransactionDate)

	bankCode := string(line[58:62])
	fmt.Printf("bankCode: %v\n", bankCode)
	if _, err := strconv.Atoi(bankCode); err != nil {
		return zengin.Header{}, errors.New("invalid bank code: contains non-numeric characters")
	}
	header.BankCode = bankCode

	header.BankName = string(line[62:77])
	fmt.Printf("bankName: %v\n", header.BankName)

	branchCode := string(line[77:80])
	fmt.Printf("branchCode: %v\n", branchCode)
	if _, err := strconv.Atoi(branchCode); err != nil {
		return zengin.Header{}, errors.New("invalid branch code: contains non-numeric characters")
	}
	header.BankCode = bankCode

	header.BranchName = string(line[80:95])
	fmt.Printf("branchName: %v\n", header.BranchName)

	accountType, err := parseAccountType(string(line[95:96]))
	fmt.Printf("accountType: %v\n", accountType)
	if err != nil {
		return zengin.Header{}, err
	}
	header.AccountType = accountType

	header.AccountNumber = string(line[96:103])
	fmt.Printf("accountNumber: %v\n", header.AccountNumber)

	return zengin.Header{}, nil
}

func parseData(line []rune) (zengin.Data, error) {
	if len(line) < zengin.MinDataLength { // Ensure the line is of expected length
		return zengin.Data{}, errors.New("data line too short")
	}

	data := zengin.Data{}

	recordType := string(line[0:1])
	if recordType != "2" {
		return zengin.Data{}, errors.New("data record type is not 2")
	}
	data.RecordType = recordType

	bankCode := string(line[1:5])
	if _, err := strconv.Atoi(bankCode); err != nil {
		return zengin.Data{}, errors.New("invalid bank code: contains non-numeric characters")
	}
	data.RecipientBankCode = bankCode

	data.RecipientBankName = string(line[5:20])
	recipientBranchCode := string(line[20:23])
	if _, err := strconv.Atoi(recipientBranchCode); err != nil {
		return zengin.Data{}, errors.New("invalid recipient branch code: contains non-numeric characters")
	}
	data.RecipientBranchCode = recipientBranchCode

	data.RecipientBranchName = string(line[23:38])

	exchangeOfficeCode := string(line[38:42]) // unused
	if _, err := strconv.Atoi(exchangeOfficeCode); err != nil {
		return zengin.Data{}, errors.New("invalid exchange office code: contains non-numeric characters")
	}
	data.ExchangeOfficeCode = exchangeOfficeCode

	accountType, err := parseAccountType(string(line[42:43]))
	if err != nil {
		return zengin.Data{}, err
	}
	data.AccountType = accountType

	accountNumber := string(line[43:50])
	if _, err := strconv.Atoi(accountNumber); err != nil {
		return zengin.Data{}, errors.New("invalid account number: contains non-numeric characters")
	}

	data.RecipientName = string(line[50:80])

	// Parse transfer amount as integer
	amount, err := strconv.ParseUint(string(line[80:90]), 10, 64)
	if err != nil {
		return zengin.Data{}, fmt.Errorf("invalid transfer amount: %v", err)
	}
	data.TransferAmount = amount

	data.NewCode = string(line[90:91]) // unused
	data.CustomerCode1 = string(line[91:101])
	data.CustomerCode2 = string(line[101:111])
	data.EDIInformation = string(line[111:131])
	data.TransferCategory = string(line[131:132]) // unused
	if _, err := strconv.Atoi(data.TransferCategory); err != nil {
		return zengin.Data{}, errors.New("invalid transfer category: contains non-numeric characters")
	}
	data.Identification = string(line[132:133])

	return data, nil
}

func parseTrailer(line []rune) (zengin.Trailer, error) {
	if len(line) < zengin.MinTrailerLength { // Ensure the line is of expected length
		return zengin.Trailer{}, errors.New("trailer line too short")
	}

	trailer := zengin.Trailer{
		RecordType: string(line[0:1]),
	}
	if trailer.RecordType != "8" {
		return zengin.Trailer{}, errors.New("invalid record type for trailer")
	}

	// Parse TotalCount as integer
	totalCount, err := strconv.Atoi(string(line[1:7]))
	if err != nil {
		return zengin.Trailer{}, fmt.Errorf("invalid total count: %v", err)
	}
	trailer.TotalCount = totalCount

	// Parse TotalAmount as integer
	totalAmount, err := strconv.ParseUint(string(line[7:19]), 10, 64)
	if err != nil {
		return zengin.Trailer{}, fmt.Errorf("invalid total amount: %v", err)
	}
	trailer.TotalAmount = totalAmount

	return trailer, nil
}

// Helper functions

func parseCategoryCode(categoryCode string) (zengin.CategoryCode, error) {
	switch categoryCode {
	case "21":
		return zengin.CategoryCodeCombination, nil
	case "11", "71":
		return zengin.CategoryCodePayment, nil
	case "12", "72":
		return zengin.CategoryCodeBonus, nil
	default:
		return zengin.CategoryCodeUndefined, errors.New("unknown category code: " + categoryCode)
	}
}

func parseSenderCode(senderCode string) (string, error) {
	if len(senderCode) != 10 {
		return "", errors.New("sender code must be 10 digits")
	}
	return senderCode, nil
}

func parseAccountType(accountType string) (zengin.AccountType, error) {
	switch accountType {
	case "1":
		return zengin.AccountTypeRegular, nil
	case "2":
		return zengin.AccountTypeChecking, nil
	case "4":
		return zengin.AccountTypeSavings, nil
	default:
		return zengin.AccountTypeUndefined, errors.New("invalid account type: " + accountType)
	}
}
