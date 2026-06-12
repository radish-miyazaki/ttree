# ttree

[English](README.md) | 日本語

ターミナルで動作するインタラクティブなツリー構造エディタです。vim 風のキー操作で階層構造を作成・編集・整理し、ASCII ツリーとしてクリップボードにコピーできます。

## デモ

<img src="vhs/demo.gif" alt="ttree demo" width="600">

## 特徴

- ターミナル上でのインタラクティブなツリー編集
- リアルタイムの ASCII ツリープレビュー
- vim 風のカーソル移動(`j`/`k`)
- Tab キーによるノードのインデント / アンインデント
- レンダリングしたツリーをクリップボードにコピー

## インストール

```bash
go install github.com/radish-miyazaki/ttree@latest
```

またはソースからビルド:

```bash
git clone https://github.com/radish-miyazaki/ttree.git
cd ttree
go build -o ttree .
```

## 使い方

```bash
./ttree
```

### キーバインド

| キー | 動作 |
|-----|------|
| `↑` / `k` | 上に移動 |
| `↓` / `j` | 下に移動 |
| `Tab` | ノードをインデント(直前の兄弟ノードの子にする) |
| `Shift+Tab` | ノードをアンインデント(親の兄弟ノードにする) |
| `Enter` | 新しい兄弟ノードを作成 |
| `Ctrl+D` | 現在のノードを削除 |
| `Ctrl+C` | ツリーをクリップボードにコピー |
| `Ctrl+Q` / `Esc` | 終了 |

### 出力例

```
├── src
│   ├── components
│   │   ├── Header.tsx
│   │   └── Footer.tsx
│   └── utils
│       └── helpers.ts
├── tests
└── README.md
```

## ライセンス

MIT
