package model

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"

	"github.com/pkg/errors"
)

var (
	fileMutex   = sync.Mutex{}
	usersFileDB = "users.db"
)

// get all users from file db
func getUsersFromFileDB() ([]*User, error) {

	fileMutex.Lock()
	usersBytes, err := ioutil.ReadFile(usersFileDB)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read from file")
	}
	fileMutex.Unlock()

	users := make([]*User, 0)

	// We unmarshal the content of the file and save it to users slice
	err = json.Unmarshal(usersBytes, &users)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal users")
	}

	return users, nil
}

// We use the file for writing (io.Writer)
func saveUsersInFileDB(users []*User) error {

	// We save it back to file by marshaling it to json
	usersBytes, err := json.Marshal(users)
	if err != nil {
		return errors.Wrap(err, "failed to marshal users")
	}

	fileMutex.Lock()
	defer fileMutex.Unlock()

	// We store the credentials in file for simplicity. We could replace it with mysql db.
	dbFile, err := os.OpenFile(usersFileDB, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return errors.Wrap(err, "failed to open file")
	}

	_, err = dbFile.Write(usersBytes)
	if err != nil {
		return errors.Wrap(err, "failed to write users to file")
	}

	return nil
}

// AddUserToFileDB adds user to the file db
func AddUserToFileDB(user *User) error {

	// Get all users
	users, err := getUsersFromFileDB()
	// In case when db is empty
	if len(users) == 0 {
		users = append(users, user)
		return saveUsersInFileDB(users)
	}
	if err != nil {
		return err
	}

	// Check that user email or phone is not registered
	for _, userDB := range users {
		if userDB.Email == user.Email {
			return errors.New("email is already registered")
		}
		if userDB.Phone == user.Phone {
			return errors.New("phone is already registered")
		}
	}

	// We add the user back to slice/array of users
	users = append(users, user)

	// We the new array of users back to file
	return saveUsersInFileDB(users)
}

// AuthenticateUserInFileDB authenticates user credentials against those in our file db
func AuthenticateUserInFileDB(user *User) error {

	usersDB, err := getUsersFromFileDB()
	if err != nil {
		return err
	}

	var userExist bool // default is false

	// Check that user exists
	for _, userDB := range usersDB {
		if userDB.Email == user.Email || userDB.Phone == user.Phone {
			// We check that password matches
			if userDB.Password != user.Password {
				return errors.New("password do not match")
			}
			userExist = true
			break
		}
	}

	// Error if they don't exist
	if !userExist {
		return errors.New("user does not exist")
	}

	return nil
}
