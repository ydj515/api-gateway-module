# API Gateway Module

[![Sync Development Standards](https://github.com/ydj515/api-gateway-module/actions/workflows/sync-dev-standards.yml/badge.svg)](https://github.com/ydj515/api-gateway-module/actions/workflows/sync-dev-standards.yml)
[![Go 1.27.0](https://img.shields.io/badge/Go-1.27.0-00ADD8?logo=go&logoColor=white)](./mise.toml)
[![Fiber v3.5](https://img.shields.io/badge/Fiber-v3.5-00ACD7?logo=go&logoColor=white)](./go.mod)
[![Uber Fx v1.24](https://img.shields.io/badge/Uber%20Fx-v1.24-000000?logo=uber&logoColor=white)](./go.mod)
[![Apache Kafka](https://img.shields.io/badge/Apache%20Kafka-Confluent%20v2-231F20?logo=apachekafka&logoColor=white)](./go.mod)

단일 YAML 설정만으로 다중 HTTP 프록시 라우터와 Kafka 프로듀서를 구성해 구동하는 Go/Fiber 기반 API 게이트웨이입니다.
YAML에 정의된 각 `app` 엔트리마다 Uber Fx로 Fiber 라우터, Resty 클라이언트, Kafka 프로듀서 생명주기를 묶어 한 프로세스에서 여러 게이트웨이를 실행할 수 있습니다.
로컬 검증용 Node.js 다운스트림 서버와 호출 스크립트를 포함해 라우팅과 서킷 브레이커 동작을 바로 재현할 수 있습니다.

**기술 스택:** Go 1.27 · Fiber v3 · Uber Fx · Resty v2 · Confluent Kafka Go v2 · gobreaker v2 · Sonic

## 이 프로젝트에서 볼 수 있는 것

- **선언적 다중 라우팅:** YAML 설정에 따라 `GET`(`query`/`url` 방식), `POST`, `PUT`, `DELETE` 엔드포인트와 공통 헤더·Bearer 토큰 주입을 동적으로 생성
- **회로 차단 및 고성능 직렬화:** `gobreaker` 기반 서킷 브레이커로 다운스트림 장애 전파를 차단하고, Bytedance Sonic으로 JSON 직렬화 처리
- **DI 생명주기 관리와 이벤트 발행:** Uber Fx 모듈(`config`, `client`, `producer`, `router`)로 멀티 포트 인스턴스와 Confluent Kafka 프로듀서를 안전하게 부트스트랩

## 요구 사항

- [mise](https://mise.jdx.dev/) — Go 버전(`1.27.0`)은 [`mise.toml`](./mise.toml)에 고정되어 있으며 `mise`가 설치합니다.
- Docker Compose & Node.js *(선택)* — 로컬 Kafka 클러스터([`docker-compose.yml`](./docker-compose.yml)) 및 다운스트림 예제 서비스([`example-backend-service`](./example-backend-service)) 실행 시 필요합니다.

## 빠른 시작

```bash
mise trust ./mise.toml              # 최초 1회, mise.toml 검토 후 신뢰
mise run bootstrap                  # 고정된 Go 1.27.0 설치 + 모듈 다운로드
cp deploy-sample.yaml deploy.yaml   # 샘플 설정을 복사해 런타임 설정 준비
mise run dev                        # 게이트웨이 실행 (기본 :8080)
```

로컬에서 다운스트림 응답까지 함께 확인하려면 별도 터미널에서 예제 서버를 실행한 뒤 검증 스크립트를 호출하세요.

```bash
# 터미널 1: 로컬 다운스트림 예제 서버 실행 (기본 :3000)
cd example-backend-service && npm install && npm start

# 터미널 2: GET/POST/PUT/DELETE 프록시 흐름 및 서킷 브레이커 검증
./scripts/curl-gateway.sh
./scripts/test-circuitbreaker.sh
```

다른 경로의 설정 파일을 사용하려면 `go run . -yamlPath=/path/to/deploy.yaml`로 실행합니다.
Kafka 클러스터 연동과 상세 라우터 확장 방법은 [기여 및 개발 가이드](./docs/contributor-guide.md)에 설명합니다.

## 주요 명령

| 명령 | 설명 |
| --- | --- |
| `mise run bootstrap` | 고정된 Go(`1.27.0`) 설치 및 Go 모듈 의존성 다운로드 |
| `mise run dev` | 선택된 설정(`deploy.yaml`)으로 게이트웨이 실행 |
| `mise run test` | Go 단위 테스트 실행 (`go test ./...`) |
| `mise run lint` | Go 정적 분석 실행 (`go vet ./...`) |
| `mise run build` | 전체 패키지 컴파일 검증 (`go build ./...`) |
| `mise run config:check` | `mise` 태스크 참조 및 순환 의존성 검사 |
| `mise run verify` | `lint` + `test` + `build` 통합 검증 |

전체 task 목록은 `mise tasks`로 확인할 수 있습니다. `dev`/`prod` 환경 프로필(`-E dev`, `-E prod`)과 개인 설정(`mise.dev.local.toml`) 등 상세 정책은 [기여 및 개발 가이드의 mise 환경 구성](./docs/contributor-guide.md#mise-환경-구성-및-프로필-정책)을 참고하세요.

## 설정 예시 (`deploy.yaml`)

[`deploy-sample.yaml`](./deploy-sample.yaml)과 동일한 스키마로 `deploy.yaml`을 작성합니다.

```yaml
apps:
  - app:
      name: "gateway_module"
      port: "8080"
      version: "1.0.0"
    http:
      base_url: "http://localhost:3000"
      router:
        - method: "GET"
          get_type: "query"   # "query" 또는 "url"
          path: "/api/query"
          variable: ["name", "age"]
          header:
            X-Request-Source: gateway
        - method: "POST"
          path: "/api/users"
          auth:
            key: "Bearer"
            token: "<token>"
    kafka:
      url: "localhost:9092"
      client_id: "gateway-producer"
      topic: "gateway-log"
      acks: "all"
      batch_time: 2
```

## 프로젝트 구조

```text
.
├── app/                        # 앱 부트스트랩, Fx DI 모듈, HTTP 클라이언트 및 라우터
│   ├── client/                 # Resty 래퍼 및 다운스트림 호출
│   ├── dependency/             # Uber Fx 모듈 (config, client, producer, router)
│   └── router/                 # HTTP 메서드별 Fiber 핸들러 및 단위 테스트
├── common/                     # Sonic JSON 핸들러 및 gobreaker 서킷 브레이커
├── config/                     # YAML 설정 파서 및 스키마 (App, HTTP, Kafka)
├── docs/                       # 기여 및 개발 가이드 (mise 환경 프로필 정책)
├── example-backend-service/    # 로컬 검증용 Express 다운스트림 샘플 서비스
├── kafka/                      # Confluent Kafka 프로듀서 래퍼
├── scripts/                    # 라우팅(curl-gateway.sh) 및 서킷 브레이커 검증 스크립트
├── types/http/                 # HTTP 메서드 및 GET 라우팅 타입 enum
├── CONTRIBUTING.md             # 기여 절차 요약
├── deploy-sample.yaml          # 게이트웨이 설정 샘플
├── docker-compose.yml          # 로컬 3노드 Kafka 클러스터 및 Kafka UI 구성
├── main.go                     # Uber Fx 애플리케이션 엔트리포인트
└── mise.toml                   # 고정 도구 버전(Go 1.27.0) 및 공통 태스크 정의
```

## 문서

- [기여 및 개발 가이드 (mise 환경 구성 정책)](./docs/contributor-guide.md)
- [저장소 작업 규칙](./AGENTS.md)

## 기여

변경 전 [기여 가이드](./CONTRIBUTING.md)를 확인해 주세요.
