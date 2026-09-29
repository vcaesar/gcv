# gcv

GoCV/OpenCV image matching: templates first, optional SIFT fallback via `Sift`. `yolo/` is a placeholder, not an implementation.

## Layout and checks

- `cv.go`: matching APIs, options, results, template/SIFT logic; `fn.go`: conversion, I/O, transforms, drawing.
- `cv_test.go`: in-memory unit tests. `examples/` and `test/main.go`: Robotgo demos requiring a GUI; the latter saves PNGs.
- Requires Go (see `go.mod`), cgo, a C/C++ compiler, and GoCV-compatible OpenCV. CI setup lives in `.github/workflows/go.yml`: Homebrew on macOS, MinGW on Windows, an Ubuntu OpenCV container on Linux.

```sh
go build -v .
go test -v .
go vet .
go test -v . -run '^TestName$'
```

For non-legacy packages, use `go test -v . ./examples ./test ./yolo` (or `go vet` with the same paths). Avoid `./...`: `v1/` is legacy. Run `gofmt` on changed Go files.

## Contracts

- `FindAllTemplate` positional options: threshold (`float32`/`float64`, `0.8`), max count (non-negative `int`, `10`), background removal, RGB, SQDIFF (all `bool`, `false`). Invalid options, empty Mats, or an oversized search yield no results.
- Matches are masked between selections. `FindAll` uses SIFT only when template matching finds nothing and `Sift` is true.
- `Result` stores center/top-left, corners, scores, and image size; SIFT uses `MaxVal` for aspect ratio and good-match count.
- Non-`C` matching APIs preserve caller-owned Mats. `FindImgMatC`, `FindAllTemplateC`, and `FindAllSiftC` close both inputs; `FindAllTemplateCS` closes only search. `FindMultiAllTemplateC` closes each search and then the shared source. Never reuse transferred Mats.
- Close locally created GoCV resources, usually with `defer`. Preserve public signatures and legacy zero/nil failure returns; handle errors in new code.

## Changes and tests

- Keep fixes focused; follow existing conventions. Use `any`, standard-library imports before third-party imports, and declaration-name comments on exported APIs.
- Add regression tests beside source using `testing` and deterministic in-memory images/Mats; cover failures and ownership as well as success.
- Use `t.Cleanup` for test Mats, check `Closed()` before `Close()`, and report close errors.
- Run focused tests, then root build/test/vet. Keep GUI demos out of automated execution. For CI changes, validate workflow syntax and native dependency setup on both cache hits and misses.
- Treat `go.mod` and the workflow as version sources; do not duplicate dependency lists here.
