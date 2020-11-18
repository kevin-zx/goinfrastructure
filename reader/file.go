package reader

import (
	"bufio"
	"io"
)

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
