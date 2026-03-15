Top-level TODOs and suggested tweaks for the climate project

Status: curated snapshot created by agent (2026-03-15)

1) CI jobs
- Add CI to run `go test ./...` on Go 1.25 and 1.26 (README mentions both).
- Add a linter/formatter step: `gofmt -s -l` and `go vet` or `staticcheck`.
- Optional: Add a TinyGo build smoke test for relevant packages.

2) Tests coverage
- Add tests/examples for packages without test files:
  - daemon (signal handling helpers) — create a smoke test that validates BreakChannel() and suppression of ^C when practicable.
  - cmd/completion — test completion output format or mark as intentionally untested in README.

3) Documentation
- Expand package-level examples for ini and paths to show file save/load and XDG behavior on Linux vs macOS.
- Add a short CONTRIBUTING.md with commit message style (matches SOUL.md guidance).

4) TinyGo compatibility
- Add notes or build tags for filesystem-dependent code used by ini/paths so consumers targeting TinyGo embedded builds know what's optional.
- Consider adding a tiny CI job that runs `tinygo build` for a minimal package (if CI runner has TinyGo available).

5) Maintenance
- Consider adding a script or Makefile targets for `fmt`, `vet`, `lint`, and `test`.
- Add a CODEOWNERS or maintainers list if multiple people will maintain the repo.

6) Misc
- Confirm license file header usage across source files if desired.
- Optionally add badges for CI and coverage to README when those exist.

