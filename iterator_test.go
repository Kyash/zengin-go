package zengin

import (
	"io"
	"strings"
	"testing"
)

func TestTransferIterator_SingleBlock(t *testing.T) {
	input := `12110110999999ｹﾝｼﾝ ﾀﾛｳ                                02242606               010               20999999
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    020ﾋﾖｳｺﾞ              19876543ｹﾝｼﾝ ｼﾖｳｼﾞ                    00000000010                    0
8000001000000000001
9`

	reader := strings.NewReader(input)
	iterator, err := NewTransferIterator(reader)
	if err != nil {
		t.Fatalf("Failed to create iterator: %v", err)
	}

	// Should have one transfer
	transfer, err := iterator.Next()
	if err != nil {
		t.Fatalf("Failed to get first transfer: %v", err)
	}

	if transfer.SenderName != "ｹﾝｼﾝ ﾀﾛｳ" {
		t.Errorf("Expected sender name 'ｹﾝｼﾝ ﾀﾛｳ', got '%s'", transfer.SenderName)
	}
	if transfer.RecipientName != "ｹﾝｼﾝ ｼﾖｳｼﾞ" {
		t.Errorf("Expected recipient name 'ｹﾝｼﾝ ｼﾖｳｼﾞ', got '%s'", transfer.RecipientName)
	}
	if transfer.Amount != 1 {
		t.Errorf("Expected amount 1, got %d", transfer.Amount)
	}

	// Should be no more transfers
	_, err = iterator.Next()
	if err != io.EOF {
		t.Errorf("Expected EOF, got %v", err)
	}
}

func TestTransferIterator_MultipleBlocks(t *testing.T) {
	input := `12110110999999ｹﾝｼﾝ ﾀﾛｳ                                02242606               010               20999999
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    020ﾋﾖｳｺﾞ              19876543ｹﾝｼﾝ ｼﾖｳｼﾞ                    00000000010                    0
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    030ｻﾝﾉﾐﾔ              29999999ｹﾝｼﾝ ﾊﾅｺ                      00000000020                    0
8000002000000000003
12110110999999ｹﾝｼﾝ ﾀﾛ                                 01142606               010               20999999
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    030ｻﾝﾉﾐﾔ              11234567ｹﾝｼﾝ ｼﾞﾛｳ                     00000000030ﾏｲﾂｷﾌﾞﾝ             0Y
8000001000000000003
9`

	reader := strings.NewReader(input)
	iterator, err := NewTransferIterator(reader)
	if err != nil {
		t.Fatalf("Failed to create iterator: %v", err)
	}

	expectedTransfers := []struct {
		senderName    string
		recipientName string
		amount        uint64
	}{
		{"ｹﾝｼﾝ ﾀﾛｳ", "ｹﾝｼﾝ ｼﾖｳｼﾞ", 1},
		{"ｹﾝｼﾝ ﾀﾛｳ", "ｹﾝｼﾝ ﾊﾅｺ", 2},
		{"ｹﾝｼﾝ ﾀﾛ", "ｹﾝｼﾝ ｼﾞﾛｳ", 3},
	}

	for i, expected := range expectedTransfers {
		transfer, err := iterator.Next()
		if err != nil {
			t.Fatalf("Failed to get transfer %d: %v", i+1, err)
		}

		if transfer.SenderName != expected.senderName {
			t.Errorf("Transfer %d: Expected sender name '%s', got '%s'", i+1, expected.senderName, transfer.SenderName)
		}
		if transfer.RecipientName != expected.recipientName {
			t.Errorf("Transfer %d: Expected recipient name '%s', got '%s'", i+1, expected.recipientName, transfer.RecipientName)
		}
		if transfer.Amount != expected.amount {
			t.Errorf("Transfer %d: Expected amount %d, got %d", i+1, expected.amount, transfer.Amount)
		}
	}

	// Should be no more transfers
	_, err = iterator.Next()
	if err != io.EOF {
		t.Errorf("Expected EOF, got %v", err)
	}
}

func TestTransferIterator_EmptyFile(t *testing.T) {
	input := ""
	reader := strings.NewReader(input)
	iterator, err := NewTransferIterator(reader)
	if err != nil {
		t.Fatalf("Failed to create iterator: %v", err)
	}

	_, err = iterator.Next()
	if err == nil {
		t.Error("Expected error for empty file, got nil")
	}
}

func TestTransferIterator_InvalidFormat(t *testing.T) {
	input := "invalid file content"
	reader := strings.NewReader(input)
	iterator, err := NewTransferIterator(reader)
	if err != nil {
		t.Fatalf("Failed to create iterator: %v", err)
	}

	_, err = iterator.Next()
	if err == nil {
		t.Error("Expected error for invalid format, got nil")
	}
}

func TestTransferIterator_HasMore(t *testing.T) {
	input := `12110110999999ｹﾝｼﾝ ﾀﾛｳ                                02242606               010               20999999
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    020ﾋﾖｳｺﾞ              19876543ｹﾝｼﾝ ｼﾖｳｼﾞ                    00000000010                    0
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    030ｻﾝﾉﾐﾔ              29999999ｹﾝｼﾝ ﾊﾅｺ                      00000000020                    0
8000002000000000003
9`

	reader := strings.NewReader(input)
	iterator, err := NewTransferIterator(reader)
	if err != nil {
		t.Fatalf("Failed to create iterator: %v", err)
	}

	// Get first transfer
	_, err = iterator.Next()
	if err != nil {
		t.Fatalf("Failed to get first transfer: %v", err)
	}

	// Should still have more in the same block
	if !iterator.HasMore() {
		t.Error("Expected HasMore() to return true after first transfer")
	}

	// Get second transfer
	_, err = iterator.Next()
	if err != nil {
		t.Fatalf("Failed to get second transfer: %v", err)
	}

	// Now try to get third transfer (should hit EOF)
	_, err = iterator.Next()
	if err != io.EOF {
		t.Errorf("Expected EOF after all transfers, got %v", err)
	}

	// After EOF, HasMore should return false
	if iterator.HasMore() {
		t.Error("Expected HasMore() to return false after EOF")
	}
}
