[![Quality code report](http://goreportcard.com/badge/dimiro1/ipe)](http://goreportcard.com/report/dimiro1/ipe)

# IPÊ

This software is written in Go - the WYSIWYG lang

# Why I wrote this software?

1. I wanted to learn Go and I needed a non trivial application;
2. I use Pusher in some projects;
3. I really like Pusher;

# Features

* Public Channels;
* Private Channels;
* Presence Channels;
* Web Hooks;
* Client events;
* Complete REST API;
* Easy instalation;
* A single binary without dependencies;
* Easy configuration;
* Protocol version 7;
* Multiple apps in the same instance;
* Drop in replacement for pusher server;

# Building

```console
$ go get github.com/dimiro1/ipe
```

or simply

```console
$ go install github.com/dimiro1/ipe
```

# How to configure?

## The server

```json
{
	"Host": ":8080",
	"SSL": false,
	"SSLHost": ":8090",
	"SSLKeyFile": "key.pem",
	"SSLCertFile": "cert.pem",
	"Apps": [
		{
			"ApplicationDisabled": false,
			"Secret": "7ad3753142a6693b25b9",
			"Key": "278d525bdf162c739803",
			"OnlySSL": false,
			"Name": "App 1",
			"AppID": "321",
			"UserEvents": true,
			"WebHooks": true,
			"URLWebHook": "http://127.0.0.1:4567/php/hook.php"
		},
		{
			"ApplicationDisabled": false,
			"Secret": "d6824d2fa32888931504",
			"Key": "c8b30f611ffb13202976",
			"OnlySSL": false,
			"Name": "App 2",
			"AppID": "123",
			"UserEvents": true,
			"WebHooks": false,
			"URLWebHook": "http://127.0.0.1:4567/php/hook.php"
		}
	]
}
```

## Libraries

### Client javascript library

```javascript
var pusher = new Pusher(APP_KEY, {
  wsHost: 'localhost',
  wsPort: 8080,
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

NodeJS

```javascript
var pusher = new Pusher({
  appId: APP_ID,
  key: APP_KEY,
  secret: APP_SECRET
  domain: 'localhost',
  port: 80
});

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

Pusher is an excelent service, their service is very reliable. I recomend for everyone.

# Where this name came from?

Here in Brazil we have this beautiful tree called [Ipê](http://en.wikipedia.org/wiki/Tabebuia_aurea), it comes in differente colors: yellow, pink, white, purple.

[I want to see pictures](https://www.flickr.com/search/?q=ipe)

# Author

Claudemiro Alves Feitosa Neto

# LICENSE

Copyright 2014, 2015, 2016 Claudemiro Alves Feitosa Neto. All rights reserved.
Use of this source code is governed by a MIT-style
license that can be found in the LICENSE file.

