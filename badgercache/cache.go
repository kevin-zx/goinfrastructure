package badgercache

import (
	"strings"
	"time"

	badger "github.com/dgraph-io/badger/v2"
)

// EntryData 是用来存取用的
type EntryData struct {
	Key  string
	Data []byte
	TTL  time.Duration
	Meta byte
}

// BadgerCache 是一个基于Badger作为底层的文件缓存
type BadgerCache interface {
	Saves([]EntryData) error
	SavesWithPrefix(eds []EntryData, prefix string) error
	Gets(keys []string) ([]EntryData, []string, error)
	GetsWithPrefix(keys []string, prefix string) ([]EntryData, []string, error)
	DeleteWithPrefix(keys []string, prefix string) error
	Delete(keys []string) error
	Close() error
}

// NewBadgerCache 新建badgercache
func NewBadgerCache(dbPath string) (BadgerCache, error) {
	db, err := badger.Open(badger.DefaultOptions(dbPath))
	if err != nil {
		return nil, err
	}
	dc := badgercache{db: db}
	go dc.periodGC()
	return &dc, nil
}

//bc *badgercache badgercache.BadgerCache
type badgercache struct {
	db *badger.DB
}

func (bc *badgercache) DeleteWithPrefix(keys []string, prefix string) error {
	return bc.Delete(assemblePrefixFromKeys(keys, prefix))
}

func (bc *badgercache) SavesWithPrefix(eds []EntryData, prefix string) error {
	eds = assemblePrefixFromEntryData(eds, prefix)
	return bc.Saves(eds)
}

func assemblePrefixFromKeys(ks []string, prefix string) []string {
	var keyWithPrefix []string
	for _, k := range ks {
		keyWithPrefix = append(keyWithPrefix, prefix+"_"+k)
	}
	return keyWithPrefix
}

func disassemblePrefixFromKeys(ks []string, prefix string) []string {
	var keyWithoutPrefix []string
	for _, k := range ks {
		keyWithoutPrefix = append(keyWithoutPrefix, strings.Replace(k, prefix+"_", "", 1))
	}
	return keyWithoutPrefix

}

func assemblePrefixFromEntryData(eds []EntryData, prefix string) []EntryData {
	for i := range eds {
		eds[i].Key = prefix + "_" + eds[i].Key
	}
	return eds
}

func disassemblePrefixFromEntryData(eds []EntryData, prefix string) []EntryData {
	for i := range eds {
		eds[i].Key = strings.Replace(eds[i].Key, prefix+"_", "", 1)
		// keyWithoutPrefix = append(keyWithoutPrefix, strings.Replace(k, prefix+"_", "", 1))
	}
	return eds

}

func (bc *badgercache) GetsWithPrefix(keys []string, prefix string) ([]EntryData, []string, error) {
	kwp := assemblePrefixFromKeys(keys, prefix)
	eds, ss, err := bc.Gets(kwp)
	if err != nil {
		return nil, nil, err
	}
	miskeys := disassemblePrefixFromKeys(ss, prefix)
	eds = disassemblePrefixFromEntryData(eds, prefix)
	return eds, miskeys, nil
}

func (bc *badgercache) Saves(eds []EntryData) error {
	err := bc.db.Update(func(txn *badger.Txn) error {
		for _, ed := range eds {
			e := badger.NewEntry([]byte(ed.Key), ed.Data)
			if ed.TTL > 0 {
				e = e.WithTTL(ed.TTL)
			}
			if ed.Meta != 0 {
				e = e.WithMeta(ed.Meta)
			}
			err := txn.SetEntry(e)
			return err
		}
		return nil
	})
	return err
}

func (bc *badgercache) Delete(keys []string) error {
	err := bc.db.Update(func(txn *badger.Txn) error {
		for _, k := range keys {
			err := txn.Delete([]byte(k))
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err

}

func (bc *badgercache) Gets(keys []string) ([]EntryData, []string, error) {
	missKeys := []string{}
	eds := []EntryData{}
	err := bc.db.View(func(txn *badger.Txn) error {
		for _, key := range keys {
			item, err := txn.Get([]byte(key))
			if err == badger.ErrKeyNotFound {
				missKeys = append(missKeys, key)
				continue
			}
			if err != nil {
				return err
			}
			ed := EntryData{Key: key, Data: []byte{}}
			item.Value(func(val []byte) error {
				ed.Data = append(ed.Data, val...)
				return nil
			})
			eds = append(eds, ed)
		}
		return nil
	})
	return eds, missKeys, err
}

func (bc *badgercache) Close() error {
	return bc.db.Close()
}

func (bc *badgercache) periodGC() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
	again:
		err := bc.db.RunValueLogGC(0.7)
		if err == nil {
			goto again
		}
	}
}
