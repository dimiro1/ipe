// Copyright 2015 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package utils

import (
	"regexp"
	"testing"
)

func BenchmarkGenerateSession(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateSessionID()
	}
}

func BenchmarkIsChannelNameValid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsChannelNameValid("hello-world")
	}
}

func TestGenerateSession(t *testing.T) {
	sessionID := GenerateSessionID()

	if matched, _ := regexp.MatchString("^\\d+\\.\\d+$", sessionID); !matched {
		t.Errorf("Must match ^\\d+\\.\\d+$, value: '%s'", sessionID)
	}
}

func TestIsValidChannelName(t *testing.T) {
	name := "#@#hhh**sasas"
	ok := IsChannelNameValid(name)

	if ok {
		t.Errorf("IsChannelNameValid(%s) == %t, wants %t", name, ok, false)
	}

	name = "private-hello"
	ok = IsChannelNameValid(name)

	if !ok {
		t.Errorf("IsChannelNameValid(%s) == %t, wants %t", name, ok, true)
	}

	name = "presence-hello"
	ok = IsChannelNameValid(name)

	if !ok {
		t.Errorf("IsChannelNameValid(%s) == %t, wants %t", name, ok, true)
	}

	name = "public"
	ok = IsChannelNameValid(name)

	if !ok {
		t.Errorf("IsChannelNameValid(%s) == %t, wants %t", name, ok, true)
	}
}

func TestIsPrivateChannel_valid(t *testing.T) {
	name := "private-hello"
	ok := IsPrivateChannel(name)

	if !ok {
		t.Errorf("IsPrivateChannel(%s) == %t, wants %t", name, ok, true)
	}
}

func TestIsPrivateChannel_invalid(t *testing.T) {
	name := "hello"
	ok := IsPrivateChannel(name)

	if ok {
		t.Errorf("IsPrivateChannel(%s) == %t, wants %t", name, ok, false)
	}
}

func TestIsIsPresenceChannel_valid(t *testing.T) {
	name := "presence-hello"
	ok := IsPresenceChannel(name)

	if !ok {
		t.Errorf("IsPresenceChannel(%s) == %t, wants %t", name, ok, true)
	}
}

func TestIsPresenceChannel_invalid(t *testing.T) {
	name := "hello"
	ok := IsPresenceChannel(name)

	if ok {
		t.Errorf("IsPresenceChannel(%s) == %t, wants %t", name, ok, false)
	}
}

func TestIsClientEvent_valid(t *testing.T) {
	name := "client-hello"
	ok := IsClientEvent(name)

	if !ok {
		t.Errorf("IsClientEvent(%s) == %t, wants %t", name, ok, true)
	}
}

func TestIsClientEvent_invalid(t *testing.T) {
	name := "hello"
	ok := IsClientEvent(name)

	if ok {
		t.Errorf("IsClientEvent(%s) == %t, wants %t", name, ok, false)
	}
}

func TestHashMAC(t *testing.T) {
	message := []byte("hello world")
	key := []byte("my super secret key")
	digest := HashMAC(message, key)

	// See: http://www.freeformatter.com/hmac-generator.html
	expected := "0811b8affc185a01e1a65b80089ebb1f7f68d287fc3b64581da9ec99136ad1db"

	if digest != expected {
		t.Errorf("HashMAC(%s, %q) == %s, wants %s", message, key, digest, expected)
	}
}
