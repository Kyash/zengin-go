package internal

import (
	"testing"

	"github.com/Kyash/zengin-go/types"
)

func TestParseHeaderSenderBankAndBranchCode(t *testing.T) {
	// 仕向銀行番号 = 2606, 仕向支店番号 = 010
	line := []rune("12110110999999ｹﾝｼﾝ ﾀﾛｳ                                02242606               010               20999999")

	header, err := parseHeader(line, types.EncodingShiftJIS)
	if err != nil {
		t.Fatalf("parseHeader returned error: %v", err)
	}

	if header.SenderBankCode != "2606" {
		t.Errorf("SenderBankCode = %q, want %q", header.SenderBankCode, "2606")
	}
	if header.SenderBranchCode != "010" {
		t.Errorf("SenderBranchCode = %q, want %q", header.SenderBranchCode, "010")
	}
}
