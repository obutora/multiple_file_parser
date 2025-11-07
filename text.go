package documentParser

import (
	"fmt"
	"io"
)

// TextParser はプレーンテキストファイルのパーサー
type TextParser struct {
	BaseParser
}

// SupportedExtensions はサポートする拡張子を返す
func (p *TextParser) SupportedExtensions() []string {
	return []string{
		".txt",          // プレーンテキスト
		".md",           // Markdown
		".markdown",     // Markdown (別拡張子)
		".mdown",        // Markdown (別拡張子)
		".mkd",          // Markdown (別拡張子)
		".mdwn",         // Markdown (別拡張子)
		".mkdn",         // Markdown (別拡張子)
		".mdtxt",        // Markdown (別拡張子)
		".mdtext",       // Markdown (別拡張子)
		".text",         // プレーンテキスト
		".log",          // ログファイル
		".csv",          // CSVファイル
		".tsv",          // TSVファイル
		".json",         // JSONファイル
		".xml",          // XMLファイル
		".yaml",         // YAMLファイル
		".yml",          // YAMLファイル (別拡張子)
		".toml",         // TOMLファイル
		".ini",          // INIファイル
		".cfg",          // 設定ファイル
		".conf",         // 設定ファイル
		".properties",   // Javaプロパティファイル
		".env",          // 環境変数ファイル
		".sh",           // シェルスクリプト
		".bash",         // Bashスクリプト
		".zsh",          // Zshスクリプト
		".fish",         // Fishスクリプト
		".py",           // Pythonコード
		".js",           // JavaScriptコード
		".ts",           // TypeScriptコード
		".jsx",          // React JSX
		".tsx",          // React TSX
		".go",           // Goコード
		".rs",           // Rustコード
		".c",            // Cコード
		".cpp",          // C++コード
		".h",            // Cヘッダーファイル
		".hpp",          // C++ヘッダーファイル
		".java",         // Javaコード
		".kt",           // Kotlinコード
		".swift",        // Swiftコード
		".rb",           // Rubyコード
		".php",          // PHPコード
		".sql",          // SQLファイル
		".html",         // HTMLファイル
		".htm",          // HTMLファイル (別拡張子)
		".css",          // CSSファイル
		".scss",         // SCSSファイル
		".sass",         // SASSファイル
		".less",         // LESSファイル
		".vue",          // Vueファイル
		".svelte",       // Svelteファイル
		".r",            // Rコード
		".R",            // Rコード (大文字)
		".m",            // MATLABコード
		".pl",           // Perlコード
		".lua",          // Luaコード
		".dart",         // Dartコード
		".scala",        // Scalaコード
		".clj",          // Clojureコード
		".ex",           // Elixirコード
		".exs",          // Elixirスクリプト
		".erl",          // Erlangコード
		".hrl",          // Erlangヘッダー
		".fs",           // F#コード
		".fsx",          // F#スクリプト
		".vb",           // Visual Basicコード
		".bas",          // BASICコード
		".pas",          // Pascalコード
		".asm",          // アセンブリコード
		".s",            // アセンブリコード
		".ps1",          // PowerShellスクリプト
		".psm1",         // PowerShellモジュール
		".bat",          // Windowsバッチファイル
		".cmd",          // Windowsコマンドファイル
		".makefile",     // Makefile
		".mk",           // Makefile
		".dockerfile",   // Dockerfile
		".gitignore",    // Gitignore
		".editorconfig", // EditorConfig
		".htaccess",     // Apache設定
		".nginx",        // Nginx設定
		".vimrc",        // Vim設定
		".zshrc",        // Zsh設定
		".bashrc",       // Bash設定
	}
}

// ParseFromBytes はバイト配列をそのまま文字列として返す
func (p *TextParser) ParseFromBytes(data []byte) (string, error) {
	// サイズ制限を設定（最大100MB）
	const maxSize = 100 * 1024 * 1024 // 100MB
	if len(data) > maxSize {
		return "", fmt.Errorf("file size %d exceeds maximum allowed size of %d bytes", len(data), maxSize)
	}
	return string(data), nil
}

// ParseFromReader はio.ReaderAtからテキストを読み込んでそのまま返す
func (p *TextParser) ParseFromReader(reader io.ReaderAt, size int64) (string, error) {
	// io.ReaderAt を io.Reader に変換
	// サイズ制限を設定（最大100MB）
	const maxSize = 100 * 1024 * 1024 // 100MB
	if size > maxSize {
		return "", fmt.Errorf("file size %d exceeds maximum allowed size of %d bytes", size, maxSize)
	}

	// バッファを作成してデータを読み込む
	buffer := make([]byte, size)
	n, err := reader.ReadAt(buffer, 0)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("error reading text file: %w", err)
	}

	// 読み込んだデータを文字列として返す
	return string(buffer[:n]), nil
}

// ParseTextToString は後方互換性のための既存メソッド
func ParseTextToString(textFilePath string) (string, error) {
	parser := &TextParser{}
	return parser.ParseFromFile(textFilePath)
}
