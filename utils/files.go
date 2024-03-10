package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func GetDirFiles(dirPath string) []string {
	dirPath = strings.Replace(dirPath,`\`,`\\`,-1)
	var fileNameList []string
	// filepath.Walk 方式会遍历所有子文件夹
	//err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
	//	if info.IsDir() {
	//		return nil
	//	}
	//	if filepath.Ext(path) != ".xml"{
	//		return nil
	//	}
	//	fileNameList = append(fileNameList, strings.Replace(path,`\`,`\\`,-1))
	//	return nil
	//})
	files, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Println("os.ReadDir 加载文件夹内容失败！")
		log.Fatal(err)
	}
	for _, file := range files {
		if !file.IsDir() {
			if filepath.Ext(file.Name()) == ".xml"{
				fileNameList = append(fileNameList,dirPath + "\\\\" + file.Name())
			}
		}
	}
	return fileNameList
}

func GetFileContent(filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(content))
	return string(content)
}