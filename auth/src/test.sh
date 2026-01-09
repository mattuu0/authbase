#!/bin/bash

# ========================================
# AuthBase Test Runner Script
# ========================================
# このスクリプトはテストを見やすく実行します

# カラーコード
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 絵文字
CHECK_MARK="✓"
CROSS_MARK="✗"
ROCKET="🚀"
PACKAGE="📦"
CLOCK="⏱️"
TROPHY="🏆"

# テスト開始時刻
START_TIME=$(date +%s)

# ========================================
# ヘルパー関数
# ========================================

print_header() {
    echo ""
    echo -e "${BLUE}╔════════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║${NC}                                                ${BLUE}║${NC}"
    echo -e "${BLUE}║${NC}  $1${BLUE}║${NC}"
    echo -e "${BLUE}║${NC}                                                ${BLUE}║${NC}"
    echo -e "${BLUE}╚════════════════════════════════════════════════╝${NC}"
    echo ""
}

print_section() {
    echo ""
    echo -e "${CYAN}═══════════════════════════════════════════════════${NC}"
    echo -e "${CYAN}  $1${NC}"
    echo -e "${CYAN}═══════════════════════════════════════════════════${NC}"
    echo ""
}

print_success() {
    echo -e "${GREEN}${CHECK_MARK} $1${NC}"
}

