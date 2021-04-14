package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	helpers "github.com/MinhL4m/blogs/api"
	auth "github.com/MinhL4m/blogs/api/auth"
	models "github.com/MinhL4m/blogs/models"
	"github.com/gorilla/mux"
)

func (a *App) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	post := models.Post{}

	limit64, err := strconv.ParseUint(query.Get("limit"), 10, 8)

	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Request")
		return
	}

	skip64, err := strconv.ParseUint(query.Get("page"), 10, 8)

	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Request")
		return
	}

	posts, err := post.FindAllPosts(a.DB, uint8(limit64), uint8(skip64))

	if err != nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Post not found")
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, posts)
}

func (a *App) GetPostByID(w http.ResponseWriter, r *http.Request) {
	post := models.Post{}

	vars := mux.Vars(r)

	pid64, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Request")
		return
	}

	postFound, err := post.FindPostById(a.DB, uint32(pid64))

	if err != nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Post not found")
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, postFound)
}

func (a *App) GetPostsByTitle(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	post := models.Post{}

	title := query.Get("title")

	if title == "" || len(title) == 0 {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Request")
		return
	}

	limit64, err := strconv.ParseUint(query.Get("limit"), 10, 8)

	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Request")
		return
	}

	skip64, err := strconv.ParseUint(query.Get("page"), 10, 8)

	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Request")
		return
	}

	postsFound, err := post.FindPostByTitle(a.DB, title, uint8(limit64), uint8(skip64))

	if err != nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Post not found")
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, postsFound)
}

func (a *App) UpdatePost(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Check if the post id in url
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Request")
		return
	}

	//Check if the auth token is valid and  get the user id from it
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Check if the post exist
	post := models.Post{}
	err = a.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Post not found")
		return
	}

	// If a user attempt to update a post not belonging to user
	if uid != post.AuthorID {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// get post data from body
	var postUpdate models.Post
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&post); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	//Also check if the request user id is equal to the one gotten from token
	if uid != postUpdate.AuthorID {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// valid new post
	postUpdate.Prepare()
	err = postUpdate.Validate()
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Unvalid Request")
		return
	}

	postUpdate.ID = post.ID //this is important to tell the model the post id to update, the other update field are set above

	postUpdated, err := postUpdate.UpdatePost(a.DB)

	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Unvalid Request")
		return
	}
	helpers.RespondWithJSON(w, http.StatusOK, postUpdated)
}

func (a *App) DeletePost(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Is a valid post id given to us?
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Unvalid Request")
		return
	}

	// Is this user authenticated?
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Check if the post exist
	post := models.Post{}
	err = a.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Is the authenticated user, the owner of this post?
	if uid != post.AuthorID {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	_, err = post.DeletePost(a.DB, pid, uid)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Unvalid Request")
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	helpers.RespondWithJSON(w, http.StatusNoContent, "")
}
