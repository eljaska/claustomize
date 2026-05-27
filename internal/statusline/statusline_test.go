package statusline

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/eljaska/claustomize/internal/block"
)

func TestGenerate_EmitsSnippetsInOrder(t *testing.T) {
	blocks := []block.Block{
		{ID: "a", Snippet: `printf 'A'`},
		{ID: "b", Snippet: `printf 'B'`},
	}
	script := Generate(blocks)

	if !strings.Contains(script, "printf 'A'") || !strings.Contains(script, "printf 'B'") {
		t.Fatalf("script missing snippets:\n%s", script)
	}
	if strings.Index(script, "printf 'A'") > strings.Index(script, "printf 'B'") {
		t.Fatalf("snippet order not preserved")
	}
}

func TestPreview_RendersSelectedBlocks(t *testing.T) {
	blocks := []block.Block{
		{ID: "model", Snippet: `printf '%s' "$(printf '%s' "$input" | jq -r '.model.display_name // empty')"`},
		{ID: "literal", Snippet: `printf 'hello'`},
	}
	got := Preview(blocks)
	if !strings.Contains(got, "Opus 4.7") {
		t.Errorf("expected preview to contain model name, got %q", got)
	}
	if !strings.Contains(got, "hello") {
		t.Errorf("expected preview to contain literal, got %q", got)
	}
	if !strings.Contains(got, " | ") {
		t.Errorf("expected separator in preview, got %q", got)
	}
}

func TestPreview_EmptySelection(t *testing.T) {
	if got := Preview(nil); !strings.Contains(got, "no blocks") {
		t.Errorf("expected empty-selection message, got %q", got)
	}
}

func TestPreview_OmitsEmptySegments(t *testing.T) {
	blocks := []block.Block{
		{ID: "a", Snippet: `printf 'A'`},
		{ID: "empty", Snippet: `printf ''`},
		{ID: "b", Snippet: `printf 'B'`},
	}
	got := Preview(blocks)
	if got != "A | B" {
		t.Errorf("expected 'A | B', got %q", got)
	}
}

func TestInstall_WritesScriptAndSettings(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	blocks := []block.Block{
		{ID: "lit", Snippet: `printf 'hi'`},
	}
	if err := Install(blocks); err != nil {
		t.Fatalf("Install: %v", err)
	}

	scriptPath := filepath.Join(tmp, ".config", "claustomize", "statusline.sh")
	info, err := os.Stat(scriptPath)
	if err != nil {
		t.Fatalf("stat script: %v", err)
	}
	if info.Mode().Perm()&0o100 == 0 {
		t.Errorf("script not executable: mode=%v", info.Mode())
	}

	settingsPath := filepath.Join(tmp, ".claude", "settings.json")
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("read settings: %v", err)
	}
	var parsed map[string]any
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("parse settings: %v", err)
	}
	sl, ok := parsed["statusLine"].(map[string]any)
	if !ok {
		t.Fatalf("statusLine missing or wrong type: %v", parsed["statusLine"])
	}
	if sl["type"] != "command" || sl["command"] != scriptPath {
		t.Errorf("unexpected statusLine: %v", sl)
	}
}

func TestInstall_PreservesExistingSettings(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	settingsPath := filepath.Join(tmp, ".claude", "settings.json")
	if err := os.MkdirAll(filepath.Dir(settingsPath), 0o755); err != nil {
		t.Fatal(err)
	}
	existing := `{"theme":"dark","statusLine":{"type":"command","command":"/old"}}`
	if err := os.WriteFile(settingsPath, []byte(existing), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Install([]block.Block{{ID: "lit", Snippet: `printf 'x'`}}); err != nil {
		t.Fatalf("Install: %v", err)
	}

	data, _ := os.ReadFile(settingsPath)
	var parsed map[string]any
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatal(err)
	}
	if parsed["theme"] != "dark" {
		t.Errorf("existing 'theme' field lost: %v", parsed)
	}
	sl := parsed["statusLine"].(map[string]any)
	if sl["command"] == "/old" {
		t.Errorf("statusLine.command not updated")
	}
}
