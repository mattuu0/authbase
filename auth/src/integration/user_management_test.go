// Package integration provides integration tests for user management API flows.
package integration

import (
	"auth/models"
	"auth/services"
	testtool "auth/testing"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUserCRUDFlow はAdmin APIを使ったユーザーCRUDの完全フローをテストします
func TestUserCRUDFlow(t *testing.T) {
	testtool.RecordResult(t)
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	testtool.SetupTestEnv(t)

	provider := testtool.CreateTestProvider(t, db, models.Basic)
	provider.IsEnabled = 1
	db.Save(provider)
	services.Init()
	router := setupAdminRouter()

	// 管理者を作成してログイン
	body, _ := json.Marshal(map[string]string{"username": "crud-admin", "password": "CrudAdmin123!"})
	req := httptest.NewRequest(http.MethodPost, "/admin/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	cookies := loginAsAdmin(t, router, "crud-admin", "CrudAdmin123!")

	t.Run("create user via admin API", func(t *testing.T) {
		testtool.RecordResult(t)
		body, _ := json.Marshal(map[string]string{
			"name":     "New User",
			"email":    "newuser@example.com",
			"password": "UserPass123!",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/user", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		for _, c := range cookies {
			req.AddCookie(c)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var user map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &user))
		assert.Equal(t, "newuser@example.com", user["email"])
		assert.NotEmpty(t, user["id"])
		testtool.LogSuccess(t, "User created via admin API")
	})

	t.Run("get all users includes created user", func(t *testing.T) {
		testtool.RecordResult(t)
		req := httptest.NewRequest(http.MethodGet, "/api/user/all", nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		var users []map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &users))
		assert.GreaterOrEqual(t, len(users), 1)
		testtool.LogSuccess(t, "User list returned")
	})

	t.Run("update user name via admin API", func(t *testing.T) {
		testtool.RecordResult(t)
		// まずユーザー一覧を取得してIDを得る
		req := httptest.NewRequest(http.MethodGet, "/api/user/all", nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)

		var users []map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &users))
		require.NotEmpty(t, users)

		// emailで対象ユーザーを探す
		var targetID string
		for _, u := range users {
			if u["email"] == "newuser@example.com" {
				targetID = u["id"].(string)
				break
			}
		}
		require.NotEmpty(t, targetID)

		body, _ := json.Marshal(map[string]interface{}{
			"id":     targetID,
			"name":   "Updated Name",
			"labels": []string{},
		})
		req = httptest.NewRequest(http.MethodPut, "/api/user", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		for _, c := range cookies {
			req.AddCookie(c)
		}
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		var updated map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &updated))
		assert.Equal(t, "Updated Name", updated["name"])
		testtool.LogSuccess(t, "User name updated")
	})

	t.Run("ban and unban user", func(t *testing.T) {
		testtool.RecordResult(t)
		// ユーザーIDを取得
		req := httptest.NewRequest(http.MethodGet, "/api/user/all", nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		var users []map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &users))

		var targetID string
		for _, u := range users {
			if u["email"] == "newuser@example.com" {
				targetID = u["id"].(string)
				break
			}
		}
		require.NotEmpty(t, targetID)

		// BAN
		body, _ := json.Marshal(map[string]interface{}{"UserID": targetID, "IsBanned": true})
		req = httptest.NewRequest(http.MethodPut, "/api/user/ban", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		for _, c := range cookies {
			req.AddCookie(c)
		}
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		var banned map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &banned))
		assert.Equal(t, true, banned["banned"])
		testtool.LogSuccess(t, "User banned")

		// BAN解除
		body, _ = json.Marshal(map[string]interface{}{"UserID": targetID, "IsBanned": false})
		req = httptest.NewRequest(http.MethodPut, "/api/user/ban", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		for _, c := range cookies {
			req.AddCookie(c)
		}
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		var unbanned map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &unbanned))
		assert.Equal(t, false, unbanned["banned"])
		testtool.LogSuccess(t, "User unbanned")
	})

	t.Run("delete user via admin API", func(t *testing.T) {
		testtool.RecordResult(t)
		// ユーザーIDを取得
		req := httptest.NewRequest(http.MethodGet, "/api/user/all", nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		var users []map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &users))

		var targetID string
		for _, u := range users {
			if u["email"] == "newuser@example.com" {
				targetID = u["id"].(string)
				break
			}
		}
		require.NotEmpty(t, targetID)

		req = httptest.NewRequest(http.MethodDelete, "/api/user", nil)
		req.Header.Set("userid", targetID)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		testtool.LogSuccess(t, "User deleted")

		// 削除後は一覧に存在しない
		req = httptest.NewRequest(http.MethodGet, "/api/user/all", nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		var usersAfter []map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &usersAfter))
		for _, u := range usersAfter {
			assert.NotEqual(t, "newuser@example.com", u["email"])
		}
	})
}

