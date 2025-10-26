package services

import (
	"context"
	"database/sql"
	"errors"
	"homedb/repository"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func Login(ctx context.Context, username, password string) (*repository.User, error) {
	// check if user exists
	db, err := sql.Open("postgres", os.Getenv("DB_STRING"))
	if err != nil {
		return nil, err
	}
	defer db.Close()

	repo := repository.New(db)

	if username == "" || password == "" {
		return nil, err
	}

	user, err := repo.GetUserByName(ctx, username)
	if err != nil {
		return nil, err
	}

	// check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}

	return &user, nil
}

func Signup(ctx context.Context, username, email, password, passwordRepeat string) (*repository.User, []error) {
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
	db, err := sql.Open("postgres", os.Getenv("DB_STRING"))
	if err != nil {
		return nil, []error{err}
	}
	defer db.Close()

	repo := repository.New(db)

	user, err := repo.CreateUser(ctx, repository.CreateUserParams{Username: username, Email: email, Password: string(hash)})
	if err != nil {
		return nil, []error{err}
	}

	return &user, nil
}
