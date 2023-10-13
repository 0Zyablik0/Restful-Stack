package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"restful_stack/models"

	"gorm.io/gorm"
)

func (app *app) MainHandler(w http.ResponseWriter, r *http.Request) {
	response := &models.MainResponse{}
	response.Message = "Hello, I am simple stack!"
	response.Code = StatusOK
	if err := json.NewEncoder(w).Encode(response); err != nil {
		app.logger.Printf("[Error] [StackTop]: %s\n", err)
	}
}

func (app *app) StackTopHandler(w http.ResponseWriter, r *http.Request) {

	user, err := app.GetUser(r)
	stackResponse := &models.StackTopResponse{}

	if err != nil {
		app.logger.Printf("[Error] [StackPush]: %s\n", err)
		stackResponse.Err = "We are not able to identify you"
		stackResponse.Code = StatusInternalError
		if err := json.NewEncoder(w).Encode(stackResponse.BaseResponse); err != nil {
			app.logger.Printf("[Error] [StackPush]: %s\n", err)
		}
		return
	}

	err = app.db.Transaction(func(tx *gorm.DB) error {
		return TransactionStackTop(tx, user, stackResponse)
	})
	if err != nil {
		stackResponse.Code = StatusInternalError
		stackResponse.Err = fmt.Sprintf("Internal error: action was not performed")
		app.logger.Printf("[Error] [StackTop]: %s\n", err)
	}
	err = json.NewEncoder(w).Encode(stackResponse)
	if err != nil {
		app.logger.Printf("[Error] [StackTop]: %s\n", err)
	}
}

func (app *app) StackPopHandler(w http.ResponseWriter, r *http.Request) {
	user, err := app.GetUser(r)
	stackResponse := &models.StackPopResponse{}

	if err != nil {
		app.logger.Printf("[Error] [StackPop]: %s\n", err)
		stackResponse.Err = "We are not able to identify you"
		stackResponse.Code = StatusInternalError
		if err := json.NewEncoder(w).Encode(stackResponse.BaseResponse); err != nil {
			app.logger.Printf("[Error] [StackPop]: %s\n", err)
		}
		return
	}

	err = app.db.Transaction(func(tx *gorm.DB) error {
		return TransactionStackPop(tx, user, stackResponse)
	})
	if err != nil {
		stackResponse.Code = StatusInternalError
		stackResponse.Err = fmt.Sprintf("Internal error: action was not performed")
		app.logger.Printf("[Error] [StackPop]: %s\n", err)
	}
	err = json.NewEncoder(w).Encode(stackResponse)
	if err != nil {
		app.logger.Printf("[Error] [StackPop]: %s\n", err)
	}

}

func (app *app) StackPushHandler(w http.ResponseWriter, r *http.Request) {
	stackRequest := &models.StackPushRequest{}
	stackResponse := &models.StackPushResponse{}
	user, err := app.GetUser(r)

	if err != nil {
		app.logger.Printf("[Error] [StackPush]: %s\n", err)
		stackResponse.Err = "We are not able to identify you"
		stackResponse.Code = StatusInternalError
		if err := json.NewEncoder(w).Encode(stackResponse.BaseResponse); err != nil {
			app.logger.Printf("[Error] [StackPush]: %s\n", err)
		}
		return
	}

	err = json.NewDecoder(r.Body).Decode(stackRequest)
	if err != nil {
		app.logger.Printf("[Error] [StackPush]: %s\n", err.Error())
		stackResponse.Err = "Wrong request syntax"
		stackResponse.Code = StatusError
		if err := json.NewEncoder(w).Encode(stackResponse.BaseResponse); err != nil {
			app.logger.Printf("[Error] [StackPush]: %s\n", err)
		}
		return
	}

	err = app.db.Transaction(func(tx *gorm.DB) error {
		return TransactionStackPush(tx, user, stackRequest, stackResponse)
	})

	if err != nil {
		stackResponse.Code = StatusInternalError
		stackResponse.Err = "Abort due to internal error"
		app.logger.Printf("[Error] [StackPush]: %s", err.Error())
	}
	err = json.NewEncoder(w).Encode(stackResponse)
	if err != nil {
		app.logger.Printf("[Error] [StackPush]: %s", err.Error())
	}
}

func (app *app) StackPushRangeHandler(w http.ResponseWriter, r *http.Request) {
	stackRequest := &models.StackPushRangeRequest{}
	stackResponse := &models.StackPushRangeResponse{}
	user, err := app.GetUser(r)

	if err != nil {
		app.logger.Printf("[Error] [StackPushRange]: %s\n", err)
		stackResponse.Err = "We are not able to identify you"
		stackResponse.Code = StatusInternalError
		if err := json.NewEncoder(w).Encode(stackResponse.BaseResponse); err != nil {
			app.logger.Printf("[Error] [StackPushRange]: %s\n", err)
		}
		return
	}

	err = json.NewDecoder(r.Body).Decode(stackRequest)
	if err != nil {
		app.logger.Printf("[Error] [StackPushRange]: %s\n", err.Error())
		stackResponse.Err = "Wrong request syntax"
		stackResponse.Code = StatusError
		if err := json.NewEncoder(w).Encode(stackResponse.BaseResponse); err != nil {
			app.logger.Printf("[Error] [StackPushRange]: %s\n", err)
		}
		return
	}

	err = app.db.Transaction(func(tx *gorm.DB) error {
		return TransactionStackPushRange(tx, user, stackRequest, stackResponse)
	})

	if err != nil {
		stackResponse.Code = StatusInternalError
		stackResponse.Err = "Abort due to internal error"
		app.logger.Printf("[Error] [StackPushRange]: %s", err.Error())
	}

	err = json.NewEncoder(w).Encode(stackResponse)
	if err != nil {
		app.logger.Printf("[Error] [StackPushRange]: %s", err.Error())
	}
}
