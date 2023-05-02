# twitter_modoki_2022
- 簡易twitter
- React (+ Typescript) + Go + SurrealDB

### 注意書き
- 個人開発なのでxx flowやWIPは使用しないで一旦main pushしてしまう
- TDDもしない

## やること
- React
  - SPAでやる？
- Go
- SurrealDB
- DDD
- docker

### フロント
- ツイートの枠作る
- アイコン読み込み
- 投稿フォーム
- 少しずつ読み込む
  - スクロールで追加

## できればやること
- Typescript
- テストコード
- CICD(Github Actions)
- 手続き版も作る
- docker-compose、こんなにenv介さなくても良くない？

## 手順
1. goサーバー
2. react
3. 疎通
4. DB導入
5. 機能を詰める

## 構築メモ
### API
- Goのローカルサーバー
  - コンテナに入って`go run main.go`
  - http://localhost:8006/ で呼べる
  - 参考: https://solomaker.club/how-to-create-go-development-environment-with-docker/
- air導入
  - `go install github.com/cosmtrek/air@latest`
  - `air`で実行

### front構築
- `npm install @mui/material @emotion/react @emotion/styled`
  - `npm notice New major version of npm available! 8.15.0 -> 9.5.1`
- `vite`は以下
```
yarn create vite
✔ Project name: … frontend
✔ Select a framework: › React
✔ Select a variant: › TypeScript + SWC
```
- 追加インストール
  - `npm install @mui/icons-material`
  - `npm install @uiw/react-split`
    - TODO: 使わなくなったので削除
  - `npm install @types/react-router-dom`
- useEffectがマウント時に2回実行される
  - StrictModeによるもの
  - 設計上問題ないのでこのままにする
  - 参考: https://qiita.com/asahina820/items/665c55594cfd55e6f14a

### DB構築
- 最初はSurrealDBにHTTPリクエストでクエリを送っていたが、SQLインジェクション対策のためライブラリを導入した

### 起動
- `docker compose up`
- go, react, DBが自動で立ち上がる
  - http://localhost:5173/ をブラウザで開くと画面が表示される

# 設計
- 起動時
  - サーバーからツイートを取得(1)
    - n個
- 追加時
  - 追加取得
- ツイート時
  - サーバーに追加
- リロード時
  - (1)

## サーバーサイド
- 全件取得
- 追加取得
- ツイート追加

### tweetテーブル
- ID
- 日付
- 本文
- userID

### userテーブル
- ID
- 名前
- 表示ID?

CREATE tweet SET
	id = 1,
  auther = 1,
  text = 'I got it.',
	created_at = time::now()
;
