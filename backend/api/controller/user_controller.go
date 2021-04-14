package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	helpers "github.com/MinhL4m/blogs/api"
	"github.com/MinhL4m/blogs/api/auth"
	models "github.com/MinhL4m/blogs/models"
	"github.com/gorilla/mux"
)

func (a *App) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&user); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}
	defer r.Body.Close()

	user.Prepare()
	err := user.Validation("")
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	userCreated, err := user.SaveUser(a.DB)

	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	// https://www.geeksforgeeks.org/http-headers-location/
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, userCreated.ID))
	helpers.RespondWithJSON(w, http.StatusCreated, userCreated)
}

func (a *App) UpdateUser(w http.ResponseWriter, r *http.Request) {

	// Get ID
	vars := mux.Vars(r)

	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	decoder := json.NewDecoder(r.Body)
	var user models.User

	if err := decoder.Decode(&user); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()

	tokenID, err := auth.ExtractTokenID(r)

	if err != nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Check if tokenId match uint32
	if tokenID != uint32(uid) {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	user.Prepare()
	err = user.Validation("update")

	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Request")
		return
	}

	updatedUser, err := user.UpdateUser(a.DB, tokenID)

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, updatedUser)
}

func (a *App) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var user models.User

	uid, err := strconv.ParseUint(vars["id"], 10, 32)

	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()

	tokenID, err := auth.ExtractTokenID(r)

	if err != nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if tokenID != uint32(uid) {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	updatedUser, err := user.UpdateUserPassword(a.DB, tokenID)

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, updatedUser)
}

func (a *App) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var user models.User

	uid, err := strconv.ParseUint(vars["id"], 10, 32)

	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()

	tokenID, err := auth.ExtractTokenID(r)

	if err != nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if tokenID != uint32(uid) {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	_, err = user.DeleteUser(a.DB, tokenID)

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	helpers.RespondWithJSON(w, http.StatusNoContent, "")
}
