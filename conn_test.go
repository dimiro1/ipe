// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"
)

func Test_New_ID(t *testing.T) {
	id := newID()

	if newID() != id+1 {
		t.Error("Every call to newID must increment the id by one")
	}
}
