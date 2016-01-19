// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"math/rand"
	"regexp"
)

// HashMAC Calculates the MAC signing with the given key and returns the hexadecimal encoded Result
func HashMAC(message, key []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expected := mac.Sum(nil)

	return hex.EncodeToString(expected)
}

// GenerateSessionID Generate a new random Hash
func GenerateSessionID() string {
	MAX := math.MaxInt64

	return fmt.Sprintf("%d.%d", rand.Intn(MAX), rand.Intn(MAX))
}

// IsChannelNameValid Verify if the channel name is valid
func IsChannelNameValid(channelName string) bool {
	matched, err := regexp.MatchString("^[A-Za-z0-9_\\-=@,.;]+$", channelName)

	if err == nil && matched {
		return true
	}

	return false
}
