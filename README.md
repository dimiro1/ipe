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

This software use the standard GOPATH structure and Godep as "dependency manager", so, install [godep](https://github.com/tools/godep) first.

Then from the project directory execute:

```console
$ godep restore
$ go install github.com/dimiro1/ipe
```

# How to configure?

## The server

```json
{
	"Host": ":8080",      // Host and Port
	"Expvar": true,       // Export variables to /debug/vars
	"User": "Username",   // Username for Expvar
	"Password": "123456", // Password for Expvar
	"Apps": [
		{
			"ApplicationDisabled": false, // Disable this app
			"OnlySSL": false,             // If this application is avaibale only in SSL
			"Secret": "APP_SECRET", // The secret key, do't tell anyone
			"Key": "APP_KEY",       // Application Key
			"Name": "APP_NAME",     // This app name
			"AppID": "APP_ID",      // This app ID
			"UserEvents": true,     // true if Client events are enabled
			"WebHooks": true,       // true if WebHooks are enabled
			"URLWebHook": "http://localhost:4567/php/hook.php" // Url for webhooks
		}
	]
}
```

## Libraries

### Client javascript library

```javascript
var pusher = new Pusher(API_KEY, {
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

Pusher is an excelent service, their service is very realible. I recomend for everyone.

# Where this name came from?

Here in Brazil we have this beatiful tree called [Ipê](http://en.wikipedia.org/wiki/Tabebuia_aurea), it comes in differente colors: yellow, pink, white, purple.

[I want to see pictures](https://www.flickr.com/search/?q=ipe)

# Author

Claudemiro Alves Feitosa Neto

# LICENSE

Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
Use of this source code is governed by a MIT-style
license that can be found in the LICENSE file.

