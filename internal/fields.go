package internal

import (
	"strings"
	"time"
)

type Encoding int

const (
	EncodingShiftJIS Encoding = iota
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

type Record struct {
	Header    Header
	Data      []Data
	Trailer   Trailer
	EndRecord EndRecord
}

type CategoryCode int

const (
	CategoryCodeCombination CategoryCode = 1
	CategoryCodePayment     CategoryCode = 2
	CategoryCodeBonus       CategoryCode = 3
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
	AccountType     string       // 1 digit
	AccountNumber   string       // 7 digits
	Dummy           string       // 17 characters (unused)
}

type Data struct {
}

type Trailer struct {
}

type EndRecord struct {
}

func IsHeader(line string) bool {
	if !strings.HasPrefix(line, "1") {
		return false
	}
	if len(line) != 120 {
		return false
	}

	return true
}
func IsData(line string) bool {
	return false
}
func IsTrailer(line string) bool {
	return false
}
func IsEndRecord(line string) bool {
	return false
}