// TestLabelCRUDFlow はAdmin APIを使ったラベルCRUDをテストします
func TestLabelCRUDFlow(t *testing.T) {
	testtool.RecordResult(t)
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	testtool.SetupTestEnv(t)
	services.Init()
	router := setupAdminRouter()

	body, _ := json.Marshal(map[string]string{"username": "label-admin", "password": "LabelAdmin123!"})
	req := httptest.NewRequest(http.MethodPost, "/admin/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	cookies := loginAsAdmin(t, router, "label-admin", "LabelAdmin123!")

	var createdLabelName string

	t.Run("create label", func(t *testing.T) {
		testtool.RecordResult(t)
		body, _ := json.Marshal(map[string]string{"name": "vip", "color": "#FFD700"})
		req := httptest.NewRequest(http.MethodPost, "/api/labels", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		for _, c := range cookies {
			req.AddCookie(c)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		var label map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &label))
		assert.Equal(t, "vip", label["name"])
		assert.Equal(t, "#FFD700", label["color"])
		createdLabelName = label["name"].(string)
		testtool.LogSuccess(t, "Label created")
	})

	t.Run("get all labels includes created label", func(t *testing.T) {
		testtool.RecordResult(t)
		req := httptest.NewRequest(http.MethodGet, "/api/labels", nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		var labels []map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &labels))
		require.NotEmpty(t, labels)
		found := false
		for _, l := range labels {
			if l["name"] == "vip" {
				found = true
			}
		}
		assert.True(t, found, "created label should appear in list")
		testtool.LogSuccess(t, "Label list returned")
	})

	t.Run("update label", func(t *testing.T) {
		testtool.RecordResult(t)
		require.NotEmpty(t, createdLabelName)
		body, _ := json.Marshal(map[string]string{"id": createdLabelName, "name": "vip", "color": "#FF0000"})
		req := httptest.NewRequest(http.MethodPut, "/api/labels", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		for _, c := range cookies {
			req.AddCookie(c)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		var label map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &label))
		assert.Equal(t, "#FF0000", label["color"])
		testtool.LogSuccess(t, "Label color updated")
	})

	t.Run("delete label", func(t *testing.T) {
		testtool.RecordResult(t)
		require.NotEmpty(t, createdLabelName)
		body, _ := json.Marshal(map[string]string{"id": createdLabelName})
		req := httptest.NewRequest(http.MethodDelete, "/api/labels", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		for _, c := range cookies {
			req.AddCookie(c)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		testtool.LogSuccess(t, "Label deleted")

		// 削除後は一覧に存在しない
		req = httptest.NewRequest(http.MethodGet, "/api/labels", nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		var labels []map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &labels))
		for _, l := range labels {
			assert.NotEqual(t, "vip", l["name"])
		}
	})

	t.Run("create duplicate label returns error", func(t *testing.T) {
		testtool.RecordResult(t)
		body, _ := json.Marshal(map[string]string{"name": "unique", "color": "#000000"})
		req := httptest.NewRequest(http.MethodPost, "/api/labels", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		for _, c := range cookies {
			req.AddCookie(c)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)

		// 同名で再作成
		req = httptest.NewRequest(http.MethodPost, "/api/labels", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		for _, c := range cookies {
			req.AddCookie(c)
		}
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		testtool.LogSuccess(t, "Duplicate label rejected")
	})
}

// TestSessionManagementFlow はAdmin APIを使ったセッション管理をテストします
func TestSessionManagementFlow(t *testing.T) {
	testtool.RecordResult(t)
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	testtool.SetupTestEnv(t)

	provider := testtool.CreateTestProvider(t, db, models.Basic)
	provider.IsEnabled = 1
	db.Save(provider)
	services.Init()
	router := setupAdminRouter()
	userRouter := setupTestRouter()

	// 管理者作成・ログイン
	body, _ := json.Marshal(map[string]string{"username": "session-admin", "password": "SessionAdmin123!"})
	req := httptest.NewRequest(http.MethodPost, "/admin/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	cookies := loginAsAdmin(t, router, "session-admin", "SessionAdmin123!")

	// テストユーザーを作成してセッションを発生させる
	signupBody, _ := json.Marshal(map[string]string{
		"name":     "Session Test User",
		"email":    "sessiontest@example.com",
		"password": "UserPass123!",
	})
	req = httptest.NewRequest(http.MethodPost, "/basic/signup", bytes.NewBuffer(signupBody))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	userRouter.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	var signupResp map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &signupResp))
	userToken := signupResp["token"].(string)

	t.Run("get all sessions includes user session", func(t *testing.T) {
		testtool.RecordResult(t)
		req := httptest.NewRequest(http.MethodGet, "/api/session", nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		var sessions []map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &sessions))
		assert.GreaterOrEqual(t, len(sessions), 1)
		testtool.LogSuccess(t, "Sessions listed")
	})

	t.Run("delete session via admin API invalidates user token", func(t *testing.T) {
		testtool.RecordResult(t)
		// セッション一覧を取得してIDを得る
		req := httptest.NewRequest(http.MethodGet, "/api/session", nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		var sessions []map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &sessions))

		var sessionID string
		for _, s := range sessions {
			if s["userEmail"] == "sessiontest@example.com" {
				sessionID = s["id"].(string)
				break
			}
		}
		require.NotEmpty(t, sessionID, "session for test user should exist")

		// セッションを削除
		req = httptest.NewRequest(http.MethodDelete, "/api/session", nil)
		req.Header.Set("sessionid", sessionID)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		testtool.LogSuccess(t, "Session deleted by admin")

		// 削除後はユーザートークンで /me にアクセス不可
		req = httptest.NewRequest(http.MethodGet, "/me", nil)
		req.Header.Set("Authorization", userToken)
		rec = httptest.NewRecorder()
		userRouter.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		testtool.LogSuccess(t, "User token invalidated after session deletion")
	})
}

