package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go-article/models"
	"go-article/util"
)

func (c Controller) CreateArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Endpoint Hit: Creating New Article")

	var article models.ArticleRequestParam
	err := json.NewDecoder(r.Body).Decode(&article)

	if err != nil {
		err := util.ResponseError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Invalid Format!",
		}
		util.SendResponseError(w, err)
		log.Println(err)
		return
	}

	if article.Body == "" || article.Description == "" || article.Title == "" {
		err := util.ResponseError{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Title, Description and Body can't be blank!",
		}
		util.SendResponseError(w, err)
		return
	}

	username, errJwt := util.CheckJwt(r)

	if errJwt != nil {
		err := util.ResponseError{
			StatusCode: 401,
			Message:    "Unauthorized!",
		}
		util.SendResponseError(w, err)
		return
	} else {
		user, err := c.DB.FindUserByUsername(username)

		if err == nil {
			user_id := user.ID
			arti, err := c.DB.CreateArticle(article, user_id)
			if err == nil {
				util.SendResponseData(w, arti)
			}
		}
	}

}

func (c Controller) GetListArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Endpoint Hit: Get List Articles")

	username, errJwt := util.CheckJwt(r)

	if errJwt != nil {
		err := util.ResponseError{
			StatusCode: 401,
			Message:    "Unauthorized!",
		}
		util.SendResponseError(w, err)
		return
	} else {
		user, err := c.DB.FindUserByUsername(username)

		if err == nil {
			user_id := user.ID
			arti, err := c.DB.GetListArticle(user_id)
			if err == nil {
				util.SendResponseData(w, arti)
			}
		}
	}
}
