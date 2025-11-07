# Multiple File Parser

複数のドキュメント形式（PDF、DOCX、PPTX、テキストファイルなど）からテキストを抽出するGoライブラリです。

## 特徴

- 🚀 **複数フォーマット対応**: PDF、DOCX、PPTX、およびテキストファイルをサポート
- 📦 **柔軟なパース方法**: ファイルパス、バイト配列、`io.ReaderAt`の3つの方法でパース可能
- 🔌 **拡張可能**: ファクトリーパターンを採用し、カスタムパーサーの追加が容易
- 🎯 **シンプルなAPI**: 統一されたインターフェースで簡単に使用可能

## サポートするフォーマット

| 形式       | 拡張子                               | 説明                                           |
| ---------- | ------------------------------------ | ---------------------------------------------- |
| PDF        | `.pdf`                               | PDFドキュメント                                |
| Word       | `.docx`                              | Microsoft Word文書                             |
| PowerPoint | `.pptx`, `.ppt`                      | Microsoft PowerPointプレゼンテーション         |
| テキスト   | `.txt`, `.md`, `.json`, `.xml`, など | プレーンテキストおよび各種ソースコードファイル |

## インストール

```bash
go get github.com/obutora/multiple_file_parser
```

## 使い方

### 基本的な使用例

```go
package main

import (
    "fmt"
    "log"
    
    service "github.com/obutora/multiple_file_parser"
)

func main() {
    // ファクトリーを作成
    factory := service.NewDocumentParserFactory()
    
    // 拡張子に基づいてパーサーを取得
    parser, err := factory.GetParser(".pdf")
    if err != nil {
        log.Fatal(err)
    }
    
    // ファイルをパース
    content, err := parser.ParseFromFile("document.pdf")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(content)
}
```

### 3つのパース方法

#### 1. ファイルパスからパース

```go
parser, _ := factory.GetParser(".pdf")
content, err := parser.ParseFromFile("document.pdf")
```

#### 2. バイト配列からパース

```go
data := []byte("テキストの内容")
parser, _ := factory.GetParser(".txt")
content, err := parser.ParseFromBytes(data)
```

#### 3. io.ReaderAtからパース

```go
file, _ := os.Open("document.pdf")
defer file.Close()

stat, _ := file.Stat()
parser, _ := factory.GetParser(".pdf")
content, err := parser.ParseFromReader(file, stat.Size())
```

### サポートされている拡張子の確認

```go
factory := service.NewDocumentParserFactory()
extensions := factory.SupportedExtensions()

for _, ext := range extensions {
    fmt.Println(ext)
}
```

### カスタムパーサーの追加

独自のパーサーを作成して登録することができます：

```go
// カスタムパーサーを実装
type CustomParser struct {
    service.BaseParser
}

func (p *CustomParser) SupportedExtensions() []string {
    return []string{".custom"}
}

func (p *CustomParser) ParseFromReader(reader io.ReaderAt, size int64) (string, error) {
    // パース処理を実装
    return "parsed content", nil
}

// ファクトリーに登録
factory := service.NewDocumentParserFactory()
factory.RegisterParser(&CustomParser{})
```

## API リファレンス

### DocumentParser インターフェース

すべてのパーサーは以下のインターフェースを実装します：

```go
type DocumentParser interface {
    ParseFromReader(reader io.ReaderAt, size int64) (string, error)
    ParseFromBytes(data []byte) (string, error)
    ParseFromFile(filePath string) (string, error)
    SupportedExtensions() []string
}
```

### DocumentParserFactory

パーサーの管理と取得を行うファクトリークラス：

- `NewDocumentParserFactory()`: 新しいファクトリーインスタンスを作成
- `GetParser(extension string)`: 拡張子に対応するパーサーを取得
- `RegisterParser(parser DocumentParser)`: カスタムパーサーを登録
- `SupportedExtensions()`: サポートされている全拡張子を取得

## サンプルコード

詳細なサンプルコードは [`examples/example.go`](examples/example.go) を参照してください。

実行方法：

```bash
cd examples
go run example.go
```

## テキストファイルでサポートされる拡張子

TextParserは以下の拡張子をサポートしています：

- **ドキュメント**: `.txt`, `.md`, `.markdown`
- **プログラミング言語**: `.go`, `.py`, `.js`, `.ts`, `.java`, `.c`, `.cpp`, など
- **設定ファイル**: `.json`, `.xml`, `.yaml`, `.yml`, `.toml`, `.ini`
- **スクリプト**: `.sh`, `.bash`, `.zsh`, `.ps1`
- **Web**: `.html`, `.css`, `.scss`, `.vue`, `.svelte`

完全なリストは [`text.go`](text.go) を参照してください。

## 依存関係

- `github.com/ledongthuc/pdf`: PDFパース用

## ライセンス

このプロジェクトは[MIT License](LICENSE)の下で公開されています。

## 貢献

プルリクエストを歓迎します。大きな変更の場合は、まずissueを開いて変更内容を議論してください。

## 開発者

GitHub: [@obutora](https://github.com/obutora)
