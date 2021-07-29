package db

import (
	"encoding/json"
	"fmt"
	"os"

	"log"

	"github.com/xujiajun/nutsdb"

	. "ytd/models"
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

	return db

}

//CloseDb .
func CloseDb() {
	db.Close()
}

func DbGetAllEntries() []GenericEntry {
	data := []GenericEntry{}

	if err := db.View(
		func(tx *nutsdb.Tx) error {
			entries, err := tx.GetAll(entriesBucketName)
			if err != nil {
				return err
			}

			for _, entry := range entries {

				genericEntry := &GenericEntry{}
				json.Unmarshal(entry.Value, genericEntry)
				data = append(data, *genericEntry)
			}

			return nil
		}); err != nil {
		log.Println(err)
	} else {
		return data
	}
	return data
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

func DbDeleteEntry(Key string) error {
	err := db.Update(
		func(tx *nutsdb.Tx) error {
			key := []byte(Key)
			if err := tx.Delete(entriesBucketName, key); err != nil {
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

func DbSaveSettingBoolValue(name string, val bool) error {
	var v string
	if val {
		v = "1"
	} else {
		v = "0"
	}
	return DbWriteSetting(name, v)
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

func DbReadSettingBoolValue(name string) (bool, error) {
	val, err := DbReadSetting(name)
	if err != nil {
		return false, err
	}
	if val == "1" {
		return true, nil
	}
	return false, nil
}
