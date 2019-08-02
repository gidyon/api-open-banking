package user

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

type usersDatabase struct {
	dbFilePath string
	mu         *sync.Mutex // guards users
	users      map[string]*User
}

// newUsersDB creates a new users db
func newUsersDB(dbFilePath string) (*usersDatabase, error) {
	if strings.TrimSpace(dbFilePath) == "" {
		dbFilePath = "users.db"
	}

	usersDB := &usersDatabase{
		dbFilePath: dbFilePath,
		mu:         &sync.Mutex{},
		users:      make(map[string]*User, 0),
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

	// load users to usersDB.users map
	return usersDB, usersDB.load()
}

// get all users
func (usersDB *usersDatabase) getUsers() map[string]*User {
	usersDB.mu.Lock()
	defer usersDB.mu.Unlock()

	return usersDB.users
}

// userExist returns true if any of user data exist and corresponding name of data
func (usersDB *usersDatabase) userExist(user *User) (bool, string) {
	usersDB.mu.Lock()
	defer usersDB.mu.Unlock()

	// Get user from map, if they don't exist, ok is false
	_, ok := usersDB.users[user.UserID]
	if !ok {
		// Check deeper
		for _, v := range usersDB.users {
			if user.Email == v.Email {
				return false, "email"
			}
			if user.Phone == v.Phone {
				return false, "phone"
			}
		}
	}

	return true, ""
}

// gets a single user
func (usersDB *usersDatabase) getUser(userID string) (*User, error) {
	usersDB.mu.Lock()
	defer usersDB.mu.Unlock()

	user, ok := usersDB.users[userID]
	if !ok {
		// Check deeper
		for _, user := range usersDB.users {
			if user.Email == userID {
				return user, nil
			}
			if user.Phone == userID {
				return user, nil
			}
			if user.UserID == userID {
				return user, nil
			}
		}

		return nil, errors.Errorf("couldn't find user with id: %v", userID)
	}

	return user, nil
}

func (usersDB *usersDatabase) addUser(user *User) error {
	_, data := usersDB.userExist(user)

	if data != "" {
		return errors.Errorf("user with %s already exist", data)
	}

	usersDB.mu.Lock()
	defer usersDB.mu.Unlock()

	usersDB.users[user.UserID] = user

	return usersDB.save()
}

func (usersDB *usersDatabase) updateUser(userID string, newUser *User) error {
	usersDB.mu.Lock()
	defer usersDB.mu.Unlock()

	user, ok := usersDB.users[userID]
	if !ok {
		return errors.New("user doesn't exist")
	}

	switch {
	case strings.TrimSpace(newUser.Email) != "":
		user.Email = newUser.Email
	case strings.TrimSpace(newUser.Phone) != "":
		user.Phone = newUser.Phone
	case strings.TrimSpace(newUser.FirstName) != "":
		user.FirstName = newUser.FirstName
	case strings.TrimSpace(newUser.Phone) != "":
		user.LastName = newUser.LastName
	case strings.TrimSpace(newUser.Password) != "":
		user.Password = newUser.Password
	}

	*usersDB.users[userID] = *newUser

	return usersDB.save()
}

// save will save the current state of the map in a file
func (usersDB *usersDatabase) save() error {
	usersBytes, err := json.Marshal(usersDB.users)
	if err != nil {
		return errors.Wrap(err, "failed to marshal users")
	}

	dbFile, err := os.OpenFile(usersDB.dbFilePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return errors.Wrap(err, "failed to open file")
	}

	_, err = dbFile.Write(usersBytes)
	if err != nil {
		return errors.Wrap(err, "failed to write users to file")
	}

	return nil
}

func (usersDB *usersDatabase) load() error {
	usersDB.mu.Lock()
	defer usersDB.mu.Unlock()

	usersBytes := make([]byte, 0)

	usersBytes, err := ioutil.ReadFile(usersDB.dbFilePath)
	if err != nil {
		return errors.Wrap(err, "failed read from file")
	}

	usersBytes = bytes.TrimSpace(usersBytes)

	if len(usersBytes) == 0 {
		usersDB.users = make(map[string]*User, 0)
		return nil
	}

	users := make(map[string]*User, 0)

	err = json.Unmarshal(usersBytes, &users)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal users")
	}

	usersDB.users = users

	return nil
}
