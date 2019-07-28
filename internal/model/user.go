package model

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"strings"
)

// User contains data for API user
type User struct {
	UserID    string `json:"user_id,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Password  string `json:"password,omitempty"`
}

// Register registers a user. Returns USerID on successful registration and error object
func (user *User) Register() (string, error) {
	// Validate user credentials. Check that credentials not empty
	switch {
	// We can do more here like making sure its a valid email
	case strings.Trim(user.Email, " ") == "":
		return "", errMissingCredential("user email")
	case strings.Trim(user.Phone, " ") == "":
		return "", errMissingCredential("user phone")
	case strings.Trim(user.Password, " ") == "":
		return "", errMissingCredential("user password")
	}

	// We create a unique user id
	user.UserID = uuid.New().String()

	// We store the credentials in file for simplicity. We could replace it with mysql db.
	err := AddUserToFileDB(user)
	if err != nil {
		return "", err
	}

	return user.UserID, nil
}

// Login authenticates user. Creates a jwt token for user on successful login and error object
func (user *User) Login() (string, error) {
	// We validate that credentials is not empty
	switch {
	// We can do more here like making sure its a valid email
	case strings.Trim(user.Email, " ") == "":
		return "", errMissingCredential("user email")
	case strings.Trim(user.Phone, " ") == "":
		return "", errMissingCredential("user phone")
	case strings.Trim(user.Password, " ") == "":
		return "", errMissingCredential("user password")
	}

	// Authenticate user credentials
	err := AuthenticateUserInFileDB(user)
	if err != nil {
		return "", err
	}

	return GenToken(context.Background(), user)
}

// Updateuser updates the user of an existing user
func (user *User) Updateuser(newUser *User) error {
	return nil
}

var (
	signingKey    = []byte("MySecretSoupRecipeOrAvengersEndGame")
	signingMethod = jwt.SigningMethodHS256
)

// JWTClaims contains JWT claims information
type JWTClaims struct {
	*User
	jwt.StandardClaims
}

// GenToken json web token
func GenToken(
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

// ParseToken parses a jwt token to claims or fail otherwise
func ParseToken(tokenString string) (*JWTClaims, error) {
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
