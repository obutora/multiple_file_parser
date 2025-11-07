package main

import (
	"fmt"
	"log"
	"os"

	service "github.com/obutora/multiple_file_parser"
)

func main() {
	// ドキュメントパーサーファクトリーを作成
	factory := service.NewDocumentParserFactory()

	// サポートされている拡張子を表示
	fmt.Println("サポートされている拡張子:")
	for _, ext := range factory.SupportedExtensions() {
		fmt.Printf("  - %s\n", ext)
	}
	fmt.Println()

	// 例1: ファイルパスからパース（PDF）
	pdfFilePath := "assets/sample.pdf"
	if err := parseFromFile(factory, pdfFilePath); err != nil {
		log.Printf("PDFファイルのパースに失敗: %v\n", err)
	}

	// 例2: バイト配列からパース（テキスト）
	exampleBytes := []byte("これはテキストファイルの内容です。\nサンプルテキスト。")
	if err := parseFromBytes(factory, ".txt", exampleBytes); err != nil {
		log.Printf("バイト配列のパースに失敗: %v\n", err)
	}

	// 例3: io.ReaderAtからパース（DOCX）
	docxFilePath := "assets/sample.docx"
	if err := parseFromReader(factory, docxFilePath); err != nil {
		log.Printf("DOCXファイルのパースに失敗: %v\n", err)
	}

	// 例4: PPTXファイルのパース
	pptxFilePath := "assets/sample.pptx"
	if err := parseFromFile(factory, pptxFilePath); err != nil {
		log.Printf("PPTXファイルのパースに失敗: %v\n", err)
	}
}

// parseFromFile はファイルパスからドキュメントをパースする例
func parseFromFile(factory *service.DocumentParserFactory, filePath string) error {
	fmt.Printf("=== ファイルからパース: %s ===\n", filePath)

	// ファイルの拡張子を取得
	ext := getFileExtension(filePath)

	// 拡張子に対応するパーサーを取得
	parser, err := factory.GetParser(ext)
	if err != nil {
		return fmt.Errorf("パーサーの取得に失敗: %w", err)
	}

	// ファイルをパース
	content, err := parser.ParseFromFile(filePath)
	if err != nil {
		return fmt.Errorf("パースに失敗: %w", err)
	}

	fmt.Printf("パース結果:\n%s\n\n", content)
	return nil
}

// parseFromBytes はバイト配列からドキュメントをパースする例
func parseFromBytes(factory *service.DocumentParserFactory, ext string, data []byte) error {
	fmt.Printf("=== バイト配列からパース (拡張子: %s) ===\n", ext)

	// 拡張子に対応するパーサーを取得
	parser, err := factory.GetParser(ext)
	if err != nil {
		return fmt.Errorf("パーサーの取得に失敗: %w", err)
	}

	// バイト配列をパース
	content, err := parser.ParseFromBytes(data)
	if err != nil {
		return fmt.Errorf("パースに失敗: %w", err)
	}

	fmt.Printf("パース結果:\n%s\n\n", content)
	return nil
}

// parseFromReader はio.ReaderAtからドキュメントをパースする例
func parseFromReader(factory *service.DocumentParserFactory, filePath string) error {
	fmt.Printf("=== io.ReaderAtからパース: %s ===\n", filePath)

	// ファイルを開く
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("ファイルのオープンに失敗: %w", err)
	}
	defer file.Close()

	// ファイルサイズを取得
	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("ファイル情報の取得に失敗: %w", err)
	}

	// ファイルの拡張子を取得
	ext := getFileExtension(filePath)

	// 拡張子に対応するパーサーを取得
	parser, err := factory.GetParser(ext)
	if err != nil {
		return fmt.Errorf("パーサーの取得に失敗: %w", err)
	}

	// ReaderAtからパース
	content, err := parser.ParseFromReader(file, stat.Size())
	if err != nil {
		return fmt.Errorf("パースに失敗: %w", err)
	}

	fmt.Printf("パース結果:\n%s\n\n", content)
	return nil
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
