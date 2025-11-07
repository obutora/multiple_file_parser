package main

import (
	"fmt"
	"log"
	"os"

	documentParser "github.com/obutora/multiple_file_parser"
)

func main() {
	// ドキュメントパーサーファクトリーを作成
	factory := documentParser.NewDocumentParserFactory()

	// サポートされている拡張子を表示
	fmt.Println("サポートされている拡張子:")
	for _, ext := range factory.SupportedExtensions() {
		fmt.Printf("  - %s\n", ext)
	}
	fmt.Println()

	// 例1: ファイルパスからパース（PDF）
	pdfFilePath := "assets/sample.pdf"
	fmt.Printf("=== ファイルからパース: %s ===\n", pdfFilePath)
	content, err := factory.ParseFromFile(pdfFilePath)
	if err != nil {
		log.Printf("PDFファイルのパースに失敗: %v\n", err)
	} else {
		fmt.Printf("パース結果:\n%s\n\n", content)
	}

	// 例2: バイト配列からパース（テキスト）
	exampleBytes := []byte("これはテキストファイルの内容です。\nサンプルテキスト。")
	fmt.Printf("=== バイト配列からパース (拡張子: .txt) ===\n")
	content, err = factory.ParseFromBytes(".txt", exampleBytes)
	if err != nil {
		log.Printf("バイト配列のパースに失敗: %v\n", err)
	} else {
		fmt.Printf("パース結果:\n%s\n\n", content)
	}

	// 例3: io.ReaderAtからパース（DOCX）
	docxFilePath := "assets/sample.docx"
	fmt.Printf("=== io.ReaderAtからパース: %s ===\n", docxFilePath)
	file, err := os.Open(docxFilePath)
	if err != nil {
		log.Printf("ファイルのオープンに失敗: %v\n", err)
	} else {
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			log.Printf("ファイル情報の取得に失敗: %v\n", err)
		} else {
			content, err = factory.ParseFromReader(".docx", file, stat.Size())
			if err != nil {
				log.Printf("DOCXファイルのパースに失敗: %v\n", err)
			} else {
				fmt.Printf("パース結果:\n%s\n\n", content)
			}
		}
	}

	// 例4: PPTXファイルのパース
	pptxFilePath := "assets/sample.pptx"
	fmt.Printf("=== ファイルからパース: %s ===\n", pptxFilePath)
	content, err = factory.ParseFromFile(pptxFilePath)
	if err != nil {
		log.Printf("PPTXファイルのパースに失敗: %v\n", err)
	} else {
		fmt.Printf("パース結果:\n%s\n\n", content)
	}
}
