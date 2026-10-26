# 기여 가이드

API Gateway Module에 기여해 주셔서 감사합니다. 변경 사항을 제출하기 전에 아래 개발 환경 설정과 검증 절차를 확인해 주세요. 상세한 저장소 규칙과 `mise` 환경 프로필 정책은 [기여 및 개발 가이드](./docs/contributor-guide.md)에 정리되어 있습니다.

## 1. 개발 환경 준비

Go 버전(`1.27.0`)과 공통 작업은 [`mise.toml`](./mise.toml)에 정의되어 있습니다.

```bash
mise trust ./mise.toml
mise run bootstrap
mise run config:check
```

- `dev`/`prod` 환경 오버레이(`-E dev`, `-E prod`)와 개인 로컬 설정(`mise.dev.local.toml`) 정책은 [기여 및 개발 가이드의 mise 환경 구성](./docs/contributor-guide.md#mise-환경-구성-및-프로필-정책)을 참고하세요.
- `mise.local.toml`은 모든 환경에서 로드되므로 개인 개발용 오버라이드는 반드시 `mise.dev.local.toml`에 작성합니다.

## 2. 코드 및 설정 작성 규칙

- **코딩 스타일:** 표준 Go 포맷(`gofmt`)과 패키지 경계(`app`, `common`, `config`, `kafka`, `types`)를 준수합니다.
- **설정 스키마 동기화:** [`config/`](./config)의 YAML 구조체를 변경한 경우 [`deploy-sample.yaml`](./deploy-sample.yaml)과 관련 문서를 함께 갱신해야 합니다.
- **커밋 메시지:** `feat:`, `fix:`, `refactor:`, `docs:`, `test:`, `chore:` 등 Conventional Commits 형식을 따릅니다.

## 3. 필수 검증

모든 변경 사항은 Pull Request를 제출하기 전에 아래 검증 명령을 통과해야 합니다.

```bash
mise run verify
```

라우터 핸들러(`app/router`), HTTP 클라이언트(`app/client`), 설정 파서(`config`)의 동작을 변경한 경우 `*_test.go` 테스트를 함께 추가하거나 업데이트해 주세요.

## 4. 참고 문서

- [기여 및 개발 가이드 (mise 환경 구성 정책)](./docs/contributor-guide.md)
- [에이전트 및 저장소 가이드](./AGENTS.md)
