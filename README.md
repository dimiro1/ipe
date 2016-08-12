[![Go Report Card](https://goreportcard.com/badge/github.com/dimiro1/ipe)](https://goreportcard.com/report/github.com/dimiro1/ipe)

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

# Download pre built binaries

You can download pre built binaries from the [releases tab](https://github.com/dimiro1/ipe/releases).

I do not have a Windows machine, so I can only distribute binaries for amd64 linux and amd64 darwin.

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

```javascript
{
	"Host": ":8080",                    // Required
	"SSL": false,                       // Required but can be false
	"SSLHost": ":4433",                 // Required if SSL is true
	"SSLKeyFile": "A key.pem file",     // Required if SSL is true
	"SSLCertFile": "A cert.pem file",   // Required if SSL is true
	"Apps": [                           // Required, A Json arrays with multiple apps
		{
			"ApplicationDisabled": false,               // Required but can be false
			"Secret": "A really secret random string",  // Required
			"Key": "A random Key string",               // Required
			"OnlySSL": false,                           // Required but can be false
			"Name": "The app name",                     // Required
			"AppID": "The app ID",                      // Required
			"UserEvents": true,                         // Required but can be false
			"WebHooks": true,                           // Required but can be false
			"URLWebHook": "Some URL to send webhooks"   // Required if WebHooks is true
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

