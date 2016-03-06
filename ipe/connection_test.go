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
		t.Errorf("c.SocketID == %s, wants %s", c.SocketID, expectedSocketID)
	}

	if c.Socket != expectedSocket {
		t.Errorf("c.Socket == %v, wants %v", c.Socket, expectedSocket)
	}

	if c.CreatedAt.IsZero() {
		t.Errorf("c.CreatedAt.IsZero() == %t, wants %t", c.CreatedAt.IsZero(), false)
	}
}
