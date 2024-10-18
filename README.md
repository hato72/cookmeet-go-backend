# CookMeet(サーバー)
![cookmeet](https://github.com/hato72/go_backend_hackathon/assets/139688965/54235b01-2da0-491e-857c-18581b70b518)

## デプロイ(自動OFF)
https://cookmeet-backend.onrender.com

## フロントエンド
ソースコード：https://github.com/hato72/CookMeet

デプロイ先：https://cook-meet.vercel.app/

## DB設計
https://free-casquette-dee.notion.site/d558148d80f742a4ac77c0bf76b4a2c9?pvs=4


## 実行方法(テスト環境)

```sh
.env.dev:

PORT=8081
POSTGRES_USER=
POSTGRES_PW=
POSTGRES_DB=
POSTGRES_PORT=
POSTGRES_HOST=
GO_ENV=dev
SECRET=<supabaseのJWT　Secretキー>
FE_URL=https://localhost:3000
```

.env.devをbackendディレクトリ直下に配置した後に以下を実行

<!-- ```sh
docker compose build

docker compose up

docker compose run --rm backend sh

go run src/migrate/migrate.go

go run src/main.go

``` -->

```sh
docker compose build

docker compose up
```

## メモ
dbイメージ　postgres latest 

バックエンドイメージ　hackathon-backend latest

Docker Composeで作ったコンテナ、イメージ、ボリューム、ネットワークを一括削除：
docker compose down -v --rmi local

postmanで確認する際はrouter.goの26行目をコメントアウトして27行目のコメントアウトを外してから実行

missing csrf ~~ というエラーが出たらX-CSRF-TOKENを設定する

cuisines/ではフォームから入力する
