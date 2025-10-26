package handlers

import (
	"database/sql"
	"homedb/repository"
	"homedb/sessions"
	"homedb/views/pages"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
)

func Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	// check if user exists
	db, err := sql.Open("postgres", os.Getenv("DB_STRING"))
	if err != nil {
		pages.Login(err).Render(r.Context(), w)
		return
	}
	defer db.Close()

	repo := repository.New(db)

	user, err := repo.GetUserByName(r.Context(), username)
	if err != nil {
		pages.Login(err).Render(r.Context(), w)
		return
	}

	// check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		pages.Login(err).Render(r.Context(), w)
		return
	}

	// session
	sessId := uuid.New()
	session := sessions.New(user.ID)

	sessions.Sessions[sessId] = session

	http.SetCookie(w, &http.Cookie{
		Name:     "session-id",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
		Secure:   false, // TODO: Change to true for prod
	})

	// redirect
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Signup(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	passwordRepeat := r.FormValue("password-repeat")

	// validation
	var errors []string

	if username == "" || email == "" || password == "" || passwordRepeat == "" {
		errors = append(errors, "All fields must be filled out!")
	}

	if password != passwordRepeat {
		errors = append(errors, "Passwords don't match!")
	}

	if len(errors) > 0 {
		pages.Signup(errors).Render(r.Context(), w)
		return
	}

	// hash
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		pages.Signup([]string{err.Error()}).Render(r.Context(), w)
		return
	}

	// save
	db, err := sql.Open("postgres", os.Getenv("DB_STRING"))
	if err != nil {
		pages.Signup([]string{err.Error()}).Render(r.Context(), w)
		return
	}
	defer db.Close()

	repo := repository.New(db)

	user, err := repo.CreateUser(r.Context(), repository.CreateUserParams{Username: username, Email: email, Password: string(hash)})
	if err != nil {
		pages.Signup([]string{err.Error()}).Render(r.Context(), w)
		return
	}

	// session
	sessId := uuid.New()
	session := sessions.New(user.ID)

	sessions.Sessions[sessId] = session

	http.SetCookie(w, &http.Cookie{
		Name:     "session-id",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
		Secure:   false, // TODO: Change to true for prod
	})

	// redirect
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
