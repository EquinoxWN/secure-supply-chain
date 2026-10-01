package pinlint

import (
	"strings"
	"testing"
)

const sha = "11d5960a326750d5838078e36cf38b85af677262"

func TestClassify(t *testing.T) {
	cases := []struct {
		ref  string
		safe bool
	}{
		{"actions/checkout@" + sha, true},
		{"github/codeql-action/upload-sarif@" + sha, true},
		{"./.github/actions/local", true},
		{"docker://alpine@sha256:" + strings.Repeat("a", 64), true},
		{"actions/checkout@v4", false},
		{"actions/checkout@main", false},
		{"actions/checkout", false},
		{"actions/checkout@" + sha[:7], false},              // short SHA is ambiguous
		{"actions/checkout@" + strings.ToUpper(sha), false}, // tags can look like SHAs; require lowercase hex
		{"docker://alpine:3.20", false},
		{"docker://alpine@sha256:abc", false},
	}
	for _, c := range cases {
		if got := classify(c.ref) == ""; got != c.safe {
			t.Errorf("classify(%q) safe=%v, want %v", c.ref, got, c.safe)
		}
	}
}

func TestCheckReportsLineNumbers(t *testing.T) {
	workflow := `name: ci
jobs:
  build:
    steps:
      - uses: actions/checkout@` + sha + ` # v4.4.0
      - uses: actions/setup-go@v5
      - name: quoted
        uses: "docker/build-push-action@v6"
      - run: echo "uses: not/an-action@v1"
`
	got, err := Check(strings.NewReader(workflow), "ci.yml")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("want 2 findings, got %v", got)
	}
	if got[0].Line != 6 || got[0].Ref != "actions/setup-go@v5" {
		t.Errorf("first finding = %v", got[0])
	}
	if got[1].Line != 8 || got[1].Ref != "docker/build-push-action@v6" {
		t.Errorf("second finding = %v", got[1])
	}
	if !strings.Contains(got[0].String(), "ci.yml:6:") {
		t.Errorf("String() = %q", got[0].String())
	}
}
