package user

import (
	"context"
	"github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"strings"
)

var db *usersDatabase

func init() {
	userDB, err := newUsersDB("users.db")
	if err != nil {
		logrus.Fatalln(err)
	}

	db = userDB
}

// User contains data for API user
type User struct {
	UserID    string `json:"user_id,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Password  string `json:"password,omitempty"`
}

// Register registers a user
func (user *User) Register() (string, error) {
	switch {
	// We can do more here like making sure its a valid email
	case strings.Trim(user.Email, " ") == "":
		return "", errMissingCredential("user email")
	case strings.Trim(user.Phone, " ") == "":
		return "", errMissingCredential("user phone")
	case strings.Trim(user.Password, " ") == "":
		return "", errMissingCredential("user password")
	}

	user.UserID = uuid.New().String()

	_, data := db.userExist(user)
	if data != "" {
		return "", errors.Errorf("user with %s already exists", data)
	}

	err := db.addUser(user)
	if err != nil {
		return "", err
	}

	// user id and nil error
	return user.UserID, nil
}

// Login authenticates user
func Login(userID, password string) (string, error) {
	switch {
	case strings.Trim(userID, " ") == "":
		return "", errMissingCredential("user id")
	case strings.Trim(password, " ") == "":
		return "", errMissingCredential("password")
	}

	user, err := db.getUser(userID)
	if err != nil {
		return "", err
	}

	if user.Password != password {
		return "", errors.New("password incorrect")
	}

	// returns jwt token and error
	return genToken(context.Background(), user)
}

// UpdateUser updates data of an existing user
func UpdateUser(userID string, newUser *User) error {
	return db.updateUser(userID, newUser)
}

// GetUser gets a user data
func GetUser(userID string) (*User, error) {
	return db.getUser(userID)
}

// JWT section:
var (
	signingKey    = []byte("MySecretSoupRecipeOrAvengersEndGame")
	signingMethod = jwt.SigningMethodHS256
)

// JWTClaims contains JWT claims information
type JWTClaims struct {
	*User
	jwt.StandardClaims
}

// genToken json web token
func genToken(
	ctx context.Context, user *User,
) (string, error) {
	token := jwt.NewWithClaims(signingMethod, JWTClaims{
		user,
		jwt.StandardClaims{
			// ExpiresAt: 1500,
			Issuer: "OpenBanking",
		},
	})

	// Generate the token ...
	return token.SignedString(signingKey)
}

// parseToken parses a jwt token to claims or fail otherwise
func parseToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return signingKey, nil
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse token with claims")
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token must be valid")
	}
	return claims, nil
}
