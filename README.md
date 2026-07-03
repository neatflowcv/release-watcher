# release-watcher

프로젝트 릴리스 정보를 추적하기 위한 Go 기반 REST API 서비스입니다.

## 요구 사항

- Go 1.26.4 이상
- PostgreSQL

## 실행 방법

먼저 PostgreSQL 데이터베이스를 준비합니다. 애플리케이션은 시작할 때 `projects` 테이블을 자동으로 생성합니다.

```sh
createdb release_watcher
```

데이터베이스 연결 문자열을 환경변수로 설정합니다.

```sh
export RELEASE_WATCHER_DATABASE_URL='postgres://localhost:5432/release_watcher?sslmode=disable'
```

인증이 필요한 데이터베이스는 다음 형식을 사용합니다.

```sh
export RELEASE_WATCHER_DATABASE_URL='postgres://USER:PASSWORD@HOST:PORT/DATABASE?sslmode=disable'
```

서비스를 실행합니다.

```sh
go run ./cmd/release-watcher
```

기본 HTTP 주소는 `:8080`입니다. 다른 주소로 실행하려면 `RELEASE_WATCHER_ADDRESS`를 설정합니다.

```sh
export RELEASE_WATCHER_ADDRESS=':9090'
go run ./cmd/release-watcher
```

## API 예시

프로젝트를 등록합니다.

```sh
curl -X POST http://localhost:8080/projects \
  -H 'Content-Type: application/json' \
  -d '{"name":"ceph","url":"https://github.com/ceph/ceph"}'
```

등록된 프로젝트 목록을 조회합니다.

```sh
curl http://localhost:8080/projects
```

## 테스트

```sh
go test ./...
```
