package file

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func GetAppData(app string) (string, error) {
	appdata, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get UserConfigDir: %v", err)
	}

	// 创建 AppData\Roaming\myapp 目录
	path := appdata + "/" + app
	path = strings.ReplaceAll(path, "\\", "/")
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}
	return path, nil
}

func WriteToFileAsJson(filename string, v interface{}, indent string, truncateIfExist bool) error {

	buf, err := json.MarshalIndent(v, "", indent)
	if err != nil {
		return err
	}
	err = WriteToFile(filename, buf, truncateIfExist)
	if err != nil {
		return err
	}
	return nil
}

func ReadFileJsonToObject(filename string, obj interface{}) error {

	buf, err := ReadFromFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buf, &obj)
	if err != nil {
		return err
	}
	return nil
}

// 递归创建文件夹
func CreateDirRecursive(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) { // not exist
		err := os.MkdirAll(filePath, os.ModePerm)
		return err
	}
	return nil
}

func WriteToFile(filename string, content []byte, truncateIfExist bool) error {
	flag := os.O_RDWR | os.O_CREATE | os.O_APPEND
	if truncateIfExist {
		flag = os.O_RDWR | os.O_CREATE | os.O_TRUNC
	}

	if IsFileNotExist(GetFilePath(filename)) {
		err := CreateDirRecursive(GetFilePath(filename))
		if err != nil {
			return fmt.Errorf("failed CreateDirRecursive, err:%v", err)
		}
	}

	fileObj, err := os.OpenFile(filename, flag, 0644)
	if err != nil {
		return err
	}
	defer fileObj.Close()

	n, err := fileObj.Write(content)
	if err != nil {
		return err
	}
	if n != len(content) {
		return errors.New("written length error")
	}
	return fileObj.Sync()
}

func WriteToFileWithFlag(filename string, content []byte, flag int) error {

	if IsFileNotExist(GetFilePath(filename)) {
		err := CreateDirRecursive(GetFilePath(filename))
		if err != nil {
			return fmt.Errorf("failed CreateDirRecursive, err:%v", err)
		}
	}

	fileObj, err := os.OpenFile(filename, flag, 0644)
	if err != nil {
		return err
	}
	defer fileObj.Close()

	n, err := fileObj.Write(content)
	if err != nil {
		return err
	}
	if n != len(content) {
		return errors.New("written length error")
	}
	return fileObj.Sync()
}

func ReadFromFile(filename string) ([]byte, error) {
	fileObj, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fileObj.Close()

	content, err := ioutil.ReadAll(fileObj)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func MoveFile(file, targetDirWithoutTargetName string) error {
	targetDir := targetDirWithoutTargetName
	err := CreateDirRecursive(targetDir)
	if err != nil {
		return err
	}

	base := path.Base(file)
	t := path.Join(targetDir, base)
	i := 1

	for {
		if IsFileNotExist(t) {
			break
		}
		fileSuffix := path.Ext(base)                         //获取文件后缀
		filenameOnly := strings.TrimSuffix(base, fileSuffix) //获取文件名
		t = path.Join(targetDir, filenameOnly+fmt.Sprintf("(%v)", i)+fileSuffix)
		i++
	}

	err = os.Rename(file, t)
	if err != nil {
		return err
	}

	return nil
}

func ListDir(folder string) ([]string, error) {
	var ret []string
	isFile, err := IsFile(folder)
	if err != nil {
		return ret, err
	}
	if isFile {
		ret = append(ret, folder)
	} else {

		files, err := ioutil.ReadDir(folder)
		if err != nil {
			return ret, err
		}
		for _, fi := range files {
			t := path.Join(folder, fi.Name())
			if fi.IsDir() {
				r, _ := ListDir(t)
				ret = append(ret, r...)
			} else {
				ret = append(ret, t)
			}
		}
	}
	return ret, nil
}

// 如果是单个文件则返回 true.
func IsFile(f string) (isFile bool, err error) {
	fi, err := os.Stat(f)
	if err != nil {
		return
	}
	isFile = !fi.IsDir()
	return
}

func IsFileExist(f string) bool {
	return !IsFileNotExist(f)
}
func IsFileNotExist(f string) bool {
	_, err := os.Stat(f)
	return os.IsNotExist(err)
}

func AddPathSepIfNeed(path string) (newPath string) {
	newPath = path
	if len(path) > 0 {
		if path[len(path)-1:] != "/" {
			newPath += "/"
		}
	} else {
		newPath += "/"
	}
	return
}

// /a/b/c.txt -> /a/b/
func GetFilePath(fullfilename string) string {
	dir, _ := filepath.Split(fullfilename)
	return dir
}

// /a/b/c.txt -> /a/b/, c.txt
func GetFilePathAndName(fullfilename string) (string, string) {
	dir, file := filepath.Split(fullfilename)
	return dir, file
}

// /a/b/c.txt -> c.txt
func GetFilename(fullfilename string) string {
	return path.Base(fullfilename)
}

// /a/b/c.txt -> c
func GetFilenameOnly(fullfilename string) string {
	return strings.TrimSuffix(fullfilename, path.Ext(fullfilename))
}

// /a/b/c.txt -> .txt
func GetFileSuffix(fullfilename string) string {
	return path.Ext(fullfilename)
}

func GenFullFileNameByAppData(appName string, fileName string, subfolders ...string) (string, error) {
	p, err := GetAppData(appName)
	if err != nil {
		return "", fmt.Errorf("failed GetAppData, err:%v", err)
	}
	fileFullName := AddPathSepIfNeed(p) + strings.Join(subfolders, "/") + "/" + fileName
	return fileFullName, nil
}

func GenPathByAppData(appName string, subfolders ...string) (string, error) {
	p, err := GetAppData(appName)
	if err != nil {
		return "", fmt.Errorf("failed GetAppData, err:%v", err)
	}
	fileFullName := AddPathSepIfNeed(p) + strings.Join(subfolders, "/")
	return fileFullName, nil
}
