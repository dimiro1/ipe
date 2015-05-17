// Copyright 2015 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package utils

import (
	"fmt"
	"regexp"
	"testing"
)

func TestGenerateSession(t *testing.T) {
	sessionID := GenerateSessionID()

	fmt.Println(sessionID)
	if matched, _ := regexp.MatchString("^\\d+\\.\\d+$", sessionID); !matched {
		t.Errorf("Must match ^\\d+\\.\\d+$, value: '%s'", sessionID)
	}
}

func TestIsValidChannelName(t *testing.T) {
	if IsChannelNameValid("#@#hhh**sasas") {
		t.Errorf("Invalid Channel Name")
	}

	if !IsChannelNameValid("private-hello") {
		t.Errorf("Must be Valid Channel Name")
	}

	if !IsChannelNameValid("presence-hello") {
		t.Errorf("Must be Valid Channel Name")
	}

	if !IsChannelNameValid("public") {
		t.Errorf("Must be Valid Channel Name")
	}
}
