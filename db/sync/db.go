package sync

import (
	"encoding/binary"
	"encoding/json"
	bolt "go.etcd.io/bbolt"
	"log"
	"path/filepath"
	"time"
	"torrsru/db/utils"
	"torrsru/models/fdb"
	"torrsru/web/global"
)

var (
	db *bolt.DB
)

func Init() {
	d, err := bolt.Open(filepath.Join(global.PWD, "torrents.db"), 0o666, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		log.Fatalln("Error open db", err)
		return
	}
	db = d
}

func GetFileTime() int64 {
	var ft int64 = -1
	err := db.View(func(tx *bolt.Tx) error {
		sets := tx.Bucket([]byte("Settings"))
		if sets == nil {
			return nil
		}
		b := sets.Get([]byte("FileTime"))
		if b != nil {
			ft = int64(binary.LittleEndian.Uint64(b))
		}
		return nil
	})
	if err != nil {
		log.Println("Error get from db:", err)
	}
	return ft
}

func SetFileTime(fileTime int64) error {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(fileTime))
	return db.Update(func(tx *bolt.Tx) error {
		sets, err := tx.CreateBucketIfNotExists([]byte("Settings"))
		if err != nil {
			return err
		}
		return sets.Put([]byte("FileTime"), b)
	})
}

func ListNames() []string {
	var ret []string
	err := db.View(func(tx *bolt.Tx) error {
		torrsb := tx.Bucket([]byte("Torrents"))
		if torrsb == nil {
			return nil
		}
		return torrsb.ForEach(func(k, v []byte) error {
			ret = append(ret, string(k))
			return nil
		})
	})
	if err != nil {
		log.Println("Error get from db:", err)
	}

	return ret
}

func ListTitles() []string {
	var ret []string
	err := db.View(func(tx *bolt.Tx) error {
		torrsb := tx.Bucket([]byte("Indexes"))
		if torrsb == nil {
			return nil
		}
		torrsb = torrsb.Bucket([]byte("ByTitle"))
		if torrsb == nil {
			return nil
		}
		return torrsb.ForEach(func(k, v []byte) error {
			ret = append(ret, string(k))
			return nil
		})
	})
	if err != nil {
		log.Println("Error get from db:", err)
	}

	return ret
}

func GetTorrentsByName(name string) []*fdb.Torrent {
	var ret []*fdb.Torrent
	err := db.View(func(tx *bolt.Tx) error {
		torrsb := tx.Bucket([]byte("Torrents"))
		if torrsb == nil {
			return nil
		}
		tnb := torrsb.Bucket([]byte(name))
		if tnb == nil {
			return nil
		}

		return tnb.ForEach(func(k, v []byte) error {
			var t *fdb.Torrent
			err := json.Unmarshal(v, &t)
			if err != nil {
				return err
			}
			ret = append(ret, t)
			return nil
		})
	})
	if err != nil {
		log.Println("Error get from db:", err)
	}

	return ret
}

func GetTorrentsByTitle(title string) []*fdb.Torrent {
	var ret []*fdb.Torrent
	index := map[string]string{}
	err := db.View(func(tx *bolt.Tx) error {
		torrsb := tx.Bucket([]byte("Indexes"))
		if torrsb == nil {
			return nil
		}
		torrsb = torrsb.Bucket([]byte("ByTitle"))
		if torrsb == nil {
			return nil
		}

		tr := torrsb.Bucket([]byte(title))
		if tr != nil {
			tr.ForEach(func(name, hash []byte) error {
				index[string(hash)] = string(name)
				return nil
			})
		}
		return nil
	})

	if err == nil {
		for hash, name := range index {
			torrs := GetTorrentsByName(name)
			for _, torr := range torrs {
				if getHash(torr.Magnet) == hash {
					ret = append(ret, torr)
				}
			}
		}
	}

	if err != nil {
		log.Println("Error get from db:", err)
	}

	return ret
}

func saveTorrent(cols []*fdb.Collection) error {
	return db.Update(func(tx *bolt.Tx) error {
		torrsb, err := tx.CreateBucketIfNotExists([]byte("Torrents"))
		if err != nil {
			return err
		}

		index, err := tx.CreateBucketIfNotExists([]byte("Indexes"))
		if err != nil {
			return err
		}

		index, err = index.CreateBucketIfNotExists([]byte("ByTitle"))
		if err != nil {
			return err
		}

		for _, col := range cols {
			tnb, err := torrsb.CreateBucketIfNotExists([]byte(col.Key))
			if err != nil {
				return err
			}
			for _, torr := range col.Value.Torrents {
				//save torrent
				hash := getHash(torr.Magnet)
				if hash == "" {
					continue
				}
				buf, err := json.Marshal(torr)
				if err != nil {
					return err
				}
				err = tnb.Put([]byte(hash), buf)
				if err != nil {
					return err
				}
				//create index
				//index: torrent title->collection name->hash
				titlid, err := index.CreateBucketIfNotExists([]byte(utils.ClearStr(torr.Title)))
				if err != nil {
					return err
				}
				err = titlid.Put([]byte(col.Key), []byte(hash))
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}
