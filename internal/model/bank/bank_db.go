package bank

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

// users banks database
type usersBanksDatabase struct {
	dbFilePath string
	mu         *sync.Mutex // guards usersBanks
	usersBanks map[string][]*Bank
}

// creates a new users banks storage
func newBankDB(dbFilePath string) (*usersBanksDatabase, error) {
	if strings.TrimSpace(dbFilePath) == "" {
		dbFilePath = "banks.db"
	}

	// create db file if it doesn't exist and close it afterwards
	f, err := os.OpenFile(dbFilePath, os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	err = f.Close()
	if err != nil {
		return nil, err
	}

	bankDB := &usersBanksDatabase{
		dbFilePath: dbFilePath,
		mu:         &sync.Mutex{},
		usersBanks: make(map[string][]*Bank, 0),
	}

	// load users to usersDB.users map
	return bankDB, bankDB.load()
}

// Add bank to a user list of banks
func (bankDB *usersBanksDatabase) addUserBank(userID string, bank *Bank) error {
	bankDB.mu.Lock()
	defer bankDB.mu.Unlock()

	userBanks, ok := bankDB.usersBanks[userID]
	// handle case where user user banks is nil or has no record
	if !ok || userBanks == nil {
		bankDB.usersBanks[userID] = make([]*Bank, 0)
		userBanks = bankDB.usersBanks[userID]
	}

	userBanks = append(userBanks, bank)

	return bankDB.save()
}

// retrieves user list of bank
func (bankDB *usersBanksDatabase) getUserBanks(userID string) ([]*Bank, error) {
	bankDB.mu.Lock()
	defer bankDB.mu.Unlock()

	banks, ok := bankDB.usersBanks[userID]
	if !ok {
		return nil, errors.New("you have no banks yet")
	}

	return banks, nil
}

// retrieves user bank
func (bankDB *usersBanksDatabase) getUserBank(userID, bankID string) (*Bank, error) {
	bankDB.mu.Lock()
	defer bankDB.mu.Unlock()

	banks, ok := bankDB.usersBanks[userID]
	if !ok {
		return nil, errors.New("you have no banks yet")
	}

	for _, userBank := range banks {
		if userBank.ID == bankID {
			return userBank, nil
		}
	}

	return nil, errors.Errorf("user has no bank with id: %s", bankID)
}

// removes a user bank from their list of banks
func (bankDB *usersBanksDatabase) removeUserBank(userID, bankID string) error {
	bankDB.mu.Lock()
	defer bankDB.mu.Unlock()

	banks, ok := bankDB.usersBanks[userID]
	if !ok {
		return errors.New("you have no banks yet")
	}

	for index, userBank := range banks {
		if userBank.ID == bankID {
			// Remove the bank using append
			banks = append(banks[:index], banks[index+1:]...)
			return bankDB.save()
		}
	}

	return errors.Errorf("user has no bank with id: %s", bankID)
}

// saves current state of usersBanks into a file db. Assumes mutex is locked
func (bankDB *usersBanksDatabase) save() error {

	// marshaling it to json first
	usersBanksBytes, err := json.Marshal(bankDB.usersBanks)
	if err != nil {
		return errors.Wrap(err, "failed to marshal users banks")
	}

	// We store the credentials in file for simplicity. We could replace it with mysql db.
	dbFile, err := os.OpenFile(bankDB.dbFilePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return errors.Wrap(err, "failed to open file")
	}

	_, err = dbFile.Write(usersBanksBytes)
	if err != nil {
		return errors.Wrap(err, "failed to write users to file")
	}

	return nil
}

// loads users banks from a file and initialize to map
func (bankDB *usersBanksDatabase) load() error {
	bankDB.mu.Lock()
	defer bankDB.mu.Unlock()

	usersBanksBytes := make([]byte, 0)

	usersBanksBytes, err := ioutil.ReadFile(bankDB.dbFilePath)
	if err != nil {
		return errors.Wrap(err, "failed read from file")
	}

	usersBanksBytes = bytes.TrimSpace(usersBanksBytes)

	if len(usersBanksBytes) == 0 {
		bankDB.usersBanks = make(map[string][]*Bank, 0)
		return nil
	}

	usersBanks := make(map[string][]*Bank, 0)

	// We unmarshal the content of the file and save it to a users map
	err = json.Unmarshal(usersBanksBytes, &usersBanks)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal data")
	}

	// assign users banks map
	bankDB.usersBanks = usersBanks

	return nil
}
