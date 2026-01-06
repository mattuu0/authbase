package testing

import (
	"fmt"
	"os"
	"testing"
)

// RecordResult records the result of a subtest to a file
func RecordResult(t *testing.T) {
	t.Cleanup(func() {
		success := !t.Failed()
		status := "✅ SUCCESS"
		if !success {
			status = "❌ FAILED"
		}

		// ファイルに書き込み
		filePath := os.Getenv("TEST_SUMMARY_FILE")
		if filePath != "" {
			f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err == nil {
				defer f.Close()
				line := fmt.Sprintf("%s|%s\n", status, t.Name())
				f.WriteString(line)
			}
		}
	})
}

// PrintFinalSummary (Deprecated) kept for compatibility but does nothing now
func PrintFinalSummary() {
	// No-op: Summary is now printed by an external tool reading the result file
}