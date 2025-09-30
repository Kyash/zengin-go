# zengin-go

**zengin-go** は、日本の銀行が金融取引を行うために一般的に使用する全銀フォーマットのテキストファイルを解析するためのGoライブラリです。

## 特徴

- 全銀フォーマットのテキストファイルを解析し、CSV形式のデータまたはすべてのフィールドを含むGo構造体として取得できます。
- UTF-8およびShift-JISの両方のエンコーディングをサポートし、他のエンコーディングもサポートする可能性があります（未テスト）。

## インターフェース

このライブラリは、バッチ処理とストリーミング処理の両方のインターフェースを提供します：

### バッチAPI
すべてのデータをメモリに読み込むことが可能な小さなファイル用：
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

### ストリーミングAPI
大きなファイルのメモリ効率的な処理用：
```go
// 振込を一つずつ処理するためのイテレータを作成
func NewTransferIterator(reader zengin.Reader) *TransferIterator

// イテレータのメソッド:
func (it *TransferIterator) Next() *types.Transfer  // 次の振込を取得
func (it *TransferIterator) HasMore() bool          // さらに振込があるかチェック
func (it *TransferIterator) Err() error             // 解析エラーを取得
```

ストリーミングAPIは、ファイル全体をメモリに読み込むことなく、一度に一つのヘッダ→データ[]→トレーラブロックを処理するため、大きなファイルに適しています。

解析可能なフィールドは [types/fields.go](./types/fields.go) にあります。


## インストール

このライブラリをインストールするには、`go get` コマンドを使用します：

```bash
go get github.com/Kyash/zengin-go
```

## 使用方法

### バッチ処理
[バッチサンプル](./samples/main.go)を参照してください。

### ストリーミング処理
大きなファイルのメモリ効率的な処理については[ストリーミングサンプル](./samples/streaming/main.go)を参照してください。

## コントリビュート

問題や機能リクエストがある場合は、イシューを作成するかプルリクエストを作ってください。

## ライセンス

このプロジェクトはMITライセンスの下でライセンスされています。詳細については [LICENSE](./LICENSE) ファイルを参照してください。
