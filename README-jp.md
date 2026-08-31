# Clean Invisible Text

> **このファイルは正本(日本語版)です。**
> 英語版(参照)は [README.md](README.md) を参照してください。

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![CI](https://github.com/y-marui/go-clean-invisible-text/actions/workflows/ci.yml/badge.svg)](https://github.com/y-marui/go-clean-invisible-text/actions/workflows/ci.yml)
[![Charter Check](https://github.com/y-marui/go-clean-invisible-text/actions/workflows/dev-charter-check.yml/badge.svg)](https://github.com/y-marui/go-clean-invisible-text/actions/workflows/dev-charter-check.yml)
[![GitHub Sponsors](https://img.shields.io/github/sponsors/y-marui?style=social)](https://github.com/sponsors/y-marui)
[![Buy Me a Coffee](https://img.shields.io/badge/Buy%20Me%20a%20Coffee-donate-yellow.svg)](https://www.buymeacoffee.com/y.marui)

UTF-8 プレーンテキスト中の危険な不可視 Unicode 文字を検出・説明し、安全にクリーニングするクロスプラットフォーム対応の Go CLI。

> **Status:** v1.0.0 リリース済み。`check`/`fix`/`explain`/`clean`、pre-commit
> パッケージング、クロスプラットフォームリリース自動化、v1.0 互換性・サポート方針
> ([ADR 0002](docs/decisions/0002-v1-compatibility-and-support-policy.md)) をすべて出荷済み
> ([ロードマップ](https://github.com/y-marui/go-clean-invisible-text/issues/1) 参照)。

## Requirements

Windows・macOS・Raspberry Pi でスタンドアロンバイナリとして動作することを目標にしている。ビルドには Go のみが必要で、Node.js や Python 等の追加ランタイムは不要。公式リリースバイナリの最小 OS バージョンと対応する Raspberry Pi アーキテクチャは [ADR 0002](docs/decisions/0002-v1-compatibility-and-support-policy.md) で定義する。ソースからビルドしたバイナリは、実際に使用した Go ツールチェーンの要件に従う。

## Setup

```bash
go install github.com/y-marui/go-clean-invisible-text/cmd/clean-invisible-text@latest
```

またはクローンしてビルド:

```bash
git clone https://github.com/y-marui/go-clean-invisible-text.git
cd go-clean-invisible-text
go build -o clean-invisible-text ./cmd/clean-invisible-text
go test ./...
```

## Usage

| コマンド | 説明 |
|---|---|
| `clean-invisible-text check FILE...` | 変更せず検出結果のみ報告する |
| `clean-invisible-text fix FILE...` | 指定ファイルを修正し、すべての変更を報告する |
| `clean-invisible-text explain FILE...` | コードポイント・Unicode 名・位置・カテゴリ・実施予定のアクションを表示する |
| `clean-invisible-text clean` | 標準入力を読み、クリーニング済みテキストを標準出力へ書き出す |

```console
$ clean-invisible-text check notes.txt
notes.txt: 2 finding(s)

$ echo "hello world" | clean-invisible-text clean
hello world
```

`check`/`fix`/`explain` に `--json` を付けると機械可読な出力になり、`fix`/`clean`
に `--keep-warnings` を付けると Warn 判定のコードポイントを除去せず残せる。どのコマンドにも
`--allow`/`--allow-file` を付けると、特定の Warn 判定コードポイントに対して監査可能な例外を
付与できる。契約の全体は [docs/cli.md](docs/cli.md) を参照。

正規の仕様は [docs/specification.md](docs/specification.md) と [docs/character-policy.md](docs/character-policy.md) を参照。作業中の議論は仕様ファイルではなく GitHub Issues で管理する。

## Documentation

- [docs/specification.md](docs/specification.md) — 機能仕様
- [docs/character-policy.md](docs/character-policy.md) — 文字ポリシー(正本)
- [docs/security-model.md](docs/security-model.md) — セキュリティモデル
- [docs/cli.md](docs/cli.md) — CLI 契約
- [docs/architecture.md](docs/architecture.md) — アーキテクチャ
- [docs/integrations/pre-commit.md](docs/integrations/pre-commit.md) — pre-commit フック契約
- [docs/release-process.md](docs/release-process.md) — リリース手順と検証方法
- [docs/decisions/](docs/decisions/) — アーキテクチャ決定記録(ADR)

## pre-commit Integration

```yaml
repos:
  - repo: https://github.com/y-marui/go-clean-invisible-text
    rev: vX.Y.Z # リリース済みタグを指定
    hooks:
      - id: clean-invisible-text-check
      - id: clean-invisible-text-fix
```

インストール済みバイナリを使う方法を含む契約の全体は [docs/integrations/pre-commit.md](docs/integrations/pre-commit.md) を参照。

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
