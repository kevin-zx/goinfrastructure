package badgercache

import (
	"os"
	"path"
	"testing"
	"time"
)

func Test_badgercache_Save(t *testing.T) {
	bc, err := NewBadgerCache(path.Join(os.TempDir(), "badgerCacheTest/"))
	if err != nil {
		panic(err)
	}
	defer bc.Close()

	eds := []EntryData{
		{Key: "stk1", Data: []byte("stv1"), TTL: time.Second * 1},
	}
	err = bc.Saves(eds)
	if err != nil {
		t.Fatal(err)
	}
	ed, miskes, err := bc.Gets([]string{"stk1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(miskes) > 0 {
		t.Errorf("miskes shoud be nil\n")
	}
	if string(ed[0].Data) != string(eds[0].Data) {
		t.Errorf("get value not equal save value, get: %s, save:%s", string(ed[0].Data), string(eds[0].Data))
	}
	time.Sleep(time.Second * 2)
	ed, miskes, err = bc.Gets([]string{"stk1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(miskes) != 1 {
		t.Errorf("miskes len shoud be 1\n")
	}
	if len(ed) > 0 {
		t.Errorf("eds shoud be nil")
	}
}
