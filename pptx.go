package documentParser

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// PPTXParser はPowerPointファイルのパーサー
type PPTXParser struct {
	BaseParser
}

// SupportedExtensions はサポートする拡張子を返す
func (p *PPTXParser) SupportedExtensions() []string {
	return []string{".pptx", ".ppt"}
}

// ParseFromFile はファイルパスからPPTXをパース
func (p *PPTXParser) ParseFromFile(filePath string) (string, error) {
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

// ParseFromBytes はバイト配列からPPTXをパース
func (p *PPTXParser) ParseFromBytes(data []byte) (string, error) {
	reader := bytes.NewReader(data)
	return p.ParseFromReader(reader, int64(len(data)))
}

// ParseFromReader はio.ReaderAtからPPTXをパース
func (p *PPTXParser) ParseFromReader(reader io.ReaderAt, size int64) (string, error) {
	r, err := zip.NewReader(reader, size)
	if err != nil {
		return "", fmt.Errorf("error reading PowerPoint: %w", err)
	}

	var allText strings.Builder
	slideNum := 1

	// 各ファイルをチェック
	for _, f := range r.File {
		// スライドファイルのみを処理
		if strings.HasPrefix(f.Name, "ppt/slides/slide") &&
			strings.HasSuffix(f.Name, ".xml") &&
			!strings.Contains(f.Name, "Layout") &&
			!strings.Contains(f.Name, "Master") {

			rc, err := f.Open()
			if err != nil {
				log.Printf("Error opening file %s: %s", f.Name, err)
				continue
			}

			// XMLをパース
			content, err := io.ReadAll(rc)
			if err != nil {
				log.Printf("Error reading file %s: %s", f.Name, err)
				rc.Close()
				continue
			}
			rc.Close()

			var slide Slide
			err = xml.Unmarshal(content, &slide)
			if err != nil {
				log.Printf("Error parsing XML for %s: %s", f.Name, err)
				continue
			}

			// テキストを抽出
			extractedText := extractTextFromSlide(slide)

			// スライド番号とテキストを追加
			allText.WriteString(fmt.Sprintf("## Slide %d\n", slideNum))
			if len(extractedText) > 0 {
				allText.WriteString(extractedText)
			} else {
				allText.WriteString("(No text found)")
			}
			allText.WriteString("\n\n")

			slideNum++
		}
	}

	return allText.String(), nil
}

// const (
// 	// pptxFilePath は、テキストを抽出するPowerPointファイルのパスです。
// 	pptxFilePath = "assets/office/AI Research.pptx"
// )

// PowerPointのXML構造を表現する構造体
type TextRun struct {
	Text string `xml:"t"`
}

type Paragraph struct {
	Runs []TextRun `xml:"r"`
}

type TextBody struct {
	Paragraphs []Paragraph `xml:"p"`
}

type Shape struct {
	TextBody TextBody `xml:"txBody"`
}

type SlideData struct {
	Shapes []Shape `xml:"spTree>sp"`
}

type Slide struct {
	SlideData SlideData `xml:"cSld"`
}

// ParsePptxToString は後方互換性のための既存メソッド
func ParsePptxToString(pptxFilePath string) string {
	parser := &PPTXParser{}
	result, err := parser.ParseFromFile(pptxFilePath)
	if err != nil {
		log.Fatalf("Error parsing PowerPoint: %s", err)
	}
	return result
}

func extractTextFromSlide(slide Slide) string {
	var result []string

	for _, shape := range slide.SlideData.Shapes {
		for _, paragraph := range shape.TextBody.Paragraphs {
			var paragraphText strings.Builder
			for _, run := range paragraph.Runs {
				if run.Text != "" {
					paragraphText.WriteString(run.Text)
				}
			}
			if paragraphText.Len() > 0 {
				result = append(result, paragraphText.String())
			}
		}
	}

	return strings.Join(result, "\n")
}
