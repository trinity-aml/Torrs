package db

import (
	"encoding/binary"
	"encoding/hex"
	bolt "go.etcd.io/bbolt"
	"log"
	"path/filepath"
	"regexp"
	"time"
	"torrsru/global"
	"torrsru/models/fdb"
)

var (
	db *bolt.DB
	re = regexp.MustCompile(`.+((19|20)\d\d)`)
)

func Init() {
	d, err := bolt.Open(filepath.Join(global.PWD, "torrents.db"), 0o666, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		log.Fatalln("Error open db", err)
		return
	}
	db = d

	err = initIndex()
	if err != nil {
		log.Fatalln("Error open indexTorrentTitle", err)
		return
	}
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

func saveTorrents(cols []*fdb.Collection) error {
	const chunkSize = 500

	batch := indexTorrentTitle.NewBatch()
	chunk := make([]*fdb.Torrent, 0, chunkSize)

	flush := func() error {
		if len(chunk) == 0 {
			return nil
		}
		err := db.Update(func(tx *bolt.Tx) error {
			torrsb, err := tx.CreateBucketIfNotExists([]byte("Torrents"))
			if err != nil {
				return err
			}
			for _, torr := range chunk {
				hash := torr.GetUnique()
				if tb := torrsb.Bucket(hash); tb != nil {
					ttmp, err := readTorrent(torrsb, hash)
					if err != nil {
						return err
					}
					torr = combineTorrents([]*fdb.Torrent{ttmp, torr})
				}
				if err := saveTorrent(torrsb, torr); err != nil {
					return err
				}
				if err := batch.Index(hex.EncodeToString(hash), torr.Title); err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
		chunk = chunk[:0]
		if err := indexTorrentTitle.Batch(batch); err != nil {
			return err
		}
		batch.Reset()
		return nil
	}

	for _, col := range cols {
		for _, torr := range col.Value.Torrents {
			chunk = append(chunk, torr)
			if len(chunk) >= chunkSize {
				if err := flush(); err != nil {
					return err
				}
			}
		}
	}
	return flush()
}
