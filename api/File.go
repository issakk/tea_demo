package api

import (
	"fmt"
	"github.com/thoas/go-funk"
	"io/fs"
	"io/ioutil"
	"strings"
)

var fileMap = map[string][]string{
	"图片":     {"png", "jpg", "jpeg", "gif"},
	"电子书":   {"pdf", "mobi", "azw3"},
	"文档资料": {"txt", "md", "doc", "docx", "xlsx", "csv", "ppt", "pptx"},
	"压缩包":   {"zip", "rar", "7z"},
	"音频":     {"mp3", "wmv", "m4a", "flac"},
	"视频":     {"mp4", "mkv", "avi"},
}

func File(path string) []string {
	files := TreeFiles(path)
	c := funk.Map(files, func(f fs.FileInfo) string {
		return f.Name()
	})

	op := make([]string, 0)
	op = c.([]string)
	return op
}

func TreeFiles(path string) []fs.FileInfo {
	dir, _ := ioutil.ReadDir(path)
	return dir
}

func readFiles(path string) map[string][]fs.FileInfo {
	m := make(map[string][]fs.FileInfo)
	dir, _ := ioutil.ReadDir(path)
	for _, i := range dir {
		name := i.Name()
		if i.IsDir() {
			continue
		}
		split := strings.Split(name, ".")
		s := split[len(split)-1]
		for i2, i3 := range fileMap {
			if funk.Contains(i3, s) {
				fmt.Println(name, "是个", i2)
				infos, ok := m[i2]
				if !ok {
					infos = make([]fs.FileInfo, 0, 0)

				}
				infos = append(infos, i)
				m[i2] = infos
			}
		}
		fmt.Println(name+" 是文件夹 ", i.IsDir())

	}
	return m
}
