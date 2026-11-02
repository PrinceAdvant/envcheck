// Package envcheck compares a .env file with a template such as .env.example
// and reports keys that are missing, empty or not documented. It is meant to
// run before a deploy, so a forgotten variable fails fast instead of crashing
// the app after it has started.
package envcheck

import (
	"bufio"
	"io"
	"os"
	"sort"
	"strings"
)

// Parse reads KEY=VALUE lines. Blank lines and lines starting with # are
// ignored, an optional "export " prefix is removed, and matching single or
// double quotes around the value are stripped.
func Parse(r io.Reader) (map[string]string, error) {
	out := map[string]string{}
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		k, v = strings.TrimSpace(k), strings.TrimSpace(v)
		if len(v) >= 2 && (v[0] == '"' || v[0] == '\'') && v[len(v)-1] == v[0] {
			v = v[1 : len(v)-1]
		}
		out[k] = v
	}
	return out, s.Err()
}

// ParseFile is Parse for a file path.
func ParseFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Parse(f)
}

// Report lists problems found by Compare. Each slice is sorted.
type Report struct {
	Missing []string // in the template, absent from the env file
	Empty   []string // present but empty
	Extra   []string // in the env file, not in the template
}

// OK reports whether there are no missing or empty keys. Extra keys are
// informational and do not fail the check.
func (r Report) OK() bool { return len(r.Missing) == 0 && len(r.Empty) == 0 }

// Compare checks env against template.
func Compare(template, env map[string]string) Report {
	var r Report
	for k := range template {
		v, ok := env[k]
		switch {
		case !ok:
			r.Missing = append(r.Missing, k)
		case v == "":
			r.Empty = append(r.Empty, k)
		}
	}
	for k := range env {
		if _, ok := template[k]; !ok {
			r.Extra = append(r.Extra, k)
		}
	}
	sort.Strings(r.Missing)
	sort.Strings(r.Empty)
	sort.Strings(r.Extra)
	return r
}
