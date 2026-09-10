package onepanel

import "testing"

// TestToken 校验 1Panel 签名：md5("1panel" + api_key + timestamp)。
func TestToken(t *testing.T) {
	got := token("1723622400", "apikey")
	want := "fc4b8539f694d44f0a4d580adde358cb"
	if got != want {
		t.Errorf("token() = %q, want %q", got, want)
	}
}

func TestRequestValidatesInput(t *testing.T) {
	c := &client{}
	if _, err := c.request(nil, "GET", "settings"); err == nil {
		t.Error("expected error for missing url")
	}
	c.Url = "https://127.0.0.1:8090"
	if _, err := c.request(nil, "GET", "settings"); err == nil {
		t.Error("expected error for missing api_key")
	}
}
