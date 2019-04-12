package tidbtest

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type fileReader struct {
	filePaths []string
}

//NewFileReader 从文件中读取数据
func NewFileReader(filePaths []string) (Reader, error) {
	for _, path := range filePaths {
		if _, err := os.Stat(path); err != nil {
			return nil, fmt.Errorf("get file error %s", err)
		}
	}
	return &fileReader{
		filePaths: filePaths,
	}, nil
}

//TODO:考虑带buffer的异步读取方式，以及是否需要保留数据
func (f *fileReader) Read() (map[string]string, error) {
	if f == nil {
		return nil, errors.New("reader is nil")
	}
	ss := make(map[string]string, 0)
	for _, path := range f.filePaths {
		f, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("open file %s error %s", path, err)
		}

		bs, err := ioutil.ReadAll(f)

		if e := f.Close(); e != nil {
			log.Printf("close file %s error %s", path, e)
		}
		if err != nil {
			return nil, fmt.Errorf("read file %s error %s", path, err)
		}
		ss[path] = string(bs)
	}
	return ss, nil
}
