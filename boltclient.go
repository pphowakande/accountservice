package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/boltdb/bolt"
)

type IBoltClient interface {
	OpenBoltDb()
	CreateAccount(account Account) (Account, error)
	QueryAccount(accountId string) (Account, error)
	Seed()
}

// Real implementation
type BoltClient struct {
	boltDB *bolt.DB
}

func (bc *BoltClient) OpenBoltDb() {
	var err error
	bc.boltDB, err = bolt.Open("accounts.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (bc *BoltClient) QueryAccount(accountId string) (Account, error) {
	// Allocate an empty Account instance we'll let json.Unmarhal populate for us in a bit.
	account := Account{}

	// Read an object from the bucket using boltDB.View
	err := bc.boltDB.View(func(tx *bolt.Tx) error {
		// Read the bucket from the DB
		b := tx.Bucket([]byte("AccountBucket"))

		// Read the value identified by our accountId supplied as []byte
		accountBytes := b.Get([]byte(accountId))
		if accountBytes == nil {
			return fmt.Errorf("No account found for " + accountId)
		}
		// Unmarshal the returned bytes into the account struct we created at
		// the top of the function
		json.Unmarshal(accountBytes, &account)

		// Return nil to indicate nothing went wrong, e.g no error
		return nil
	})
	// If there were an error, return the error
	if err != nil {
		return Account{}, err
	}
	// Return the Account struct and nil as error.
	return account, nil
}

func (bc *BoltClient) CreateAccount(account Account) (Account, error) {

	// Serialize the struct to JSON
	jsonBytes, _ := json.Marshal(account)

	// Write the data to the AccountBucket
	bc.boltDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("AccountBucket"))
		err := b.Put([]byte(account.Id), jsonBytes)
		return err
	})

	return account, nil
}

// Start seeding accounts
func (bc *BoltClient) Seed() {
	bc.initializeBucket()
	bc.seedAccounts()
}

// Creates an "AccountBucket" in our BoltDB. It will overwrite any existing bucket of the same name.
func (bc *BoltClient) initializeBucket() {
	bc.boltDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("AccountBucket"))
		if err != nil {
			return fmt.Errorf("create bucket failed: %s", err)
		}
		return nil
	})
}

// Seed (n) make-believe account objects into the AcountBucket bucket.
func (bc *BoltClient) seedAccounts() {

	total := 10
	for i := 0; i < total; i++ {

		// Generate a key 10000 or larger
		key := strconv.Itoa(10 + i)

		// Create an instance of our Account struct
		acc := Account{
			Id:   key,
			Name: "Person_" + strconv.Itoa(i),
		}

		// Serialize the struct to JSON
		jsonBytes, _ := json.Marshal(acc)

		// Write the data to the AccountBucket
		bc.boltDB.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("AccountBucket"))
			err := b.Put([]byte(key), jsonBytes)
			return err
		})
	}
	//fmt.Printf("Seeded %v fake accounts...\n", total)
}
