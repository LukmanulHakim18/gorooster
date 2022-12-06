# Gorooster

Gorooster is an event scheduler that can schedule the events you want.
For now, Gorooster only provides events in the form of API requests.

For ease of use it would be nice to use [gorooster-client](https://"git.bluebird.id/mybb/gorooster.git) as a remote to the server. But if you want to use rest-api go ahead we provide that.

# Features

- Set event
- Get event
- Update release event
- Update data event
- Delete event

# Instalation

Gorooster requires a Go version with modules support.

```git
git clone https://"git.bluebird.id/mybb/gorooster.git
```

So make sure instal the dependency in your local:

```go
go mod tidy
```

# Setup

Gorooster Requires **redis 6** to handle scheduling built in, create an .env file and add the following code:

```env
REDIS_SERVER_IP= localhost:6379
REDIS_SERVER_PASSWORD=
REDIS_SELECT_DB= 14
```

This service runs on the default port: `1407` but if you want to change it, just add the following code to the .env file

```env
RUNING_PORT=: 1407
```

This service also has a retry fire event mode, if the endpoint that is scheduled to be requested returns a response code not 2xx. Then this mode will retry hit until it succeeds or until reach the `RETRY_COUNT`.

```env
RETRY_MODE= true
RETRY_COUNT= 3
```

You can set in env your log file location

```.env
LOG_PATH=/your/directory/logs.log
```

# Quickstart

Start your server to listen for requests from rest-api and listen for event which should fire

```go
go run main.go
```

After server running and you want to make event and maintain event, make request with example like below

## Create event

1. Endpoint
   ```
   {base-url}/event/{event-key}/{event-relese-in}
   ```

- base-url : server host and port default `localhost:1407`
- event-key : uniq string and can not contain `: (colon)`
- event-relese-in : time to release event in, with format `1h30m20s`

2. Methode POST
3. Header

   ```
   X-CLIENT-NAME:POSTMANT-CLIENT
   Accept-Encoding:application/json
   ```

4. Body
   ```json
   {
     "Name": "cancel_order",
     "id": "901ec8dc-8de2-448c-b64c-6f0bc49cabff",
     "type": "api_event",
     "job_data": {
       "endpoint": "https://foo.id/bar",
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
   ```

## Get event

1. Endpoint

   ```
   {base-url}/event/{event-key}
   ```

2. Methode GET
3. Header

   ```
   X-CLIENT-NAME:POSTMANT-CLIENT
   Accept-Encoding:application/json
   ```

## Update release event

1. Endpoint

   ```
   {base-url}/event/{event-key}/{event-relese-in}
   ```

2. Methode `PUT`
3. Header

   ```
   X-CLIENT-NAME:POSTMANT-CLIENT
   Accept-Encoding:application/json
   ```

## Update data event

1. Endpoint

   ```
   {base-url}/event/{event-key}
   ```

2. Methode `PUT`
3. Header

   ```
   X-CLIENT-NAME:POSTMANT-CLIENT
   Accept-Encoding:application/json
   ```

4. Body
   ```json
   {
     "Name": "cancel_order_update",
     "id": "901ec8dc-8de2-448c-b64c-6f0bc49cabff",
     "type": "api_event",
     "job_data": {
       "endpoint": "https://foo.id/bar",
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
   ```

## Delete event

1. Endpoint
   ```
   {base-url}/event/{event-key}
   ```
2. Methode `DELETE`
3. Header

   ```
   X-CLIENT-NAME:POSTMANT-CLIENT
   Accept-Encoding:application/json
   ```
