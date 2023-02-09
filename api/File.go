package api

import (
	"fmt"
	"github.com/thoas/go-funk"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"
)

var fileMap = map[string][]string{
	"图片":    {"png", "jpg", "jpeg", "gif"},
	"电子书":   {"pdf", "mobi", "azw3", "epub"},
	"文档资料":  {"txt", "md", "doc", "docx", "xlsx", "csv", "ppt", "pptx"},
	"压缩包":   {"zip", "rar", "7z"},
	"音频":    {"mp3", "wmv", "m4a", "flac"},
	"视频":    {"mp4", "mkv", "avi"},
	"exe程序": {"exe"},
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
		split := strings.Split(name, ".")
		s := split[len(split)-1]
		for i2, i3 := range fileMap {
			if funk.Contains(i3, s) {
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
	backupPath := path + "/" + "备份"
	if createDirIfNotExist(backupPath) {
		return
	}

	for k, v := range m {
		tempPath := path + "/" + k
		createDirIfNotExist(tempPath)
		for _, i := range v {
			file, _ := ioutil.ReadFile(path + "/" + i.Name())
			err := ioutil.WriteFile(tempPath+"/"+i.Name(), file, os.ModePerm)
			if err != nil {
				fmt.Println(err)
				return
			}
			err = ioutil.WriteFile(backupPath+"/"+i.Name(), file, os.ModePerm)
			if err != nil {
				fmt.Println(err)
				return
			}
			err = os.Remove(path + "/" + i.Name())
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
