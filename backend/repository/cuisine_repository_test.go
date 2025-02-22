package repository

import (
	"testing"
	"time"

	"backend/model"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetAllCuisines(t *testing.T) {
	db := SetupTestDB()
	defer CleanupTestDB(db)

	repo := NewCuisineRepository(db)
	user := CreateTestUser(db)

	// テストデータを作成
	cuisines := []model.Cuisine{
		{
			Title:     "Test Cuisine 1",
			URL:       "https://example.com/1",
			UserId:    user.ID,
			CreatedAt: time.Now(),
		},
		{
			Title:     "Test Cuisine 2",
			URL:       "https://example.com/2",
			UserId:    user.ID,
			CreatedAt: time.Now().Add(time.Hour),
		},
	}

	for _, c := range cuisines {
		assert.NoError(t, repo.CreateCuisine(&c))
	}

	var fetchedCuisines []model.Cuisine
	err := repo.GetAllCuisines(&fetchedCuisines, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(fetchedCuisines))
	assert.Equal(t, "Test Cuisine 1", fetchedCuisines[0].Title)
	assert.Equal(t, "Test Cuisine 2", fetchedCuisines[1].Title)
}

func TestGetCuisineById(t *testing.T) {
	db := SetupTestDB()
	defer CleanupTestDB(db)

	repo := NewCuisineRepository(db)
	user := CreateTestUser(db)

	// テスト料理を作成
	cuisine := model.Cuisine{
		Title:  "Test Cuisine",
		URL:    "https://example.com",
		UserId: user.ID,
	}
	assert.NoError(t, repo.CreateCuisine(&cuisine))

	testCases := []struct {
		name    string
		userId  uint
		wantErr bool
	}{
		{
			name:    "正常なケース",
			userId:  user.ID,
			wantErr: false,
		},
		{
			name:    "存在しないユーザーID",
			userId:  user.ID + 1,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var fetchedCuisine model.Cuisine
			err := repo.GetCuisineById(&fetchedCuisine, tc.userId, cuisine.ID)
			if tc.wantErr {
				assert.Error(t, err, gorm.ErrRecordNotFound)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, cuisine.Title, fetchedCuisine.Title)
				assert.Equal(t, cuisine.URL, fetchedCuisine.URL)
				assert.Equal(t, cuisine.UserId, fetchedCuisine.UserId)
			}
		})
	}
}

func TestCreateCuisine(t *testing.T) {
	db := SetupTestDB()
	defer CleanupTestDB(db)

	repo := NewCuisineRepository(db)
	user := CreateTestUser(db)

	testCases := []struct {
		name    string
		cuisine model.Cuisine
		wantErr bool
	}{
		{
			name: "正常なケース",
			cuisine: model.Cuisine{
				Title:  "New Cuisine",
				URL:    "https://example.com",
				UserId: user.ID,
			},
			wantErr: false,
		},
		{
			name: "タイトルなし",
			cuisine: model.Cuisine{
				URL:    "https://example.com",
				UserId: user.ID,
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.CreateCuisine(&tc.cuisine)
			if tc.wantErr {
				t.Logf("Error: %v", err)
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, tc.cuisine.ID)
			}
		})
	}
}

func TestDeleteCuisine(t *testing.T) {
	db := SetupTestDB()
	defer CleanupTestDB(db)

	repo := NewCuisineRepository(db)
	user := CreateTestUser(db)

	cuisine := model.Cuisine{
		Title:  "Test Cuisine",
		URL:    "https://example.com",
		UserId: user.ID,
	}
	assert.NoError(t, repo.CreateCuisine(&cuisine))

	testCases := []struct {
		name    string
		userId  uint
		wantErr bool
	}{
		{
			name:    "正常なケース",
			userId:  user.ID,
			wantErr: false,
		},
		{
			name:    "存在しないユーザーID",
			userId:  user.ID + 1,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.DeleteCuisine(tc.userId, cuisine.ID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// 削除を確認
				var count int64
				db.Model(&model.Cuisine{}).Where("id = ?", cuisine.ID).Count(&count)
				assert.Equal(t, int64(0), count)
			}
		})
	}
}

func TestSettingCuisine(t *testing.T) {
	db := SetupTestDB()
	defer CleanupTestDB(db)

	repo := NewCuisineRepository(db)
	user := CreateTestUser(db)

	// テスト料理を作成
	cuisine := model.Cuisine{
		Title:  "Test Cuisine",
		URL:    "https://example.com",
		UserId: user.ID,
	}
	assert.NoError(t, repo.CreateCuisine(&cuisine))

	iconURL := "https://example.com/icon.png"
	newURL := "https://example.com/new"

	testCases := []struct {
		name    string
		update  model.Cuisine
		wantErr bool
	}{
		{
			name: "アイコンURLの更新",
			update: model.Cuisine{
				ID:      cuisine.ID,
				UserId:  user.ID,
				IconUrl: &iconURL,
			},
			wantErr: false,
		},
		{
			name: "URLの更新",
			update: model.Cuisine{
				ID:     cuisine.ID,
				UserId: user.ID,
				URL:    newURL,
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.SettingCuisine(&tc.update)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// 更新を確認
				var updated model.Cuisine
				db.First(&updated, cuisine.ID)

				if tc.update.IconUrl != nil {
					assert.Equal(t, *tc.update.IconUrl, *updated.IconUrl)
				}
				if tc.update.URL != "" {
					assert.Equal(t, tc.update.URL, updated.URL)
				}
			}
		})
	}
}
