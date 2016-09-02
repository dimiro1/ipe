// Copyright 2016 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ipe

import "testing"

func Benchmark_memdb_GetAppByAppID(b *testing.B) {
	db := newMemdb()
	db.AddApp(&app{AppID: "123456", Name: "Example"})
	db.AddApp(&app{AppID: "654321", Name: "Example2"})
	db.AddApp(&app{AppID: "678901", Name: "Example3"})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		db.GetAppByAppID("123456")
	}
}

func Test_db_GetAppByAppID(t *testing.T) {
	app := &app{AppID: "123456", Name: "Example"}

	db := newMemdb()
	db.AddApp(app)

	a, err := db.GetAppByAppID("123456")

	if err != nil {
		t.Errorf("GetAppByAppID(%q) == %q, want %q", "123456", a, app)
	}
}

func Test_db_GetAppByAppID__error(t *testing.T) {
	app := &app{AppID: "123456", Name: "Example"}

	db := newMemdb()
	db.AddApp(app)

	a, err := db.GetAppByAppID("not-found")

	if err == nil {
		t.Errorf("GetAppByAppID(%q) == %q, want %q", "123456", a, app)
	}
}

func Test_db_GetAppByKey(t *testing.T) {
	app := &app{AppID: "123456", Name: "Example", Key: "654321"}

	db := newMemdb()
	db.AddApp(app)

	a, err := db.GetAppByKey("654321")

	if err != nil {
		t.Errorf("GetAppByKey(%q) == %q, want %q", "654321", a, app)
	}
}

func Test_db_GetAppByKey__error(t *testing.T) {
	app := &app{AppID: "123456", Name: "Example", Key: "654321"}

	db := newMemdb()
	db.AddApp(app)

	a, err := db.GetAppByKey("not-found")

	if err == nil {
		t.Errorf("GetAppByKey(%q) == %q, want %v", "not-found", a, nil)
	}
}
