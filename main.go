package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type exitCode int

const (
	exitCodeOk exitCode = iota
	exitCodeErrArgs
)

func main() {

	//runtime.GOOS
	/*
		code, err := run(os.Args[1:])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(int(code))
	*/
	createDir("test")
	zipDownload("https://code.visualstudio.com/sha/download?build=stable&os=win32-x64-archive", "test")
}

/*
func run(args []string) (exitCode, error) {


	if err != nil {
		return exitCodeErrArgs, err
	}
	return exitCodeOk, nil

}
*/
//https://github.com/oneclick/rubyinstaller2/releases/download/RubyInstaller-3.1.3-1/rubyinstaller-devkit-3.1.3-1-x64.exe

func vscodeInstall(osType string) error {
	return nil
}

func rubyInstall(osType string) error {
	url := ""
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch file %s: %w", url, err)
	}
	defer resp.Header.Clone()

	return nil

}

func json2map(fetchedJson string) (map[string]interface{}, error) {

	var mapData map[string]interface{}
	err := json.Unmarshal([]byte(fetchedJson), &mapData)
	if err != nil {
		return nil, fmt.Errorf("failed to json convert: %w", err)
	}

	return mapData, nil
}

func createDir(destDir string) error {
	//ファイル存在確認
	_, err := os.Stat(destDir)
	if err == nil {
		//ファイルが存在していた場合削除する
		err := os.RemoveAll(destDir)
		if err != nil {
			return fmt.Errorf("failed to remove directory %s: %w", destDir, err)
		}
	}
	err = os.Mkdir(destDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %w", destDir, err)
	}

	return nil

}

func zipDownload(url string, destDir string) error {

	//zipDonwload
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch file %s: %w", url, err)
	}
	defer resp.Header.Clone()

	//メモリ上にzipBinaryを展開
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to open file in memory :%w", err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(bodyBytes), int64(len(bodyBytes)))
	if err != nil {
		return fmt.Errorf("failed to create zipReader :%w", err)
	}

	for _, file := range zipReader.File {
		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		destPath := filepath.Join(destDir, file.Name)

		//出力がディレクトリの場合
		if file.FileInfo().IsDir() {
			err = os.MkdirAll(destPath, file.Mode())
			if err != nil {
				return fmt.Errorf("failed to create directory %s :%w", destPath, err)
			}
			continue
		}

		//出力がファイルの場合
		dstFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return fmt.Errorf("failed to open file %s :%w", destPath, err)
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, fileReader)
		if err != nil {
			return fmt.Errorf("failed to write file %s :%w", destPath, err)
		}

	}
	return nil

}
