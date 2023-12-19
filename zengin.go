package zengin

import (
	zengin "Kyash/zengin/internal"
	"bufio"
	"errors"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"io"
	"log"
	"strconv"
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
		line := scanner.Text()

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
				line = scanner.Text()

				if !zengin.IsData(line) {
					if !zengin.IsTrailer(line) {
						return nil, errors.New("unexpected record type: " + line)
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
			return nil, errors.New("unknown record type: " + line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if !zengin.IsEndRecord(scanner.Text()) {
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

func parseHeader(line string, encoding zengin.Encoding) (zengin.Header, error) {
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
	if encoding == zengin.EncodingShiftJIS && encodingType != "0" {
		return zengin.Header{}, errors.New("unsupported encoding type: " + encodingType)
	}
	header.EncodingType = encodingType

	senderCode, err := parseSenderCode(line[4:14])
	if err != nil {
		return zengin.Header{}, errors.New("invalid sender code: " + err.Error())
	}
	header.SenderCode = senderCode

	senderName := line[14:54] // Trim spaces for names
	header.SenderName = strings.TrimSpace(senderName)

	transactionDate := line[54:58]
	header.TransactionDate = transactionDate

	bankCode := line[58:62]
	if _, err := strconv.Atoi(bankCode); err != nil {
		return zengin.Header{}, errors.New("invalid bank code: contains non-numeric characters")
	}
	header.BankCode = bankCode

	header.BankName = strings.TrimSpace(line[62:77])

	branchCode := line[77:80]
	if _, err := strconv.Atoi(branchCode); err != nil {
		return zengin.Header{}, errors.New("invalid branch code: contains non-numeric characters")
	}
	header.BankCode = bankCode

	header.BranchName = strings.TrimSpace(line[80:95])

	header.AccountType = line[95:96]

	header.AccountNumber = line[96:103]

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
