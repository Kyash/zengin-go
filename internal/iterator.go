package internal

import (
	"errors"
	"io"

	"github.com/Kyash/zengin-go/types"
)

type TransferIterator struct {
	engine *ParsingEngine

	// Current block state
	currentIndex int
	transfers    []types.Transfer
	err          error
	done         bool
}

func NewTransferIterator(reader Reader) (*TransferIterator, error) {
	engine, err := NewParsingEngine(reader)
	if err != nil {
		return nil, err
	}

	return &TransferIterator{
		engine: engine,
	}, nil
}

func (t *TransferIterator) Next() (types.Transfer, error) {
	if t.err != nil {
		return types.Transfer{}, t.err
	}

	if t.done {
		return types.Transfer{}, io.EOF
	}

	// If we have transfers from current block, return next one
	if t.currentIndex < len(t.transfers) {
		transfer := t.transfers[t.currentIndex]
		t.currentIndex++
		return transfer, nil
	}

	// Need to parse next block
	err := t.parseNextBlock()
	if err != nil {
		t.err = err
		if err == io.EOF {
			t.done = true
		}
		return types.Transfer{}, err
	}

	// Return first transfer from newly parsed block
	if len(t.transfers) > 0 {
		transfer := t.transfers[0]
		t.currentIndex = 1
		return transfer, nil
	}

	// Should not reach here
	return types.Transfer{}, errors.New("no transfers available")
}

func (t *TransferIterator) HasMore() bool {
	if t.err != nil || t.done {
		return false
	}

	// Check if we have more transfers in current block
	if t.currentIndex < len(t.transfers) {
		return true
	}

	// If we don't have transfers buffered, we need to try parsing
	// We can't easily peek ahead without consuming the scanner,
	// so we'll be conservative and return true if no error occurred yet
	return !t.done
}

func (t *TransferIterator) Err() error {
	if t.err == io.EOF {
		return nil
	}
	return t.err
}

type iteratorHandler struct {
	transfers   []types.Transfer
	blockParsed bool
}

func (i *iteratorHandler) HandleBlock(header types.Header, data []types.Data, trailer types.Trailer) error {
	var err error
	i.transfers, err = createTransfers(header, data, trailer)
	if err != nil {
		return err
	}
	i.blockParsed = true
	return nil
}

func (i *iteratorHandler) ShouldContinue() bool {
	return !i.blockParsed
}

func (t *TransferIterator) parseNextBlock() error {
	// Reset block state
	t.transfers = nil
	t.currentIndex = 0

	handler := &iteratorHandler{}
	err := t.engine.ParseWithHandler(handler)
	if err != nil {
		if err == io.EOF {
			t.done = true
		}
		return err
	}

	if !handler.blockParsed {
		// No block was parsed, which means we've reached the end
		t.done = true
		return io.EOF
	}

	t.transfers = handler.transfers
	return nil
}