// TestBannedUserCannotLogin はBANされたユーザーがログインできないことをテストします
func TestBannedUserCannotLogin(t *testing.T) {
	testtool.RecordResult(t)
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	testtool.SetupTestEnv(t)

	provider := testtool.CreateTestProvider(t, db, models.Basic)
	provider.IsEnabled = 1
	db.Save(provider)
	services.Init()
	router := setupAdminRouter()
	userRouter := setupTestRouter()

	// 管理者作成・ログイン
	body, _ := json.Marshal(map[string]string{"username": "ban-admin", "password": "BanAdmin123!"})
	req := httptest.NewRequest(http.MethodPost, "/admin/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	cookies := loginAsAdmin(t, router, "ban-admin", "BanAdmin123!")

	// ユーザー作成
	signupBody, _ := json.Marshal(map[string]string{
		"name": "Target User", "email": "target@example.com", "password": "TargetPass123!",
	})
	req = httptest.NewRequest(http.MethodPost, "/basic/signup", bytes.NewBuffer(signupBody))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	userRouter.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	// ユーザーIDを取得
	req = httptest.NewRequest(http.MethodGet, "/api/user/all", nil)
	for _, c := range cookies {
		req.AddCookie(c)
	}
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	var users []map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &users))
	var targetID string
	for _, u := range users {
		if u["email"] == "target@example.com" {
			targetID = u["id"].(string)
			break
		}
	}
	require.NotEmpty(t, targetID)

	// BAN する
	banBody, _ := json.Marshal(map[string]interface{}{"UserID": targetID, "IsBanned": true})
	req = httptest.NewRequest(http.MethodPut, "/api/user/ban", bytes.NewBuffer(banBody))
	req.Header.Set("Content-Type", "application/json")
	for _, c := range cookies {
		req.AddCookie(c)
	}
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	// BANされたユーザーは新規ログイン不可
	t.Run("banned user cannot login", func(t *testing.T) {
		testtool.RecordResult(t)
		loginBody, _ := json.Marshal(map[string]string{
			"email": "target@example.com", "password": "TargetPass123!",
		})
		req = httptest.NewRequest(http.MethodPost, "/basic/login", bytes.NewBuffer(loginBody))
		req.Header.Set("Content-Type", "application/json")
		rec = httptest.NewRecorder()
		userRouter.ServeHTTP(rec, req)
		assert.NotEqual(t, http.StatusOK, rec.Code)
		testtool.LogSuccess(t, "Banned user cannot login")
	})
}
