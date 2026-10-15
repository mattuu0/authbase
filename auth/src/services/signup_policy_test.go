package services

import (
	"auth/models"
	testtool "auth/testing"
	"net/http"
	"testing"
)

// TestSignupAllowed は ALLOW_SIGNUP の値ごとに、登録の可否が期待どおりになることを確認する
func TestSignupAllowed(t *testing.T) {
	cases := []struct {
		name  string
		value string
		set   bool
		want  bool
	}{
		{"未設定は許可(従来どおり)", "", false, true},
		{"空文字は許可", "", true, true},
		{"true は許可", "true", true, true},
		{"false は拒否", "false", true, false},
		{"FALSE(大文字)は拒否", "FALSE", true, false},
		{"前後の空白があっても拒否", "  false ", true, false},
		{"0 は拒否", "0", true, false},
		{"no は拒否", "no", true, false},
		{"off は拒否", "off", true, false},
		{"1 は許可", "1", true, true},
		{"想定外の値は許可(誤設定で全員が登録不能にならない側に倒す)", "maybe", true, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.set {
				t.Setenv("ALLOW_SIGNUP", c.value)
			} else {
				t.Setenv("ALLOW_SIGNUP", "") // t.Setenv で、テスト後に元へ戻す
			}
			if got := SignupAllowed(); got != c.want {
				t.Errorf("ALLOW_SIGNUP=%q: 期待 %v, 実際 %v", c.value, c.want, got)
			}
		})
	}
}

// TestCreateBasicUser_ALLOW_SIGNUPがfalseなら403で拒否する は、登録が禁止されているとき、
// DB に触れる前に 403 を返し、ユーザーを作らないことを確認する(DB 未接続でも動く)
func TestCreateBasicUser_ALLOW_SIGNUPがfalseなら403で拒否する(t *testing.T) {
	t.Setenv("ALLOW_SIGNUP", "false")

	token, result := CreateBasicUser(CreateBasicUserArgs{Name: "x", Email: "x@example.com", Password: "SecurePass123!"})

	if result.Success || result.Code != http.StatusForbidden {
		t.Errorf("403 で拒否されることを期待しましたが、実際: success=%v code=%d", result.Success, result.Code)
	}
	if token != "" {
		t.Errorf("拒否されたのにトークンが返っています: %q", token)
	}
	if result.Error == nil || result.Error.Error() != "signup is disabled" {
		t.Errorf("エラーが signup is disabled ではありません: %v", result.Error)
	}
}

// TestCreateUser_ALLOW_SIGNUPがfalseでも管理者はユーザーを作成できる は、
// 公開の登録を止めても、管理画面(管理者認証つきの POST /api/user)からの作成は影響を受けないことを確認する
func TestCreateUser_ALLOW_SIGNUPがfalseでも管理者はユーザーを作成できる(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	testtool.SetupTestEnv(t)

	provider := testtool.CreateTestProvider(t, db, models.Basic)
	provider.IsEnabled = 1
	db.Save(provider)
	Init()

	t.Setenv("ALLOW_SIGNUP", "false")

	// 公開の登録は拒否される
	_, publicResult := CreateBasicUser(CreateBasicUserArgs{Name: "Public", Email: "public@example.com", Password: "SecurePass123!"})
	if publicResult.Success || publicResult.Code != http.StatusForbidden {
		t.Fatalf("公開の登録は 403 で拒否されることを期待しましたが、実際: success=%v code=%d", publicResult.Success, publicResult.Code)
	}

	// 管理者による作成は成功する
	createdUser, err := CreateUser(CreateBasicUserArgs{Name: "Created By Admin", Email: "admin-created@example.com", Password: "SecurePass123!"})
	if err != nil {
		t.Fatalf("管理者によるユーザー作成が失敗しました: %v", err)
	}
	if createdUser.Email != "admin-created@example.com" {
		t.Errorf("作成されたユーザーのメールが期待と異なります: %s", createdUser.Email)
	}

	// 拒否された公開の登録では、ユーザーが作られていない
	if _, lookup := models.GetUserByEmail("public@example.com"); lookup.IsExists {
		t.Error("拒否されたはずの公開の登録で、ユーザーが作られています")
	}
}
