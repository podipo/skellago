package main

import (
	"podipo.com/skellago/be"

	"example.com/api/cms"
)

/*
Call's Skella's NewTestAPI and then adds the example APIs
*/
func NewTestAPI() (*be.TestAPI, error) {
	api, err := be.NewTestAPI()
	if err != nil {
		return nil, err
	}
	api.API.AddResource(NewEchoResource(), true)
	api.API.AddResource(cms.NewLogsResource(), true)
	api.API.AddResource(cms.NewLogResource(), true)
	return api, err
}
