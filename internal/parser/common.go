package parser

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	namespaceRe = regexp.MustCompile(`(?m)^\s*namespace\s+([^;]+);`)
	classRe     = regexp.MustCompile(`(?m)^\s*class\s+([A-Za-z_][A-Za-z0-9_]*)`)
	methodRe    = regexp.MustCompile(`(?m)^\s*public\s+function\s+([A-Za-z_][A-Za-z0-9_]*)\s*\(`)
)

func readFile(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func parseNamespace(content string) string {
	m := namespaceRe.FindStringSubmatch(content)
	if len(m) < 2 {
		return ""
	}
	return strings.TrimSpace(m[1])
}

func parseClass(content, path string) string {
	m := classRe.FindStringSubmatch(content)
	if len(m) >= 2 {
		return strings.TrimSpace(m[1])
	}
	base := filepath.Base(path)
	return strings.TrimSuffix(base, filepath.Ext(base))
}
