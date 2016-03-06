package ipe

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	testApp *app
	ctx     *applicationContext
)

func init() {
	testApp = newTestApp()

	channel := newChannel("presence-c1")
	testApp.AddChannel(channel)
	testApp.AddChannel(newChannel("c2"))
	testApp.AddChannel(newChannel("private-c3"))

	conn := newConnection("123.456", mockSocket{})
	testApp.Subscribe(channel, conn, "{}")

	conn = newConnection("321.654", mockSocket{})
	testApp.Subscribe(channel, conn, "{}")

	db := newMemdb()
	db.AddApp(testApp)

	ctx = &applicationContext{DB: db}
}

// All Channels
func Test_getChannels_all(t *testing.T) {

	appID := testApp.AppID

	p := map[string]string{}
	p["app_id"] = appID

	r, _ := http.NewRequest("GET", fmt.Sprintf("/apps/%s/channels", appID), nil)
	w := httptest.NewRecorder()

	getChannels(ctx, params(p), w, r)

	if w.Code != http.StatusOK {
		t.Errorf("w.Code == %d, wants %d", w.Code, http.StatusOK)
	}

	data := make(map[string]interface{})
	json.Unmarshal(w.Body.Bytes(), &data)

	channels := data["channels"].(map[string]interface{})

	if len(channels) != 3 {
		t.Errorf("len(%q) == %d, want %d", channels, len(channels), 3)
	}
}

func Test_getChannels_filter_by_presence_prefix(t *testing.T) {
	appID := testApp.AppID

	p := map[string]string{}
	p["app_id"] = appID

	r, _ := http.NewRequest("GET", fmt.Sprintf("/apps/%s/channels?filter_by_prefix=presence-", appID), nil)
	w := httptest.NewRecorder()

	getChannels(ctx, params(p), w, r)

	if w.Code != http.StatusOK {
		t.Errorf("w.Code == %d, wants %d", w.Code, http.StatusOK)
	}

	data := make(map[string]interface{})
	json.Unmarshal(w.Body.Bytes(), &data)

	channels := data["channels"].(map[string]interface{})

	if len(channels) != 1 {
		t.Errorf("len(%q) == %d, want %d", channels, len(channels), 1)
	}
}

// Only presence channels and user_count
func Test_getChannels_filter_by_presence_prefix_and_user_count(t *testing.T) {

	appID := testApp.AppID

	p := map[string]string{}
	p["app_id"] = appID

	r, _ := http.NewRequest("GET", fmt.Sprintf("/apps/%s/channels?filter_by_prefix=presence-&info=user_count", appID), nil)
	w := httptest.NewRecorder()

	getChannels(ctx, params(p), w, r)

	if w.Code != http.StatusOK {
		t.Errorf("w.Code == %d, wants %d", w.Code, http.StatusOK)
	}

	data := make(map[string]interface{})
	json.Unmarshal(w.Body.Bytes(), &data)

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

// User count only alowed in Presence channels
func Test_getChannels_filter_by_private_prefix_and_info_user_count(t *testing.T) {
	appID := testApp.AppID

	p := map[string]string{}
	p["app_id"] = appID

	r, _ := http.NewRequest("GET", fmt.Sprintf("/apps/%s/channels?filter_by_prefix=private-&info=user_count", appID), nil)
	w := httptest.NewRecorder()

	getChannels(ctx, params(p), w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("w.Code == %d, wants %d", w.Code, http.StatusBadRequest)
	}
}

func Test_getChannels_filter_by_public_prefix(t *testing.T) {
	appID := testApp.AppID

	p := map[string]string{}
	p["app_id"] = appID

	r, _ := http.NewRequest("GET", fmt.Sprintf("/apps/%s/channels?filter_by_prefix=public-", appID), nil)
	w := httptest.NewRecorder()

	getChannels(ctx, params(p), w, r)

	if w.Code != http.StatusOK {
		t.Errorf("w.Code == %d, wants %d", w.Code, http.StatusOK)
	}

	data := make(map[string]interface{})

	json.Unmarshal(w.Body.Bytes(), &data)

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

	p := map[string]string{}
	p["app_id"] = appID

	r, _ := http.NewRequest("GET", fmt.Sprintf("/apps/%s/channels?filter_by_prefix=private-", appID), nil)
	w := httptest.NewRecorder()

	getChannels(ctx, params(p), w, r)

	if w.Code != http.StatusOK {
		t.Errorf("w.Code == %d, wants %d", w.Code, http.StatusOK)
	}

	data := make(map[string]interface{})

	json.Unmarshal(w.Body.Bytes(), &data)

	channels := data["channels"].(map[string]interface{})

	if len(channels) != 1 {
		t.Errorf("len(%q) == %d, want %d", channels, len(channels), 1)
	}

	_, exists := channels["private-c3"]

	if !exists {
		t.Errorf("!exists == %t, want %t", !exists, false)
	}
}
