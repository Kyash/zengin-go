package internal

import (
	"bufio"
	"errors"
	"fmt"
	"io"

	"github.com/Kyash/zengin-go/types"
)

type BlockHandler interface {
	HandleBlock(header types.Header, data []types.Data, trailer types.Trailer) error
	ShouldContinue() bool
}

type ParsingEngine struct {
	scanner  *bufio.Scanner
	encoding types.Encoding
	state    ParseState
}

func NewParsingEngine(reader Reader) (*ParsingEngine, error) {
	scanner, encoding, err := guessEncoding(reader)
	if err != nil {
		return nil, err
	}

	return &ParsingEngine{
		scanner:  scanner,
		encoding: encoding,
		state:    StateUnknown,
	}, nil
}

func (p *ParsingEngine) ParseWithHandler(handler BlockHandler) error {
	var header types.Header
	var data []types.Data
	var trailer types.Trailer

	for p.scanner.Scan() {
		line := preprocessLine(p.scanner.Text())
		if len(line) == 0 {
			continue
		}

		recordProcessed, err := p.processRecord(line, &header, &data, &trailer, handler)
		if err != nil {
			return err
		}

		if recordProcessed && !handler.ShouldContinue() {
			return nil
		}

		if p.state == StateEnd {
			break
		}
	}

	if err := p.scanner.Err(); err != nil {
		return err
	}

	if p.state != StateEnd {
		return errors.New("unexpected end of file")
	}

	return io.EOF
}

func preprocessLine(text string) []rune {
	line := []rune(text)
	if len(line) == 0 {
		return line
	}

	// Remove BOM if exists
	if len(line) >= 1 && line[0] == '\ufeff' {
		line = line[1:]
	}

	return line
}

func (p *ParsingEngine) processRecord(line []rune, header *types.Header, data *[]types.Data, trailer *types.Trailer, handler BlockHandler) (bool, error) {
	switch {
	case types.IsHeader(line):
		if p.state == StateData || p.state == StateEnd {
			return false, errors.New("found record with missing trailer")
		}
		var err error
		*header, err = parseHeader(line, p.encoding)
		if err != nil {
			return false, fmt.Errorf("error parsing header: %w", err)
		}
		p.state = StateHeader

	case types.IsData(line):
		if p.state != StateHeader && p.state != StateData {
			return false, errors.New("data record found before header")
		}
		dataRecord, err := parseData(line)
		if err != nil {
			return false, fmt.Errorf("error parsing data record: %w", err)
		}
		*data = append(*data, dataRecord)
		p.state = StateData

	case types.IsTrailer(line):
		if p.state != StateData && p.state != StateHeader {
			return false, errors.New("trailer record found before header")
		}
		var err error
		*trailer, err = parseTrailer(line)
		if err != nil {
			return false, fmt.Errorf("error parsing trailer record: %w", err)
		}

		// Handle the completed block
		err = handler.HandleBlock(*header, *data, *trailer)
		if err != nil {
			return false, err
		}

		// Reset data for next block
		*data = (*data)[:0]
		p.state = StateTrailer
		return true, nil // Block processed

	case types.IsEndRecord(line):
		if p.state != StateTrailer {
			return false, errors.New("end record found before trailer")
		}
		p.state = StateEnd
		return false, io.EOF // Signal end of file

	default:
		// Some programs seem to put invisible characters, just ignore them
	}

	return false, nil
}
