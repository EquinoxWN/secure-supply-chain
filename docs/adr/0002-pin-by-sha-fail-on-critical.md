# ADR 0002: Pin everything by immutable reference; gate only on critical vulnerabilities

- **Status:** Accepted

## Context

Two choices shape the M1 pipeline. First, how third-party code is referenced: tags such as
`actions/checkout@v4` and `golang:1.23-alpine` can be moved to new content at any time, which is
how several real supply-chain attacks spread. Second, which vulnerability severity blocks a
build: blocking on everything makes the pipeline permanently red because of unfixable base-image
findings, and blocking on nothing makes the scan decorative.

## Decision

- Every third-party action is pinned to a full 40-character commit SHA, with the release in a
  trailing comment. `pinlint` enforces this in `make lint`, so an unpinned reference fails CI.
- Base images are pinned by `@sha256` digest in the Dockerfile.
- Grype fails the build on **critical** findings. High and lower findings are counted in the job
  summary and kept in the uploaded report for triage.

## Consequences

- A moved tag or a compromised action release cannot change what CI runs without a visible diff.
- Pins go stale, so security fixes arrive only through reviewed pull requests (Dependabot).
- Some high-severity findings can ship. That is a deliberate, documented risk, revisited when
  the image is signed and deployed in M2 and M3.
