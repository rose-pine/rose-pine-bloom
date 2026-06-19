<p align="center">
    <img src="https://github.com/rose-pine/rose-pine-theme/raw/main/assets/icon.png" width="80" />
    <h2 align="center">Rosé Pine Bloom</h2>
</p>

<p align="center">All natural pine, faux fur and a bit of soho vibes for the classy minimalist</p>

---

Bloom is an opinionated theme generator, matching the Rosé Pine style guide.

## Install

### Homebrew

```sh
brew install rose-pine/tap/bloom
```

### Other methods

```sh
# Goblin (zero-dependency)
curl -sf http://goblin.run/github.com/rose-pine/rose-pine-bloom | OUT=bloom sh

# Go
go install github.com/rose-pine/rose-pine-bloom@latest

# Arch Linux (community AUR)
yay -S rose-pine-bloom
```

Pre-built binaries are also available on the [releases page](https://github.com/rose-pine/rose-pine-bloom/releases).

## Usage

Create a template:

```yaml
# template.yaml
name: $name
background: $base
foreground: $rose
```

Build it:

```sh
bloom build template.yaml
```

If you already have a theme, convert it with `bloom init theme.json`. You can also build from a directory: `bloom build templates/`.

## Templates

### Variables

By default, variables are prefixed with `$`.

| Variable       | Resolves to                                                                  |
| -------------- | ---------------------------------------------------------------------------- |
| `$id`          | `rose-pine`, `rose-pine-moon`, `rose-pine-dawn`                              |
| `$name`        | `Rosé Pine`, `Rosé Pine Moon`, `Rosé Pine Dawn`                              |
| `$appearance`  | `dark`, `dark`, `light`                                                      |
| `$description` | All natural pine, faux fur and a bit of soho vibes for the classy minimalist |

Every colour in the [Rosé Pine palette](https://rosepinetheme.com/palette) is available as a variable — `$base`, `$surface`, `$overlay`, `$muted`, `$subtle`, `$text`, `$love`, `$gold`, `$rose`, `$pine`, `$foam`, `$iris`, `$highlightLow`, `$highlightMed`, `$highlightHigh`. Control opacity by appending a value, e.g. `$love/10` for 10% opacity.

### Accents

Using `$accent` generates variants for each accent colour. The accent name is appended to the filename, e.g. `rose-pine-gold.yaml`.

| Variable      | Description                         |
| ------------- | ----------------------------------- |
| `$accent`     | Accent colour value                 |
| `$onaccent`   | Contrasting foreground colour       |
| `$accentname` | Lowercase accent name (e.g. `gold`) |

### Variant values

For variant-specific values, use the `$(main|moon|dawn)` syntax. Variables are also allowed inside the variant values.

- `priority: $(10|20|30)` → `priority: 10` in rose-pine, `20` in rose-pine-moon, `30` in rose-pine-dawn
- `background: $($rose|$pine|$gold)` → `background: #ebbcba` in rose-pine, `#3e8fb0` in rose-pine-moon, `#ea9d34` in rose-pine-dawn

## Options

### Prefix

Change variable prefix:

```sh
bloom build template.json --prefix @
```

### Output

Change the output destination:

```sh
bloom build template.json --out themes
```

### Format

Specify one of the supported formats:

```sh
bloom build template.json --format <format>
```

Available formats:

| Name        | Example              |
| ----------- | -------------------- |
| `hex`       | `#ebbcba`            |
| `hsl`       | `hsl(2, 55%, 83%)`   |
| `hsl-css`   | `hsl(2deg 55% 83%)`  |
| `hsl-array` | `[2, 0.55, 0.83]`    |
| `rgb`       | `rgb(235, 188, 186)` |
| `rgb-css`   | `rgb(235 188 186)`   |
| `rgb-array` | `[235, 188, 186]`    |
| `ansi`      | `235;188;186`        |

Commas and spaces can be removed by passing `--no-commas` and `--no-spaces`. Decorators (#, rgb(), hsl(), brackets) can be removed by passing `--plain`.

## Contributing

We welcome and appreciate contributions of any kind. Please create an issue for any proposed changes.
