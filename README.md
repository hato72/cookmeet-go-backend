# CookMeet(サーバー)
![cookmeet](https://github.com/hato72/go_backend_hackathon/assets/139688965/54235b01-2da0-491e-857c-18581b70b518)

## フロントエンド
ソースコード：https://github.com/hato72/CookMeet

デプロイ先：https://cook-meet.vercel.app/

## AI部分
ソースコード：https://github.com/hato72/cookmeet-recommend-recipes


## 実行方法(ローカル環境)

```sh
.env.dev:

PORT=8080
POSTGRES_USER=hato
POSTGRES_PW=hato72
POSTGRES_DB=hato
POSTGRES_PORT=5432
POSTGRES_HOST=db
SECRET=uu5pveql
GO_ENV=dev
API_DOMAIN=localhost
FE_URL=http://localhost:3000
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

## test
```
make test
```
or
```
make test-div test=
```

## lint

```
make lint
```

## メモ
dbイメージ　postgres latest 

バックエンドイメージ　hackathon-backend latest

Docker Composeで作ったコンテナ、イメージ、ボリューム、ネットワークを一括削除：
docker compose down -v --rmi local

postmanで確認する際はrouter.goの26行目をコメントアウトして27行目のコメントアウトを外してから実行

missing csrf ~~ というエラーが出たらX-CSRF-TOKENを設定する

cuisines/ではフォームから入力する
