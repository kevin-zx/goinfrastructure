package reader

import (
	"encoding/csv"
	"io"
)

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
