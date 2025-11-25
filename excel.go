package documentParser

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
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

func (p *ExcelParser) ParseFromBytes(data []byte) (string, error) {
	reader := bytes.NewReader(data)
	return p.ParseFromReader(reader, int64(len(data)))
}

func (p *ExcelParser) ParseFromReader(reader io.ReaderAt, size int64) (string, error) {
	f, err := excelize.OpenReader(io.NewSectionReader(reader, 0, size))
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	sheetList := f.GetSheetList()

	var buf strings.Builder
	for _, sheet := range sheetList {
		buf.WriteString(fmt.Sprintf("# Sheet %s\n", sheet))

		rows, err := f.Rows(sheet)
		if err != nil {
			log.Printf("failed to get rows: %v\n", err)
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

		buf.WriteString("\n---\n\n")
	}

	if len(buf.String()) == 0 {
		return "", fmt.Errorf("no data found")
	}

	return buf.String(), nil
}
