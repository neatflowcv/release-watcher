# ADR-0001: PostgreSQL 클라이언트로 `pgx`를 사용한다

## 상태

Accepted

## 배경

release-watcher는 프로젝트 정보를 저장하고 조회하기 위해 PostgreSQL과 통신해야 한다.

PostgreSQL 클라이언트 라이브러리는 다음 후보를 확인했다.

| 후보 | 확인한 최신 버전 | 출시 날짜 | 판단 |
| --- | --- | --- | --- |
| `pgx` | `v5.10.0` | 2026-06-03 | 선택 |
| `lib/pq` | `v1.12.3` | 2026-04-03 | 보류 |
| `bun pgdriver` | `v1.2.18` | 2026-02-28 | 보류 |
| `go-pg/pg` | `v10.15.1` | 2026-05-29 | 거부 |

후보 모듈 경로는 각각 `github.com/jackc/pgx/v5`,
`github.com/lib/pq`, `github.com/uptrace/bun/driver/pgdriver`,
`github.com/go-pg/pg/v10`이다.

## 결정

PostgreSQL 클라이언트 라이브러리로 `pgx`를 사용한다.

## 근거

- 확인한 후보 중 `pgx`의 출시 날짜가 오늘(2026-06-29)에 가장 가깝다.
- 이 서비스는 PostgreSQL 서버하고만 통신하면 된다.
- 범용 데이터베이스 추상화보다 PostgreSQL에 적당하게 맞춰 쓸 수 있는 클라이언트가 충분하다.
- `go-pg/pg`는 사용하지 않는다. 프로젝트가 스스로 maintenance mode라고 선언한 라이브러리는 굳이 채택하지 않는다.
- `lib/pq`와 `bun pgdriver`는 강하게 끌리는 선택지가 아니다. 현재 요구에는 `pgx`면 충분하다.

## 결과

- PostgreSQL 저장소 구현은 `pgx` 기반으로 작성한다.
- ORM 도입은 기본 방향으로 보지 않는다.
- `go-pg/pg`는 이후 후보에서도 제외한다.
