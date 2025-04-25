package controller

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"backend/model"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCuisineUsecase struct {
	mock.Mock
}

// 以下のメソッドは、mock.Mockを埋め込んでいるため、自動的にモック化される
// モック化したいメソッドをオーバーライド
func (m *mockCuisineUsecase) GetAllCuisines(userID uint) ([]model.CuisineResponse, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.CuisineResponse), args.Error(1)
}

func (m *mockCuisineUsecase) GetCuisineByID(userID uint, cuisineID uint) (model.CuisineResponse, error) {
	args := m.Called(userID, cuisineID)
	return args.Get(0).(model.CuisineResponse), args.Error(1)
}

func (m *mockCuisineUsecase) DeleteCuisine(userID uint, cuisineID uint) error {
	args := m.Called(userID, cuisineID)
	return args.Error(0)
}

// AddCuisineメソッドのシグネチャを変更
func (m *mockCuisineUsecase) AddCuisine(cuisine model.Cuisine, iconURL *string, url string, title string) (model.CuisineResponse, error) {
	args := m.Called(cuisine, iconURL, url, title)
	return args.Get(0).(model.CuisineResponse), args.Error(1)
}

// SetCuisineメソッドも修正が必要
func (m *mockCuisineUsecase) SetCuisine(cuisine model.Cuisine, iconFile *multipart.FileHeader, url string, title string, userID uint, cuisineID uint) (model.CuisineResponse, error) {
	args := m.Called(cuisine, iconFile, url, title, userID, cuisineID)
	return args.Get(0).(model.CuisineResponse), args.Error(1)
}

// Echo のコンテキストとモックユースケース、そしてテスト対象の Cuisine Controller を初期化
func setupCuisineTest(_ *testing.T) (*echo.Echo, *mockCuisineUsecase, ICuisineController) {
	e := echo.New()
	mockUsecase := new(mockCuisineUsecase)
	controller := NewCuisineController(mockUsecase)
	return e, mockUsecase, controller
}

// JWT トークンを生成
func createJWTToken(userID float64) *jwt.Token {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	return token
}

