package surge

import (
	"encoding/json"
	"fmt"
	"os"

	"log"

	"github.com/xujiajun/nutsdb"
)

const entriesBucketName = "entriesBucket"
const settingBucketName = "settingsBucket"

var db *nutsdb.DB

//InitializeDb initializes db
func InitializeDb() *nutsdb.DB {
	var err error
	opt := nutsdb.DefaultOptions

	fmt.Println(os.Getenv("HOME") + string(os.PathSeparator) + ".ytd")
	opt.Dir = os.Getenv("HOME") + string(os.PathSeparator) + ".ytd"
	db, err = nutsdb.Open(opt)
	if err != nil {
		log.Panic(err)
	}

	// dbGetAllEntries()

	return db

}

//CloseDb .
func CloseDb() {
	db.Close()
}

func DbWriteEntry(Key string, value interface{}) error {
	err := db.Update(
		func(tx *nutsdb.Tx) error {

			keyBytes := []byte(Key)
			valueBytes, err := json.Marshal(value)
			if err != nil {
				return err
			}

			if err := tx.Put(entriesBucketName, keyBytes, valueBytes, 0); err != nil {
				return err
			}
			return nil
		})
	return err
}

//DbWriteSetting .
func DbWriteSetting(Name string, value string) error {
	err := db.Update(
		func(tx *nutsdb.Tx) error {

			keyBytes := []byte(Name)
			valueBytes := []byte(value)

			if err := tx.Put(settingBucketName, keyBytes, valueBytes, 0); err != nil {
				return err
			}
			return nil
		})
	return err
}

//DbReadSetting .
func DbReadSetting(Name string) (string, error) {
	result := ""
	key := []byte(Name)

	if err := db.View(
		func(tx *nutsdb.Tx) error {
			bytes, err := tx.Get(settingBucketName, key)
			if err != nil {
				return err
			}

			result = string(bytes.Value)

			return err
		}); err != nil {
		return "", err
	}

	return result, nil
}
