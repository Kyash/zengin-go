# zengin-go

**zengin-go** は、日本の銀行が金融取引を行うために一般的に使用する全銀フォーマットのテキストファイルを解析するためのGoライブラリです。

## 特徴

- 全銀フォーマットのテキストファイルを解析し、CSV形式のデータまたはすべてのフィールドを含むGo構造体として取得できます。
- UTF-8およびShift-JISの両方のエンコーディングをサポートし、他のエンコーディングもサポートする可能性があります（未テスト）。

## インターフェース

このライブラリは、全銀フォーマットのテキストファイルをCSV形式のデータまたはGo構造体に変換するために使用できます。

```go
// 全銀フォーマットファイルを解析し、すべてのフィールドを含む行を返します
func Parse(reader zengin.Reader) ([]types.Transfer, error)

// 全銀フォーマットファイルを解析し、以下のフィールド名を持つCSV形式のテーブルを返します
// SenderName,TransferDate,BankCode,BranchCode,AccountType,AccountNumber,AccountName,Amount
func ToCSV(reader zengin.Reader) ([][]string, error)

// 全銀フォーマットファイルを解析し、以下のフィールド名を持つCSV形式のテーブルを返します
// 振込名義人, 振込日, 金融機関コード, 支店コード, 科目, 口座番号, 口座名義人, 金額
func ToCSVJa(reader zengin.Reader) ([][]string, error)
```

解析可能なフィールドは [types/fields.go](./types/fields.go) にあります。


## インストール

このライブラリをインストールするには、`go get` コマンドを使用します：

```bash
go get github.com/Kyash/zengin-go
```

## 使用方法

[サンプル](./samples/main.go)を参照にしてください。

## コントリビュート

問題や機能リクエストがある場合は、イシューを作成するかプルリクエストを作ってください。

## ライセンス

このプロジェクトはMITライセンスの下でライセンスされています。詳細については [LICENSE](./LICENSE) ファイルを参照してください。
