SHELL := /bin/bash

.PHONY: all fmt vet staticcheck tools run build

# デフォルト: すべて実行
all: fmt vet staticcheck

# gofmt による整形
fmt:
	go fmt ./...

# go vet による静的解析
vet:
	go vet ./...

# staticcheck による静的解析
staticcheck:
	@command -v staticcheck >/dev/null || { \
		echo "staticcheck が見つかりません。インストール: go install honnef.co/go/tools/cmd/staticcheck@latest"; \
		exit 1; \
	}
	staticcheck ./...

# 開発環境にツールをインストール
tools:
	go install honnef.co/go/tools/cmd/staticcheck@latest

# アプリを実行（引数は ARGS で指定）
ARGS ?=
run:
	go run . $(ARGS)

# ビルド（出力先は BINARY で指定）
BINARY ?= bin/someday
build:
	@mkdir -p $(dir $(BINARY))
	go build -o $(BINARY) .
