package statusline

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/eljaska/claustomize/internal/slot"
	"github.com/eljaska/claustomize/internal/statusline/blocks"
)

func singleStyle(id, snippet string) *blocks.Block {
	return &blocks.Block{
		ID:     id,
		Name:   id,
		Styles: []blocks.Style{{ID: "only", Name: "only", Snippet: snippet}},
	}
}

// build constructs a slot list from a sequence of filled blocks, inserting
// the required empties to satisfy the invariant.
func build(items ...*blocks.Block) slot.List {
	l := slot.New()
	for _, b := range items {
		l, _ = l.Fill(len(l)-1, b, 0)
	}
	return l
}

func TestGenerate_EmitsSnippetsInOrder(t *testing.T) {
	a := singleStyle("a", `printf 'A'`)
	b := singleStyle("b", `printf 'B'`)
	script := Generate(build(a, b))

	if !strings.Contains(script, "printf 'A'") || !strings.Contains(script, "printf 'B'") {
		t.Fatalf("script missing snippets:\n%s", script)
	}
	if strings.Index(script, "printf 'A'") > strings.Index(script, "printf 'B'") {
		t.Fatalf("snippet order not preserved")
	}
}

func TestPreview_RendersFilledBlocks(t *testing.T) {
	model := singleStyle("model", `printf '%s' "$(printf '%s' "$input" | jq -r '.model.display_name // empty')"`)
	lit := singleStyle("lit", `printf 'hello'`)
	got := Preview(build(model, lit))
	if !strings.Contains(got, "Opus 4.7") {
		t.Errorf("expected preview to contain model name, got %q", got)
	}
	if !strings.Contains(got, "hello") {
		t.Errorf("expected preview to contain literal, got %q", got)
	}
}

func TestPreview_NoFilledSlots(t *testing.T) {
	if got := Preview(slot.New()); got != "" {
		t.Errorf("expected empty preview for empty slots, got %q", got)
	}
}

func TestPreview_NoAutomaticSeparator(t *testing.T) {
	a := singleStyle("a", `printf 'A'`)
	b := singleStyle("b", `printf 'B'`)
	got := Preview(build(a, b))
	if got != "AB" {
		t.Errorf("expected concatenation 'AB' with no auto-separator, got %q", got)
	}
}

func TestPreview_UsesExplicitSeparatorBlock(t *testing.T) {
	a := singleStyle("a", `printf 'A'`)
	sep := singleStyle("sep", `printf ' | '`)
	b := singleStyle("b", `printf 'B'`)
	got := Preview(build(a, sep, b))
	if got != "A | B" {
		t.Errorf("expected 'A | B', got %q", got)
	}
}

func TestInstall_WritesScriptAndSettings(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	if err := Install(build(singleStyle("lit", `printf 'hi'`))); err != nil {
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

	if err := Install(build(singleStyle("lit", `printf 'x'`))); err != nil {
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
