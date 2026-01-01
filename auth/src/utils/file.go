package utils

import "os"

func CheckExistFile(path string) bool {
	// ファイルが存在するかチェック
	_, err := os.Stat(path)

	return err == nil
}