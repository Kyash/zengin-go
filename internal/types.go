package internal

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
	Header
	Data
	Trailer
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
	AccountTypeUndefined AccountType = 0
	AccountTypeRegular   AccountType = 1
	AccountTypeChecking  AccountType = 2
	AccountTypeSavings   AccountType = 4
)

type NewCode int // 新規コード

const (
	CodeFirstTransfer  NewCode = 1 // 第 1 回振込分
	CodeUpdateTransfer         = 2 // 変更分(被仕向銀行・支店、預金種目・口座番号)    //
	CodeOther                  = 0
	CodeUndefined              = -1
)
