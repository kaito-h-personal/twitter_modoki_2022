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

## メモ
- Goのローカルサーバー
  - コンテナに入って`go run main.go` # TODO: 起動時に実行するようにする
  - http://localhost:8006/ で入れる
  - 参考: https://solomaker.club/how-to-create-go-development-environment-with-docker/
- front構築
  - `npm install @mui/material @emotion/react @emotion/styled`
    - `npm notice New major version of npm available! 8.15.0 -> 9.5.1`
  - `viteは以下`
```
yarn create vite
✔ Project name: … frontend
✔ Select a framework: › React
✔ Select a variant: › TypeScript + SWC
```
- air導入
  - `go install github.com/cosmtrek/air@latest` # TODO: 起動時に実行するようにする
  - `air`で実行
    - http://localhost:8006/ 呼べる
- DB導入
  - 疎通確認は以下
```

DATA="INFO FOR DB;"
curl --request POST --header "Accept: application/json" --header "NS: test" --header "DB: test" --user "root:pasuwado" --data "${DATA}" http://localhost:8009/sql
[{"time":"68.129µs","status":"OK","result":{"dl":{},"dt":{},"sc":{},"tb":{}}}]
```
- 追加インストール
  - `npm install @mui/icons-material`
  - `npm install @uiw/react-split`
    - TODO: 使わなくなったので削除
- apiコンテナ内で`curl db:8009`で疎通確認
  - `http://db:8009/sql`でも可？
- `curl -k -L -s --compressed POST --header "Accept: application/json" --header "NS: test" --header "DB: test" --user "root:pasuwado" --data "INFO FOR DB;" http://db:8
000/sql`
- `curl -X POST --header "Accept: application/json" --header "NS: test" --header "DB: test" --user "root:pasuwado" --data "INFO FOR DB;" http://db:8000/sql`

### 起動
- `docker compose up`
- reactは何もしなくても起動する
  - http://localhost:5173/
- goは起動しないので以下のコマンドを実行
  - `docker compose exec api sh`
    - bashは無い
  - `air`

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
