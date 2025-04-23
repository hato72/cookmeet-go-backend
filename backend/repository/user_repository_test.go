package repository

import (
	"testing"

	"backend/model"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	db := SetupTestDB()
	defer CleanupTestDB(db)

	repo := NewUserRepository(db)

	testCases := []struct {
		name    string
		user    *model.User
		wantErr bool
	}{
		{
			name: "正常なユーザー作成",
			user: &model.User{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "重複するメールアドレス",
			user: &model.User{
				Name:     "Another User",
				Email:    "test@example.com", // 既に使用されているメールアドレス
				Password: "password456",
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.CreateUser(tc.user)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, tc.user.ID)
			}
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	db := SetupTestDB()
	defer CleanupTestDB(db)

	repo := NewUserRepository(db)

	// テストユーザーを作成
	testUser := &model.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}
	if err := repo.CreateUser(testUser); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	testCases := []struct {
		name     string
		email    string
		wantUser bool
		wantErr  bool
	}{
		{
			name:     "存在するユーザーの取得",
			email:    "test@example.com",
			wantUser: true,
			wantErr:  false,
		},
		{
			name:     "存在しないユーザーの取得",
			email:    "notfound@example.com",
			wantUser: false,
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var user model.User
			err := repo.GetUserByEmail(&user, tc.email)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.email, user.Email)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	db := SetupTestDB()
	defer CleanupTestDB(db)

	repo := NewUserRepository(db)

	// テストユーザーを作成
	testUser := &model.User{
		Name:     "Original Name",
		Email:    "original@example.com",
		Password: "password123",
	}
	if err := repo.CreateUser(testUser); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	testCases := []struct {
		name    string
		update  *model.User
		wantErr bool
	}{
		{
			name: "名前の更新",
			update: &model.User{
				ID:   testUser.ID,
				Name: "Updated Name",
			},
			wantErr: false,
		},
		{
			name: "メールアドレスの更新",
			update: &model.User{
				ID:    testUser.ID,
				Email: "updated@example.com",
			},
			wantErr: false,
		},
		{
			name: "パスワードの更新",
			update: &model.User{
				ID:       testUser.ID,
				Password: "newpassword123",
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.UpdateUser(tc.update)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// 更新されたユーザー情報を取得して確認
				var updatedUser model.User
				db.First(&updatedUser, testUser.ID)

				if tc.update.Name != "" {
					assert.Equal(t, tc.update.Name, updatedUser.Name)
				}
				if tc.update.Email != "" {
					assert.Equal(t, tc.update.Email, updatedUser.Email)
				}
				if tc.update.Password != "" {
					assert.Equal(t, tc.update.Password, updatedUser.Password)
				}
			}
		})
	}
}