func TestGetAllCuisines(t *testing.T) {
	e, mockUsecase, controller := setupCuisineTest(t)

	testCases := []struct {
		name         string
		userID       float64
		mockResponse []model.CuisineResponse
		mockError    error
		expectStatus int
	}{
		{
			name:   "正常な取得",
			userID: 1,
			mockResponse: []model.CuisineResponse{
				{
					ID:        1,
					Title:     "Test Cuisine 1",
					URL:       "https://example.com/1",
					UserID:    1,
					CreatedAt: time.Now(),
				},
				{
					ID:        2,
					Title:     "Test Cuisine 2",
					URL:       "https://example.com/2",
					UserID:    1,
					CreatedAt: time.Now(),
				},
			},
			mockError:    nil,
			expectStatus: http.StatusOK,
		},
		{
			name:         "データなし",
			userID:       2,
			mockResponse: []model.CuisineResponse{},
			mockError:    nil,
			expectStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/cuisines", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("user", createJWTToken(tc.userID))

			mockUsecase.On("GetAllCuisines", uint(tc.userID)).Return(tc.mockResponse, tc.mockError)

			err := controller.GetAllCuisines(c) // テスト対象のメソッドを実行
			assert.NoError(t, err)
			assert.Equal(t, tc.expectStatus, rec.Code) // レスポンスのステータスコードが期待通りか確認

			if tc.expectStatus == http.StatusOK { //モックが期待通りの結果を返す場合
				var response []model.CuisineResponse
				if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				assert.Equal(t, len(tc.mockResponse), len(response))
			}

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestGetCuisineByID(t *testing.T) {
	e, mockUsecase, controller := setupCuisineTest(t)

	testCases := []struct {
		name         string
		userID       float64
		cuisineID    string
		mockResponse model.CuisineResponse
		mockError    error
		expectStatus int
	}{
		{
			name:      "正常な取得",
			userID:    1,
			cuisineID: "1",
			mockResponse: model.CuisineResponse{
				ID:        1,
				Title:     "Test Cuisine",
				URL:       "https://example.com",
				UserID:    1,
				CreatedAt: time.Now(),
			},
			mockError:    nil,
			expectStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/cuisines/:id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("cuisineID")
			c.SetParamValues(tc.cuisineID)
			c.Set("user", createJWTToken(tc.userID))

			mockUsecase.On("GetCuisineByID", uint(tc.userID), uint(1)).Return(tc.mockResponse, tc.mockError)

			err := controller.GetCuisineByID(c)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectStatus, rec.Code)

			if tc.expectStatus == http.StatusOK {
				var response model.CuisineResponse
				if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				assert.Equal(t, tc.mockResponse.ID, response.ID)
				assert.Equal(t, tc.mockResponse.Title, response.Title)
			}

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestDeleteCuisine(t *testing.T) {
	e, mockUsecase, controller := setupCuisineTest(t)

	testCases := []struct {
		name         string
		userID       float64
		cuisineID    string
		mockError    error
		expectStatus int
	}{
		{
			name:         "正常な削除",
			userID:       1,
			cuisineID:    "1",
			mockError:    nil,
			expectStatus: http.StatusNoContent,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/cuisines/:id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("cuisineID")
			c.SetParamValues(tc.cuisineID)
			c.Set("user", createJWTToken(tc.userID))

			mockUsecase.On("DeleteCuisine", uint(tc.userID), uint(1)).Return(tc.mockError)

			err := controller.DeleteCuisine(c)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectStatus, rec.Code)

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestAddCuisine(t *testing.T) {
	e, mockUsecase, controller := setupCuisineTest(t)

	testCases := []struct {
		name         string
		userID       float64
		title        string
		url          string
		mockResponse model.CuisineResponse
		mockError    error
		expectStatus int
	}{
		{
			name:   "正常な追加",
			userID: 1,
			title:  "New Cuisine",
			url:    "https://example.com/new",
			mockResponse: model.CuisineResponse{
				ID:        1,
				Title:     "New Cuisine",
				URL:       "https://example.com/new",
				UserID:    1,
				CreatedAt: time.Now(),
			},
			mockError:    nil,
			expectStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			if err := writer.WriteField("title", tc.title); err != nil {
				t.Fatalf("Failed to write title field: %v", err)
			}
			if err := writer.WriteField("url", tc.url); err != nil {
				t.Fatalf("Failed to write url field: %v", err)
			}
			if err := writer.Close(); err != nil {
				t.Fatalf("Failed to close writer: %v", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/cuisines", body)
			req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("user", createJWTToken(tc.userID))

			// モックの設定を修正: 型チェックのみではなく、任意の値を受け入れるように変更
			mockUsecase.On("AddCuisine",
				mock.AnythingOfType("model.Cuisine"),
				mock.AnythingOfType("*string"), // nilであるかどうかにかかわらず任意の*string型を受け入れる
				tc.url,
				tc.title,
			).Return(tc.mockResponse, tc.mockError)

			err := controller.AddCuisine(c)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectStatus, rec.Code)

			if tc.expectStatus == http.StatusOK {
				var response model.CuisineResponse
				if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				assert.Equal(t, tc.mockResponse.Title, response.Title)
				assert.Equal(t, tc.mockResponse.URL, response.URL)
			}

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestSetCuisine(t *testing.T) {
	e, mockUsecase, controller := setupCuisineTest(t)

	testCases := []struct {
		name         string
		userID       float64
		cuisineID    string
		title        string
		url          string
		mockGetRes   model.CuisineResponse
		mockSetRes   model.CuisineResponse
		mockError    error
		expectStatus int
	}{
		{
			name:      "正常な更新",
			userID:    1,
			cuisineID: "1",
			title:     "Updated Cuisine",
			url:       "https://example.com/updated",
			mockGetRes: model.CuisineResponse{
				ID:        1,
				Title:     "Original Cuisine",
				URL:       "https://example.com/original",
				UserID:    1,
				CreatedAt: time.Now(),
			},
			mockSetRes: model.CuisineResponse{
				ID:        1,
				Title:     "Updated Cuisine",
				URL:       "https://example.com/updated",
				UserID:    1,
				CreatedAt: time.Now(),
			},
			mockError:    nil,
			expectStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			if err := writer.WriteField("title", tc.title); err != nil {
				t.Fatalf("Failed to write title field: %v", err)
			}
			if err := writer.WriteField("url", tc.url); err != nil {
				t.Fatalf("Failed to write url field: %v", err)
			}
			if err := writer.Close(); err != nil {
				t.Fatalf("Failed to close writer: %v", err)
			}

			req := httptest.NewRequest(http.MethodPut, "/cuisines/:id", body)
			req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("cuisineID")
			c.SetParamValues(tc.cuisineID)
			c.Set("user", createJWTToken(tc.userID))

			mockUsecase.On("GetCuisineByID", uint(tc.userID), uint(1)).Return(tc.mockGetRes, nil)
			mockUsecase.On("SetCuisine",
				mock.AnythingOfType("model.Cuisine"),
				(*multipart.FileHeader)(nil),
				tc.url,
				tc.title,
				uint(tc.userID),
				uint(1),
			).Return(tc.mockSetRes, tc.mockError)

			err := controller.SetCuisine(c)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectStatus, rec.Code)

			if tc.expectStatus == http.StatusOK {
				var response model.CuisineResponse
				if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				assert.Equal(t, tc.mockSetRes.Title, response.Title)
				assert.Equal(t, tc.mockSetRes.URL, response.URL)
			}

			mockUsecase.AssertExpectations(t)
		})
	}
}
