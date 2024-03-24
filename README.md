# Cubism Go

[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen?style=flat-square)](/LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/aethiopicuschan/cubism-go.svg)](https://pkg.go.dev/github.com/aethiopicuschan/cubism-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/aethiopicuschan/cubism-go)](https://goreportcard.com/report/github.com/aethiopicuschan/cubism-go)
[![CI](https://github.com/aethiopicuschan/cubism-go/actions/workflows/ci.yaml/badge.svg)](https://github.com/aethiopicuschan/cubism-go/actions/workflows/ci.yaml)

cubism-goは[Live2D Cubism SDK](https://www.live2d.com/sdk/about/)の非公式版のGolang実装です。[ebitengine/purego](https://github.com/ebitengine/purego)を用いているため扱いやすいです。

## インストール

```bash
go get -u github.com/aethiopicuschan/cubism-go
```

## 動作に必要なもの

- cubism coreの動的ライブラリ
- Live2Dモデル

## 使い方

exampleディレクトリにサンプルコードがあります。おおむね全ての機能を利用したものとなっているので、[Go Reference](https://pkg.go.dev/github.com/aethiopicuschan/cubism-go)と合わせて参照してください。

また、描画の実装として `render/ebitengine` パッケージがあります。
これにより、[Ebitegine](https://ebitengine.org/)を用いたプロジェクトで簡単に利用することができます。もちろん、自身で実装した `renderer` を使うことも可能です。

また、音声の再生のための実装をいくつか用意しています。

- `sound/normal`
  - 一番素直と思われる実装
- `sound/delay`
  - 音声ファイルの読み込みやデコード等を再生時まで遅延させる実装
- `sound/disabled`
  - 音声再生を無効化する実装

こちらも自身で実装することが可能です。
