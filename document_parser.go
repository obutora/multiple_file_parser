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

// parseFromFileCommon は各パーサーで利用可能な共通実装
func parseFromFileCommon(p DocumentParser, filePath string) (string, error) {
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

// parseFromBytesCommon は各パーサーで利用可能な共通実装
func parseFromBytesCommon(p DocumentParser, data []byte) (string, error) {
	reader := bytes.NewReader(data)
	return p.ParseFromReader(reader, int64(len(data)))
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

	excelParser := &ExcelParser{}
	for _, ext := range excelParser.SupportedExtensions() {
		factory.parsers[ext] = excelParser
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

// ParseFromFile はファイルパスからドキュメントをパースする
// ファイルの拡張子を自動的に検出し、適切なパーサーを使用する
func (f *DocumentParserFactory) ParseFromFile(filePath string) (string, error) {
	ext := getFileExtension(filePath)
	parser, err := f.GetParser(ext)
	if err != nil {
		return "", fmt.Errorf("failed to get parser: %w", err)
	}

	content, err := parser.ParseFromFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to parse file: %w", err)
	}

	return content, nil
}

// ParseFromBytes はバイト配列からドキュメントをパースする
func (f *DocumentParserFactory) ParseFromBytes(ext string, data []byte) (string, error) {
	parser, err := f.GetParser(ext)
	if err != nil {
		return "", fmt.Errorf("failed to get parser: %w", err)
	}

	content, err := parser.ParseFromBytes(data)
	if err != nil {
		return "", fmt.Errorf("failed to parse bytes: %w", err)
	}

	return content, nil
}

// ParseFromReader はio.ReaderAtからドキュメントをパースする
func (f *DocumentParserFactory) ParseFromReader(ext string, reader io.ReaderAt, size int64) (string, error) {
	parser, err := f.GetParser(ext)
	if err != nil {
		return "", fmt.Errorf("failed to get parser: %w", err)
	}

	content, err := parser.ParseFromReader(reader, size)
	if err != nil {
		return "", fmt.Errorf("failed to parse from reader: %w", err)
	}

	return content, nil
}

// PageSeparatedParser はページやシートごとに分割してパースするインターフェース
type PageSeparatedParser interface {
	DocumentParser
	// ParseWithPages はio.ReaderAtからドキュメントをパースし、ページ/シートごとのマップを返す
	ParseWithPages(reader io.ReaderAt, size int64) (map[string]string, error)
}

// ParseFromFileWithPages はファイルパスからドキュメントをパースし、可能な場合はページ/シートごとに分割して返す
func (f *DocumentParserFactory) ParseFromFileWithPages(filePath string) (map[string]string, error) {
	ext := getFileExtension(filePath)
	parser, err := f.GetParser(ext)
	if err != nil {
		return nil, fmt.Errorf("failed to get parser: %w", err)
	}

	if p, ok := parser.(PageSeparatedParser); ok {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			return nil, fmt.Errorf("failed to get file stats: %w", err)
		}

		return p.ParseWithPages(file, stat.Size())
	}

	// PageSeparatedParserを実装していない場合は通常パースを行い、全体を一つの要素として返す
	content, err := parser.ParseFromFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	return map[string]string{"Content": content}, nil
}

// ParseFromBytesWithPages はバイト配列からドキュメントをパースし、可能な場合はページ/シートごとに分割して返す
func (f *DocumentParserFactory) ParseFromBytesWithPages(ext string, data []byte) (map[string]string, error) {
	parser, err := f.GetParser(ext)
	if err != nil {
		return nil, fmt.Errorf("failed to get parser: %w", err)
	}

	if p, ok := parser.(PageSeparatedParser); ok {
		reader := bytes.NewReader(data)
		return p.ParseWithPages(reader, int64(len(data)))
	}

	content, err := parser.ParseFromBytes(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bytes: %w", err)
	}

	return map[string]string{"Content": content}, nil
}

// ParseFromReaderWithPages はio.ReaderAtからドキュメントをパースし、可能な場合はページ/シートごとに分割して返す
func (f *DocumentParserFactory) ParseFromReaderWithPages(ext string, reader io.ReaderAt, size int64) (map[string]string, error) {
	parser, err := f.GetParser(ext)
	if err != nil {
		return nil, fmt.Errorf("failed to get parser: %w", err)
	}

	if p, ok := parser.(PageSeparatedParser); ok {
		return p.ParseWithPages(reader, size)
	}

	content, err := parser.ParseFromReader(reader, size)
	if err != nil {
		return nil, fmt.Errorf("failed to parse from reader: %w", err)
	}

	return map[string]string{"Content": content}, nil
}

// getFileExtension はファイルパスから拡張子を取得
func getFileExtension(filePath string) string {
	for i := len(filePath) - 1; i >= 0; i-- {
		if filePath[i] == '.' {
			return filePath[i:]
		}
	}
	return ""
}
