// Copyright (c) 2026 TAGtastic contributors
// SPDX-License-Identifier: MIT
//
// This file is part of TAGtastic and is licensed under the MIT License.

// next-codename prints the alphabetically next Crayola colour that has not
// yet been used as a TAGtastic release codename. The colour list comes from
// the embedded `crayola_colors` theme (internal/data/themes.yaml); used
// codenames are detected by parsing CHANGELOG.md.
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aenawi/tagtastic/internal/data"
)

const codenameThemeID = "crayola_colors"

func main() {
	root, err := repoRoot()
	if err != nil {
		fatal(err)
	}

	colors, err := loadCodenameTheme()
	if err != nil {
		fatal(err)
	}

	used, err := loadUsedCodenames(filepath.Join(root, "CHANGELOG.md"))
	if err != nil {
		fatal(err)
	}

	for _, color := range colors {
		if _, ok := used[color]; ok {
			continue
		}
		fmt.Println(color)
		return
	}

	fatal(fmt.Errorf("no available codenames left in theme %q", codenameThemeID))
}

func repoRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return wd, nil
}

// loadCodenameTheme returns the item display names from the embedded
// crayola_colors theme in declaration order (which is alphabetical in the
// shipped themes.yaml).
func loadCodenameTheme() ([]string, error) {
	repo, err := data.NewEmbeddedThemeRepository()
	if err != nil {
		return nil, fmt.Errorf("load themes: %w", err)
	}
	theme, err := repo.GetThemeByName(codenameThemeID)
	if err != nil {
		return nil, fmt.Errorf("load %q theme: %w", codenameThemeID, err)
	}

	names := make([]string, 0, len(theme.Items))
	for _, item := range theme.Items {
		name := strings.TrimSpace(item.Name)
		if name == "" {
			continue
		}
		names = append(names, name)
	}
	return names, nil
}

func loadUsedCodenames(path string) (map[string]struct{}, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]struct{}{}, nil
		}
		return nil, err
	}
	defer func() { _ = file.Close() }()

	used := make(map[string]struct{})
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if idx := strings.Index(line, "– \""); idx != -1 {
			fragment := line[idx+len("– \""):]
			if end := strings.Index(fragment, "\""); end != -1 {
				codename := fragment[:end]
				codename = strings.TrimSpace(codename)
				if codename != "" {
					used[codename] = struct{}{}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return used, nil
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
