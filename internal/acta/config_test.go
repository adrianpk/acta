package acta

import (
	"path/filepath"
	"testing"
)

func TestConfigFromHome(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		home    string
		want    Config
		wantErr bool
	}{
		{
			name: "user paths",
			home: "/home/tester",
			want: Config{
				ConfigRoot:  filepath.FromSlash("/home/tester/.config/acta"),
				FilesRoot:   filepath.FromSlash("/home/tester/.acta"),
				JournalRoot: filepath.FromSlash("/home/tester/.acta/events"),
				PendingPath: filepath.FromSlash("/home/tester/.acta/pending.json"),
			},
		},
		{
			name:    "relative home",
			home:    "home/tester",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := ConfigFromHome(test.home)
			if test.wantErr {
				if err == nil {
					t.Fatal("expected an error")
				}

				return
			}

			if err != nil {
				t.Fatalf("ConfigFromHome() error = %v", err)
			}

			if got != test.want {
				t.Fatalf("ConfigFromHome() = %#v, want %#v", got, test.want)
			}
		})
	}
}

func TestDefaultConfigIncludesRuntimeSocket(t *testing.T) {
	home := t.TempDir()
	runtimeDirectory := filepath.Join(t.TempDir(), "runtime")
	t.Setenv("HOME", home)
	t.Setenv("XDG_RUNTIME_DIR", runtimeDirectory)

	config, err := DefaultConfig()
	if err != nil {
		t.Fatalf("DefaultConfig() error = %v", err)
	}

	if config.FilesRoot != filepath.Join(home, ".acta") {
		t.Fatalf("FilesRoot = %q", config.FilesRoot)
	}

	wantRuntimeRoot := filepath.Join(runtimeDirectory, "acta")
	if config.RuntimeRoot != wantRuntimeRoot || config.SocketPath != filepath.Join(wantRuntimeRoot, "acta.sock") {
		t.Fatalf("runtime config = %#v", config)
	}
}

func TestDefaultConfigRequiresRuntimeDirectory(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("XDG_RUNTIME_DIR", "")

	if _, err := DefaultConfig(); err == nil {
		t.Fatal("expected runtime directory error")
	}
}
