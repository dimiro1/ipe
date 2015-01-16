// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
)

// HashMAC Calculates the MAC signing with the given key and returns the hexadecimal encoded Result
func HashMAC(message, key []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expected := mac.Sum(nil)

	return hex.EncodeToString(expected)
}

// Generate a new random Hash
func RandomHash() string {
	b := make([]byte, 25)

	if _, err := rand.Read(b); err != nil {
		panic("websockets: Could not generate a random session ID")
	}

	return base32.StdEncoding.EncodeToString(b)
}
