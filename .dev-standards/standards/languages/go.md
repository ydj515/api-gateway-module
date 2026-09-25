# Go Guidelines & Standards

Go 코드는 다른 언어의 클래스 구조를 옮긴 코드가 아니라 작은 package, 명시적인 값과
에러, 단순한 제어 흐름으로 읽혀야 합니다. 추상화보다 호출자가 이해할 수 있는 구체적인
동작을 먼저 만들고, 실제 교체 지점이 확인될 때만 interface를 도입합니다.

## 1. Go다운 기본 원칙

- `go.mod`의 `go`와 `toolchain` 지시문, `go.work` 및 실제 선택된 toolchain을 확인합니다.
  더 최신 toolchain이 설치돼 있어도 module의 최소 지원 버전을 넘는 문법·표준 API를
  임의로 사용하지 않습니다. 아래 `any` 예시는 Go 1.18 이상을 전제로 합니다.
- 해당 Go 버전의 release note를 기준으로 언어 의미와 표준 library 변경을 확인합니다.
  loop variable이나 timer 동작처럼 버전에 영향을 받는 코드는 관련 테스트를 실행하고,
  `gofmt`, `go vet`와 linter가 사용하는 toolchain의 호환성도 확인합니다.

- 복잡한 계층보다 조합과 작은 함수를 사용합니다. embedding은 상속이 아닙니다.
- 성공 경로는 왼쪽에 유지하고 오류와 예외 조건은 guard clause로 일찍 반환합니다.
- 가능한 타입은 유용한 zero value를 갖게 설계합니다. 불변식 검증이 필요한 경우에만
  `NewXxx` 생성자를 둡니다.
- package 이름은 짧고 구체적으로 짓고 `util`, `common`, `manager` 같은 의미 없는
  이름과 `order.OrderService` 같은 stutter를 피합니다.
- 의존성은 package 전역 변수나 service locator가 아니라 함수 또는 struct 필드로
  명시적으로 전달합니다.

## 2. 타입과 API 설계

- struct는 데이터와 동작을 함께 표현하되 getter/setter를 기계적으로 만들지 않습니다.
  Go getter는 보통 `GetName()`이 아니라 `Name()`입니다.
- interface는 구현 package가 아니라 사용하는 package에서 필요한 최소 메서드로
  정의합니다. 함수는 interface를 받을 수 있지만 특별한 이유 없이 interface를 반환하지
  않습니다.
- 값 의미가 있는 작고 복사 가능한 타입은 value receiver를, 공유 상태를 변경하거나
  복사하면 안 되는 타입은 pointer receiver를 사용합니다. 한 타입에서 receiver 방식을
  임의로 섞지 않습니다.
- slice와 map을 전달할 때 호출자와 피호출자 중 누가 수정 권한을 갖는지 정합니다.
  공개 API에서 `nil`과 빈 slice가 직렬화 결과를 바꾸면 계약으로 문서화합니다.
- 함수 옵션 패턴은 선택 항목이 많고 앞으로 확장될 공개 API에만 사용합니다. 단순한
  필수 인자를 숨기기 위한 옵션은 만들지 않습니다.

```go
package order

import "context"

type Repository interface {
	Find(ctx context.Context, id ID) (Order, error)
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}
```

## 3. 제어 흐름과 컬렉션

- `if err != nil`을 정상적인 제어 흐름으로 취급하고 중첩을 줄입니다.
- 짧은 선언은 값의 사용 범위를 좁힐 때 사용하며, 바깥 변수를 실수로 shadowing하지
  않는지 확인합니다.
- 단순 반복과 누적에는 명시적인 `for`를 사용합니다. generic helper가 실제 중복과
  오류를 줄이지 않는다면 map/filter 스타일 추상화를 만들지 않습니다.
- `defer`는 자원 획득 직후 배치합니다. 반복문에서 다수 자원을 열 때는 반복 본문을
  함수로 분리하여 defer가 함수 종료 시 실행되게 합니다.

## 4. comma-ok와 error 선택

