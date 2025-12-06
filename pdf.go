package documentParser

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/ledongthuc/pdf"
)

// PDFParser はPDFファイルのパーサー
type PDFParser struct {
	BaseParser
}

// SupportedExtensions はサポートする拡張子を返す
func (p *PDFParser) SupportedExtensions() []string {
	return []string{".pdf"}
}

// ParseFromFile はファイルパスからPDFをパース
func (p *PDFParser) ParseFromFile(filePath string) (string, error) {
	return parseFromFileCommon(p, filePath)
}

// ParseFromBytes はバイト配列からPDFをパース
func (p *PDFParser) ParseFromBytes(data []byte) (string, error) {
	return parseFromBytesCommon(p, data)
}

// ParseFromReader はio.ReaderAtからPDFをパース
func (p *PDFParser) ParseFromReader(reader io.ReaderAt, size int64) (string, error) {
	pdfReader, err := pdf.NewReader(reader, size)
	if err != nil {
		return "", fmt.Errorf("error reading PDF: %w", err)
	}

	var result strings.Builder

	// 全てのページからテキストを抽出
	numPages := pdfReader.NumPage()
	for i := 1; i <= numPages; i++ {
		page := pdfReader.Page(i)
		if page.V.IsNull() {
			continue
		}

		// ページ番号を追加
		result.WriteString(fmt.Sprintf("## Page %d\n\n", i))

		// ページからテキストを抽出し、連結
		var pageTexts []string
		texts := page.Content().Text
		for _, text := range texts {
			cleanedText := strings.TrimSpace(text.S)
			if cleanedText != "" {
				pageTexts = append(pageTexts, cleanedText)
			}
		}

		// ページ内のテキストを結合してサニタイズ
		if len(pageTexts) > 0 {
			pageContent := strings.Join(pageTexts, " ")
			sanitizedContent := sanitizeText(pageContent)
			result.WriteString(sanitizedContent)
		}

		result.WriteString("\n\n")
	}

	return result.String(), nil
}

// ParsePdfToString は後方互換性のための既存メソッド
func ParsePdfToString(pdfFilePath string) string {
	parser := &PDFParser{}
	result, err := parser.ParseFromFile(pdfFilePath)
	if err != nil {
		log.Fatalf("Error parsing PDF: %s", err)
	}
	return result
}

// sanitizeText はテキストから余分な空白文字を除去し、正規化する
func sanitizeText(text string) string {
	// 文字化け文字（置換文字）を除去
	text = strings.ReplaceAll(text, "�", "")

	// 全角英数字を半角に変換
	text = convertFullWidthToHalfWidth(text)

	// 日本語文字間の不要なスペースを除去（ひらがな、カタカナ、漢字の間）
	text = removeJapaneseSpaces(text)

	// タブ文字のみスペースに置換（改行は維持）
	text = strings.ReplaceAll(text, "\t", " ")
	text = strings.ReplaceAll(text, "\r\n", "\n") // Windows形式の改行を統一
	text = strings.ReplaceAll(text, "\r", "\n")   // Mac形式の改行を統一

	// 複数の連続スペースを単一のスペースに変換
	for strings.Contains(text, "  ") {
		text = strings.ReplaceAll(text, "  ", " ")
	}

	// 前後の空白を削除
	return strings.TrimSpace(text)
}

// removeJapaneseSpaces は日本語文字間の不要なスペースを除去する
func removeJapaneseSpaces(text string) string {
	var result strings.Builder
	runes := []rune(text)

	for i, r := range runes {
		// スペース文字の場合
		if r == ' ' {
			// 前後の文字をチェック
			prevIsJapanese := i > 0 && isJapaneseChar(runes[i-1])
			nextIsJapanese := i < len(runes)-1 && isJapaneseChar(runes[i+1])

			// 日本語文字間のスペースは除去
			if prevIsJapanese && nextIsJapanese {
				continue
			}
		}
		result.WriteRune(r)
	}

	return result.String()
}

// isJapaneseChar は文字が日本語（ひらがな、カタカナ、漢字）かどうかを判定する
func isJapaneseChar(r rune) bool {
	// ひらがな: U+3040-U+309F
	// カタカナ: U+30A0-U+30FF
	// CJK統合漢字: U+4E00-U+9FAF
	return (r >= 0x3040 && r <= 0x309F) ||
		(r >= 0x30A0 && r <= 0x30FF) ||
		(r >= 0x4E00 && r <= 0x9FAF)
}

// convertFullWidthToHalfWidth は全角英数字を半角に変換する
func convertFullWidthToHalfWidth(text string) string {
	var result strings.Builder

	for _, r := range text {
		// 全角英数字の範囲: U+FF01-U+FF5E
		if r >= 0xFF01 && r <= 0xFF5E {
			// 全角文字から対応する半角文字に変換
			halfWidth := r - 0xFF00 + 0x0020
			result.WriteRune(halfWidth)
		} else {
			// 全角文字でない場合はそのまま
			result.WriteRune(r)
		}
	}

	return result.String()
}
