// Copyright 2014 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import "testing"

func Test_New_ID(t *testing.T) {
	_, id := newID()

	if _, i := newID(); i != id+1 {
		t.Errorf("Every call to newID must increment the id by one, got: %s", i)
	}
}
