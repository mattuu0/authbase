#!/bin/bash

# エラーが発生したら即座に終了
set -e

# プロジェクト名や出力ファイルの設定
COVERAGE_OUT="coverage.out"
COVERAGE_HTML="coverage.html"

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

# 実行する処理の分岐
case "$1" in
    test|all|"")
        echo "Running all tests..."
        go test -v ./...
        ;;
    unit)
        echo "Running unit tests..."
        go test -v ./models/... ./services/... ./middlewares/... ./utils/...
        ;;
    integration)
        echo "Running integration tests..."
        go test -v ./integration/...
        ;;
    coverage)
        echo "Generating coverage report..."
        go test -v -coverprofile=$COVERAGE_OUT ./...
        go tool cover -html=$COVERAGE_OUT -o $COVERAGE_HTML
        echo "Coverage report generated: $COVERAGE_HTML"
        ;;
    verbose)
        echo "Running tests in verbose mode..."
        go test -v -race ./...
        ;;
    models)
        echo "Testing models package..."
        go test -v ./models/...
        ;;
    services)
        echo "Testing services package..."
        go test -v ./services/...
        ;;
    middlewares)
        echo "Testing middlewares package..."
        go test -v ./middlewares/...
        ;;
    benchmark)
        echo "Running benchmark tests..."
        go test -bench=. -benchmem ./...
        ;;
    clean-cache)
        echo "Cleaning test cache..."
        go clean -testcache
        ;;
    clean)
        echo "Cleaning up..."
        rm -f $COVERAGE_OUT $COVERAGE_HTML
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        echo "Error: Unknown command '$1'"
        show_help
        exit 1
        ;;
esac