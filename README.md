[![Go Report Card](https://goreportcard.com/badge/github.com/dimiro1/ipe)](https://goreportcard.com/report/github.com/dimiro1/ipe)

Try browsing [the code on Sourcegraph](https://sourcegraph.com/github.com/dimiro1/ipe)!

# IPÊ

An open source Pusher server implementation compatible with Pusher client libraries written in Go.

# Why I wrote this software?

1. I wanted to learn Go and I needed a non trivial application;
2. I use Pusher in some projects;
3. I really like Pusher;
4. I was using Pusher on some projects behind a firewall;

# Features

* Public Channels;
* Private Channels;
* Presence Channels;
* Web Hooks;
* Client events;
* Complete REST API;
* Easy installation;
* A single binary without dependencies;
* Easy configuration;
* Protocol version 7;
* Multiple apps in the same instance;
* Drop in replacement for pusher server;
* Dynamically add clients via api
* Get all clients via api

# Download pre built binaries

You can download pre built binaries from the [releases tab](https://github.com/dimiro1/ipe/releases).

# Building

```console
$ go get github.com/dimiro1/ipe
```

or simply

```console
$ go install github.com/dimiro1/ipe
```

Building from a local checkout

```console
$ git clone https://github.com/dimiro1/ipe.git
$ cd ipe/cmd
$ go build -o ipe
```

# How to configure?

## The server

```yaml

---
host: ":8080"
profiling: false
ssl:
  enabled: false
  host: ":4343"
  key_file: "key.pem"
  cert_file: "cert.pem"
apps:
  - name: "Sample Application"
    enabled: true
    only_ssl: false
    key: "278d525bdf162c739803"
    secret: "${APP_SECRET}" # Expand env vars
    app_id: "1"
    user_events: true
    webhooks:
      enabled: true
      url: "http://127.0.0.1:5000/hook"

```

## Libraries

### Client javascript library

```javascript
let pusher = new Pusher(APP_KEY, {
  wsHost: 'localhost',
  wsPort: 8080,
  wssPort: 4433,    // Required if encrypted is true
  encrypted: false, // Optional. the application must use only SSL connections
  enabledTransports: ["ws", "flash"],
  disabledTransports: ["flash"]
});
```

### Client server libraries

Ruby

```ruby
Pusher.host = 'localhost'
Pusher.port = 8080
```

PHP

```php
$pusher = new Pusher(APP_KEY, APP_SECRET, APP_ID, DEBUG, "http://localhost", "8080");
```

Go

```go
package main

import "github.com/pusher/pusher-http-go"

func main() {
	client := pusher.Client{
        AppId:  "APP_ID",
        Key:    "APP_KEY",
        Secret: "APP_SECRET",
        Host:   ":8080",
    }
	
	// use the client
}
```

NodeJS

```javascript
let pusher = new Pusher({
  appId: APP_ID,
  key: APP_KEY,
  secret: APP_SECRET
  domain: 'localhost',
  port: 80
});

```

# API Access

### To create client api with credentials call

```
{host}/app/create
```

with Authorization header token (via env)

Example 

```
curl --location --request POST 'http://127.0.0.1:8080/app/create' \
--header 'Content-Type: application/json' \
--header 'Authorization: 7kAWg6ADHx' \
--data-raw '{
	"name" : "Testing 5",
  "appID" : "5",
  "key" : "testing-5",
  "secret" : "testing-5",
  "onlySSL" : false,
  "webHookURL":"https://webhook.site/e4cf00c5-c379-4c19-8eb4-a32e6002728b5"
}'
```


responds

```
{
    "Message": "Testing 5 Successfully added to pusher server!",
    "apps": [
        {
            "Name": "Sample Application",
            "AppID": "1",
            "Key": "",
            "Secret": "",
            "OnlySSL": false,
            "Enabled": true,
            "UserEvents": true,
            "WebHooks": true,
            "URLWebHook": "http://0.0.0.0:5000/hook"
        },
        {
            "Name": "Testing 5",
            "AppID": "5",
            "Key": "testing-5",
            "Secret": "testing-5",
            "OnlySSL": false,
            "Enabled": true,
            "UserEvents": true,
            "WebHooks": true,
            "URLWebHook": "https://webhook.site/e4cf00c5-c379-4c19-8eb4-a32e6002728b5"
        }
    ]
}
```


### To get all clients

```
{host}/app/all
```

Example

```
curl --location --request GET 'http://127.0.0.1:8080/apps/all' \
--header 'Authorization: 7kAWg6ADHx'
```

Respose will be

```
{
    [
        {
            "Name": "Sample Application",
            "AppID": "1",
            "Key": "",
            "Secret": "",
            "OnlySSL": false,
            "Enabled": true,
            "UserEvents": true,
            "WebHooks": true,
            "URLWebHook": "http://0.0.0.0:5000/hook"
        },
        {
            "Name": "Testing 5",
            "AppID": "5",
            "Key": "testing-5",
            "Secret": "testing-5",
            "OnlySSL": false,
            "Enabled": true,
            "UserEvents": true,
            "WebHooks": true,
            "URLWebHook": "https://webhook.site/e4cf00c5-c379-4c19-8eb4-a32e6002728b5"
        }
    ]
}
```




# Logging

This software uses the [glog](https://github.com/golang/glog) library

for more information about logging type the following in console.

```console
$ ipe -h
```

# When use this software?

* When you are offline;
* When you want to control your infrastructure;
* When you do not want to have external dependencies;
* When you want extend the protocol;

# Contributing.

Feel free to fork this repo.

# Pusher

Pusher is an excellent service, their service is very reliable. I recommend for everyone.

# Where this name came from?

Here in Brazil we have this beautiful tree called [Ipê](http://en.wikipedia.org/wiki/Tabebuia_aurea), it comes in differente colors: yellow, pink, white, purple.

[I want to see pictures](https://www.flickr.com/search/?q=ipe)

# Author

Claudemiro Alves Feitosa Neto

# LICENSE

Copyright 2014, 2018 Claudemiro Alves Feitosa Neto. All rights reserved.
Use of this source code is governed by a MIT-style
license that can be found in the LICENSE file.

