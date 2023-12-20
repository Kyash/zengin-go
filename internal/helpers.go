package internal

import (
	"bufio"
	"errors"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"io"
	"log"
	"strconv"
	"time"
)

type Reader interface {
	io.Reader
}

func guessEncoding(file Reader) (*bufio.Scanner, Encoding, error) {
	var encoding Encoding
	reader := bufio.NewReader(file)
	peekBytes, err := reader.Peek(1024)
	if err != nil && err != io.EOF {
		log.Fatal("couldn't read from file", "error", err)
		return nil, EncodingUndefined, err
	}

	// Ignore "certain" (3rd value), as during testing it was always false, even though it correctly detects utf-8.
	_, name, _ := charset.DetermineEncoding(peekBytes, "")

	var scanner *bufio.Scanner
	switch name {
	case "utf-8":
		encoding = EncodingUTF8
		scanner = bufio.NewScanner(reader)
	// Shift-JIS can't be reliably detected, so we'll assume it's Shift-JIS if it's not UTF-8.
	default:
		encoding = EncodingShiftJIS
		scanner = bufio.NewScanner(transform.NewReader(reader, japanese.ShiftJIS.NewDecoder()))
	}

	return scanner, encoding, nil
}

func parseCategoryCode(categoryCode string) (CategoryCode, error) {
	switch categoryCode {
	case "21":
		return CategoryCodeCombination, nil
	case "11", "71":
		return CategoryCodePayment, nil
	case "12", "72":
		return CategoryCodeBonus, nil
	default:
		return CategoryCodeUndefined, errors.New("unknown category code: " + categoryCode)
	}
}

func parseSenderCode(senderCode string) (string, error) {
	if len(senderCode) != 10 {
		return "", errors.New("sender code must be 10 digits")
	}
	return senderCode, nil
}

func parseAccountType(accountType string) (AccountType, error) {
	switch accountType {
	case "1":
		return AccountTypeRegular, nil
	case "2":
		return AccountTypeChecking, nil
	case "4":
		return AccountTypeSavings, nil
	default:
		return AccountTypeUndefined, errors.New("invalid account type: " + accountType)
	}
}

func parseNewCode(accountType string) (NewCode, error) {
	switch accountType {
	case "1":
		return CodeFirstTransfer, nil
	case "2":
		return CodeUpdateTransfer, nil
	case "0":
		return CodeOther, nil
	default:
		return CodeUndefined, errors.New("invalid account type: " + accountType)
	}
}

func parseDate(date string) (string, error) {
	_, err := time.Parse("0102", date)
	if err != nil {
		return "", nil
	}
	return date, err
}
func parseBankCode(bankCode string) (string, error) {
	if _, err := strconv.Atoi(bankCode); err != nil {
		return "", errors.New("invalid bank code: contains non-numeric characters")
	}
	if len(bankCode) != 4 {
		return "", errors.New("bank code must be 4 digits")
	}
	return bankCode, nil
}

func parseBranchCode(bankCode string) (string, error) {
	if _, err := strconv.Atoi(bankCode); err != nil {
		return "", errors.New("invalid branch code: contains non-numeric characters")
	}
	if len(bankCode) != 3 {
		return "", errors.New("branch code must be 3 digits")
	}
	return bankCode, nil
}
