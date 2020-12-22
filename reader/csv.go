package reader

import (
	"encoding/csv"
	"io"
	"os"
)

// CsvFromFileWithBreak 从文件中读取，允许中断 true为中断信号
func CsvFromFileWithBreak(name string, f func(rcs []string) bool) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()
	return CsvWithBreak(file, f)
}

// CsvWithBreak 从csv io.Reader 中读取，允许中断 true为中断信号
func CsvWithBreak(r io.Reader, f func(rcs []string) bool) error {
	cr := csv.NewReader(r)
	for {
		rcs, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		f(rcs)
	}
	return nil
}

// Csv 读取解析csv内容
func Csv(r io.Reader, f func(rcs []string)) error {
	cr := csv.NewReader(r)
	for {
		rcs, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		f(rcs)
	}
	return nil
}
