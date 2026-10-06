package router

import (
	"api-gateway-module/app/client"
	"api-gateway-module/config"
	gatewayhttp "api-gateway-module/types/http"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouterForwardsConfiguredGet(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/items" || r.URL.Query().Get("name") != "sample" {
			t.Errorf("unexpected downstream URL: %s", r.URL)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	t.Cleanup(backend.Close)

	var cfg config.App
	cfg.App.Name = "test"
	cfg.Http.BaseUrl = backend.URL
	cfg.Http.Router = []config.Router{{
		Method: gatewayhttp.GET, GetType: gatewayhttp.QUERY,
		Path: "/items", Variable: []string{"name"},
	}}
	downstream := client.NewHttpClient(cfg, nil)
	router := NewRouter(cfg, map[string]*client.HttpClient{"test": downstream})
	response, err := router.engine.Test(httptest.NewRequest(http.MethodGet, "/items?name=sample", nil))
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", response.StatusCode)
	}
	var body string
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body != `{"ok":true}` {
		t.Fatalf("body = %q", body)
	}
}
