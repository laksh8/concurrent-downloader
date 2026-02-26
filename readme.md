# Concurrent downloader

A production-oriented CLI tool that downloads multiple files concurrently using Go’s concurrency primitives.

## Design Document

**See /docs/concurrent-downloader-design.md for full architecture details.**

## Usage

### Build

```bash
go build -o downloader ./cmd/downloader
```

---

### Basic Execution

```bash
./downloader -limit 5 -output ./downloads https://example.com/file1.jpg https://example.com/file2.jpg
```

---

### Flags

| Flag       | Description                                  | Default |
| ---------- | -------------------------------------------- | ------- |
| `-limit`   | Maximum number of concurrent downloads       | `10`    |
| `-output`  | Destination directory for downloaded files   | `./`    |
| `-keep`  | Keep partially downloaded files on failure | `false` |
| `-timeout` | Global timeout for all downloads             | `30s`   |

---

### Example

Download multiple files with controlled concurrency:

```bash
./downloader \
  -limit 10 \
  -output ./downloads \
  -delete \
  https://example.com/a.png \
  https://example.com/b.png \
  https://example.com/c.png
```

---

### Behavior

* Downloads run concurrently up to the specified limit.
* Files are streamed directly to disk (no full buffering in memory).
* If `-keep` is enabled, partial files are kept on failure.
* If the global timeout is reached, in-flight downloads are cancelled.

---

### Exit Codes

| Code | Meaning                              |
| ---- | ------------------------------------ |
| `0`  | All downloads completed successfully |
| `1`  | One or more downloads failed         |
| `2`  | Invalid input or configuration       |