- map 조회, type assertion, channel receive처럼 언어가 제공하는 두 번째 `ok` 값은
  “값이 없음”, “해당 타입이 아님”, “channel이 닫힘”이 정상적으로 확인 가능한 상태이고
  별도 원인 정보가 필요 없을 때 사용합니다.
- parsing, I/O, validation, transaction처럼 실패 이유와 작업 문맥, 재시도 가능 여부를
  호출자가 알아야 하면 `error`를 반환합니다. 단순히 성공 여부만 남기는 `bool`로 실패
  원인을 버리지 않습니다.
- map의 key 부재가 domain 실패라면 map 경계에서는 `ok`로 확인하고 application/domain
  API 경계에서 의미 있는 sentinel/typed error로 변환합니다.
- `(T, bool, error)`는 “정상적인 부재”, “값 존재”, “조회 자체 실패” 세 상태가 모두 실제
  계약일 때만 사용합니다. 같은 상태를 `bool`과 `error`로 중복 표현하지 않습니다.
- 짧은 `if value, ok := ...; ok`는 값의 scope를 좁힐 때 사용합니다. 이후에도 값이나 상태가
  필요하면 읽기 어려운 `else`를 늘리지 말고 먼저 선언합니다.

```go
func regionOf(attributes map[string]any) (string, error) {
	value, ok := attributes["region"]
	if !ok {
		return "", ErrRegionMissing
	}

	region, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("region has type %T: %w", value, ErrInvalidRegion)
	}
	return region, nil
}
```

첫 번째 `ok`는 map 부재를 구분하고, 두 번째 실패는 외부 계약의 잘못된 값이므로 원인을
보존하는 domain error로 변환합니다. `error`가 `ok`보다 우월한 것이 아니라 필요한 정보량이
다릅니다.

## 5. 에러는 값으로 다룬다

- 오류에는 작업 문맥을 붙이고 원인은 `%w`로 보존합니다.
- 호출자가 분기해야 하는 실패는 sentinel error 또는 구체적인 error type으로
  표현하고 `errors.Is`와 `errors.As`로 판별합니다. 문자열을 비교하지 않습니다.
- 같은 오류를 모든 계층에서 반복 로깅하지 않습니다. 처리할 수 없으면 문맥을 붙여
  반환하고 프로세스나 요청 경계에서 한 번 기록합니다.
- `panic`은 일반 입력 오류나 외부 시스템 실패에 사용하지 않습니다. 복구할 수 없는
  초기화 실패나 깨진 내부 불변식으로 한정합니다.
- `Close`, `Flush`, transaction commit처럼 결과가 중요한 정리 작업의 오류를 버리지
  않습니다.

```go
order, err := s.repository.Find(ctx, id)
if err != nil {
	return Order{}, fmt.Errorf("find order %s: %w", id, err)
}
return order, nil
```

## 6. domain-oriented package 적용

Go 언어 선택만으로 architecture를 자동 적용하지 않습니다. 업무 기능 중심 구성이 필요한
서비스는 `domain-oriented`를 명시적으로 선택합니다.

```yaml
languages: [go]
architectures: [domain-oriented]
```

이 profile을 Go로 구현할 때는 bounded context 아래에 Java식 계층 package를 모두
만들기보다 작은 core package와 inbound/outbound adapter package로 평탄화합니다.

```text
cmd/api/
└─ main.go                         # composition root
internal/
├─ order/
│  ├─ order.go                     # entity, value와 domain rule
│  ├─ create.go                    # application operation, params/result
│  ├─ errors.go                    # caller가 분기하는 domain/application error
│  ├─ repository.go                # order package가 소비하는 작은 port
│  ├─ http/
│  │  ├─ handler.go                # inbound adapter와 error mapping
│  │  └─ contract.go               # request/response wire type
│  └─ postgres/
│     └─ repository.go             # outbound persistence adapter
└─ payment/
```

- 의존 방향은 `main → adapter → order`입니다. `order` package는 HTTP, database와
  framework package를 import하지 않습니다.
