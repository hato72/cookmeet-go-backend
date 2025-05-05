package usecase

import (
	"backend/model"
	"backend/validator"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockCuisineRepository はCuisineRepositoryのモック
type MockCuisineRepository struct {
	mock.Mock
}

type MockCuisineValidator struct {
	mock.Mock
}

func (m *MockCuisineRepository) GetAllCuisines(cuisines *[]model.Cuisine, userID uint) error {
	args := m.Called(cuisines, userID)
	if args.Get(0) != nil {
		*cuisines = args.Get(0).([]model.Cuisine)
	}
	return args.Error(1)
}

func (m *MockCuisineRepository) GetCuisineByID(cuisine *model.Cuisine, userID uint, cuisineID uint) error {
	args := m.Called(cuisine, userID, cuisineID)
	return args.Error(0)
}

func (m *MockCuisineRepository) CreateCuisine(cuisine *model.Cuisine) error {
	args := m.Called(cuisine)
	return args.Error(0)
}

func (m *MockCuisineRepository) DeleteCuisine(userID uint, cuisineID uint) error {
	args := m.Called(userID, cuisineID)
	return args.Error(0)
}

func (m *MockCuisineRepository) SettingCuisine(cuisine *model.Cuisine) error {
	args := m.Called(cuisine)
	return args.Error(0)
}

func (m *MockCuisineValidator) CuisineValidate(cuisine model.Cuisine) error {
	args := m.Called(cuisine)
	return args.Error(0)
}

func TestGetAllCuisines(t *testing.T) {
	// モックの準備
	mockRepo := new(MockCuisineRepository)
	validator := validator.NewCuisineValidator()
	usecase := NewCuisineUsecase(mockRepo, validator)

	UserID := uint(1)
	now := time.Now()
	mockCuisines := []model.Cuisine{
		{
			ID:        1,
			Title:     "Test Cuisine 1",
			URL:       "http://example.com/1",
			CreatedAt: now,
			UpdatedAt: now,
			UserID:    UserID,
		},
		{
			ID:        2,
			Title:     "Test Cuisine 2",
			URL:       "http://example.com/2",
			CreatedAt: now,
			UpdatedAt: now,
			UserID:    UserID,
		},
	}

	// モックの振る舞いを設定
	mockRepo.On("GetAllCuisines", mock.AnythingOfType("*[]model.Cuisine"), UserID).
		Run(func(args mock.Arguments) {
			cuisines := args.Get(0).(*[]model.Cuisine)
			*cuisines = mockCuisines
		}).
		Return(mockCuisines, nil)

	// テスト実行
	cuisines, err := usecase.GetAllCuisines(UserID)

	// アサーション
	assert.NoError(t, err)
	assert.Len(t, cuisines, 2)
	assert.Equal(t, mockCuisines[0].Title, cuisines[0].Title)
	assert.Equal(t, mockCuisines[0].URL, cuisines[0].URL)
	mockRepo.AssertExpectations(t)
}

func TestGetCuisineByID(t *testing.T) {
	// モックの準備
	mockRepo := new(MockCuisineRepository)
	validator := validator.NewCuisineValidator()
	usecase := NewCuisineUsecase(mockRepo, validator)

	UserID := uint(1)
	cuisineID := uint(1)
	now := time.Now()
	mockCuisine := model.Cuisine{
		ID:        cuisineID,
		Title:     "Test Cuisine",
		URL:       "http://example.com",
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    UserID,
	}

	// モックの振る舞いを設定
	mockRepo.On("GetCuisineByID", mock.AnythingOfType("*model.Cuisine"), UserID, cuisineID).
		Run(func(args mock.Arguments) {
			cuisine := args.Get(0).(*model.Cuisine)
			*cuisine = mockCuisine
		}).Return(nil) // エラーの型を修正

	// テスト実行
	cuisine, err := usecase.GetCuisineByID(UserID, cuisineID)

	// アサーション
	assert.NoError(t, err)
	assert.Equal(t, mockCuisine.Title, cuisine.Title)
	assert.Equal(t, mockCuisine.URL, cuisine.URL)
	mockRepo.AssertExpectations(t)
}

func TestDeleteCuisine(t *testing.T) {
	mockRepo := new(MockCuisineRepository)
	mockValidator := new(MockCuisineValidator)
	cu := NewCuisineUsecase(mockRepo, mockValidator)

	tests := []struct {
		name      string
		userID    uint
		cuisineID uint
		mockSetup func()
		wantErr   error
	}{
		{
			name:      "正常に削除できる場合",
			userID:    1,
			cuisineID: 1,
			mockSetup: func() {
				mockRepo.On("GetCuisineByID", mock.AnythingOfType("*model.Cuisine"), uint(1), uint(1)).
					Run(func(args mock.Arguments) {
						cuisine := args.Get(0).(*model.Cuisine)
						cuisine.ID = 1
						cuisine.UserID = 1
					}).Return(nil)

				mockRepo.On("DeleteCuisine", uint(1), uint(1)).Return(nil)
			},
			wantErr: nil,
		},
		{
			name:      "料理が存在しない場合",
			userID:    1,
			cuisineID: 999,
			mockSetup: func() {
				mockRepo.On("GetCuisineByID", mock.AnythingOfType("*model.Cuisine"), uint(1), uint(999)).
					Return(gorm.ErrRecordNotFound)
			},
			wantErr: ErrCuisineNotFound,
		},
		{
			name:      "権限がない場合",
			userID:    2,
			cuisineID: 1,
			mockSetup: func() {
				mockRepo.On("GetCuisineByID", mock.AnythingOfType("*model.Cuisine"), uint(2), uint(1)).
					Run(func(args mock.Arguments) {
						cuisine := args.Get(0).(*model.Cuisine)
						cuisine.ID = 1
						cuisine.UserID = 1 // 別のユーザーのCuisine
					}).Return(nil)
			},
			wantErr: ErrUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックをリセット
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil

			// モックの設定
			tt.mockSetup()

			// テスト実行
			err := cu.DeleteCuisine(tt.userID, tt.cuisineID)

			// アサーション
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			// モックの呼び出しを検証
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAddCuisine(t *testing.T) {
	// モックの準備
	mockRepo := new(MockCuisineRepository)
	validator := validator.NewCuisineValidator()
	usecase := NewCuisineUsecase(mockRepo, validator)

	cuisine := model.Cuisine{
		Title:  "Test Cuisine",
		URL:    "http://example.com",
		UserID: 1,
	}

	// モックの振る舞いを設定
	mockRepo.On("CreateCuisine", mock.AnythingOfType("*model.Cuisine")).Return(nil)

	// テスト実行
	response, err := usecase.AddCuisine(cuisine, nil, cuisine.URL, cuisine.Title)

	// アサーション
	assert.NoError(t, err)
	assert.Equal(t, cuisine.Title, response.Title)
	assert.Equal(t, cuisine.URL, response.URL)
	mockRepo.AssertExpectations(t)
}
