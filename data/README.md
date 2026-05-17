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

## Release codename source

TAGtastic distinguishes two codename jobs:

1. **End-user codenames** — `tagtastic generate --theme <id>` produces a
   codename for the user's own release. The user picks the theme.
2. **TAGtastic's own release codenames** — `make codename` (and the
   release helper at [`cmd/tools/release`](../cmd/tools/release/)) picks
   the alphabetically next item from a fixed *release theme* that has
   not yet appeared in [`CHANGELOG.md`](../CHANGELOG.md). That codename
   is stamped on the GitHub release title and `git tag` annotation.

The release theme is set by the `codenameThemeID` constant in both
[`cmd/tools/next-codename/main.go`](../cmd/tools/next-codename/main.go)
and [`cmd/tools/release/main.go`](../cmd/tools/release/main.go).
Current value: **`arabian_birds`** (switched at v0.2.0). Previous
releases used `crayola_colors`.

Both tools read the chosen theme from the embedded copy of this file at
runtime — there is no separate JSON snapshot. The YAML catalogue is the
single source of truth.

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
