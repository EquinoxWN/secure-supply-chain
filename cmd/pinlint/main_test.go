package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOwnWorkflowsArePinned(t *testing.T) {
	n, err := run(filepath.Join("..", "..", ".github", "workflows"))
	if err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatalf("%d unpinned action reference(s) in this repository's workflows", n)
	}
}

func TestUnpinnedWorkflowIsCounted(t *testing.T) {
	dir := t.TempDir()
	wf := "jobs:\n  a:\n    steps:\n      - uses: actions/checkout@v4\n"
	if err := os.WriteFile(filepath.Join(dir, "bad.yml"), []byte(wf), 0o644); err != nil {
		t.Fatal(err)
	}
	if n, err := run(dir); err != nil || n != 1 {
		t.Fatalf("got n=%d err=%v, want 1 finding", n, err)
	}
}

func TestEmptyDirectoryIsAnError(t *testing.T) {
	if _, err := run(t.TempDir()); err == nil {
		t.Fatal("a directory without workflows must be reported, not silently pass")
	}
}
