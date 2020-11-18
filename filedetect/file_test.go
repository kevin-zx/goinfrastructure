package filedetect

import (
	"fmt"
	"os"
	"path"
	"testing"
	"time"
)

func TestFileIsExist(t *testing.T) {
	notExistFile := path.Join(os.TempDir(), fmt.Sprintf("%d.notexistfile", time.Now().Unix()))
	e, err := FileIsExist(notExistFile)
	if err != nil {
		t.Fatal(err)
	}
	if e {
		fmt.Println("not exist file detect wrong")
		t.Fail()
	}

	existFile := path.Join(os.TempDir(), fmt.Sprintf("%d.existfile", time.Now().Unix()))
	f, err := os.Create(existFile)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	defer func() {
		os.Remove(existFile)
	}()
	e, err = FileIsExist(existFile)
	if err != nil {
		t.Fatal(err)
	}
	if !e {
		fmt.Println("exist file detect wrong")
		t.Fail()
	}
}
