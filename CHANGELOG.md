# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed
- **Release helper now commits `.tagtastic.yaml` alongside `CHANGELOG.md` and `VERSION`.** The previous flow updated `used_codenames` on disk after `commitRelease` had already run, leaving the audit-trail change in the working tree. `cmd/tools/release` now writes `.tagtastic.yaml` first (via `updateRepoConfig`) and stages it as part of the release commit (via a new `releaseCommitPaths` helper that gracefully omits the file when it doesn't exist).
- Backfilled `.tagtastic.yaml` with the `0.2.1: Arabian Green Bee-eater` entry that the helper failed to commit during the v0.2.1 cut.

### Added
- Two regression tests in `cmd/tools/release/main_test.go` covering `releaseCommitPaths` for the present and absent cases.

## [0.2.1] – "Arabian Green Bee-eater" – 2026-05-17

### Fixed
- **Release helper no longer writes `vUnreleased` into CHANGELOG compare URLs.** `cmd/tools/release` used to sort existing link refs descending and pick `[Unreleased]` (which sorts after every `[0.x.y]`) as the previous version, producing broken links like `.../compare/vUnreleased...v0.2.0`. The helper now filters `[Unreleased]:` out before picking the previous tag and skips any existing ref for the version being released, so re-runs produce exactly one ref per version.
- **Release helper preserves the blank line before `## [Unreleased]`** when rewriting `CHANGELOG.md`. The previous join collapsed the preamble onto the heading (`...Semantic Versioning).## [Unreleased]`), which broke markdown rendering and had to be hand-patched during the v0.2.0 cut.
- Patched the two stale `vUnreleased` lines in `CHANGELOG.md` for v0.2.0 and v0.2.0-beta.1; all link refs now point to real previous tags.

### Changed
- **Pinned gosec to v2.23.0** (latest version compatible with the Go 1.24 floor in `go.mod`). The Makefile exposes the pin as `GOSEC_VERSION` and CI inherits it via `make security`. Eliminates the local-vs-CI drift that caused a surprise CI failure during the v0.2.0 cut (CI's `@latest` install introduced rule G703 that local gosec didn't have).
- Bumped GitHub Actions to Node 24 majors: `actions/checkout` v4→v6, `actions/setup-go` v5→v6, `goreleaser/goreleaser-action` v6→v7. Addresses the Node 20 deprecation warnings that surfaced in v0.2.0 CI logs.

### Added
- Three regression tests in `cmd/tools/release/main_test.go` covering the link-ref dedup, previous-tag selection, and preamble-newline preservation. Coverage on `cmd/tools/release` rose from 20.9% to 35.1%.
- Five `#nosec` annotations (2× G115 for `os.Stdout.Fd()` int conversion, 3× G705 for JSON error formatting to stderr) suppressing false positives that surfaced when pinning gosec to v2.23.0. All match the existing project convention with one-line justifications.

## [0.2.0] – "Arabian Babbler" – 2026-05-17

### ⚠️ Breaking changes
- **Default theme switched from `crayola_colors` to `arabian_birds`** across the CLI flag defaults (`generate --theme`, `list --theme`), the config-package default (`internal/config/config.Default()`), the release helper's `.tagtastic.yaml` seed, and this repo's own `.tagtastic.yaml`. Users who relied on `tagtastic generate` (no `--theme`) returning Crayola colours must now pass `--theme crayola_colors` explicitly. The `crayola_colors` theme is unchanged and still ships embedded in the binary.

### Added
- Five new themes covering Arabian Peninsula and Arab Gulf wildlife: `arabian_mammals` (31 items), `arabian_birds` (62 items), `arabian_trees` (31 items), `arabian_reptiles` (24 items), and `arabian_marine` (30 items). Each item carries an English display name, dash-slug aliases, and a description containing the scientific name and Arabic vernacular.
- New `data/README.md` explaining the master → sync → embed flow and why both `data/themes.yaml` and `internal/data/themes.yaml` exist.

### Changed
- **Release codename theme switched from `crayola_colors` to `arabian_birds`.** Starting with v0.2.0, TAGtastic releases are named after Arabian Peninsula bird species. Earlier releases (v0.1.0-alpha.1 through v0.2.0-beta.1) used Crayola crayon colours; the legacy `crayola_colors` theme remains available to end users via `--theme crayola_colors`.
- Consolidated the Crayola colour catalogue: `cmd/tools/next-codename` and `cmd/tools/release` now read directly from the embedded `crayola_colors` theme in `internal/data/themes.yaml` instead of a parallel `data/crayola.json` snapshot.
- README "Themes" section expanded into a 9-row table; Credits section split into "Tooling" and "Theme data sources" with full attribution for the Arabian wildlife pack (NCW, EAD, UAE Atlas, UAE Flora, Flora of Arabia, Oman Open Data, Fujairah Research Centre, DDCR, Wikipedia CC BY-SA 4.0, IUCN, GBIF, Avibase, AVONET).
- `themes` JSON golden snapshot regenerated to include the five new theme IDs.

### Removed
- `data/crayola.json` — the colour catalogue lives in `data/themes.yaml` as the single source of truth.

## [0.2.0-beta.1] – "Asparagus" – 2026-01-04

### Added

- GitHub Actions CI workflow for automated quality checks on PRs and commits
- Security scanning with gosec (runs on every PR and commit)
- `make security` target for local security scanning
- CI workflow badge in README.md
- Security considerations documentation in README.md

### Changed

- Updated Go version requirement to 1.24.0 (required by dependencies)
- Updated README.md with CI/CD integration examples and security scanning details

## [0.1.1-beta.1] – "Aquamarine" – 2026-01-04

### Added
- Placeholder version entry for this release.

## [0.1.0-beta.2] – "Apricot" – 2026-01-04

### Added
- N/A

### Changed
- N/A

### Fixed
- Release helper now handles flags after the version argument correctly in dry-run mode.

### Deprecated
- N/A

### Removed
- N/A

### Security
- N/A

## [0.1.0-beta.1] – "Almond" – 2026-01-03

### Added
- Initial CLI scaffold and project documentation.
- Release helper tool to prepare releases, update files, and tag versions.
- Repo-local config support with precedence (`--config-path`, `TAGTASTIC_CONFIG`, `./.tagtastic.yaml`).
- `generate --record` to write selected codenames into repo config.
- CI-friendly flags: `--quiet`, `--json-errors`, and `--dry-run` (where applicable).
- Banner asset and repo badges (Go Report Card, release status, license, Go version).
- Quality checks via `make quality` and Go Report Card guidance.
- Expanded documentation with real-world release workflows and CI usage examples.

### Changed
- Shell output now uses aliases (slug style) with a safe fallback.
- JSON output is pretty-printed for readability.
- Banner and help behavior tuned for interactive vs. CI usage.
- Codename lookup prefers git tags, then config, then changelog.

### Fixed
- Changelog reference links corrected for beta and unreleased entries.

## [0.1.0-alpha.1] – "Antique Brass" – 2026-01-01

### Added
- Placeholder version entry for the first alpha release.







[Unreleased]: https://github.com/infravillage/tagtastic/compare/v0.2.1...HEAD
[0.2.1]: https://github.com/infravillage/tagtastic/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/infravillage/tagtastic/compare/v0.2.0-beta.1...v0.2.0
[0.2.0-beta.1]: https://github.com/infravillage/tagtastic/compare/v0.1.1-beta.1...v0.2.0-beta.1
[0.1.1-beta.1]: https://github.com/infravillage/tagtastic/compare/v0.1.0-beta.2...v0.1.1-beta.1
[0.1.0-beta.2]: https://github.com/infravillage/tagtastic/compare/v0.1.0-beta.1...v0.1.0-beta.2
[0.1.0-beta.1]: https://github.com/infravillage/tagtastic/compare/v0.1.0-alpha.1...v0.1.0-beta.1
[0.1.0-alpha.1]: https://github.com/infravillage/tagtastic/releases/tag/v0.1.0-alpha.1