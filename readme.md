# Concurrent downloader

A production-oriented CLI tool that downloads multiple files concurrently using Go’s concurrency primitives.

[ In progress ]

Bounded worker pool (configurable concurrency)
Context-based timeout and cancellation
Streaming downloads (io.Copy) — no full file buffering
Atomic file writes (temp file + rename)
Error aggregation and proper exit codes
Unit and integration tests using standard library only
Race-detector clean (go test -race)
