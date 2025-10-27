# API Gateway Module

단일 YAML 정의만으로 HTTP 라우터와 Kafka 프로듀서를 동시에 구성하는 Go/Fiber 기반 API 게이트웨이입니다. 각 `app` 엔트리는 Uber Fx를 통해 라우터, Resty 클라이언트, Kafka 프로듀서를 묶어 다중 게이트웨이를 한 프로세스에서 구동할 수 있게 합니다.

## 사전 요구사항
- Go 1.24 (go.mod에 명시된 최소 버전) 이상이 설치되어 있어야 합니다.

## 주요 기능
- GET/POST/PUT/DELETE 라우트를 YAML 설정에 따라 동적으로 생성하고, GET은 `query`/`url` 두 타입을 지원.
- Resty 기반 HTTP 클라이언트에 Base URL, 헤더, 인증 토큰 등을 YAML로 주입.
- Confluent Kafka 프로듀서를 설정에 맞춰 자동 생성해 다운스트림 호출 로그나 팬아웃 처리 가능.
- Uber Fx 라이프사이클 관리로 모든 라우터와 의존성을 안전하게 부트스트랩.
- Bytedance Sonic 기반 JSON 핸들러로 빠른 직렬화와 공통 에러 로깅 제공.
- Node 기반 예제 서비스와 `scripts/curl-gateway.sh` 스크립트를 통해 전체 라우팅 흐름을 손쉽게 검증.

## 프로젝트 구조
```
.
├── app/              # 앱 부트스트랩, DI 모듈, 라우터 등록
│   ├── app.go
│   ├── client/       # Resty 래퍼 + Kafka 프로듀서 연결
│   ├── dependency/   # Fx 모듈 (config/client/producer/router)
│   └── router/       # HTTP 메서드별 핸들러
├── common/           # JSON 핸들러 등 공용 유틸
├── config/           # YAML 스키마 정의 (App/HTTP/Kafka)
├── kafka/            # confluent-kafka-go 기반 프로듀서
├── types/http/       # HTTP 메서드 및 GET 타입 enum
├── example-service/  # 로컬 검증용 Node 기반 다운스트림 샘플 서비스
├── scripts/          # 게이트웨이 호출을 돕는 보조 스크립트
├── deploy.yaml       # 실제 실행 시 사용하는 설정 파일
├── deploy-sample.yaml# 샘플 설정 (필수 참고)
├── go.mod / go.sum
└── main.go           # Fx 엔트리포인트
```

## 구성 방법
런타임 설정은 기본적으로 `./deploy.yaml`에서 읽습니다. 다른 파일을 쓰려면 `-yamlPath=/path/to/file` 플래그를 전달하세요. **반드시 `deploy-sample.yaml`을 참고해 동일한 구조로 `deploy.yaml`을 만들어야** 파서가 정상 동작합니다.

```yaml
apps:
  - app:
      name: gateway_module
      port: "8080"
      version: "1.0.0"
    http:
      base_url: "http://localhost:3000"
      router:
        - method: "GET"
          get_type: "query"   # 또는 "url"
          path: "/api/query"
          variable: ["name", "age"]
          header:
            X-Request-Source: gateway
        - method: "GET"
          get_type: "url"
          path: "/api/user/:id"
        - method: "POST"
          path: "/api/create"
          header:
            Content-Type: application/json
          auth:
            key: "Bearer"
            token: "<token>"
        - method: "PUT"
          path: "/api/update/:id"
          header:
            Content-Type: application/json
        - method: "DELETE"
          path: "/api/delete/:id"
          header:
            Content-Type: application/json
    kafka:
      url: "localhost:9092"
      client_id: "gateway-producer"
      topic: "gateway-log"
      acks: "all"
      batch_time: 2
```

핵심 필드
- `http.router[].get_type`: `query`는 요청의 쿼리 파라미터를 그대로 전달, `url`은 들어온 Path를 그대로 프록시.
- `http.router[].auth`, `header`: Resty 요청에 공통 헤더/인증 정보를 삽입.
- `kafka`: 각 앱 전용 Kafka 프로듀서 설정. 필요 시 `app/client`에서 `producer.SendEvent` 호출 로직을 확장하세요.

## 실행 방법
1. Go 1.24 이상을 설치하고 최초 한 번 `go mod download`.
2. `deploy-sample.yaml`을 복사해 `deploy.yaml`을 만들고 환경에 맞게 수정.
3. (선택) 로컬 검증을 위해 `example-service` 디렉터리에서 예제 서비스를 실행:
   ```bash
   cd example-backend-service
   npm install
   npm start
   ```
   기본 포트는 `3000`이며, 필요 시 `PORT` 환경 변수를 지정하세요.
4. 게이트웨이 실행:
   ```bash
   go run . -yamlPath=./deploy.yaml
   ```
5. 설정된 각 `app.port`마다 Fiber 인스턴스가 하나씩 구동되고, YAML에 정의한 경로가 그대로 마운트됩니다.

바이너리 빌드:
```bash
go build -o bin/api-gateway
./bin/api-gateway -yamlPath=/path/to/deploy.yaml
```

## 로컬 검증 예시
- 샘플 서비스와 게이트웨이가 구동 중일 때 `scripts/curl-gateway.sh`를 실행하면 GET/POST/PUT/DELETE 요청 흐름을 한 번에 확인할 수 있습니다.
  ```bash
  ./scripts/curl-gateway.sh             # 기본 게이트웨이 주소 http://localhost:8080
  ./scripts/curl-gateway.sh http://localhost:9090  # 포트가 다를 때
  ```
- 스크립트는 각 요청에 대한 응답 바디와 HTTP 코드를 함께 출력하므로 라우팅/헤더 전달 여부를 빠르게 확인할 수 있습니다.

## 확장 가이드
- **새 다운스트림 서비스**: `apps` 배열에 항목을 추가하면 Fx가 자동으로 새 클라이언트/라우터/프로듀서를 생성합니다.
- **커스텀 미들웨어**: `app/router/router.go`에서 Fiber 미들웨어(CORS, 로깅 등)를 등록하세요.
- **HTTP 메서드 추가**: `app/router`의 패턴을 따라 새 파일을 만들고 `types/http` enum을 확장합니다.
- **Kafka 이벤트 구조화**: `kafka/producer.go` 또는 `app/client`에서 직렬화 로직을 추가해 로그/메트릭을 보강합니다.

## 사용 라이브러리
- [Fiber v2](https://github.com/gofiber/fiber) – HTTP 서버
- [Resty](https://github.com/go-resty/resty) – 외부 API 호출
- [Uber Fx](https://github.com/uber-go/fx) – DI 및 라이프사이클
- [Confluent Kafka Go](https://github.com/confluentinc/confluent-kafka-go) – Kafka 프로듀서
- [Bytedance Sonic](https://github.com/bytedance/sonic) – 고성능 JSON 직렬화
