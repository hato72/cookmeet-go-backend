# バックエンドAPIサービス

## 開発環境のセットアップ

### 必要条件
- Go 1.22以上
- Docker
- Docker Compose
- PostgreSQL

### セットアップ手順

1. リポジトリのクローン
```bash
git clone <repository-url>
cd backend
```

2. 依存関係のインストール
```bash
go mod download
```

3. 環境変数の設定
```bash
cp .env.dev .env
```

4. Dockerコンテナの起動
```bash
docker-compose up -d
```

## テスト実行

### テスト環境のセットアップ

1. テスト用のデータベース作成
```bash
docker-compose up -d test-db
```

2. テスト用の環境変数設定
```bash
cp .env.test .env
```

3. テスト用のディレクトリ作成
```bash
mkdir -p user_images/icons
mkdir -p cuisine_images/cuisine_icons
```

### テストの実行

すべてのテストを実行:
```bash
go test -v ./...
```

特定のパッケージのテストを実行:
```bash
go test -v ./db/...
go test -v ./usecase/...
go test -v ./repository/...
go test -v ./controller/...
```

カバレッジレポートの生成:
```bash
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## テストを避ける
コミット(プッシュ)時のメッセージに [Notest] を含める


## CI/CD

GitHub Actionsを使用して以下の自動化を実施:

- プルリクエスト時のテスト実行
- メインブランチへのマージ時のテスト実行とデプロイ

### CI/CD環境変数

GitHub Secretsに以下の環境変数を設定する必要があります:

- `DOCKER_USERNAME`: DockerHubのユーザー名
- `DOCKER_PASSWORD`: DockerHubのアクセストークン

## プロジェクト構成

```
backend/
├── controller/     # HTTPリクエストハンドラー
├── db/            # データベース接続管理
├── model/         # データモデル
├── repository/    # データアクセス層
├── router/        # ルーティング設定
├── usecase/       # ビジネスロジック
├── validator/     # バリデーション
└── testutil/      # テストユーティリティ
```

## API エンドポイント

### ユーザー関連
- `POST /signup` - ユーザー登録
- `POST /login` - ログイン
- `PUT /users` - ユーザー情報更新

### 料理関連
- `GET /cuisines` - 料理一覧取得
- `GET /cuisines/:id` - 料理詳細取得
- `POST /cuisines` - 料理追加
- `PUT /cuisines/:id` - 料理更新
- `DELETE /cuisines/:id` - 料理削除