package services

import (
	"context"
	"testing"

	"git.bluebird.id/mybb/gorooster/database"
)

func TestMapping(t *testing.T) {
	es := `
	{
		"Name": "cancel order",
		"id": "901ec8dc-8de2-448c-b64c-6f0bc49cabff",
		"type": "api_event",
		"job_data": {
		  "endpoint": "https://jsonplaceholder.typicode.com/posts/1",
		  "data": null,
		  "method": "GET",
		  "headers": [
			{
			  "key": "Token",
			  "value": "b77d808805559c2fa028add373b661a3"
			},
			{
			  "key": "App-Version",
			  "value": "6.0.0"
			},
			{
			  "key": "Device-Id",
			  "value": "e60c90b865524f76"
			},
			{
			  "key": "Content-Type",
			  "value": "application/json"
			}
		  ]
		}
	  }
	`
	RedisClient := database.GetRedisClient()
	NewEventMapper().CreateEvent(context.Background(), RedisClient, es)
	// t.Log(res)
}
