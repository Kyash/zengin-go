package internal

import (
	"strings"
	"time"
)

type Encoding int

const (
	MinHeaderLength  = 80 // until "仕向支店番号"
	MinDataLength    = 91 // until "新規コード"
	MinTrailerLength = 19 // until "dummy"
	MinEndLength     = 1  // until "dummy"
)

const (
	EncodingUndefined Encoding = iota
	EncodingShiftJIS
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
	CategoryCodeUndefined CategoryCode = iota
	CategoryCodeCombination
	CategoryCodePayment
	CategoryCodeBonus
)

type AccountType int

const (
	AccountTypeUndefined AccountType = iota
	AccountTypeRegular
	AccountTypeChecking
	AccountTypeSavings
)

type NewCode int // 新規コード

const (
	CodeFirstTransfer  NewCode = 1 // 第 1 回振込分
	CodeUpdateTransfer         = 2 // 変更分(被仕向銀行・支店、預金種目・口座番号)    //
	CodeOther                  = 0
	CodeUndefined              = -1
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
	NewCode             NewCode     // 1 digit (unused)
	// Next 20 characters can be used for CustomerCode1&2, or EDIInformation
	Extra            string // 20 characters
	TransferCategory string // 1 digit (unused)
	EdiPresent       bool   // 1 character, if "Y", EDIInformation is used
	Dummy            string // 7 characters (unused)
}

type Trailer struct {
	RecordType  string // 1 digit
	TotalCount  int    // 6 digits
	TotalAmount uint64 // 12 digits
	Dummy       string // 101 characters (unused)
}

// helper functions

func IsHeader(line []rune) bool {
	if len(line) < MinHeaderLength || !strings.HasPrefix(string(line), "1") {
		return false
	}
	return true
}
func IsData(line []rune) bool {
	if len(line) < MinDataLength || string(line[0:1]) != "2" {
		return false
	}
	return true
}
func IsTrailer(line []rune) bool {
	if len(line) < MinTrailerLength || string(line[0:1]) != "8" {
		return false
	}
	return true
}
func IsEndRecord(line []rune) bool {
	if len(line) < MinEndLength || string(line[0:1]) != "9" {
		return false
	}
	return true
}
