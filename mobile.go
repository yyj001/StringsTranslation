package main

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type LanguageStrings struct {
	language string
	strings  map[string]string
}

type StringLineXML struct {
	XMLName    xml.Name `xml:"string"`
	StringName string   `xml:"name,attr"`
	Value      string   `xml:",chardata"`
}

type StringFileXML struct {
	XMLName    xml.Name        `xml:"resources"`
	StringLine []StringLineXML `xml:",innerxml"`
}

type StringJsonElement struct {
	Lan string
	Str string
}

func GenerateStringsFile(context *gin.Context) {
	var stringRecords []TableStrings
	db.Find(&stringRecords)
	var languageMaps map[string]map[string]string
	languageMaps = make(map[string]map[string]string)
	for _, strRecord := range stringRecords {
		if len(strRecord.Name) == 0 {
			continue
		}
		var stringLans []StringJsonElement
		err := json.Unmarshal([]byte(strRecord.TranslateStrs), &stringLans)
		if err != nil {
			OnError("解析json失败", context)
			return
		}
		for _, stringLan := range stringLans {
			if languageMaps[stringLan.Lan] == nil {
				languageMaps[stringLan.Lan] = make(map[string]string)
			}
			var languageMap = languageMaps[stringLan.Lan]
			languageMap[strRecord.Name] = stringLan.Str
		}
	}

	for lan, languageMap := range languageMaps {
		var dirName = "values"
		if lan != "en" {
			dirName = "values-" + lan
		}
		// 生成文件夹
		var dirPath = "strings/" + dirName
		os.MkdirAll(dirPath, os.ModePerm)
		// 生成xml
		var stringXML = StringFileXML{}
		for name, str := range languageMap {
			var line = StringLineXML{StringName: name, Value: str}
			stringXML.StringLine = append(stringXML.StringLine, line)
		}
		// 写文件
		output, err := xml.MarshalIndent(stringXML, "", "    ")
		if err != nil {
			fmt.Println(err)
			return
		}
		file, _ := os.Create(dirPath + "/strings.xml")
		defer file.Close()
		file.Write([]byte(xml.Header))
		file.Write(output)
		file.Close()
	}
	// 压缩
	ZipFile(FILE_DIR, FILE_NAME)
	return
}

func ZipFile(src_dir string, zip_file_name string) error {
	dir, err := ioutil.ReadDir(src_dir)
	if err != nil {
		return err
	}
	if len(dir) == 0 {
		return nil
	}
	// 预防：旧文件无法覆盖  删除路径下所有的文件
	os.RemoveAll(zip_file_name)

	// 创建：zip文件
	zipfile, _ := os.Create(zip_file_name)
	defer zipfile.Close()

	// 打开：zip文件
	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	// 遍历路径信息
	filepath.Walk(src_dir, func(path string, info os.FileInfo, _ error) error {

		// 如果是源路径，提前进行下一个遍历
		if path == src_dir {
			return nil
		}
		// 获取：文件头信息
		header, _ := zip.FileInfoHeader(info)

		header.Name = strings.TrimPrefix(path, src_dir+`\`)

		// 判断：文件是不是文件夹
		if info.IsDir() {
			header.Name += `/`
		} else {
			// 设置：zip的文件压缩算法
			header.Method = zip.Deflate
		}

		// 创建：压缩包头部信息
		writer, _ := archive.CreateHeader(header)
		if !info.IsDir() {
			file, _ := os.Open(path)
			defer file.Close()
			io.Copy(writer, file)
		}

		return nil
	})
	return nil
}

const FILE_NAME = "strings.jar"
const FILE_DIR = "./strings"

var isNeedUpdateFile = true

func DownloadStringsFile(context *gin.Context) {
	// 生成jar
	isFileExist, _ := PathExists(FILE_NAME)
	// 缓存更新或者文件不存在，重新生成jar
	if isNeedUpdateFile || !isFileExist {
		GenerateStringsFile(context)
		isNeedUpdateFile = false
		fmt.Println("重新生成" + FILE_NAME)
	}
	file, _ := os.Open(FILE_NAME)
	defer file.Close()
	fileHeader := make([]byte, 512)
	file.Read(fileHeader)
	fileStat, _ := file.Stat()
	context.Writer.Header().Set("Content-Disposition", "attachment; filename="+FILE_NAME)
	context.Writer.Header().Set("Content-Type", http.DetectContentType(fileHeader))
	context.Writer.Header().Set("Content-Length", strconv.FormatInt(fileStat.Size(), 10))
	file.Seek(0, 0)
	io.Copy(context.Writer, file)
	return
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
