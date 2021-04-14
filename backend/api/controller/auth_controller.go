package controller

import (
	"encoding/json"
	"net/http"

	helpers "github.com/MinhL4m/blogs/api"
	auth "github.com/MinhL4m/blogs/api/auth"
	authHelper "github.com/MinhL4m/blogs/helpers"
	"github.com/MinhL4m/blogs/models"
	"golang.org/x/crypto/bcrypt"
)

func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	user := models.User{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		helpers.RespondWithError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	user.Prepare()
	err := user.Validation("login")
	if err != nil {
		helpers.RespondWithError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	token, err := a.SignIn(user.Email, user.Password)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helpers.RespondWithJSON(w, http.StatusOK, token)
}

func (a *App) SignIn(email, password string) (string, error) {

	var err error

	user := models.User{}

	err = a.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", err
	}
	err = authHelper.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CreateToken(user.ID)
}
