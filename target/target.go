package target

import (
	"bufio"
	"github.com/bitly/go-simplejson"
	"io"
	"io/ioutil"
	"os"
)

type Target interface {
	GetFilePath() (path string, err error)
}

func Write(t Target, content []byte) (n int, err error) {
	path, err := t.GetFilePath()
	if err != nil {
		return
	}

	err = os.Remove(path)
	if err != nil {
		return
	}

	file, err := os.Create(path)
	if err != nil {
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	n, err = writer.Write(content)
	return
}

func Backup(t Target, ext string) (backupPath string, err error) {
	filepath, err := t.GetFilePath()
	if err != nil {
		return "", err
	}
	backupPath = filepath + "." + ext
	_, err = copyFile(filepath, backupPath)
	return
}

func Restore(t Target, ext string) (backupPath string, err error) {
	filepath, err := t.GetFilePath()
	if err != nil {
		return "", err
	}
	backupPath = filepath + "." + ext

	_, err = copyFile(backupPath, filepath)
	if err != nil {
		return
	}

	err = os.Remove(backupPath)
	return
}

func copyFile(src, dst string) (n int64, err error) {
	srcFile, err := os.OpenFile(src, os.O_RDONLY, os.FileMode(0400))
	if err != nil {
		return
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return
	}
	defer dstFile.Close()

	reader := bufio.NewReader(srcFile)
	writer := bufio.NewWriter(dstFile)
	defer writer.Flush()

	n, err = io.Copy(writer, reader)
	if err != nil {
		return
	}
	return
}

func readJson(path string) (j *simplejson.Json, err error) {
	file, err := os.OpenFile(path, os.O_RDONLY, os.FileMode(0400))
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	return simplejson.NewJson(content)
}
