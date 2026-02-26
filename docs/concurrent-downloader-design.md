
# Concurrent Downloader

## 1. Problem Statement

Design a concurrent file downloader that retrieves multiple URLs efficiently while maintaining bounded resource usage, supporting cancellation, and handling failures safely.

## 2. Goals

* Efficient parallel downloads
* Bounded concurrency
* Deterministic resource usage
* Graceful cancellation
* Safe partial failure handling
* Clean shutdown without goroutine or FD leaks

## 3. Non-Goals

* Distributed execution across multiple machines
* Persistent job storage
* Automatic resume of partial downloads
* Advanced retry orchestration
* UI or progress dashboard

## 4. Requirements

### 4.1 Functional

* **FR-1:** Accept a list of URLs as input.
* **FR-2:** Download each URL to a specified directory.
* **FR-3:** Limit concurrent downloads to a configurable worker count.
* **FR-4:** Support cancellation via context.
* **FR-5:** Surface download errors to caller.
* **FR-6:** Optionally delete partially downloaded files on failure.

### 4.2 Non-Functional

* **NFR-1:** Memory usage must scale with worker count, not total URLs.
* **NFR-2:** No unbounded goroutine creation.
* **NFR-3:** All file descriptors must be closed deterministically.
* **NFR-4:** System must tolerate slow or hanging servers via timeouts.
* **NFR-5:** Shutdown must not leave orphaned workers.

## 5. Proposed Design

Architecture:

Input -> Job Producer -> Buffered Channel -> Worker Pool (N) -> HTTP Download -> File Write

Key decisions:

* Fixed-size worker pool for bounded concurrency.
* Shared `http.Client` for connection reuse.
* Streaming (`io.Copy`) from response body to file to avoid buffering entire file.
* Context propagation for cancellation.
* Error channel to collect worker errors.
* WaitGroup for lifecycle synchronization.

Concurrency ceiling:

```
Max Concurrent Downloads = Worker Count
```

## 6. Alternatives Considered

**A1: Spawn 1 goroutine per URL**
Rejected: Unbounded concurrency risks FD exhaustion and memory pressure.


**A2: Dynamic worker scaling**
Deferred: Adds complexity; static pool sufficient for CLI scope.

## 7. Failure Modes

* **Slow server:** Mitigated via request timeouts.
* **Partial write:** Optional file deletion.
* **Network drop:** Error propagated to caller.
* **Worker panic:** Could add recover block (optional improvement).
* **Channel blockage:** Bounded buffer provides backpressure.

## 8. Tradeoffs

* Static worker pool simplifies reasoning but may not maximize throughput under variable latency.
* Fail-fast simplifies control flow but may stop valid downloads.
* Streaming reduces memory use but complicates retry mid-transfer.

## 9. Observability

Current:

* Error reporting to caller.

Potential improvements:

* Structured logs
* Download duration metrics
* Throughput measurement



