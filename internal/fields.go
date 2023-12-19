package internal

import (
	"strings"
	"time"
)

type Encoding int

const (
	HeaderLength  = 120 // 103?
	DataLength    = 120 // 106?
	TrailerLength = 120 // 19?
	EndLength     = 120 // 1?
)

const (
	EncodingUndefined Encoding = iota
	EncodingShiftJIS  Encoding = iota
	EncodingUTF8
)

type Transfer struct {
	SenderName    string
	TransferDate  time.Time
	BankCode      string
	BranchCode    string
	AccountType   string
	AccountNumber string
	AccountName   string
	Amount        uint64
}

type CategoryCode int

const (
	CategoryCodeUndefined   CategoryCode = iota
	CategoryCodeCombination CategoryCode = iota
	CategoryCodePayment     CategoryCode = iota
	CategoryCodeBonus       CategoryCode = iota
)

type AccountType int

const (
	AccountTypeUndefined AccountType = iota
	AccountTypeRegular   AccountType = iota
	AccountTypeChecking  AccountType = iota
	AccountTypeSavings   AccountType = iota
)

type Header struct {
	RecordType      string       // 1 digit
	CategoryCode    CategoryCode // 2 digits
	EncodingType    string       // 1 digit
	SenderCode      string       // 10 digits
	SenderName      string       // 40 characters
	TransactionDate string       // 4 digits (MMDD)
	BankCode        string       // 4 digits
	BankName        string       // 15 characters
	BranchCode      string       // 3 digits
	BranchName      string       // 15 characters
	AccountType     AccountType  // 1 digit
	AccountNumber   string       // 7 digits
	Dummy           string       // 17 characters (unused)
}

type Data struct {
	RecordType          string      // 1 digit
	RecipientBankCode   string      // 4 digits
	RecipientBankName   string      // 15 characters
	RecipientBranchCode string      // 3 digits
	RecipientBranchName string      // 15 characters
	ExchangeOfficeCode  string      // 4 digits (unused)
	AccountType         AccountType // 1 digit
	AccountNumber       string      // 7 digits
	RecipientName       string      // 30 characters
	TransferAmount      uint64      // 10 digits
	NewCode             string      // 1 digit (unused)
	CustomerCode1       string      // 10 characters
	CustomerCode2       string      // 10 characters
	EDIInformation      string      // 20 characters
	TransferCategory    string      // 1 digit (unused)
	Identification      string      // 1 character
	Dummy               string      // 7 characters (unused)
}

type Trailer struct {
	RecordType  string // 1 digit
	TotalCount  int    // 6 digits
	TotalAmount uint64 // 12 digits
	Dummy       string // 101 characters (unused)
}

// helper functions

func IsHeader(line string) bool {
	if !strings.HasPrefix(line, "1") || len(line) != HeaderLength {
		return false
	}
	return true
}
func IsData(line string) bool {
	if len(line) != DataLength || line[0:1] != "2" {
		return false
	}
	return true
}
func IsTrailer(line string) bool {
	if len(line) < TrailerLength || line[0:1] != "8" {
		return false
	}
	return true
}
func IsEndRecord(line string) bool {
	if len(line) != EndLength || line[0:1] != "9" {
		return false
	}
	return true
}
