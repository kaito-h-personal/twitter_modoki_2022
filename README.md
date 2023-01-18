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

## できればやること
- Typescript
- テストコード
- CICD(Github Actions)
- 手続き版も作る

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
```
yarn create vite
✔ Project name: … frontend
✔ Select a framework: › React
✔ Select a variant: › TypeScript + SWC
```
- air導入
  - `go install github.com/cosmtrek/air@latest`
  - `air`で実行
    - http://localhost:8006/ 呼べる
