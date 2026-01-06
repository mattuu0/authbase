package models_test

import (
	testtool "auth/testing"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// テスト実行
	code := m.Run()

	// 最後に集計を表示
	testtool.PrintFinalSummary()

	os.Exit(code)
}