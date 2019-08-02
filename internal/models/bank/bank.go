package bank

import (
	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"strings"
)

var db *usersBanksDatabase

func init() {
	bankDB, err := newBankDB("banks.db")
	if err != nil {
		logrus.Fatalln(err)
	}

	db = bankDB
}

// Bank is a bank resource ...
type Bank struct {
	ID                 string
	FullName           string
	ShortName          string
	LogoURL            string
	WebsiteURL         string
	SwiftBIC           string
	NationalIdentifier string
	BankRouting        struct {
		Scheme  string
		Address string
	}
}

// AddBank adds bank to user list of banks
func AddBank(userID string, bank *Bank) error {
	if bank == nil {
		return errors.New("cannot add nil bank")
	}
	return db.addUserBank(userID, bank)
}

// RemoveBank removes a bank from user list of banks
func RemoveBank(userID, bankID string) error {
	// Validate input
	switch {
	case strings.TrimSpace(userID) == "":
		return errors.New("missing user id")
	case strings.TrimSpace(bankID) == "":
		return errors.New("missing bank id")
	}

	return db.removeUserBank(userID, bankID)
}

// GetBank gets information of a bank for a user
func GetBank(userID, bankID string) (*Bank, error) {
	// Validate input
	switch {
	case strings.TrimSpace(userID) == "":
		return nil, errors.New("missing user id")
	case strings.TrimSpace(bankID) == "":
		return nil, errors.New("missing bank id")
	}

	return db.getUserBank(userID, bankID)
}

// GetBanks gets list of banks for the user
func GetBanks(userID string) ([]*Bank, error) {
	// Validate input
	switch {
	case strings.TrimSpace(userID) == "":
		return nil, errors.New("missing user id")
	}

	return db.getUserBanks(userID)
}
