package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

var version = "0.1.0"

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

// データモデル
type Hurdle struct {
	Name  string `json:"name"`
	Level string `json:"level"`
}

type Item struct {
	Description     string   `json:"description"`
	Hurdles         []string `json:"hurdles,omitempty"`
	HurdlesDetailed []Hurdle `json:"hurdles_detailed,omitempty"`
	Significance    string   `json:"significance,omitempty"`
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
	for i, it := range items {
		if i > 0 {
			// ビジュアル上の区切りとして空行を入れない（README の例に合わせて連続出力）
		}
		fmt.Printf("- %s\n", it.Description)
		// 詳細ハードルがあればそれを優先
		if len(it.HurdlesDetailed) > 0 {
			parts := make([]string, 0, len(it.HurdlesDetailed))
			for _, h := range it.HurdlesDetailed {
				if strings.TrimSpace(h.Level) != "" {
					parts = append(parts, fmt.Sprintf("%s (%s)", h.Name, h.Level))
				} else {
					parts = append(parts, h.Name)
				}
			}
			fmt.Printf("  %s\n", strings.Join(parts, ", "))
			continue
		}

		// 通常の hurdles + significance
		left := strings.Join(it.Hurdles, ", ")
		right := strings.TrimSpace(it.Significance)
		if left == "" && right == "" {
			// 追加情報なし
			continue
		}
		if left != "" && right != "" {
			// Usage の例に合わせて (Significance) を括弧で併記
			fmt.Printf("  %s (%s)\n", left, right)
		} else if left != "" {
			fmt.Printf("  %s\n", left)
		} else {
			fmt.Printf("  %s\n", right)
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
    fmt.Println("    例: { \"description\": \"...\", \"hurdles\": [\"知識\"], \"significance\": \"中\" }")
    fmt.Println("    または hurdles_detailed: [{ name, level }] で '時間 (大), アイデア (中)' の表記に対応。")
}
