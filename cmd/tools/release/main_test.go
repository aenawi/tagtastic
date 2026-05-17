// Copyright (c) 2026 InfraVillage
// SPDX-License-Identifier: MIT
//
// This file is part of TAGtastic and is licensed under the MIT License.
package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUpdateRepoConfigCreatesFile(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, ".tagtastic.yaml")

	if err := updateRepoConfig(path, "Almond", "0.1.0-beta.1"); err != nil {
		t.Fatalf("updateRepoConfig failed: %v", err)
	}

	payload, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	content := string(payload)

	if !strings.Contains(content, "default_theme: arabian_birds") {
		t.Fatalf("expected default_theme to be set to arabian_birds")
	}
	if !strings.Contains(content, "0.1.0-beta.1: Almond") {
		t.Fatalf("expected codename to be recorded")
	}
}

func TestResolveConfigPathOverride(t *testing.T) {
	tmp := t.TempDir()
	override := filepath.Join(tmp, "custom.yaml")

	path, err := resolveConfigPath(tmp, override)
	if err != nil {
		t.Fatalf("resolveConfigPath failed: %v", err)
	}
	if path == "" || !strings.HasSuffix(path, "custom.yaml") {
		t.Fatalf("expected override path, got %q", path)
	}
}

func TestEnsureSemVerForward(t *testing.T) {
	if err := ensureSemVerForward("0.1.1", "0.1.0"); err != nil {
		t.Fatalf("expected forward version, got error: %v", err)
	}
	if err := ensureSemVerForward("0.1.0", "0.1.0"); err == nil {
		t.Fatalf("expected error for non-forward version")
	}
}

func TestBumpVersion(t *testing.T) {
	version, err := bumpVersion("0.1.0-beta.2", "patch")
	if err != nil {
		t.Fatalf("bumpVersion failed: %v", err)
	}
	if version != "0.1.1" {
		t.Fatalf("expected 0.1.1, got %q", version)
	}
}

func TestResolvePreReleaseVersionAutoNum(t *testing.T) {
	tags := []string{"v0.1.1-beta.1", "v0.1.1-beta.2", "v0.1.1-rc.1"}
	version, err := resolvePreReleaseVersionWithTags("0.1.1", "beta", 0, tags)
	if err != nil {
		t.Fatalf("resolvePreReleaseVersionWithTags failed: %v", err)
	}
	if version != "0.1.1-beta.3" {
		t.Fatalf("expected 0.1.1-beta.3, got %q", version)
	}
}

func TestResolvePreReleaseVersionManualNum(t *testing.T) {
	version, err := resolvePreReleaseVersionWithTags("0.1.1", "rc", 5, nil)
	if err != nil {
		t.Fatalf("resolvePreReleaseVersionWithTags failed: %v", err)
	}
	if version != "0.1.1-rc.5" {
		t.Fatalf("expected 0.1.1-rc.5, got %q", version)
	}
}

func TestResolvePreReleaseVersionInvalidLabel(t *testing.T) {
	if _, err := resolvePreReleaseVersionWithTags("0.1.1", "preview", 0, nil); err == nil {
		t.Fatalf("expected error for invalid prerelease label")
	}
}

// TestReleaseCommitPaths_IncludesTagtasticYaml is a regression test for the
// missing .tagtastic.yaml in release commits: commitRelease used to stage
// only CHANGELOG.md and VERSION, leaving the used_codenames audit-trail
// update uncommitted in the working tree. Verifies that when .tagtastic.yaml
// exists at the repo root it is included in the list of paths to stage.
func TestReleaseCommitPaths_IncludesTagtasticYaml(t *testing.T) {
	tmp := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmp, ".tagtastic.yaml"), []byte("default_theme: x\n"), 0o644); err != nil {
		t.Fatalf("seed config: %v", err)
	}

	paths := releaseCommitPaths(tmp)

	want := map[string]bool{"CHANGELOG.md": false, "VERSION": false, ".tagtastic.yaml": false}
	for _, p := range paths {
		if _, ok := want[p]; ok {
			want[p] = true
		}
	}
	for p, found := range want {
		if !found {
			t.Fatalf("expected releaseCommitPaths to include %q, got %v", p, paths)
		}
	}
}

// TestReleaseCommitPaths_OmitsTagtasticYamlWhenAbsent guards the inverse:
// projects without a repo-local .tagtastic.yaml (e.g. user-global config
// via --config) must not have a non-existent path forced into git add.
func TestReleaseCommitPaths_OmitsTagtasticYamlWhenAbsent(t *testing.T) {
	tmp := t.TempDir() // no .tagtastic.yaml seeded

	paths := releaseCommitPaths(tmp)

	for _, p := range paths {
		if p == ".tagtastic.yaml" {
			t.Fatalf("releaseCommitPaths must not include .tagtastic.yaml when the file is absent; got %v", paths)
		}
	}
	if len(paths) < 2 {
		t.Fatalf("expected at least CHANGELOG.md and VERSION in paths, got %v", paths)
	}
}

