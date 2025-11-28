package documentParser

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
			err = func() error {
				defer rc.Close()
				decoder := xml.NewDecoder(rc)
				inBody := false
				for {
					t, err := decoder.Token()
					if err == io.EOF {
						break
					}
					if err != nil {
						return fmt.Errorf("error parsing XML: %w", err)
					}

					switch se := t.(type) {
					case xml.StartElement:
						if se.Name.Local == "body" {
							inBody = true
						}
						if inBody {
							if se.Name.Local == "p" {
								var p DocxParagraph
								if err := decoder.DecodeElement(&p, &se); err != nil {
									return err
								}
								text := extractTextFromParagraph(p)
								if text != "" {
									allText.WriteString(text + "\n")
								}
							} else if se.Name.Local == "tbl" {
								var tbl DocxTable
								if err := decoder.DecodeElement(&tbl, &se); err != nil {
									return err
								}
								allText.WriteString(extractTextFromTable(tbl))
							}
						}
					case xml.EndElement:
						if se.Name.Local == "body" {
							inBody = false
						}
					}
				}
				return nil
			}()

			if err != nil {
				return "", err
			}

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

// テーブル構造体
type DocxTable struct {
	Rows []DocxTableRow `xml:"tr"`
}

type DocxTableRow struct {
	Cells []DocxTableCell `xml:"tc"`
}

type DocxTableCell struct {
	Paragraphs []DocxParagraph `xml:"p"`
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

func extractTextFromParagraph(p DocxParagraph) string {
	var paragraphText strings.Builder
	for _, run := range p.Runs {
		paragraphText.WriteString(run.Text.Content)
	}
	return paragraphText.String()
}

func extractTextFromTable(tbl DocxTable) string {
	var sb strings.Builder
	for _, row := range tbl.Rows {
		rowTexts := []string{}
		for _, cell := range row.Cells {
			for _, p := range cell.Paragraphs {
				text := extractTextFromParagraph(p)
				if text != "" {
					rowTexts = append(rowTexts, text)
				}
			}
		}
		if len(rowTexts) > 0 {
			sb.WriteString(strings.Join(rowTexts, "\t") + "\n")
		}
	}
	return sb.String()
}
