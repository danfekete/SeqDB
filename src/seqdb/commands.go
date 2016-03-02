package seqdb
import (
	"github.com/boltdb/bolt"
	"fmt"
	"bytes"
	"encoding/binary"
	"log"
	"errors"
)

type SeqPointer struct {
	BucketName string
	SequenceName string
	lock *BucketLocks
}

func (s SeqPointer) Set(v uint64) {

	Db.Update(func(tx *bolt.Tx) error {

		bucket, err := tx.CreateBucketIfNotExists([]byte(s.BucketName))

		if err != nil {
			log.Printf("%s bucket cannot be created: %v\r\n", s.BucketName, err)
			return err
		}

		buf := new(bytes.Buffer)
		err = binary.Write(buf, binary.LittleEndian, v)

		if err != nil {
			log.Printf("Cannot convert value %d to binary: %v\r\n", v, err)
			return err
		}

		err = bucket.Put([]byte(s.SequenceName), buf.Bytes())

		if err != nil {
			log.Printf("Cannot put data to sequence %s: %v\r\n", s.SequenceName, err)
			return err
		}

		log.Printf("[%s %s] <- %d", s.BucketName, s.SequenceName, v)
		return nil
	})

}

func (s SeqPointer) Get() uint64 {

	var val uint64;

	Db.View(func(tx *bolt.Tx) error {

		// Try loading the bucket
		bucket := tx.Bucket([]byte(s.BucketName))

		// Bucket is not found
		if bucket == nil {
			log.Printf("Bucket %s not found", s.BucketName)
			return errors.New("Bucket not found")
		}

		// Read the data from the sequence
		b := bucket.Get([]byte(s.SequenceName))

		if b == nil {
			log.Printf("Sequence %s is not available\r\n", s.SequenceName)
			return errors.New("Sequence not available for reading")
		}

		buf := bytes.NewReader(b)
		err := binary.Read(buf, binary.LittleEndian, &val)

		if err != nil {
			log.Printf("Cannot convert to binary while reading: %v\r\n", err)
			return errors.New("Binary read fail")
		}

		log.Printf("[%s %s] -> %d", s.BucketName, s.SequenceName, val)
		return nil
	})

	return val
}

func (s SeqPointer) Inc() uint64 {
	var val uint64

	Db.Batch(func(tx *bolt.Tx) error {
		// Get the given sequence from bucket
		bucket := tx.Bucket([]byte(s.BucketName))

		if bucket == nil {
			// Bucket is not found
			log.Printf("Bucket %s not found", s.BucketName)
			return fmt.Errorf("Bucket %s not found", s.BucketName)
		}

		// Convert the binary to Uint64
		bufR := bytes.NewReader(bucket.Get([]byte(s.SequenceName)))
		err := binary.Read(bufR, binary.LittleEndian, &val)

		if err != nil {
			log.Printf("Cannot convert to binary while reading: %v\r\n", err)
			return err
		}

		// Increment variable
		val++

		// Put it back
		bufW := new(bytes.Buffer)
		err = binary.Write(bufW, binary.LittleEndian, val)

		if err != nil {
			log.Printf("Cannot convert value %d to binary: %v\r\n", val, err)
			return err
		}

		err = bucket.Put([]byte(s.SequenceName), bufW.Bytes())

		if err != nil {
			log.Printf("Cannot put data to sequence %s: %v\r\n", s.SequenceName, err)
			return err
		}

		log.Printf("[%s %s] -> %d", s.BucketName, s.SequenceName, val)
		return nil
	})

	return val
}

func (s SeqPointer) Lock() {
	s.lock.WaitAndSet(s.BucketName)
}

func (s SeqPointer) Unlock() {
	s.lock.RemoveLock(s.BucketName)
}