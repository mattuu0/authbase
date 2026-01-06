#!/bin/bash

# プロジェクト名や出力ファイルの設定
COVERAGE_OUT="coverage.out"
COVERAGE_HTML="coverage.html"
export TEST_SUMMARY_FILE="$(pwd)/test_results.log"

# ヘルプメッセージを表示する関数
show_help() {
    echo "Usage: ./test.sh [command]"
    echo ""
    echo "Available commands:"
    echo "  all             - Run all tests (default)"
    echo "  unit            - Run unit tests only"
    echo "  integration     - Run integration tests only"
    echo "  coverage        - Generate coverage report"
    echo "  verbose         - Run tests with race detector"
    echo "  models          - Test models package"
    echo "  services        - Test services package"
    echo "  middlewares     - Test middlewares package"
    echo "  benchmark       - Run benchmark tests"
    echo "  clean-cache     - Clean test cache"
    echo "  clean           - Remove coverage files"
    echo "  help            - Show this help message"
}

# 以前の結果ファイルを削除
rm -f "$TEST_SUMMARY_FILE"

EXIT_CODE=0

# 実行する処理の分岐
case "$1" in
    test|all|"")
        echo "Running all tests..."
        go test -v ./...
        EXIT_CODE=$?
        ;;
    unit)
        echo "Running unit tests..."
        go test -v ./models/... ./services/... ./middlewares/... ./utils/...
        EXIT_CODE=$?
        ;;
    integration)
        echo "Running integration tests..."
        go test -v ./integration/...
        EXIT_CODE=$?
        ;;
    coverage)
        echo "Generating coverage report..."
        go test -v -coverprofile=$COVERAGE_OUT ./...
        EXIT_CODE=$?
        if [ $EXIT_CODE -eq 0 ]; then
            go tool cover -html=$COVERAGE_OUT -o $COVERAGE_HTML
            echo "Coverage report generated: $COVERAGE_HTML"
        fi
        ;;
    verbose)
        echo "Running tests in verbose mode..."
        go test -v -race ./...
        EXIT_CODE=$?
        ;;
    models)
        echo "Testing models package..."
        go test -v ./models/...
        EXIT_CODE=$?
        ;;
    services)
        echo "Testing services package..."
        go test -v ./services/...
        EXIT_CODE=$?
        ;;
    middlewares)
        echo "Testing middlewares package..."
        go test -v ./middlewares/...
        EXIT_CODE=$?
        ;;
    benchmark)
        echo "Running benchmark tests..."
        go test -bench=. -benchmem ./...
        EXIT_CODE=$?
        ;;
    clean-cache)
        echo "Cleaning test cache..."
        go clean -testcache
        EXIT_CODE=$?
        ;;
    clean)
        echo "Cleaning up..."
        rm -f $COVERAGE_OUT $COVERAGE_HTML
        EXIT_CODE=0
        ;;
    help|--help|-h)
        show_help
        EXIT_CODE=0
        ;;
    *)
        echo "Error: Unknown command '$1'"
        show_help
        exit 1
        ;;
esac

# 最後に集計結果を表示
echo ""
if [ -f "$TEST_SUMMARY_FILE" ]; then
    go run cmd/summary/main.go
fi

exit $EXIT_CODE