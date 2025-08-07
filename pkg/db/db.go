package db

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"

	bolt "go.etcd.io/bbolt"
)

const DBName string = "projects"
const DBTestName string = "projects_test"

var (
	ErrOpenDB          = errors.New("error opening database")
	ErrCreateBucket    = errors.New("error creating bucket")
	ErrWriteToDB       = errors.New("error writing to database")
	ErrBucketNotFound  = errors.New("bucket not found")
	ErrProjectNotFound = errors.New("project not found")
	ErrDeleteFromDB    = errors.New("error deleting from database")
	ErrKeyNotFound     = errors.New("key not found in db")
	ErrListAllRecords  = errors.New("error listing all records")
	ErrClearDB         = errors.New("error clearing database")
	ErrDBNameEmpty     = errors.New("dbname cannot be empty")
)

// WriteToDB writes the data to the specified bucket in the database
func WriteToDB(dbname string, data map[string]string, bucketName string) error {
	dbLoc, err := GetDBLoc(dbname)
	if err != nil {
		log.Printf("%v : %v \n", ErrOpenDB, err)
		return errors.Join(ErrOpenDB, err)
	}
	db, _ := bolt.Open(dbLoc, 0o600, nil) // create the database if it doesn't exist then open it
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return errors.Join(ErrCreateBucket, err)
		}
		for k, v := range data {
			err = bucket.Put([]byte(k), []byte(v))
			if err != nil {
				return errors.Join(ErrWriteToDB, err)
			}
		}
		return nil
	})
	return err
}

func DeleteFromDb(dbname, key, bucketName string) error {
	dbLoc, err := GetDBLoc(dbname)
	if err != nil {
		log.Printf("%v : %v \n", ErrOpenDB, err)
		return errors.Join(ErrOpenDB, err)
	}
	db, _ := bolt.Open(dbLoc, 0o600, nil) // create the database if it doesn't exist then open it
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return ErrBucketNotFound
		}
		return bucket.Delete([]byte(key))
	})
	return err
}

// getDBLoc returns the path to the database file, creating the directory if it doesn't exist
func GetDBLoc(dbname string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	if dbname == "" {
		return "", ErrDBNameEmpty
	}
	dbname = dbname + ".db"
	dbPath := filepath.Join(usr.HomeDir, ".local", "share", "pman", dbname)
	if _, err := os.Stat(filepath.Dir(dbPath)); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(dbPath), 0o755)
		if err != nil {
			return "", err
		}
	}
	return dbPath, nil
}

// GetRecord returns the value of the key from the specified bucket, and error if it does not exist
func GetRecord(dbname, key, bucketName string) (string, error) {
	var rec string
	dbLoc, err := GetDBLoc(dbname)
	if err != nil {
		log.Printf("%v : %v \n", ErrOpenDB, err)
		return "", err
	}
	db, _ := bolt.Open(dbLoc, 0o600, nil)

	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			// log.Println("bucket not found error", bucketName,"for key")
			return ErrBucketNotFound
		}
		v := bucket.Get([]byte(key))
		if v == nil {
			// log.Println("key not found error", key)
			return ErrKeyNotFound
		}
		rec = string(v)
		return nil
	})
	if err != nil {
		return "", err
	}
	return rec, nil
}

var GetAllRecords = func(dbname, bucketName string) (map[string]string, error) {
	dbLoc, err := GetDBLoc(dbname)
	if err != nil {
		log.Printf("%v : %v \n", ErrOpenDB, err)

		return map[string]string{}, err
	}
	db, _ := bolt.Open(dbLoc, 0o600, nil)
	defer db.Close()
	records := make(map[string]string)
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			fmt.Print("Database not found \nThis could be because no project dir has been initialized yet \n")
			return ErrBucketNotFound
		}
		err := bucket.ForEach(func(k, v []byte) error {
			records[string(k)] = string(v)
			return nil
		})
		if err != nil {
			return errors.Join(ErrListAllRecords, err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return records, nil
}

// UpdateRec updates the value of the key in the specified bucket, usually used to update the status of a project
func UpdateRec(dbname, key, val, bucketName string) error {
	dbLoc, err := GetDBLoc(dbname)
	if err != nil {
		log.Print(err)
		return err
	}
	db, _ := bolt.Open(dbLoc, 0o600, nil)

	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return ErrBucketNotFound
		}
		v := bucket.Get([]byte(key))
		if v == nil {
			return ErrProjectNotFound
		}
		err := bucket.Put([]byte(key), []byte(val))
		return err
	})
	return err
}

func DeleteDb(dbname string) error {
	dbLoc, err := GetDBLoc(dbname)
	if err != nil {
		return err
	}
	return os.Remove(dbLoc)
}
