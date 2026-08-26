# Clean Invisible Text

> **このファイルは正本(日本語版)です。**
> 英語版(参照)は [README.md](README.md) を参照してください。

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![CI](https://github.com/y-marui/go-clean-invisible-text/actions/workflows/ci.yml/badge.svg)](https://github.com/y-marui/go-clean-invisible-text/actions/workflows/ci.yml)
[![Charter Check](https://github.com/y-marui/go-clean-invisible-text/actions/workflows/dev-charter-check.yml/badge.svg)](https://github.com/y-marui/go-clean-invisible-text/actions/workflows/dev-charter-check.yml)
[![GitHub Sponsors](https://img.shields.io/github/sponsors/y-marui?style=social)](https://github.com/sponsors/y-marui)
[![Buy Me a Coffee](https://img.shields.io/badge/Buy%20Me%20a%20Coffee-donate-yellow.svg)](https://www.buymeacoffee.com/y.marui)

UTF-8 プレーンテキスト中の危険な不可視 Unicode 文字を検出・説明し、安全にクリーニングするクロスプラットフォーム対応の Go CLI。

> **Status:** 仕様とロードマップの段階。CLI コマンドはまだ実装されていない(検出・クリーニングエンジンは `internal/cleaner` に実装済み)。

## Requirements

Windows・macOS・Raspberry Pi でスタンドアロンバイナリとして動作することを目標にしている。ビルドには Go のみが必要で、Node.js や Python 等の追加ランタイムは不要。

## Setup

```bash
git clone https://github.com/y-marui/go-clean-invisible-text.git
cd go-clean-invisible-text
go build ./...
go test ./...
```

## Usage

CLI コマンドは未実装。計画しているコマンド一覧:

~~~console
clean-invisible-text check FILE...
clean-invisible-text fix FILE...
clean-invisible-text explain FILE...
clean-invisible-text clean
~~~

正規の仕様は [docs/specification.md](docs/specification.md) と [docs/character-policy.md](docs/character-policy.md) を参照。作業中の議論は仕様ファイルではなく GitHub Issues で管理する。

## Documentation

- [docs/specification.md](docs/specification.md) — 機能仕様
- [docs/character-policy.md](docs/character-policy.md) — 文字ポリシー(正本)
- [docs/security-model.md](docs/security-model.md) — セキュリティモデル
- [docs/cli.md](docs/cli.md) — CLI 契約
- [docs/architecture.md](docs/architecture.md) — アーキテクチャ
- [docs/decisions/](docs/decisions/) — アーキテクチャ決定記録(ADR)

## Alfred Integration

別リポジトリの [alfred-clean-invisible-text](https://github.com/y-marui/alfred-clean-invisible-text) がこの CLI を Alfred 向けにパッケージ化する予定。

## Contributing

[CONTRIBUTING.md](CONTRIBUTING.md) を参照。

## Changelog

[CHANGELOG.md](CHANGELOG.md) を参照。

## License

[MIT](LICENSE)

---
*この文書には英語版 [README.md](README.md) があります。編集時は同一コミットで更新してください。*
