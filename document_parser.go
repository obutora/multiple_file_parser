package documentParser

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

// DocumentParser はドキュメントをパースするインターフェース
type DocumentParser interface {
	// ParseFromReader はio.ReaderAtからドキュメントをパース
	ParseFromReader(reader io.ReaderAt, size int64) (string, error)

	// ParseFromBytes はバイト配列からドキュメントをパース
	ParseFromBytes(data []byte) (string, error)

	// ParseFromFile はファイルパスからドキュメントをパース（後方互換性）
	ParseFromFile(filePath string) (string, error)

	// SupportedExtensions はサポートする拡張子を返す
	SupportedExtensions() []string
}

// BaseParser は共通処理を提供する基底構造体
type BaseParser struct{}

// ParseFromBytes のデフォルト実装（ReaderAtを使う実装にフォールバック）
func (p *BaseParser) ParseFromBytes(data []byte) (string, error) {
	reader := bytes.NewReader(data)
	return p.ParseFromReader(reader, int64(len(data)))
}

// ParseFromFile のデフォルト実装
func (p *BaseParser) ParseFromFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to get file stats: %w", err)
	}

	return p.ParseFromReader(file, stat.Size())
}

// ParseFromReader は各パーサーで実装が必要
func (p *BaseParser) ParseFromReader(reader io.ReaderAt, size int64) (string, error) {
	return "", fmt.Errorf("ParseFromReader not implemented")
}

// DocumentParserFactory はファイル拡張子に基づいてパーサーを返す
type DocumentParserFactory struct {
	parsers map[string]DocumentParser
}

// NewDocumentParserFactory はファクトリーを初期化
func NewDocumentParserFactory() *DocumentParserFactory {
	factory := &DocumentParserFactory{
		parsers: make(map[string]DocumentParser),
	}

	// パーサーを登録
	pptxParser := &PPTXParser{}
	for _, ext := range pptxParser.SupportedExtensions() {
		factory.parsers[ext] = pptxParser
	}

	pdfParser := &PDFParser{}
	for _, ext := range pdfParser.SupportedExtensions() {
		factory.parsers[ext] = pdfParser
	}

	docxParser := &DOCXParser{}
	for _, ext := range docxParser.SupportedExtensions() {
		factory.parsers[ext] = docxParser
	}

	textParser := &TextParser{}
	for _, ext := range textParser.SupportedExtensions() {
		factory.parsers[ext] = textParser
	}

	return factory
}

// GetParser は拡張子に対応するパーサーを返す
func (f *DocumentParserFactory) GetParser(extension string) (DocumentParser, error) {
	ext := strings.ToLower(extension)
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	parser, ok := f.parsers[ext]
	if !ok {
		return nil, fmt.Errorf("unsupported file extension: %s", extension)
	}

	return parser, nil
}

// RegisterParser はカスタムパーサーを登録
func (f *DocumentParserFactory) RegisterParser(parser DocumentParser) {
	for _, ext := range parser.SupportedExtensions() {
		f.parsers[ext] = parser
	}
}

// SupportedExtensions はファクトリでサポートされる全ての拡張子を返す
func (f *DocumentParserFactory) SupportedExtensions() []string {
	var extensions []string
	for ext := range f.parsers {
		extensions = append(extensions, ext)
	}
	return extensions
}
