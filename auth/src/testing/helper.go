// Package testing provides enhanced test utilities with better logging and output.
package testing

import (
	"fmt"
	"testing"
	"time"
)

// TestLogger はテスト用の詳細ログを提供します
type TestLogger struct {
	t         *testing.T
	testName  string
	startTime time.Time
}

// NewTestLogger は新しいテストロガーを作成します
func NewTestLogger(t *testing.T, testName string) *TestLogger {
	t.Helper()
	logger := &TestLogger{
		t:         t,
		testName:  testName,
		startTime: time.Now(),
	}
	logger.LogStart()
	return logger
}

// LogStart はテスト開始をログ出力します
func (l *TestLogger) LogStart() {
	l.t.Helper()
	l.log("🚀", "START", "")
}

// LogStep はテストステップをログ出力します
func (l *TestLogger) LogStep(step string) {
	l.t.Helper()
	elapsed := time.Since(l.startTime).Round(time.Millisecond)
	l.log("📍", "STEP", fmt.Sprintf("%s [%v]", step, elapsed))
}

// LogSuccess は成功をログ出力します
func (l *TestLogger) LogSuccess(message string) {
	l.t.Helper()
	elapsed := time.Since(l.startTime).Round(time.Millisecond)
	l.log("✅", "PASS", fmt.Sprintf("%s [%v]", message, elapsed))
}

// LogWarning は警告をログ出力します
func (l *TestLogger) LogWarning(message string) {
	l.t.Helper()
	l.log("⚠️", "WARN", message)
}

// LogError はエラーをログ出力します
func (l *TestLogger) LogError(message string) {
	l.t.Helper()
	l.log("❌", "FAIL", message)
}

// LogInfo は情報をログ出力します
func (l *TestLogger) LogInfo(message string) {
	l.t.Helper()
	l.log("ℹ️", "INFO", message)
}

// log は内部ログ出力関数です
func (l *TestLogger) log(emoji, level, message string) {
	l.t.Helper()
	if message == "" {
		fmt.Printf("\n%s [%s] %s\n", emoji, level, l.testName)
	} else {
		fmt.Printf("%s [%s] %s: %s\n", emoji, level, l.testName, message)
	}
}

// Finish はテスト終了をログ出力します
func (l *TestLogger) Finish() {
	l.t.Helper()
	elapsed := time.Since(l.startTime).Round(time.Millisecond)
	l.log("🏁", "END", fmt.Sprintf("completed in %v", elapsed))
	fmt.Println(string([]rune{'-'}[0]) + "─────────────────────────────────────────────────────")
}

// AssertNoError はエラーがないことをアサートし、詳細ログを出力します
func (l *TestLogger) AssertNoError(err error, context string) {
	l.t.Helper()
	if err != nil {
		l.LogError(fmt.Sprintf("%s: %v", context, err))
		l.t.Fatalf("%s: %v", context, err)
	} else {
		l.LogSuccess(fmt.Sprintf("%s: no error", context))
	}
}

// AssertEqual は値が等しいことをアサートします
func (l *TestLogger) AssertEqual(expected, actual interface{}, context string) {
	l.t.Helper()
	if expected != actual {
		l.LogError(fmt.Sprintf("%s: expected %v, got %v", context, expected, actual))
		l.t.Fatalf("%s: expected %v, got %v", context, expected, actual)
	} else {
		l.LogSuccess(fmt.Sprintf("%s: values match", context))
	}
}

// AssertNotEmpty は値が空でないことをアサートします
func (l *TestLogger) AssertNotEmpty(value string, context string) {
	l.t.Helper()
	if value == "" {
		l.LogError(fmt.Sprintf("%s: value is empty", context))
		l.t.Fatalf("%s: value is empty", context)
	} else {
		l.LogSuccess(fmt.Sprintf("%s: value is not empty", context))
	}
}

// AssertTrue は条件が真であることをアサートします
func (l *TestLogger) AssertTrue(condition bool, context string) {
	l.t.Helper()
	if !condition {
		l.LogError(fmt.Sprintf("%s: condition is false", context))
		l.t.Fatalf("%s: condition is false", context)
	} else {
		l.LogSuccess(fmt.Sprintf("%s: condition is true", context))
	}
}

// PrintSection はセクションヘッダーを出力します
func PrintTestSection(name string) {
	fmt.Println("\n" + string([]rune{'═'}[0]) + "═══════════════════════════════════════════════════")
	fmt.Printf("  %s\n", name)
	fmt.Println(string([]rune{'═'}[0]) + "═══════════════════════════════════════════════════\n")
}