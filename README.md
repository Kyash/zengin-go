# zengin-go

[Japanese](./README.ja.md)

**zengin-go** is a Go library designed to parse Zengin format text files, 
which are commonly used by Japanese banks to conduct financial transactions.

## Features

- Parses Zengin format text files (全銀フォーマット) and get CSV-like data or all the fields as a go struct.
- Supports both UTF-8 and Shift-JIS encodings and possibly other encodings (not tested).

## Interface

The library provides both batch and streaming interfaces:

### Batch APIs
For smaller files where loading all data into memory is acceptable:
```go
// Parse Zengin format file and return rows with all fields
func Parse(reader zengin.Reader) ([]types.Transfer, error)

// Parse Zengin format file and return a csv like table with field names as below:
// SenderName,TransferDate,BankCode,BranchCode,AccountType,AccountNumber,AccountName,Amount
func ToCSV(reader zengin.Reader) ([][]string, error)

// Parse Zengin format file and return a csv like table with field names as below:
// 振込名義人,振込日,金融機関コード,支店コード,科目,口座番号,口座名義人,金額
func ToCSVJa(reader zengin.Reader) ([][]string, error)
```

### Streaming API
For memory-efficient processing of large files:
```go
// Create an iterator for processing transfers one by one
func NewTransferIterator(reader zengin.Reader) *TransferIterator

// Iterator methods:
func (it *TransferIterator) Next() *types.Transfer  // Get next transfer
func (it *TransferIterator) HasMore() bool          // Check if more transfers available
func (it *TransferIterator) Err() error             // Get any parsing errors
```

The streaming API processes one Header→Data[]→Trailer block at a time without loading the entire file into memory, making it suitable for large files.

Parsable fields can be found in [types/fields.go](./types/fields.go).


## Installation

To install the library, use the `go get` command: 

```bash
go get github.com/Kyash/zengin-go
```

## Usage

### Batch Processing
See [batch sample](./samples/main.go)

### Streaming Processing
See [streaming sample](./samples/streaming/main.go) for memory-efficient processing of large files

## Contributing

If there are any issues or feature requests, please create an issue or a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.