// TestUpdateChangelog_PreservesBlankLineBeforeUnreleased is a regression test
// for the preamble-newline bug: splitChangelog's preamble ends in a single
// '\n', which when concatenated against unreleasedTemplate produced
// "Semantic Versioning].## [Unreleased]" on one line. Verifies the rewritten
// file has a blank line between the preamble and the [Unreleased] heading.
func TestUpdateChangelog_PreservesBlankLineBeforeUnreleased(t *testing.T) {
	const seed = `# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Feature X.

## [0.1.0] – "Almond" – 2026-01-01

### Added
- Initial release.

[Unreleased]: https://github.com/infravillage/tagtastic/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/infravillage/tagtastic/releases/tag/v0.1.0
`

	tmp := t.TempDir()
	path := filepath.Join(tmp, "CHANGELOG.md")
	if err := os.WriteFile(path, []byte(seed), 0o644); err != nil {
		t.Fatalf("seed changelog: %v", err)
	}

	if err := updateChangelog(path, "0.2.0", "Arabian Babbler", "2026-05-17"); err != nil {
		t.Fatalf("updateChangelog failed: %v", err)
	}

	out, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read updated changelog: %v", err)
	}
	got := string(out)

	// The bug produced "Semantic Versioning).## [Unreleased]" on one line.
	if strings.Contains(got, ").## [Unreleased]") {
		t.Fatalf("preamble collapsed onto Unreleased heading without blank line:\n%s", got)
	}
	// Positive assertion: there must be a blank line between the preamble's
	// last sentence and the next ## heading.
	const wantSeparator = "v2.0.0.html).\n\n## [Unreleased]"
	if !strings.Contains(got, wantSeparator) {
		t.Fatalf("expected blank line before ## [Unreleased]; output:\n%s", got)
	}
	// Sanity: the new release header for v0.2.0 must be present.
	if !strings.Contains(got, `## [0.2.0] – "Arabian Babbler" – 2026-05-17`) {
		t.Fatalf("expected new release header in output:\n%s", got)
	}
}

// TestUpdateReferenceLines_PrevVersionIsRealTag is a regression test for the
// "vUnreleased" bug: updateReferenceLines used to pick the lexically-largest
// existing ref as the previous version, which selected "[Unreleased]" and
// produced compare URLs like ".../compare/vUnreleased...vX.Y.Z". The fix
// filters [Unreleased] out before sorting.
func TestUpdateReferenceLines_PrevVersionIsRealTag(t *testing.T) {
	existing := []string{
		"[Unreleased]: https://github.com/infravillage/tagtastic/compare/v0.2.0-beta.1...HEAD",
		"[0.2.0-beta.1]: https://github.com/infravillage/tagtastic/compare/v0.1.1-beta.1...v0.2.0-beta.1",
		"[0.1.1-beta.1]: https://github.com/infravillage/tagtastic/compare/v0.1.0-beta.2...v0.1.1-beta.1",
	}

	refs := updateReferenceLines(existing, "0.2.0")

	joined := strings.Join(refs, "\n")
	if strings.Contains(joined, "vUnreleased") {
		t.Fatalf("refs must not contain literal vUnreleased, got:\n%s", joined)
	}
	wantRelease := "[0.2.0]: https://github.com/infravillage/tagtastic/compare/v0.2.0-beta.1...v0.2.0"
	if !strings.Contains(joined, wantRelease) {
		t.Fatalf("expected new release ref %q in output, got:\n%s", wantRelease, joined)
	}
	wantUnreleased := "[Unreleased]: https://github.com/infravillage/tagtastic/compare/v0.2.0...HEAD"
	if !strings.Contains(joined, wantUnreleased) {
		t.Fatalf("expected updated unreleased ref %q in output, got:\n%s", wantUnreleased, joined)
	}
}

// TestUpdateReferenceLines_DedupesUnreleasedAndCurrentVersion guards against
// the accumulating-duplicates bug: re-running the helper on a changelog that
// already contains the target version's ref or stale [Unreleased] refs must
// produce exactly one of each.
func TestUpdateReferenceLines_DedupesUnreleasedAndCurrentVersion(t *testing.T) {
	existing := []string{
		"[Unreleased]: https://github.com/infravillage/tagtastic/compare/v0.2.0...HEAD",
		"[Unreleased]: https://github.com/infravillage/tagtastic/compare/v0.2.0-beta.1...HEAD",
		"[0.2.0]: https://github.com/infravillage/tagtastic/compare/vUnreleased...v0.2.0",
		"[0.2.0-beta.1]: https://github.com/infravillage/tagtastic/compare/v0.1.1-beta.1...v0.2.0-beta.1",
	}

	refs := updateReferenceLines(existing, "0.2.0")

	unreleasedCount := 0
	versionCount := 0
	for _, line := range refs {
		if strings.HasPrefix(line, "[Unreleased]:") {
			unreleasedCount++
		}
		if strings.HasPrefix(line, "[0.2.0]:") {
			versionCount++
		}
	}
	if unreleasedCount != 1 {
		t.Fatalf("expected exactly 1 [Unreleased]: ref, got %d:\n%s", unreleasedCount, strings.Join(refs, "\n"))
	}
	if versionCount != 1 {
		t.Fatalf("expected exactly 1 [0.2.0]: ref, got %d:\n%s", versionCount, strings.Join(refs, "\n"))
	}
}
