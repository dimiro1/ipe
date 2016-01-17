// Copyright 2016 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ipe

import (
	"testing"

	"github.com/gorilla/websocket"
)

func TestNewConnection(t *testing.T) {
	expectedSocketID := "socketID"
	expectedSocket := &websocket.Conn{}

	c := newConnection(expectedSocketID, expectedSocket)

	if c.SocketID != expectedSocketID {
		t.Errorf("Expected: %s but got %s", expectedSocketID, c.SocketID)
	}

	if c.Socket != expectedSocket {
		t.Errorf("Expected: %+v but got %+v", expectedSocket, c.Socket)
	}

	if c.CreatedAt.IsZero() {
		t.Errorf("Expected %s to not be zero", c.CreatedAt)
	}
}
