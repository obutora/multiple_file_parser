package service

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"strings"
)

// DOCXParser はWordファイルのパーサー
type DOCXParser struct {
	BaseParser
}

// SupportedExtensions はサポートする拡張子を返す
func (p *DOCXParser) SupportedExtensions() []string {
	return []string{".docx", ".doc"}
}

// ParseFromReader はio.ReaderAtからDOCXをパース
func (p *DOCXParser) ParseFromReader(reader io.ReaderAt, size int64) (string, error) {
	r, err := zip.NewReader(reader, size)
	if err != nil {
		return "", fmt.Errorf("error reading Word file: %w", err)
	}

	var allText strings.Builder

	// word/document.xmlファイルを探す
	for _, f := range r.File {
		if f.Name == "word/document.xml" {
			rc, err := f.Open()
			if err != nil {
				return "", fmt.Errorf("error opening file %s: %w", f.Name, err)
			}

			// XMLをパース
			content, err := io.ReadAll(rc)
			if err != nil {
				rc.Close()
				return "", fmt.Errorf("error reading file %s: %w", f.Name, err)
			}
			rc.Close()

			var document DocxDocument
			err = xml.Unmarshal(content, &document)
			if err != nil {
				return "", fmt.Errorf("error parsing XML for %s: %w", f.Name, err)
			}

			// テキストを抽出
			extractedText := extractTextFromDocument(document)
			allText.WriteString(extractedText)

			break // document.xmlは1つなので見つけたら終了
		}
	}

	return allText.String(), nil
}

// WordのXML構造を表現する構造体
type DocxText struct {
	Content string `xml:",chardata"`
}

type DocxRun struct {
	Text DocxText `xml:"t"`
}

type DocxParagraph struct {
	Runs []DocxRun `xml:"r"`
}

type DocxBody struct {
	Paragraphs []DocxParagraph `xml:"p"`
}

type DocxDocument struct {
	Body DocxBody `xml:"body"`
}

// ParseDocxToString は後方互換性のための既存メソッド
func ParseDocxToString(docxFilePath string) string {
	parser := &DOCXParser{}
	result, err := parser.ParseFromFile(docxFilePath)
	if err != nil {
		log.Fatalf("Error parsing Word: %s", err)
	}
	return result
}

func extractTextFromDocument(document DocxDocument) string {
	var result []string

	for _, paragraph := range document.Body.Paragraphs {
		var paragraphText strings.Builder

		for _, run := range paragraph.Runs {
			if run.Text.Content != "" {
				paragraphText.WriteString(run.Text.Content)
			}
		}

		if paragraphText.Len() > 0 {
			result = append(result, paragraphText.String())
		} else {
			// 空の段落も改行として扱う
			result = append(result, "")
		}
	}

	return strings.Join(result, "\n")
}
