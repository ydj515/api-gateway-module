# golangci-lint Guidelines

이 문서는 Go 프로젝트에서 golangci-lint를 재현 가능한 품질 게이트로 사용하는 기준입니다.

## 설정과 버전

- golangci-lint 실행 버전과 Go toolchain을 `mise.toml` 같은 프로젝트 설정에 고정합니다.
- version 2 schema의 `.golangci.yml`을 저장소 루트에서 관리하고 `golangci-lint config verify`로
  실행 버전과 설정의 호환성을 확인합니다.
- 기본 `standard` 집합에서 시작하고 추가 linter는 검출하려는 실패 유형과 도입 비용을
  설명한 뒤 최소한으로 활성화합니다. `all`은 릴리스마다 규칙 집합이 변하므로 기본값으로
  사용하지 않습니다.
- 시작점은 `templates/golangci-lint/.golangci.yml.example`이며 소비할 때 루트의
  `.golangci.yml`로 복사하고 모든 `REPLACE_ME` 버전을 확정합니다.

## 실행과 예외

- CI lint는 `golangci-lint run`으로 formatter 차이와 linter 위반을 함께 검출하며 `--fix`를
  사용해 working tree를 암묵적으로 변경하지 않습니다.
- formatter 적용은 개발자가 `golangci-lint fmt`로 명시적으로 수행하고 변경된 diff를
  검토합니다.
- generated source는 표준 `// Code generated ... DO NOT EDIT.` 표시가 있는 경우에만
  제외하도록 strict mode를 사용합니다.
- `//nolint:<linter> // reason`은 구체적인 linter 이름과 필요한 이유를 함께 기록합니다.
  file 또는 package 전체 예외는 최소 범위 예외로 표현할 수 없을 때만 허용합니다.
- `modules-download-mode: readonly`로 lint 실행이 `go.mod`와 `go.sum`을 조용히 변경하지
  않게 합니다.

## 최소 검증

```sh
golangci-lint config verify
golangci-lint run
go test ./...
```

Go와 golangci-lint를 고정하는 mise 시작점은 `templates/mise/go/mise.toml.example`을
사용합니다. 설정 형식은 [golangci-lint configuration 문서](https://golangci-lint.run/docs/configuration/file/)를 기준으로 합니다.
