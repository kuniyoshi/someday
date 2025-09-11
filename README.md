# Someday

Someday is a wish list tool likes TODO list but, more positive purpose.

## Usage

```
$ someday
- Copy on Write を実装する (構造共有)
  知識 (中)
- タイピングゲームを作る
  時間 (大), アイデア (中)
- sketch スクリプトの再利用性を高める
  時間 (小)
$
```

## Design

## Screen

List is formatted below.

- Description
  {Hurdle    Significance}{1..3}  // 複数行のペア

e.g.,

- Copy on Write を実装する (構造共有)
  知識      中

## 値域

Description:
  - 40 文字
Hurdle:
  - 知識
  - 時間
  - その他
Significance:
  - 大
  - 中
  - 小

## JSON 形式

- 配列、または `{ "items": [] }`
- 各要素は以下の形式:

```
{
  "description": "...",
  "hurdles": [
    { "type": "時間", "significance": "大" },
    { "type": "アイデア", "significance": "中" }
  ]
}
```

後方互換はありません。`hurdles` のみを使用してください。
