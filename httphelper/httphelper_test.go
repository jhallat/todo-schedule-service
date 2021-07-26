package httphelper

import "testing"

func TestParseUrl(t *testing.T) {

	var result, _ = parseUrl("localhost:8080/api/test/1/subtest/ok", "test/:id/subtest/:label")
	if result["id"] != "1" {
		t.Errorf("Invalid value from 'id', expected '%s', actual '%s'", "1", result["id"])
	}
	if result["label"] != "ok" {
		t.Errorf("Invalid value from 'label', expected '%s', actual '%s'", "ok", result["label"])
	}

}
