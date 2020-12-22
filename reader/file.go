package reader

import (
	"bufio"
	"io"
	"os"
)

// ReadLinesFromFileWithBreak 从文件中读取行 允许中断 true为中断信号
func ReadLinesFromFileWithBreak(name string, f func(line []byte) bool) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()
	return ReadLinesWithBreak(file, f)
}

// ReadLines 按行读取
func ReadLines(r io.Reader, f func(line []byte)) error {
	br := bufio.NewReader(r)
	for {
		data, err := br.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		f(data)
	}
	return nil
}

// ReadLinesWithBreak 按行读取 允许中断 true为中断信号
func ReadLinesWithBreak(r io.Reader, f func(line []byte) bool) error {
	br := bufio.NewReader(r)
	for {
		data, err := br.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		b := f(data)
		if b {
			break
		}

	}
	return nil
}
