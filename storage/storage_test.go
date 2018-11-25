// Copyright 2016 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package storage

import (
	"ipe/app"
	"testing"
)

func Benchmark_memdb_GetAppByAppID(b *testing.B) {
	storage := NewInMemory()
	_ = storage.AddApp(&app.Application{AppID: "123456", Name: "Example"})
	_ = storage.AddApp(&app.Application{AppID: "654321", Name: "Example2"})
	_ = storage.AddApp(&app.Application{AppID: "678901", Name: "Example3"})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = storage.GetAppByAppID("123456")
	}
}

func Test_db_GetAppByAppID(t *testing.T) {
	_app := &app.Application{AppID: "123456", Name: "Example"}

	storage := NewInMemory()
	_ = storage.AddApp(_app)

	a, err := storage.GetAppByAppID("123456")

	if err != nil {
		t.Errorf("GetAppByAppID(%q) == %+v, want %+v", "123456", a, _app)
	}
}

func Test_db_GetAppByAppID__error(t *testing.T) {
	_app := &app.Application{AppID: "123456", Name: "Example"}

	storage := NewInMemory()
	_ = storage.AddApp(_app)

	a, err := storage.GetAppByAppID("not-found")

	if err == nil {
		t.Errorf("GetAppByAppID(%q) == %+v, want %+v", "123456", a, _app)
	}
}

func Test_db_GetAppByKey(t *testing.T) {
	_app := &app.Application{AppID: "123456", Name: "Example", Key: "654321"}

	storage := NewInMemory()
	_ = storage.AddApp(_app)

	a, err := storage.GetAppByKey("654321")

	if err != nil {
		t.Errorf("GetAppByKey(%q) == %+v, want %+v", "654321", a, _app)
	}
}

func Test_db_GetAppByKey__error(t *testing.T) {
	_app := &app.Application{AppID: "123456", Name: "Example", Key: "654321"}

	storage := NewInMemory()
	_ = storage.AddApp(_app)

	a, err := storage.GetAppByKey("not-found")

	if err == nil {
		t.Errorf("GetAppByKey(%q) == %+v, want %+v", "not-found", a, nil)
	}
}
