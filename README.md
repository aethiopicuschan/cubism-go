# Cubism Go

[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen?style=flat-square)](/LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/aethiopicuschan/cubism-go.svg)](https://pkg.go.dev/github.com/aethiopicuschan/cubism-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/aethiopicuschan/cubism-go)](https://goreportcard.com/report/github.com/aethiopicuschan/cubism-go)
[![CI](https://github.com/aethiopicuschan/cubism-go/actions/workflows/ci.yaml/badge.svg)](https://github.com/aethiopicuschan/cubism-go/actions/workflows/ci.yaml)

cubism-go is an unofficial Golang implementation of the [Live2D Cubism SDK](https://www.live2d.com/sdk/about/). It leverages [ebitengine/purego](https://github.com/ebitengine/purego), making it easy to use.

## Installation

```bash
go get -u github.com/aethiopicuschan/cubism-go
```

## Requirements

- Dynamic library for Cubism Core
- Live2D model

## Usage

Sample code is available in the `example` directory. It demonstrates the use of almost all functionalities, so please refer to it alongside the [Go Reference](https://pkg.go.dev/github.com/aethiopicuschan/cubism-go).

Additionally, there is a `renderer/ebitengine` package for rendering implementations.
This package enables seamless integration with projects using [Ebiten](https://ebitengine.org/). Of course, you can also use your custom `renderer`.

Moreover, there are several implementations available for audio playback:

- `sound/normal`
  - A straightforward implementation
- `sound/delay`
  - An implementation that defers loading and decoding of audio files until playback
- `sound/disabled`
  - An implementation that disables audio playback

You can also implement your own version of these.

## Development

For `pre-commit` hooks, we use [lefthook](https://github.com/evilmartians/lefthook). The configured tools include:

- [staticcheck](https://staticcheck.dev)
- [typos](https://github.com/crate-ci/typos)

If you're using [Homebrew](https://brew.sh/), you can install all the necessary tools with the following command:

```sh
brew install lefthook staticcheck typos-cli
```

After that, run `lefthook install` to enable the hooks.

The same checks are also performed using GitHub Actions.
