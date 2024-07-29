package internal

import (
	"bufio"
	"errors"
	"github.com/Kyash/zengin-go/types"
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

func guessEncoding(file Reader) (*bufio.Scanner, types.Encoding, error) {
	var encoding types.Encoding
	reader := bufio.NewReader(file)
	peekBytes, err := reader.Peek(1024)
	if err != nil && err != io.EOF {
		log.Fatal("couldn't read from file", "error", err)
		return nil, types.EncodingUndefined, err
	}

	// Ignore "certain" (3rd value), as during testing it was always false, even though it correctly detects utf-8.
	_, name, _ := charset.DetermineEncoding(peekBytes, "")

	var scanner *bufio.Scanner
	switch name {
	case "utf-8":
		encoding = types.EncodingUTF8
		scanner = bufio.NewScanner(reader)
	// Shift-JIS can't be reliably detected, so we'll assume it's Shift-JIS if it's not UTF-8.
	default:
		encoding = types.EncodingShiftJIS
		scanner = bufio.NewScanner(transform.NewReader(reader, japanese.ShiftJIS.NewDecoder()))
	}

	return scanner, encoding, nil
}

func parseCategoryCode(categoryCode string) (types.CategoryCode, error) {
	switch categoryCode {
	case "21":
		return types.CategoryCodeCombination, nil
	case "11", "71":
		return types.CategoryCodePayment, nil
	case "12", "72":
		return types.CategoryCodeBonus, nil
	default:
		return types.CategoryCodeUndefined, errors.New("unknown category code: " + categoryCode)
	}
}

func parseSenderCode(senderCode string) (string, error) {
	if len(senderCode) != 10 {
		return "", errors.New("sender code must be 10 digits")
	}
	return senderCode, nil
}

func parseAccountType(accountType string) (types.AccountType, error) {
	switch accountType {
	case "1":
		return types.AccountTypeRegular, nil
	case "2":
		return types.AccountTypeChecking, nil
	case "4":
		return types.AccountTypeSavings, nil
	default:
		return types.AccountTypeUndefined, errors.New("invalid account type: " + accountType)
	}
}

func parseNewCode(accountType string) (types.NewCode, error) {
	switch accountType {
	case "1":
		return types.CodeFirstTransfer, nil
	case "2":
		return types.CodeUpdateTransfer, nil
	case "0":
		return types.CodeOther, nil
	default:
		return types.CodeUndefined, errors.New("invalid account type: " + accountType)
	}
}

func parseDate(date string) (string, error) {
	_, err := time.Parse("0102", date)
	if err != nil {
		return "", err
	}
	return date, nil
}
func parseBankCode(bankCode string) (string, error) {
	if len(bankCode) != 4 {
		return "", errors.New("bank code must be 4 digits")
	}
	if _, err := strconv.Atoi(bankCode); err != nil {
		return "", errors.New("invalid bank code: contains non-numeric characters: " +
			bankCode + "Error: " + err.Error())
	}
	return bankCode, nil
}

func parseBranchCode(branchCode string) (string, error) {
	if len(branchCode) != 3 {
		return "", errors.New("branch code must be 3 digits")
	}
	if _, err := strconv.Atoi(branchCode); err != nil {
		return "", errors.New("invalid branch code: contains non-numeric characters: " + branchCode)
	}
	return branchCode, nil
}
