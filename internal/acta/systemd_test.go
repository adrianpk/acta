package acta

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestScheduledVerificationAssets(t *testing.T) {
	t.Parallel()

	service := readRepositoryAsset(t, "acta-verify.service")
	timer := readRepositoryAsset(t, "acta-verify.timer")

	serviceValid := strings.Contains(service, "ExecStart=/usr/bin/env acta verify --scheduled") &&
		strings.Contains(service, "Requires=acta.socket") &&
		strings.Contains(service, "Type=oneshot")
	if !serviceValid {
		t.Fatalf("verification service = %q", service)
	}

	timerValid := strings.Contains(timer, "OnCalendar=*-*-* 03:00:00") &&
		strings.Contains(timer, "Persistent=true") &&
		strings.Contains(timer, "Unit=acta-verify.service")
	if !timerValid {
		t.Fatalf("verification timer = %q", timer)
	}
}

func readRepositoryAsset(t *testing.T, name string) string {
	t.Helper()

	path := filepath.Join("..", "..", "assets", "systemd", "user", name)

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read asset %q: %v", name, err)
	}

	return string(content)
}
