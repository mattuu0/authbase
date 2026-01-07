# ========================================
# AuthBase Testing Makefile
# ========================================
# このMakefileはテストを簡単かつ見やすく実行するためのものです

.PHONY: all test test-unit test-integration test-coverage clean help
.PHONY: test-models test-services test-middlewares test-utils
.PHONY: test-user test-session test-label test-auth test-jwt
.PHONY: watch test-verbose test-short test-summary

# デフォルトカラー設定
BLUE := \033[0;34m
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

# ========================================
# メインターゲット
# ========================================

all: banner test

# バナー表示
banner:
	@echo "$(BLUE)"
	@echo "╔════════════════════════════════════════════════╗"
	@echo "║                                                ║"
	@echo "║         AuthBase Testing Suite                ║"
	@echo "║                                                ║"
	@echo "╚════════════════════════════════════════════════╝"
	@echo "$(NC)"

# すべてのテストを実行（リアルタイム出力）
test: banner
	@echo "$(GREEN)▶ Running all tests...$(NC)"
	@go test -v -count=1 ./... 2>&1 | \
		sed 's/PASS/$(GREEN)✓ PASS$(NC)/g' | \
		sed 's/FAIL/$(RED)✗ FAIL$(NC)/g' | \
		sed 's/^=== RUN/$(BLUE)=== RUN$(NC)/g'

# ========================================
# カテゴリ別テスト
# ========================================

# 単体テストのみ
test-unit: banner
	@echo "$(GREEN)▶ Running unit tests...$(NC)"
	@go test -v -count=1 ./models/... ./services/... ./middlewares/... ./utils/... 2>&1 | \
		sed 's/PASS/$(GREEN)✓ PASS$(NC)/g' | \
		sed 's/FAIL/$(RED)✗ FAIL$(NC)/g'

# 統合テストのみ
test-integration: banner
	@echo "$(GREEN)▶ Running integration tests...$(NC)"
	@go test -v -count=1 ./integration/... 2>&1 | \
		sed 's/PASS/$(GREEN)✓ PASS$(NC)/g' | \
		sed 's/FAIL/$(RED)✗ FAIL$(NC)/g'

# ========================================
# パッケージ別テスト
# ========================================

# モデル層テスト
test-models:
	@echo "$(BLUE)▶ Testing models package...$(NC)"
	@go test -v -count=1 ./models/... 2>&1 | \
		sed 's/PASS/$(GREEN)✓ PASS$(NC)/g' | \
		sed 's/FAIL/$(RED)✗ FAIL$(NC)/g'

# サービス層テスト
test-services:
	@echo "$(BLUE)▶ Testing services package...$(NC)"
	@go test -v -count=1 ./services/... 2>&1 | \
		sed 's/PASS/$(GREEN)✓ PASS$(NC)/g' | \
		sed 's/FAIL/$(RED)✗ FAIL$(NC)/g'

# ミドルウェアテスト
test-middlewares:
	@echo "$(BLUE)▶ Testing middlewares package...$(NC)"
	@go test -v -count=1 ./middlewares/... 2>&1 | \
		sed 's/PASS/$(GREEN)✓ PASS$(NC)/g' | \
		sed 's/FAIL/$(RED)✗ FAIL$(NC)/g'

# ユーティリティテスト
test-utils:
	@echo "$(BLUE)▶ Testing utils package...$(NC)"
	@go test -v -count=1 ./utils/... 2>&1 | \
		sed 's/PASS/$(GREEN)✓ PASS$(NC)/g' | \
		sed 's/FAIL/$(RED)✗ FAIL$(NC)/g'

# ========================================
# 機能別テスト
# ========================================

# ユーザー関連テスト
test-user:
	@echo "$(YELLOW)▶ Testing user functionality...$(NC)"
	@go test -v -count=1 -run "TestUser|TestCreate|TestGet|TestUpdate|TestDelete" ./models/... 2>&1 | \
		sed 's/PASS/$(GREEN)✓ PASS$(NC)/g' | \
		sed 's/FAIL/$(RED)✗ FAIL$(NC)/g'

# セッション関連テスト
test-session:
	@echo "$(YELLOW)▶ Testing session functionality...$(NC)"
	@go test -v -count=1 -run "TestSession|TestNewSession|TestGetSession|TestDeleteSession" ./... 2>&1 | \
		sed 's/PASS/$(GREEN)✓ PASS$(NC)/g' | \
		sed 's/FAIL/$(RED)✗ FAIL$(NC)/g'

# ラベル関連テスト
test-label:
	@echo "$(YELLOW)▶ Testing label functionality...$(NC)"
	@go test -v -count=1 -run "TestLabel|TestAddLabel|TestRemoveLabel" ./... 2>&1 | \
		sed 's/PASS/$(GREEN)✓ PASS$(NC)/g' | \
		sed 's/FAIL/$(RED)✗ FAIL$(NC)/g'

# 認証関連テスト
test-auth:
	@echo "$(YELLOW)▶ Testing authentication...$(NC)"
	@go test -v -count=1 -run "TestLogin|TestSignup|TestAuth" ./... 2>&1 | \
		sed 's/PASS/$(GREEN)✓ PASS$(NC)/g' | \
		sed 's/FAIL/$(RED)✗ FAIL$(NC)/g'

# JWT関連テスト
test-jwt:
	@echo "$(YELLOW)▶ Testing JWT...$(NC)"
	@go test -v -count=1 -run "TestJWT|TestToken|TestAccess" ./services/... 2>&1 | \
		sed 's/PASS/$(GREEN)✓ PASS$(NC)/g' | \
		sed 's/FAIL/$(RED)✗ FAIL$(NC)/g'

# ========================================
# テストオプション
# ========================================

# 短縮版（高速）
test-short:
	@echo "$(GREEN)▶ Running tests in short mode...$(NC)"
	@go test -short -v -count=1 ./... 2>&1 | \
		sed 's/PASS/$(GREEN)✓ PASS$(NC)/g' | \
		sed 's/FAIL/$(RED)✗ FAIL$(NC)/g'

# サマリーのみ表示
test-summary:
	@echo "$(GREEN)▶ Running tests (summary only)...$(NC)"
	@go test -count=1 ./... 2>&1 | \
		grep -E "(PASS|FAIL|^ok|^FAIL)" | \
		sed 's/PASS/$(GREEN)✓ PASS$(NC)/g' | \
		sed 's/FAIL/$(RED)✗ FAIL$(NC)/g' | \
		sed 's/^ok/$(GREEN)✓ ok$(NC)/g'

# レースコンディションチェック
test-race:
	@echo "$(YELLOW)▶ Running tests with race detector...$(NC)"
	@go test -race -v -count=1 ./... 2>&1 | \
		sed 's/PASS/$(GREEN)✓ PASS$(NC)/g' | \
		sed 's/FAIL/$(RED)✗ FAIL$(NC)/g'

# 詳細モード
test-verbose:
	@echo "$(GREEN)▶ Running tests in verbose mode...$(NC)"
	@go test -v -count=1 -race ./... 2>&1 | \
		sed 's/PASS/$(GREEN)✓ PASS$(NC)/g' | \
		sed 's/FAIL/$(RED)✗ FAIL$(NC)/g' | \
		sed 's/^=== RUN/$(BLUE)=== RUN$(NC)/g' | \
		sed 's/^--- PASS/$(GREEN)--- PASS$(NC)/g' | \
		sed 's/^--- FAIL/$(RED)--- FAIL$(NC)/g'

# ファイル監視モード（要: entr コマンド）
watch:
	@echo "$(YELLOW)▶ Watching for file changes...$(NC)"
	@echo "$(YELLOW)  Press Ctrl+C to stop$(NC)"
	@find . -name "*.go" | entr -c make test-summary

# ========================================
# カバレッジ
# ========================================

# カバレッジレポート生成
test-coverage:
	@echo "$(GREEN)▶ Generating coverage report...$(NC)"
	@go test -v -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Coverage report generated: coverage.html$(NC)"
	@go tool cover -func=coverage.out | grep total | \
		awk '{print "$(BLUE)Total Coverage: $(GREEN)" $$3 "$(NC)"}'

# カバレッジをコンソールに表示
coverage-console:
	@echo "$(GREEN)▶ Calculating coverage...$(NC)"
	@go test -coverprofile=coverage.out ./... > /dev/null 2>&1
	@go tool cover -func=coverage.out | \
		sed 's/^/  /' | \
		sed 's/\t/  /g'
	@echo ""
	@go tool cover -func=coverage.out | grep total | \
		awk '{print "$(BLUE)Total Coverage: $(GREEN)" $$3 "$(NC)"}'

# ========================================
# ベンチマーク
# ========================================

# ベンチマークテスト
benchmark:
	@echo "$(YELLOW)▶ Running benchmark tests...$(NC)"
	@go test -bench=. -benchmem -run=^Benchmark ./...

# ベンチマーク比較
benchmark-compare:
	@echo "$(YELLOW)▶ Running benchmark (saving to bench.txt)...$(NC)"
	@go test -bench=. -benchmem -run=^Benchmark ./... | tee bench.txt

# ========================================
# クリーンアップ
# ========================================

# テストキャッシュクリア
test-clean:
	@echo "$(YELLOW)▶ Cleaning test cache...$(NC)"
	@go clean -testcache
	@echo "$(GREEN)✓ Test cache cleared$(NC)"

# カバレッジファイル削除
clean:
	@echo "$(YELLOW)▶ Cleaning up...$(NC)"
	@rm -f coverage.out coverage.html bench.txt
	@echo "$(GREEN)✓ Cleanup complete$(NC)"

# すべてクリーン
clean-all: clean test-clean
	@echo "$(GREEN)✓ All cleaned$(NC)"

# ========================================
# ヘルプ
# ========================================

help:
	@echo "$(BLUE)"
	@echo "╔════════════════════════════════════════════════╗"
	@echo "║        AuthBase Testing Commands               ║"
	@echo "╚════════════════════════════════════════════════╝"
	@echo "$(NC)"
	@echo "$(GREEN)Main Commands:$(NC)"
	@echo "  make test              - Run all tests with color output"
	@echo "  make test-unit         - Run unit tests only"
	@echo "  make test-integration  - Run integration tests only"
	@echo ""
	@echo "$(GREEN)Package Tests:$(NC)"
	@echo "  make test-models       - Test models package"
	@echo "  make test-services     - Test services package"
	@echo "  make test-middlewares  - Test middlewares package"
	@echo "  make test-utils        - Test utils package"
	@echo ""
	@echo "$(GREEN)Feature Tests:$(NC)"
	@echo "  make test-user         - Test user functionality"
	@echo "  make test-session      - Test session management"
	@echo "  make test-label        - Test label operations"
	@echo "  make test-auth         - Test authentication"
	@echo "  make test-jwt          - Test JWT operations"
	@echo ""
	@echo "$(GREEN)Test Options:$(NC)"
	@echo "  make test-short        - Fast test execution"
	@echo "  make test-summary      - Show summary only"
	@echo "  make test-race         - Run with race detector"
	@echo "  make test-verbose      - Detailed output"
	@echo "  make watch             - Watch mode (requires entr)"
	@echo ""
	@echo "$(GREEN)Coverage:$(NC)"
	@echo "  make test-coverage     - Generate HTML coverage report"
	@echo "  make coverage-console  - Show coverage in console"
	@echo ""
	@echo "$(GREEN)Benchmark:$(NC)"
	@echo "  make benchmark         - Run benchmarks"
	@echo "  make benchmark-compare - Save benchmark results"
	@echo ""
	@echo "$(GREEN)Cleanup:$(NC)"
	@echo "  make test-clean        - Clear test cache"
	@echo "  make clean             - Remove coverage files"
	@echo "  make clean-all         - Clean everything"
	@echo ""
	@echo "$(GREEN)Help:$(NC)"
	@echo "  make help              - Show this help message"
	@echo ""