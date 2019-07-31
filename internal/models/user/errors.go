package user

import "github.com/pkg/errors"

func errMissingCredential(cred string) error {
	return errors.Errorf("missing credential: %s", cred)
}

func errFailedToRegisterUser(err error) error {
	return errors.Errorf("failed to register user: %s", err)
}