print_error() {
    echo -e "${RED}${CROSS_MARK} $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# テスト実行関数
run_test() {
    local test_name=$1
    local test_path=$2
    
    echo -e "${MAGENTA}${PACKAGE} Testing: ${test_name}${NC}"
    echo ""
    
    # テストを実行してリアルタイム出力
    go test -v -count=1 ${test_path} 2>&1 | while IFS= read -r line; do
        # カラー付けと絵文字追加
        if [[ $line == *"PASS"* ]]; then
            echo -e "${GREEN}${CHECK_MARK} ${line}${NC}"
        elif [[ $line == *"FAIL"* ]]; then
            echo -e "${RED}${CROSS_MARK} ${line}${NC}"
        elif [[ $line == *"RUN"* ]]; then
            echo -e "${CYAN}${ROCKET} ${line}${NC}"
        elif [[ $line == *"---"* ]]; then
            echo -e "${YELLOW}    ${line}${NC}"
        else
            echo "    ${line}"
        fi
    done
    
    # テスト結果を取得
    local exit_code=${PIPESTATUS[0]}
    
    echo ""
    if [ $exit_code -eq 0 ]; then
        print_success "${test_name} - All tests passed!"
        return 0
    else
        print_error "${test_name} - Some tests failed!"
        return 1
    fi
}

# 経過時間を計算
calc_elapsed() {
    local end_time=$(date +%s)
    local elapsed=$((end_time - START_TIME))
    echo "${elapsed}s"
}

# ========================================
# メイン処理
# ========================================

main() {
    print_header "      AuthBase Testing Suite      "
    
    print_info "Test started at: $(date '+%Y-%m-%d %H:%M:%S')"
    echo ""
    
    # テストカウンター
    local total_tests=0
    local passed_tests=0
    local failed_tests=0
    
    # モデル層テスト
    print_section "📊 Models Layer Tests"
    
    if run_test "User Create" "./models/ -run TestCreateUser"; then
        ((passed_tests++))
    else
        ((failed_tests++))
    fi
    ((total_tests++))
    
    if run_test "User Retrieve" "./models/ -run TestGetUser"; then
        ((passed_tests++))
    else
        ((failed_tests++))
    fi
    ((total_tests++))
    
    if run_test "User Update" "./models/ -run TestUpdateUser"; then
        ((passed_tests++))
    else
        ((failed_tests++))
    fi
    ((total_tests++))
    
    if run_test "User Delete" "./models/ -run TestDeleteUser"; then
        ((passed_tests++))
    else
        ((failed_tests++))
    fi
    ((total_tests++))
    
    if run_test "User Labels" "./models/ -run TestLabel"; then
        ((passed_tests++))
    else
        ((failed_tests++))
    fi
    ((total_tests++))
    
    # サービス層テスト
    print_section "⚙️  Services Layer Tests"
    
    if run_test "User Signup" "./services/ -run TestSignup"; then
        ((passed_tests++))
    else
        ((failed_tests++))
    fi
    ((total_tests++))
    
    if run_test "User Login" "./services/ -run TestLogin"; then
        ((passed_tests++))
    else
        ((failed_tests++))
    fi
    ((total_tests++))
    
    # 最終結果
    print_section "📈 Test Summary"
    
    local elapsed=$(calc_elapsed)
    echo -e "${BLUE}${CLOCK} Total Time: ${elapsed}${NC}"
    echo -e "${BLUE}${PACKAGE} Total Test Suites: ${total_tests}${NC}"
    echo -e "${GREEN}${CHECK_MARK} Passed: ${passed_tests}${NC}"
    echo -e "${RED}${CROSS_MARK} Failed: ${failed_tests}${NC}"
    echo ""
    
    # Markdown 出力
    local summary_file="TEST_RESULT.md"
    echo "# 🧪 Auth Service Test Report" > $summary_file
    echo "" >> $summary_file
    echo "## 📊 Summary" >> $summary_file
    echo "| Metric | Value |" >> $summary_file
    echo "| :--- | :--- |" >> $summary_file
    echo "| ⏱️ Total Time | ${elapsed} |" >> $summary_file
    echo "| 📦 Total Suites | ${total_tests} |" >> $summary_file
    echo "| ✅ Passed | ${passed_tests} |" >> $summary_file
    echo "| ❌ Failed | ${failed_tests} |" >> $summary_file
    
    # 成功率計算
    if [ $total_tests -gt 0 ]; then
        local success_rate=$((passed_tests * 100 / total_tests))
        echo -e "${CYAN}Success Rate: ${success_rate}%${NC}"
        echo "| 📈 Success Rate | ${success_rate}% |" >> $summary_file
        
        if [ $success_rate -eq 100 ]; then
            echo ""
            echo -e "${GREEN}${TROPHY}${TROPHY}${TROPHY} Perfect! All tests passed! ${TROPHY}${TROPHY}${TROPHY}${NC}"
            echo "" >> $summary_file
            echo "### 🎉 Result: Perfect!" >> $summary_file
            echo "All tests passed successfully." >> $summary_file
        else
            echo "" >> $summary_file
            echo "### ⚠️ Result: Needs Attention" >> $summary_file
            echo "Some tests failed. Please check the logs." >> $summary_file
        fi
    fi
    
    echo "" >> $summary_file
    echo "---" >> $summary_file
    echo "*Generated at: $(date '+%Y-%m-%d %H:%M:%S')*" >> $summary_file
    
    echo ""
    print_info "Test report generated: ${summary_file}"
    print_info "Test finished at: $(date '+%Y-%m-%d %H:%M:%S')"
    
    # 失敗があれば終了コード1を返す
    if [ $failed_tests -gt 0 ]; then
        exit 1
    fi
}

# ========================================
# スクリプト実行
# ========================================

# 引数処理
case "${1:-all}" in
    "models")
        print_header "      Models Tests Only      "
        run_test "Models Package" "./models/..."
        ;;
    "services")
        print_header "      Services Tests Only      "
        run_test "Services Package" "./services/..."
        ;;
    "integration")
        print_header "      Integration Tests Only      "
        run_test "Integration Package" "./integration/..."
        ;;
    "quick")
        print_header "      Quick Test Run      "
        print_info "Running tests in short mode..."
        go test -short -v ./... 2>&1 | sed 's/PASS/✓ PASS/g' | sed 's/FAIL/✗ FAIL/g'
        ;;
    "watch")
        print_header "      Watch Mode      "
        print_warning "Watching for file changes... (Press Ctrl+C to stop)"
        find . -name "*.go" | entr -c $0 quick
        ;;
    "all"|*)
        main
        ;;
esac