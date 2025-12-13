package documentParser

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/xuri/excelize/v2"
)

type ExcelParser struct {
	BaseParser
}

func (p *ExcelParser) SupportedExtensions() []string {
	return []string{".xlsx", ".xls"}
}

func (p *ExcelParser) ParseFromFile(filePath string) (string, error) {
	return parseFromFileCommon(p, filePath)
}

func (p *ExcelParser) ParseFromBytes(data []byte) (string, error) {
	return parseFromBytesCommon(p, data)
}

// sheetContent はシート名と内容を保持する構造体
type sheetContent struct {
	name    string
	content string
}

// extractSheets はExcelファイルから全シートの内容を抽出する
func (p *ExcelParser) extractSheets(reader io.ReaderAt, size int64) ([]sheetContent, error) {
	f, err := excelize.OpenReader(io.NewSectionReader(reader, 0, size))
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	sheetList := f.GetSheetList()
	var results []sheetContent

	for _, sheet := range sheetList {
		var buf strings.Builder

		rows, err := f.Rows(sheet)
		if err != nil {
			log.Printf("failed to get rows for sheet %s: %v\n", sheet, err)
			continue
		}

		for rows.Next() {
			row, err := rows.Columns()
			if err != nil {
				log.Printf("failed to get row: %v\n", err)
				continue
			}
			buf.WriteString(fmt.Sprintf("%v\n", strings.Join(row, " | ")))
		}

		results = append(results, sheetContent{
			name:    sheet,
			content: buf.String(),
		})
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no data found")
	}

	return results, nil
}

func (p *ExcelParser) ParseFromReader(reader io.ReaderAt, size int64) (string, error) {
	sheets, err := p.extractSheets(reader, size)
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	for _, sheet := range sheets {
		buf.WriteString(fmt.Sprintf("# Sheet %s\n", sheet.name))
		buf.WriteString(sheet.content)
		buf.WriteString("\n---\n\n")
	}

	if buf.Len() == 0 {
		return "", fmt.Errorf("no data found")
	}

	return buf.String(), nil
}

// ParseWithPages はシートごとに内容を分けてマップ形式で返す
func (p *ExcelParser) ParseWithPages(reader io.ReaderAt, size int64) (map[string]string, error) {
	sheets, err := p.extractSheets(reader, size)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, sheet := range sheets {
		result[sheet.name] = sheet.content
	}

	return result, nil
}
