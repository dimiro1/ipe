// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

var channelValidationRegex = regexp.MustCompile("^[A-Za-z0-9_\\-=@,.;]+$")

// HashMAC Calculates the MAC signing with the given key and returns the hexadecimal encoded Result
func HashMAC(message, key []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expected := mac.Sum(nil)

	return hex.EncodeToString(expected)
}

// GenerateSessionID Generate a new random Hash
func GenerateSessionID() string {
	return fmt.Sprintf("%d.%d", rand.Intn(math.MaxInt32), rand.Intn(math.MaxInt32))
}

// IsChannelNameValid Verify if the channel name is valid
func IsChannelNameValid(channelName string) bool {
	return channelValidationRegex.MatchString(channelName)
}

// IsPrivateChannel Verify if the channel name represents a private channel
func IsPrivateChannel(channelName string) bool {
	return strings.HasPrefix(channelName, "private-")
}

// IsPresenceChannel Verify if the channel name represents a presence channel
func IsPresenceChannel(channelName string) bool {
	return strings.HasPrefix(channelName, "presence-")
}

// IsClientEvent Verify if the event name represents a client event type
func IsClientEvent(event string) bool {
	return strings.HasPrefix(event, "client-")
}

type ValueType string

const (
	String ValueType = "string"
	Int    ValueType = "int"
)

type UserID struct {
	strVal    string
	intVal    int64
	valueType ValueType
}

func (u *UserID) UnmarshalJSON(bytes []byte) error {
	var (
		err error
	)

	if bytes[0] == '"' {
		err = json.Unmarshal(bytes, &u.strVal)
		u.valueType = String
	} else {
		err = json.Unmarshal(bytes, &u.intVal)
		u.valueType = Int
	}

	return err
}

func (u *UserID) GetID() string {
	if u.valueType == Int {
		return strconv.FormatInt(u.intVal, 10)
	}

	return u.strVal
}
