package services

import (
	"context"
	"errors"
	"homedb/repository"

	"golang.org/x/crypto/bcrypt"
)

func Login(ctx context.Context, repo *repository.Queries, username, password string) (*repository.User, error) {
	if username == "" || password == "" {
		return nil, errors.New("all fields must be filled out")
	}

	user, err := repo.GetUserByName(ctx, username)
	if err != nil {
		return nil, err
	}

	// check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, errors.New("passwords dont match")
		}
		return nil, err
	}

	return &user, nil
}

func Signup(ctx context.Context, repo *repository.Queries, username, email, password, passwordRepeat string) (*repository.User, []error) {
	// validation
	var errs []error

	if username == "" || email == "" || password == "" || passwordRepeat == "" {
		errs = append(errs, errors.New("all fields must be filled out"))
	}

	if password != passwordRepeat {
		errs = append(errs, errors.New("passwords don't match"))
	}

	if len(errs) > 0 {
		return nil, errs
	}

	// hash
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, []error{err}
	}

	// save
	user, err := repo.CreateUser(ctx, repository.CreateUserParams{Username: username, Email: email, Password: string(hash)})
	if err != nil {
		return nil, []error{err}
	}

	return &user, nil
}
