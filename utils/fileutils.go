package utils

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// 判断文件或文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil // 存在
	}
	if os.IsNotExist(err) {
		return false, nil // 不存在
	}
	return false, err // 不确定
}

// 复制文件夹及文件夹下的文件
func CopyDir(srcPath string, destPath string) error {
	//检测目录正确性
	if srcInfo, err := os.Stat(srcPath); err != nil {
		fmt.Println(err.Error())
		return err
	} else {
		if !srcInfo.IsDir() {
			e := errors.New("srcPath不是一个正确的目录！")
			fmt.Println(e.Error())
			return e
		}
	}
	if destInfo, err := os.Stat(destPath); err != nil {
		fmt.Println(err.Error())
		return err
	} else {
		if !destInfo.IsDir() {
			e := errors.New("destInfo不是一个正确的目录！")
			fmt.Println(e.Error())
			return e
		}
	}

	//加上拷贝时间:不用可以去掉
	//destPath = destPath + "_" + time.Now().Format("20060102150405")

	err := filepath.Walk(srcPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		path = strings.Replace(path, "\\", "/", -1)
		destNewPath := strings.Replace(path, srcPath, destPath, -1)
		if f.IsDir() {
			// dest创建对应的目录
			err = os.MkdirAll(destNewPath, 0777)
			return err
		} else {
			// 复制文件
			// fmt.Println("复制文件:" + path + " 到 " + destNewPath)
			CopyFile(path, destNewPath)
		}
		return nil
	})
	if err != nil {
		fmt.Printf(err.Error())
	}
	return err
}

// 复制文件
func CopyFile(src, dest string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}
	defer srcFile.Close()
	//分割path目录
	destSplitPathDirs := strings.Split(dest, "/")

	//检测时候存在目录
	destSplitPath := ""
	for index, dir := range destSplitPathDirs {
		if index < len(destSplitPathDirs)-1 {
			destSplitPath = destSplitPath + dir + "/"
			b, _ := PathExists(destSplitPath)
			if b == false {
				fmt.Println("创建目录:" + destSplitPath)
				//创建目录
				err := os.Mkdir(destSplitPath, os.ModePerm)
				if err != nil {
					fmt.Println(err)
					return 0, err
				}
			}
		}
	}
	dstFile, err := os.Create(dest)
	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}
	defer dstFile.Close()

	return io.Copy(dstFile, srcFile)
}

// 解压
func Unzip(archive, target string) error {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(target, 0777); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return nil
}

// 强制删除文件
func ForceDelFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		err = exec.Command("cmd", "DEL", "/S", "/Q", "names", path).Start()
		if err != nil {
			return err
		}
	}
	return nil
}
