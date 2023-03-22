package api

import (
	"fmt"
	"github.com/samber/lo"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

var fileMap = map[string][]string{
	"图片":    {"png", "jpg", "jpeg", "gif", "bmp", "tiff", "svg", "ico"},
	"电子书":   {"pdf", "mobi", "azw3", "epub", "djvu", "cbz", "cbr"},
	"文档资料":  {"txt", "md", "doc", "docx", "xlsx", "csv", "ppt", "pptx", "rtf", "odt", "ods", "odp"},
	"压缩包":   {"zip", "rar", "7z", "tar", "gz", "bz2", "xz"},
	"音频":    {"mp3", "wmv", "m4a", "flac", "wav", "ogg", "aac", "m4b"},
	"视频":    {"mp4", "mkv", "mov", "flv", "wmv", "rmvb"},
	"exe程序": {"exe", "sh", "bat", "py", "jar", "app", "dmg", "msi", "apk", "ipa"},
}

func File(path string) []string {
	files := TreeFiles(path)
	c := lo.Map(files, func(f fs.FileInfo, _ int) string {
		return f.Name()
	})

	//op := make([]string, 0)
	//op = c.([]string)
	return c
}

func TreeFiles(path string) []fs.FileInfo {
	dir, _ := ioutil.ReadDir(path)
	return dir
}

func Drop(path string) {
	copyFiles(readFiles(path), path)

}

func readFiles(path string) map[string][]fs.FileInfo {
	m := make(map[string][]fs.FileInfo)
	dir, _ := ioutil.ReadDir(path)
	for _, i := range dir {
		name := i.Name()
		if i.IsDir() {
			continue
		}
		ext := filepath.Ext(name)
		// 去掉扩展名前面的点号
		ext = ext[1:]
		for i2, i3 := range fileMap {
			if lo.Contains(i3, ext) {
				//fmt.Println(name, "是个", i2)
				infos, ok := m[i2]
				if !ok {
					infos = make([]fs.FileInfo, 0, 0)

				}
				infos = append(infos, i)
				m[i2] = infos
			}
		}
		//fmt.Println(name+" 是文件夹 ", i.IsDir())

	}
	return m
}
func fileIsExisted(filename string) bool {
	existed := true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		existed = false
	}
	return existed
}
func copyFiles(m map[string][]fs.FileInfo, path string) {
	backupPath := filepath.Join(path, "备份")
	if createDirIfNotExist(backupPath) {
		return
	}

	for k, v := range m {
		tempPath := filepath.Join(path, k)
		createDirIfNotExist(tempPath)
		for _, i := range v {
			file, _ := ioutil.ReadFile(filepath.Join(path, i.Name()))

			err := ioutil.WriteFile(filepath.Join(tempPath, i.Name()), file, os.ModePerm)
			if err != nil {
				fmt.Println(err)
				return
			}
			err = ioutil.WriteFile(filepath.Join(backupPath, i.Name()), file, os.ModePerm)
			if err != nil {
				fmt.Println(err)
				return
			}
			err = os.Remove(filepath.Join(path, i.Name()))
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}

}

func createDirIfNotExist(backupPath string) bool {
	if !fileIsExisted(backupPath) {
		err := os.Mkdir(backupPath, fs.ModeDir)
		if err != nil {
			fmt.Println(err)
			return true
		}
	}
	return false
}
