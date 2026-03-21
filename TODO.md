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


-- Proposed new subpackages & tools (2026-03-20)

A. completions (refactor suggestion)
- Purpose: centralise generation of shell completion scripts (bash, zsh, fish, pwsh) for an arg.Options instance.
- Note: there's already arg/completion.go inside the arg package. Options:
  - Keep it in arg (low friction) OR extract to a top-level package `completions` if you want reuse across unrelated packages, extra features, tests and examples.
- API sketch:
  - func Generate(opt *arg.Options, shell string) (string, error)
  - func WriteFile(opt *arg.Options, shell, path string) error
- Effort: small. Low risk.

B. suggest (typo correction / suggestion)
- Purpose: when an unknown option/command is encountered, compute nearest matches and suggest them ("did you mean ...?").
- Where: top-level package `suggest` or `arg/suggest`. I recommend a small top-level `suggest` package so it can be reused by other packages.
- How it might work (internals):
  - Use Damerau–Levenshtein (costs for transposition) for edit distance.
  - Token-aware scoring: split `--some-flag` into tokens, prefer prefix and substring matches, boost matches against long names first, treat short aliases specially.
  - Support case-insensitive compare, alias resolution, and thresholding (only show suggestions with score above a minimum).
  - Optionally use heuristics: prefer same prefix, penalise distance more for short names.
- API sketch:
  - func Suggest(query string, candidates []string, max int) []string
  - func BestMatch(query string, candidates []string) (string, float64)
  - func ForOptions(opt *arg.Options, query string, max int) []string
- Integration point: call suggest.ForOptions in arg.Parse when an unknown option is found and include suggestions in the error/help text.
- Effort: small. Low risk.

C. config (env + ini + flags merger)
- Purpose: merge configuration from INI files, environment variables (with optional prefix), and flags (parsed via arg.Options) into a unified typed struct or map.
- Precedence (recommended): flags > env > ini > defaults.
- API sketch:
  - type Loader struct { Opt *arg.Options; IniPath string; EnvPrefix string }
  - func (l *Loader) Load() (map[string]any, error)
  - func (l *Loader) Populate(target any) error // populate a struct via tags `ini:"" env:"" flag:""`
- Features to consider: type conversion, slices, nested structs, required fields, default tags, and validation hooks.
- Where: top-level package `config`.
- Effort: medium. Worthwhile and highly reuseable by CLI apps.

D. validators
- Purpose: reusable validators and constraint combinators (range, regex, one-of, mutually-exclusive, required-if).
- API sketch:
  - type Validator func(value any) error
  - func Range(min, max float64) Validator
  - func Regex(re *regexp.Regexp) Validator
  - func OneOf(allowed ...string) Validator
  - func MutuallyExclusive(names ...string) Validator
  - func RequiredIf(name string, cond func() bool) Validator
- Integration: opt.SetValidator(name string, v Validator) or a central Validate(opt *arg.Options) call.
- Where: top-level `validators` package.
- Effort: small–medium. Very useful for real CLI programs.

E. helpfmt (alternate help renderers)
- Purpose: separate help formatting into backends: plain text (current), markdown, manpage, and compact synopsis.
- API sketch:
  - type Formatter interface { Format(opt *arg.Options) string }
  - NewTextFormatter(width int) Formatter
  - NewMarkdownFormatter() Formatter
  - NewManFormatter(section int) Formatter
- Effort: medium. No API break if PrintHelp continues to call the default formatter.

F. testhelpers (capture output)
- Purpose: small helpers for tests: CaptureStdout, CaptureStderr, RunWithArgs, etc.
- Effort: trivial. Low risk.

G. mangen (external tool)
- Purpose: generate manpage stubs (roff/groff) or markdown manpages for programs that use climate/arg.
- Two approaches:
  1) Static analysis (recommended first pass):
     - Use go/parser and go/ast to find calls to arg.New(...) and SetOption / SetCommand in source code.
     - Extract literals for appname, description, option names, help text, placeholders, and defaults.
     - Limitations: dynamic option construction (variables, loops, conditionals, fmt.Sprintf etc.) won't be fully captured.
  2) Runtime metadata export (robust):
     - Add a small opt.Metadata() helper in arg that returns a serializable descriptor (JSON) of the application options/commands.
     - Run the target binary with a special flag (e.g. `--_dump_arg_meta=json`) and have mangen invoke the binary to get full, accurate metadata.
     - This approach handles dynamic option registration and computed help text.
- Suggested CLI for mangen:
  - mangen generate --source ./cmd/myapp -o docs/myapp.1 --format roff
  - mangen scan --pkg ./... --out manpages/  (static mode)
  - mangen from-binary --binary ./bin/myapp --out docs/
- Output: roff manpage stubs, markdown stubs, and optional README snippets.
- Effort: medium–high depending on chosen approach. Runtime metadata approach is more robust and easier to implement reliably if arg exposes metadata.
- UX suggestion: add an opt.Metadata() or opt.DumpMetadata(writer) in arg to make this tool trivial and exact.


Next steps (if you want me to act now)
- I can append/merge these notes into TODO.md (done — content added), but I will not create branches or commits beyond this file unless you tell me to.
- If you want, I can scaffold any of these packages (create directories, initial code, tests, README) and run `go test ./...` — say which one to start with.
- For mangen I can prototype a small tool that either (A) parses source for literal-registered options, or (B) uses a metadata hook; tell me which approach you prefer.

Questions for you
- Do you want `suggest` to live under `arg` (arg/suggest) or as a reuseable top-level `suggest` package?
- For `config` precedence do you agree flags > env > ini > defaults?
- For `mangen`, do you prefer static analysis-only initially, or should we add a small metadata-export hook to arg to enable robust generation?

No commits were made beyond updating TODO.md as requested; I did not create branches or push changes. Let me know how you want to proceed.