- interface는 구현 package가 아니라 사용하는 package에 필요한 최소 method로 둡니다.
  구현이 하나이고 대체 경계가 없다면 concrete type을 사용합니다.
- HTTP request/response는 `<feature>/http/contract.go`, use case params/result는
  `<feature>/<operation>.go`, caller가 분기할 error는 `<feature>/errors.go`에 둡니다.
- driver/client 오류는 adapter 안에서 `%w`로 원인을 보존하고, 업무 의미가 안정된 실패만
  domain/application error로 변환합니다.
- 모든 경계에 `dto`, `mapper`, `exception` package와 `Dto` suffix type을 만들지 않습니다.
  type은 실제 wire 호환성이나 operation 의미가 달라질 때 분리합니다.

library처럼 애플리케이션 package 구조가 필요하지 않으면 architecture를 생략합니다.
선택한 framework의 자동 기본값까지 끄려는 경우에만 `architectures: []`를 사용합니다.

## 7. Context와 동시성

- 요청 범위 I/O 함수는 `context.Context`를 첫 번째 인자로 받고 struct에 저장하지
  않습니다. `nil` context를 전달하지 않습니다.
- 동기 코드로 충분하면 goroutine을 만들지 않습니다. 모든 goroutine에는 소유자,
  종료 조건, cancellation 경로가 있어야 합니다.
- channel은 데이터 전달과 소유권 이전이 목적일 때 사용합니다. 단순 공유 상태 보호는
  `sync.Mutex`가 더 명확할 수 있습니다.
- channel을 닫는 책임은 일반적으로 값을 보내는 쪽에 둡니다. 수신자가 더 이상 읽지
  않을 수 있는 경로에서 sender가 영구 block되지 않게 cancellation을 전달합니다.
- 반복 변수 capture, timer/ticker 정리, goroutine 내부 오류 전파를 테스트합니다.

## 8. 주석과 Go Doc

- exported package와 식별자는 선언 이름으로 시작하는 완전한 문장으로 계약을 설명합니다.
- 구현을 줄 단위로 번역하지 말고 호출자가 알아야 하는 특수 조건, 동시성 안전성,
  오류와 복잡도를 기록합니다.
- deprecated API는 `Deprecated:` 문단으로 대체 API와 전환 방법을 명시합니다.

```go
// Load returns the order visible to the current tenant.
// It returns ErrNotFound when id is unknown or belongs to another tenant.
func (s *Service) Load(ctx context.Context, id ID) (Order, error) {
	// 구현 내용
}
```

## 9. 테스트와 품질 게이트

- 여러 입력과 경계값은 table-driven test와 명확한 subtest 이름으로 검증합니다.
- mock 호출 횟수보다 공개 동작, 반환 오류와 상태 변화를 우선 검증합니다.
- 동시성 코드 변경에는 종료, cancellation과 race 조건을 포함합니다.

```sh
test -z "$(gofmt -l .)"
go vet ./...
go test ./...
go test -race ./...
```

`-race` 비용이 큰 저장소는 기본 테스트와 별도 task로 두되 CI 적용 범위를 문서화합니다.
도구 시작점은 `templates/mise/go/`와 `templates/golangci-lint/`에 있습니다.

## 10. 다른 언어 습관을 옮기지 않는다

- class hierarchy, annotation 기반 DI, exception 중심 흐름을 Go에 재현하지 않습니다.
- 모든 구현체 앞에 interface를 만들거나 `IRepository`, `OrderRepositoryImpl`처럼 Java식
  이름을 사용하지 않습니다.
- context 없는 background goroutine, fluent builder 남용, 의미 없는 wrapper type을
  “아키텍처 일관성”만을 위해 추가하지 않습니다.

참고: [Effective Go](https://go.dev/doc/effective_go),
[Go Code Review Comments](https://go.dev/wiki/CodeReviewComments),
[Go Doc Comments](https://go.dev/doc/comment)

버전별 설정 참고: [Go Toolchains](https://go.dev/doc/toolchain),
[Go 1.18 Release Notes](https://go.dev/doc/go1.18)
