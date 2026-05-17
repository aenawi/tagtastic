# `data/` — Theme master data

This folder is the **human-editable source** of TAGtastic's theme catalogue.
It contains one file:

- [`themes.yaml`](themes.yaml) — every theme TAGtastic ships (Crayola colors,
  birds, cities, landmarks, and the Arabian wildlife pack).

## Why this folder exists

Go's `//go:embed` directive can only read files **inside the same package
directory**. So [`internal/data/repository.go`](../internal/data/repository.go)
embeds [`internal/data/themes.yaml`](../internal/data/themes.yaml), not this
file.

To keep theme data discoverable at the repo root (where datasets normally
live), the project maintains two copies:

```text
data/themes.yaml             ← edit here (master)
        │
        │  make sync-themes  (copies the file verbatim)
        ▼
internal/data/themes.yaml    ← embedded into the binary by //go:embed
        │
        │  go build
        ▼
bin/tagtastic                ← what end users run
```

The compiled binary depends **only on the embedded copy**.

## Editing themes

1. Edit `data/themes.yaml`.
2. Run `make sync-themes` (or `go run ./cmd/tools/sync-themes`).
3. Run `make build && make test` to validate.

Schema shape (see [`internal/data/types.go`](../internal/data/types.go)):

```yaml
version: "1.0"
themes:
  your_theme:
    id: your_theme
    name: "Your Theme Name"
    description: "Short description."
    category: "Category"
    items:
      - name: "Item One"
        aliases: ["item-one", "alt-slug"]
        description: "Short label or fact."
```

Rules:

- Theme IDs are **snake_case** and unique across the file.
- The first alias must equal the dash-lowercase slug of `name` so the
  `shell` output formatter produces a predictable env-var value.
- `category` is free-text but kept consistent with sibling themes
  (`Colors`, `Nature`, `Places`).

## Release codename source (`crayola_colors`)

TAGtastic uses one theme for two jobs:

1. **End-user codenames** — `tagtastic generate --theme crayola_colors`
   produces a codename for the user's release.
2. **TAGtastic's own release codenames** — `make codename` (and the
   release helper at [`cmd/tools/release`](../cmd/tools/release/)) picks the
   alphabetically next colour that has not yet appeared in
   [`CHANGELOG.md`](../CHANGELOG.md). That codename is stamped on the
   GitHub release title and `git tag` annotation.

Both jobs read from the same `crayola_colors` theme inside this file (via
the embedded copy at runtime). There is no separate Crayola JSON file —
the catalogue is the single source of truth.

## Adding the Arabian wildlife pack

The Arabian themes shipped under `arabian_mammals`, `arabian_birds`,
`arabian_trees`, `arabian_reptiles`, and `arabian_marine` were imported
from the corpora documented at
[`docs/arabian-wildlife-datasets/`](../../docs/arabian-wildlife-datasets/).

To refresh them, regenerate from the rich dataset
(`docs/arabian-wildlife-datasets/datasets/data/biodiversity/arabian_wildlife.json`)
and re-merge the per-group blocks into `themes.yaml`. See
[`docs/arabian-wildlife-datasets/01-tagtastic-integration.md`](../../docs/arabian-wildlife-datasets/01-tagtastic-integration.md)
for the spec-shaped procedure.
