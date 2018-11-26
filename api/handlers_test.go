package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"

	"ipe/app"
	channel2 "ipe/channel"
	"ipe/connection"
	"ipe/mocks"
	"ipe/storage"
)

var (
	testApp  *app.Application
	database storage.Storage
	id       = 0
)

func newTestApp() *app.Application {
	a := app.NewApplication("Test", strconv.Itoa(id), "123", "123", false, false, true, false, "")
	id++

	return a
}

func init() {
	testApp = newTestApp()

	channel := channel2.New("presence-c1")
	testApp.AddChannel(channel)
	testApp.AddChannel(channel2.New("c2"))
	testApp.AddChannel(channel2.New("private-c3"))

	conn := connection.New("123.456", mocks.MockSocket{})
	_ = testApp.Subscribe(channel, conn, "{}")

	conn = connection.New("321.654", mocks.MockSocket{})
	_ = testApp.Subscribe(channel, conn, "{}")

	_storage := storage.NewInMemory()
	_ = _storage.AddApp(testApp)

	database = _storage
}

// All channels
func Test_getChannels_all(t *testing.T) {
	appID := testApp.AppID

	r, _ := http.NewRequest("GET", fmt.Sprintf("/apps/%s/channels", appID), nil)
	r = mux.SetURLVars(r, map[string]string{
		"app_id": appID,
	})
	w := httptest.NewRecorder()

	handler := &GetChannels{database}
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("w.Code == %d, wants %d", w.Code, http.StatusOK)
	}

	data := make(map[string]interface{})
	_ = json.Unmarshal(w.Body.Bytes(), &data)

	channels := data["channels"].(map[string]interface{})

	if len(channels) != 3 {
		t.Errorf("len(%q) == %d, want %d", channels, len(channels), 3)
	}
}

func Test_getChannels_filter_by_presence_prefix(t *testing.T) {
	appID := testApp.AppID

	r, _ := http.NewRequest("GET", fmt.Sprintf("/apps/%s/channels?filter_by_prefix=presence-", appID), nil)
	r = mux.SetURLVars(r, map[string]string{
		"app_id": appID,
	})
	w := httptest.NewRecorder()

	handler := &GetChannels{database}
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("w.Code == %d, wants %d", w.Code, http.StatusOK)
	}

	data := make(map[string]interface{})
	_ = json.Unmarshal(w.Body.Bytes(), &data)

	channels := data["channels"].(map[string]interface{})

	if len(channels) != 1 {
		t.Errorf("len(%q) == %d, want %d", channels, len(channels), 1)
	}
}

// Only presence channels and user_count
func Test_getChannels_filter_by_presence_prefix_and_user_count(t *testing.T) {
	appID := testApp.AppID

	r, _ := http.NewRequest("GET", fmt.Sprintf("/apps/%s/channels?filter_by_prefix=presence-&info=user_count", appID), nil)
	r = mux.SetURLVars(r, map[string]string{
		"app_id": appID,
	})
	w := httptest.NewRecorder()

	handler := &GetChannels{database}
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("w.Code == %d, wants %d", w.Code, http.StatusOK)
	}

	data := make(map[string]interface{})
	_ = json.Unmarshal(w.Body.Bytes(), &data)

	channels := data["channels"].(map[string]interface{})

	if len(channels) != 1 {
		t.Errorf("len(%q) == %d, want %d", channels, len(channels), 1)
	}

	c, exists := channels["presence-c1"]

	if !exists {
		t.Errorf("!exists == %t, want %t", !exists, false)
	}

	_channel := c.(map[string]interface{})

	if _channel["user_count"] != float64(1) {
		t.Errorf("_channel['user_count'] == %f, want %d", _channel["user_count"], 1)
	}
}

// User count only allowed in Presence channels
func Test_getChannels_filter_by_private_prefix_and_info_user_count(t *testing.T) {
	appID := testApp.AppID

	r, _ := http.NewRequest("GET", fmt.Sprintf("/apps/%s/channels?filter_by_prefix=private-&info=user_count", appID), nil)
	r = mux.SetURLVars(r, map[string]string{
		"app_id": appID,
	})
	w := httptest.NewRecorder()

	handler := &GetChannels{database}
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("w.Code == %d, wants %d", w.Code, http.StatusBadRequest)
	}
}

func Test_getChannels_filter_by_public_prefix(t *testing.T) {
	appID := testApp.AppID

	r, _ := http.NewRequest("GET", fmt.Sprintf("/apps/%s/channels?filter_by_prefix=public-", appID), nil)
	r = mux.SetURLVars(r, map[string]string{
		"app_id": appID,
	})
	w := httptest.NewRecorder()

	handler := &GetChannels{database}
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("w.Code == %d, wants %d", w.Code, http.StatusOK)
	}

	data := make(map[string]interface{})

	_ = json.Unmarshal(w.Body.Bytes(), &data)

	channels := data["channels"].(map[string]interface{})

	if len(channels) != 1 {
		t.Errorf("len(%q) == %d, want %d", channels, len(channels), 1)
	}

	_, exists := channels["c2"]

	if !exists {
		t.Errorf("!exists == %t, want %t", !exists, false)
	}
}

func Test_getChannels_filter_by_private_prefix(t *testing.T) {
	appID := testApp.AppID

	r, _ := http.NewRequest("GET", fmt.Sprintf("/apps/%s/channels?filter_by_prefix=private-", appID), nil)
	r = mux.SetURLVars(r, map[string]string{
		"app_id": appID,
	})
	w := httptest.NewRecorder()

	handler := &GetChannels{database}
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("w.Code == %d, wants %d", w.Code, http.StatusOK)
	}

	data := make(map[string]interface{})

	_ = json.Unmarshal(w.Body.Bytes(), &data)

	channels := data["channels"].(map[string]interface{})

	if len(channels) != 1 {
		t.Errorf("len(%q) == %d, want %d", channels, len(channels), 1)
	}

	_, exists := channels["private-c3"]

	if !exists {
		t.Errorf("!exists == %t, want %t", !exists, false)
	}
}
