package main

import (
	"net/http"
	"restful_stack/models"
)

func (app *app) GetUser(r *http.Request) (*models.User, error) {
	return &models.User{Id: 1}, nil
}
