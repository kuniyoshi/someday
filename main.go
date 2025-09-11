package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

var version = "0.4.0"

func main() {
	// 引数なし: 既定ファイルから読み込み表示
	if len(os.Args) == 1 {
		if err := runWithFile(defaultFilePath()); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	// サブコマンド処理
	switch os.Args[1] {
	case "help", "-h", "--help":
		printUsage()
	case "version", "-v", "--version":
		fmt.Println(version)
	case "-f", "--file":
		path := ""
		if len(os.Args) >= 3 {
			path = os.Args[2]
		}
		if path == "" {
			fmt.Fprintln(os.Stderr, "ファイルパスを指定してください: someday -f <path>")
			os.Exit(2)
		}
		if err := runWithFile(path); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(2)
	}
}

// 既定のファイルパスを返す
func defaultFilePath() string {
	if p := os.Getenv("SOMEDAY_FILE"); strings.TrimSpace(p) != "" {
		return p
	}
	return "someday.json"
}

// 入力ファイルを読み込み、表示する
func runWithFile(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("ファイルを読み込めませんでした (%s)。\nJSON ファイルを指定してください。例: someday -f data.json または SOMEDAY_FILE=someday.json\n%v", path, err)
	}
	if err := renderFromJSON(strings.NewReader(string(b))); err != nil {
		return fmt.Errorf("JSON のパースに失敗しました。配列または { \"items\": [] } を使用してください。\n%w", err)
	}
	return nil
}

// データモデル: ハードルと重要度のペア
type Hurdle struct {
	Type         string `json:"type"`
	Significance string `json:"significance"`
}

type Item struct {
	Description string   `json:"description"`
	Hurdles     []Hurdle `json:"hurdles"`
}

type ItemsFile struct {
	Items []Item `json:"items"`
}

// JSON を読み込んで README のスクリーンに近い形式で表示
func renderFromJSON(contentReader io.Reader) error {
	// contentReader は strings.NewReader から渡される
	// まず全部読み込む
	b, _ := io.ReadAll(contentReader)

	// 1) 配列として試行
	var list []Item
	if err := json.Unmarshal(b, &list); err == nil {
		printItems(list)
		return nil
	}

	// 2) { items: [] } として試行
	var wrapped ItemsFile
	if err := json.Unmarshal(b, &wrapped); err == nil {
		printItems(wrapped.Items)
		return nil
	}

	return fmt.Errorf("JSON のパースに失敗しました")
}

// 余分な関数は不要

func printItems(items []Item) {
	for _, it := range items {
		fmt.Printf("- %s\n", it.Description)
		// hurdles のみを表示
		for _, h := range it.Hurdles {
			lh := strings.TrimSpace(h.Type)
			rh := strings.TrimSpace(h.Significance)
			if lh != "" && rh != "" {
				fmt.Printf("  %s  %s\n", lh, rh)
			} else if lh != "" {
				fmt.Printf("  %s\n", lh)
			} else if rh != "" {
				fmt.Printf("  %s\n", rh)
			}
		}
	}
}

func printUsage() {
	fmt.Println("someday - 前向きなウィッシュリスト CLI")
	fmt.Println()
	fmt.Println("使い方:")
	fmt.Println("  someday                      既定 JSON (someday.json) を読み込み表示")
	fmt.Println("  someday -f <path>            指定 JSON を読み込み表示")
	fmt.Println("  SOMEDAY_FILE=<path> someday  環境変数で JSON 指定")
	fmt.Println("  someday help            ヘルプを表示")
	fmt.Println("  someday version         バージョンを表示")
	fmt.Println()
	fmt.Println("入力形式:")
	fmt.Println("  - JSON のみ: 配列、または { \"items\": [] } で渡す。")
	fmt.Println("    形式: { description, hurdles: [{ type, significance }] }")
}
