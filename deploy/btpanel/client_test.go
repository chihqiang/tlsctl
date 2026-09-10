package btpanel

import "testing"

// TestSign 校验宝塔开放式接口签名：md5(timestamp + md5(api_key))。
func TestSign(t *testing.T) {
	got := sign(1723622400, "apikey")
	want := "4066384c49342da772fabc24c207e322"
	if got != want {
		t.Errorf("sign() = %q, want %q", got, want)
	}
}

// TestRequestValidatesInput 校验 request 在发起网络请求前先校验必填项。
func TestRequestValidatesInput(t *testing.T) {
	c := &client{}
	if _, err := c.request(nil, "ssl/cert/save_cert"); err == nil {
		t.Error("expected error for missing url")
	}
	c.Url = "https://127.0.0.1:8888"
	if _, err := c.request(nil, "ssl/cert/save_cert"); err == nil {
		t.Error("expected error for missing api_key")
	}
}
