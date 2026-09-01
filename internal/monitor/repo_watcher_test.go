package monitor_test

import (
	"strings"
	"testing"

	"github.com/Lance52259/doc-draft/internal/monitor"
)

func TestToHTTPSURL(t *testing.T) {
	if got := monitor.ToHTTPSURL("acme/docs", ""); got != "https://github.com/acme/docs.git" {
		t.Fatalf("got %s", got)
	}
	got := monitor.ToHTTPSURL("acme/docs", "tok")
	if !strings.Contains(got, "x-access-token:tok@") {
		t.Fatalf("got %s", got)
	}
}
