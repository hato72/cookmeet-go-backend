package controller

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"backend/model"
	"backend/usecase"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserUsecase struct {
	mock.Mock
}

func (m *mockUserUsecase) SignUp(user model.User) (model.UserResponse, error) {
	args := m.Called(user)
	return args.Get(0).(model.UserResponse), args.Error(1)
}

func (m *mockUserUsecase) Login(user model.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *mockUserUsecase) Update(user model.User, newEmail string, newName string, newPassword string, iconFile *multipart.FileHeader) (model.UserResponse, error) {
	args := m.Called(user, newEmail, newName, newPassword, iconFile)
	return args.Get(0).(model.UserResponse), args.Error(1)
}

func TestSignUp(t *testing.T) {
	// Echoのインスタンスを作成
	e := echo.New()

	testCases := []struct {
		name         string
		inputJSON    string
		mockResponse model.UserResponse
		mockError    error
		expectStatus int
	}{
		{
			name:      "正常なサインアップ",
			inputJSON: `{"name":"Test User","email":"test@example.com","password":"password123"}`,
			mockResponse: model.UserResponse{
				ID:    1,
				Name:  "Test User",
				Email: "test@example.com",
			},
			mockError:    nil,
			expectStatus: http.StatusCreated,
		},
		{
			name:         "無効なJSONリクエスト",
			inputJSON:    `{"invalid_json":`,
			mockResponse: model.UserResponse{},
			mockError:    nil,
			expectStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// モックを設定
			mockUsecase := new(mockUserUsecase)
			controller := NewUserController(mockUsecase)

			// リクエストを作成
			req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBufferString(tc.inputJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// モックの期待値を設定
			if tc.expectStatus == http.StatusCreated {
				var user model.User
				if err := json.Unmarshal([]byte(tc.inputJSON), &user); err != nil {
					t.Fatalf("Failed to unmarshal test input: %v", err)
				}
				mockUsecase.On("SignUp", user).Return(tc.mockResponse, tc.mockError)
			}

			// テスト対象の関数を実行
			err := controller.SignUp(c)

			// アサーション
			if tc.expectStatus != http.StatusBadRequest {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expectStatus, rec.Code)

			if tc.expectStatus == http.StatusCreated {
				var response model.UserResponse
				if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				assert.Equal(t, tc.mockResponse.ID, response.ID)
				assert.Equal(t, tc.mockResponse.Name, response.Name)
				assert.Equal(t, tc.mockResponse.Email, response.Email)
			}

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestLogin(t *testing.T) {
	e := echo.New()
	os.Setenv("API_DOMAIN", "localhost")

	testCases := []struct {
		name         string
		inputJSON    string
		mockToken    string
		mockError    error
		expectStatus int
	}{
		{
			name:         "正常なログイン",
			inputJSON:    `{"name":"Test User","email":"test@example.com","password":"password123"}`,
			mockToken:    "valid.jwt.token",
			mockError:    nil,
			expectStatus: http.StatusOK,
		},
		{
			name:         "無効な認証情報",
			inputJSON:    `{"name":"Test User","email":"test@example.com","password":"wrongpassword"}`,
			mockToken:    "",
			mockError:    usecase.ErrUserNotFound,
			expectStatus: http.StatusNotFound,
		},
		{
			name:         "パスワードが間違っている",
			inputJSON:    `{"name":"Test User","email":"test@example.com","password":"wrongpassword"}`,
			mockToken:    "",
			mockError:    usecase.ErrInvalidPassword, // usecase で定義しているパスワード不一致エラー
			expectStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUsecase := new(mockUserUsecase)
			controller := NewUserController(mockUsecase)

			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(tc.inputJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			var user model.User
			if err := json.Unmarshal([]byte(tc.inputJSON), &user); err != nil {
				t.Fatalf("Failed to unmarshal test input: %v", err)
			}
			// if tc.expectStatus == http.StatusOK {
			// 	var user model.User
			// 	json.Unmarshal([]byte(tc.inputJSON), &user)
			// 	mockUsecase.On("Login", user).Return(tc.mockToken, tc.mockError)
			// }

			err := controller.Login(c)

			// ログイン成功時のみトークンを確認
			if tc.expectStatus == http.StatusOK {
				assert.NoError(t, err)
				cookies := rec.Result().Cookies()
				assert.Equal(t, 1, len(cookies))
				assert.Equal(t, "token", cookies[0].Name)
				assert.Equal(t, tc.mockToken, cookies[0].Value)
			}

			assert.Equal(t, tc.expectStatus, rec.Code)
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestLogout(t *testing.T) {
	e := echo.New()
	os.Setenv("API_DOMAIN", "localhost")

	t.Run("ログアウト処理", func(t *testing.T) {
		mockUsecase := new(mockUserUsecase)
		controller := NewUserController(mockUsecase)

		req := httptest.NewRequest(http.MethodPost, "/logout", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := controller.Logout(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		cookies := rec.Result().Cookies()
		assert.Equal(t, 1, len(cookies))
		assert.Equal(t, "token", cookies[0].Name)
		assert.Equal(t, "", cookies[0].Value)
		assert.True(t, cookies[0].Expires.Before(time.Now())) //現在時刻を指定
	})
}

func TestUpdate(t *testing.T) {
	e := echo.New()

	// テスト用のJWTトークンを作成
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = float64(1)

	testCases := []struct {
		name         string
		setupRequest func() (*http.Request, *httptest.ResponseRecorder)
		mockSetup    func(*mockUserUsecase)
		expectStatus int
	}{
		{
			name: "プロフィール更新（アイコンなし）",
			setupRequest: func() (*http.Request, *httptest.ResponseRecorder) {
				body := new(bytes.Buffer)
				writer := multipart.NewWriter(body)

				// エラーチェックを追加
				if err := writer.WriteField("name", "Updated Name"); err != nil {
					t.Fatalf("failed to write name field: %v", err)
				}
				if err := writer.WriteField("email", "new@example.com"); err != nil {
					t.Fatalf("failed to write email field: %v", err)
				}

				// Close()のエラーもチェック
				if err := writer.Close(); err != nil {
					t.Fatalf("failed to close writer: %v", err)
				}

				req := httptest.NewRequest(http.MethodPut, "/users/update", body)
				req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
				return req, httptest.NewRecorder()
			},
			mockSetup: func(m *mockUserUsecase) {
				m.On("Update",
					//mock.AnythingOfType("model.User"),
					mock.MatchedBy(func(user model.User) bool {
						return user.ID == 1
					}),
					"new@example.com",
					"Updated Name",
					"",
					mock.Anything,
				).Return(model.UserResponse{
					ID:    1,
					Name:  "Updated Name",
					Email: "new@example.com",
				}, nil)
			},
			expectStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUsecase := new(mockUserUsecase)
			controller := NewUserController(mockUsecase)

			req, rec := tc.setupRequest()
			c := e.NewContext(req, rec)
			c.Set("user", token)

			tc.mockSetup(mockUsecase)

			err := controller.Update(c)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectStatus, rec.Code)

			if tc.expectStatus == http.StatusOK {
				var response model.UserResponse
				if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
					t.Fatalf("レスポンスのUnmarshalに失敗: %v", err)
				}
				assert.Equal(t, uint(1), response.ID)
			}

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestCsrfToken(t *testing.T) {
	e := echo.New()

	t.Run("CSRFトークン取得", func(t *testing.T) {
		mockUsecase := new(mockUserUsecase)
		controller := NewUserController(mockUsecase)

		req := httptest.NewRequest(http.MethodGet, "/csrf", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("csrf", "test-csrf-token")

		err := controller.CsrfToken(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]string
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Fatalf("CSRFトークンレスポンスのUnmarshalに失敗: %v", err)
		}
		assert.Equal(t, "test-csrf-token", response["csrf_token"])
	})
}
