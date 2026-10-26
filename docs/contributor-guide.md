# 기여 및 개발 가이드

이 문서는 `api-gateway-module` 저장소의 로컬 개발 환경, `mise` 환경 프로필 정책, 아키텍처 확장 방법 및 검증 절차를 안내합니다.

## 요구 사항 및 초기 설정

Go 버전(`1.27.0`)과 공통 작업은 [`mise.toml`](../mise.toml)에 고정되어 있습니다.

```bash
mise trust ./mise.toml   # 최초 1회, task와 overlay를 검토한 뒤 신뢰
mise run bootstrap       # 고정된 버전의 도구(Go 1.27.0) 설치 후 의존성 준비
mise run config:check    # task 참조·순환 검사, 앱 실행 없음
mise run verify          # 프로젝트 검증 (lint, test, build)
mise -E dev run verify   # 개발 프로필 기준 검증
```

## mise 환경 구성 및 프로필 정책

`mise.toml`은 도구·공통 task, `mise.dev.toml`/`mise.prod.toml`은 공유 환경 선택을 담당합니다.

### 1. 환경 프로필 (`APP_ENV`)
- 기본 실행은 `APP_ENV=local`, `-E dev`는 개발 overlay(`gateway_config = "deploy.yaml"`), `-E prod`는 운영 설정(`gateway_config = "deploy.prod.yaml"`) 선택입니다. 환경 선택 자체가 배포나 서비스 시작을 수행하지 않습니다.
- `APP_ENV`는 공통 환경 식별자입니다. 애플리케이션의 기존 실행·배포 설정은 유지합니다.
- `dev` 태스크는 `guard:local`을 통해 `local` 또는 `dev` 환경에서만 실행되도록 보호됩니다.
- `-E local`은 사용하지 않습니다.

### 2. 개인 로컬 설정 및 `.env` 관리
- **개인 개발 설정 (`mise.dev.local.toml`)**: [`mise.dev.local.toml.example`](../mise.dev.local.toml.example)을 검토해 `mise.dev.local.toml`로 복사합니다. `.env.dev.local`을 만든 뒤 `env._.file`을 활성화하면 `dev`에서만 읽습니다. 기존 개인 파일을 덮어쓰지 않습니다.
- **전역 로컬 오버라이드 (`mise.local.toml`)**: `mise.local.toml`은 **모든 환경**에서 로드됩니다. `prod` checkout에 개인 override나 개발 dotenv를 두지 않습니다.

### 3. 버전 고정 및 Lockfile 정책
- `mise`는 개발 도구의 정확한 버전 고정과 프로필 분리에 사용합니다. `[settings] lockfile = false`로 도구 lock 생성을 끄고 `mise.lock`은 관리하지 않습니다. 공통 표준의 `mise.lock` 지침보다 이 저장소의 정책을 우선합니다. 설치 파일까지 고정해야 하는 요구가 생기면 다시 도입합니다.
- 도구 버전은 `mise.toml`의 `[tools]`에서 관리하며 `dev`/`prod`에서도 같은 버전을 사용합니다. 로컬과 CI는 `mise install` 또는 같은 정확한 버전의 setup action으로 도구를 준비합니다. `go.sum`, `package-lock.json` 등 애플리케이션 의존성 잠금과 검증 옵션은 유지합니다.
- 지원 OS는 각 도구와 실행 스크립트의 호환성에 따릅니다. `mise.lock`을 사용하지 않는다고 Windows 실행까지 보장되는 것은 아닙니다.
- `mise run bootstrap`은 프로젝트 초기 설정입니다. OS package·dotfile·서비스를 관리하는 `mise bootstrap`은 개인 머신 설정에서 별도로 채택합니다.

## 로컬 통합 실행 및 검증

### 1. 런타임 설정 (`deploy.yaml`)
게이트웨이는 기본적으로 `./deploy.yaml`을 읽으며, `-yamlPath=/path/to/file` 플래그로 경로를 변경할 수 있습니다. [`deploy-sample.yaml`](../deploy-sample.yaml)을 복사해 `deploy.yaml`을 작성하세요.

- `http.router[].get_type`: `query`는 지정된 `variable` 쿼리 파라미터를 전달하고, `url`은 요청 Path를 그대로 프록시합니다.
- `http.router[].auth`, `header`: Resty 요청에 공통 헤더와 인증 토큰을 주입합니다.
- `kafka`: 각 앱 인스턴스 전용 Confluent Kafka 프로듀서 설정입니다.

### 2. 로컬 인프라 및 샘플 다운스트림 기동 (선택)
Kafka 클러스터와 다운스트림 예제 서버를 함께 띄워 전체 흐름을 검증할 수 있습니다.

```bash
# 1) 로컬 Kafka 클러스터 기동 (필요 시 .env-sample을 .env로 복사)
cp .env-sample .env
docker compose up -d

# 2) Node.js 예제 다운스트림 서버 실행 (기본 포트 3000)
cd example-backend-service
npm install
npm start
```

### 3. 라우팅 및 서킷 브레이커 테스트
게이트웨이(`mise run dev`)가 실행 중일 때 보조 스크립트로 동작을 확인할 수 있습니다.

```bash
# GET/POST/PUT/DELETE 라우팅 일괄 호출 (기본 http://localhost:8080)
./scripts/curl-gateway.sh

# 다운스트림 중지 상태에서 서킷 브레이커 OPEN 전환 확인
./scripts/test-circuitbreaker.sh http://localhost:8080
```

## 확장 가이드

- **새 다운스트림 서비스 추가:** `deploy.yaml`의 `apps` 배열에 항목을 추가하면 Uber Fx가 해당 앱의 라우터, Resty 클라이언트, Kafka 프로듀서를 자동 생성합니다.
- **커스텀 미들웨어 등록:** [`app/router/router.go`](../app/router/router.go)에서 Fiber 미들웨어(CORS, Recover, 로깅 등)를 구성합니다.
- **HTTP 메서드 확장:** [`app/router/`](../app/router)의 핸들러 패턴에 따라 새 메서드를 추가하고 [`types/http/root.go`](../types/http/root.go) enum을 확장합니다.
- **서킷 브레이커 및 직렬화 조정:** [`common/circuit.go`](../common/circuit.go)에서 `gobreaker` 임계값을 조정하고, [`common/json.go`](../common/json.go)에서 Bytedance Sonic 기반 직렬화 동작을 관리합니다.
- **Kafka 이벤트 발행 확장:** [`kafka/producer.go`](../kafka/producer.go) 또는 [`app/client/`](../app/client)에서 `SendEvent` 호출 및 메시지 스키마를 확장합니다.

## 테스트 및 품질 검증

변경 사항을 제출하기 전에 아래 검증 명령을 실행해 모든 검사가 통과하는지 확인하세요.

| 명령 | 설명 |
| --- | --- |
| `mise run lint` | `go vet ./...` 정적 분석 실행 |
| `mise run test` | `go test ./...` 단위 테스트 실행 |
| `mise run build` | `go build ./...` 전체 패키지 컴파일 확인 |
| `mise run verify` | `lint` + `test` + `build` 통합 품질 게이트 실행 |
