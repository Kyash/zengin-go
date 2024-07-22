# zengin-go

[Japanese](./README.ja.md)

**zengin-go** is a Go library designed to parse Zengin format text files, 
which are commonly used by Japanese banks to conduct financial transactions.

## Features

- Parses Zengin format text files (全銀フォーマット) and get CSV-like data or all the fields as a go struct.
- Supports both UTF-8 and Shift-JIS encodings and possibly other encodings (not tested).

## Interface

The library can be used to convert Zengin format text files into CSV-like data or Go structs.
```go
// Parse Zengin format file and return rows with all fields
func Parse(reader zengin.Reader) ([]types.Transfer, error)

// Parse Zengin format file and return a csv like table with field names as below:
// SenderName,TransferDate,BankCode,BranchCode,AccountType,AccountNumber,AccountName,Amount
func ToCSV(reader zengin.Reader) ([][]string, error)

// Parse Zengin format file and return a csv like table with field names as below:
// 振込名義人,振込日,金融機関コード,支店コード,科目,口座番号,口座名義人,金額
func ToCSVJa(reader zengin.Reader) ([][]string, error) {
```

Parsable fields can be found in [types/fields.go](./types/fields.go).


## Installation

To install the library, use the `go get` command: 

```bash
go get github.com/Kyash/zengin-go
```

## Usage

See [sample](./samples/main.go)

## Contributing

If there are any issues or feature requests, please create an issue or a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.