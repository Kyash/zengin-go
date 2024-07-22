package types

import (
	"strings"
)

type Header struct {
	RecordType          string       // 1 digit
	CategoryCode        CategoryCode // 2 digits
	EncodingType        string       // 1 digit
	SenderCode          string       // 10 digits
	SenderName          string       // 40 characters
	TransferDate        string       // 4 digits (MMDD)
	SenderBankCode      string       // 4 digits
	SenderBankName      string       // 15 characters
	SenderBranchCode    string       // 3 digits
	SenderBranchName    string       // 15 characters
	SenderAccountType   AccountType  // 1 digit
	SenderAccountNumber string       // 7 digits
	Dummy               string       // 17 characters (unused)
}

type Data struct {
	RecordType             string      // 1 digit
	RecipientBankCode      string      // 4 digits
	RecipientBankName      string      // 15 characters
	RecipientBranchCode    string      // 3 digits
	RecipientBranchName    string      // 15 characters
	ExchangeOfficeCode     string      // 4 digits (unused)
	RecipientAccountType   AccountType // 1 digit
	RecipientAccountNumber string      // 7 digits
	RecipientName          string      // 30 characters
	Amount                 uint64      // 10 digits
	NewCode                NewCode     // 1 digit (unused)
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
