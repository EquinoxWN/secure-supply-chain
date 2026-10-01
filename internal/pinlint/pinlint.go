// Package pinlint finds GitHub Actions references that are not pinned to an immutable commit.
package pinlint

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

var (
	usesLine = regexp.MustCompile(`^\s*(?:-\s*)?uses:\s*["']?([^"'\s#]+)`)
	pinned   = regexp.MustCompile(`^[\w.-]+/[\w.-]+(?:/[\w./-]+)?@[0-9a-f]{40}$`)
	digest   = regexp.MustCompile(`^docker://[^@]+@sha256:[0-9a-f]{64}$`)
)

// Finding is one unpinned reference.
type Finding struct {
	File string
	Line int
	Ref  string
	Why  string
}

// String formats a finding like a compiler error.
func (f Finding) String() string {
	return fmt.Sprintf("%s:%d: %s: %s", f.File, f.Line, f.Ref, f.Why)
}

// classify returns why ref is unsafe, or "" when it is pinned or local.
func classify(ref string) string {
	switch {
	case strings.HasPrefix(ref, "./"):
		return "" // local action, versioned with this repository
	case strings.HasPrefix(ref, "docker://"):
		if digest.MatchString(ref) {
			return ""
		}
		return "container action must be pinned by @sha256 digest"
	case pinned.MatchString(ref):
		return ""
	case !strings.Contains(ref, "@"):
		return "missing version; pin to a full 40-character commit SHA"
	}
	return "tag or branch is mutable; pin to a full 40-character commit SHA"
}

// Check scans one workflow file and returns every unpinned `uses:` reference.
func Check(r io.Reader, name string) ([]Finding, error) {
	var out []Finding
	sc := bufio.NewScanner(r)
	for n := 1; sc.Scan(); n++ {
		m := usesLine.FindStringSubmatch(sc.Text())
		if m == nil {
			continue
		}
		if why := classify(m[1]); why != "" {
			out = append(out, Finding{File: name, Line: n, Ref: m[1], Why: why})
		}
	}
	return out, sc.Err()
}
