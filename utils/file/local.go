package file

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"prismx_cli/utils/logger"
	"strings"
)

// FilesList 获取目录下指定后缀文件
// 如： yaml、json
func FilesList(fileDir string, v string) (s []string) {
	filepath.Walk(fileDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.Contains(path, "lib/exploits/nuclei") || strings.Contains(path, "lib\\exploits\\nuclei") {
				return filepath.SkipDir
			}
			if strings.HasSuffix(path, "."+v) {
				s = append(s, path)
			}
			return nil
		})
	return s
}

// FilesEmbedList Embed获取目录下的指定后缀类型文件
func FilesEmbedList(emb embed.FS, fileDir string, v string) (s []string) {
	dirList, err := emb.ReadDir(fileDir)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	for _, item := range dirList {
		if item.IsDir() {
			FilesEmbedList(emb, fileDir+item.Name(), v)
		} else {
			if strings.HasSuffix(item.Name(), "."+v) {
				s = append(s, fmt.Sprintf("%s/%s", fileDir, item.Name()))
			}
		}
	}
	return s
}
