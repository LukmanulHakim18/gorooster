package repositories

import "testing"

func TestDoJobGet(t *testing.T) {
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
	repo := NewJobAPI()
	if err := repo.DoJob(es); err != nil {
		t.Log(err)
		t.Fail()
	}
}
func TestDoJobPost(t *testing.T) {
	es := `
	{
		"Name": "cancel order",
		"id": "901ec8dc-8de2-448c-b64c-6f0bc49cabff",
		"type": "api_event",
		"job_data": {
			"endpoint": "https://jsonplaceholder.typicode.com/posts/1",
		  "data": {
			"customer_id": "BB00546345",
			"location": {
			  "latitude": -6.246131139461152,
			  "longitude": 16.82597713520182
			}
		  },
		  "method": "POST",
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
	repo := NewJobAPI()
	if err := repo.DoJob(es); err != nil {
		t.Log(err)
		t.Fail()
	}
}
