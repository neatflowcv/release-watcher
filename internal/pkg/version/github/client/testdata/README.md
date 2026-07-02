# GitHub 클라이언트 테스트 데이터

이 디렉터리는 `https://github.com/ceph/ceph`에서 수집한 GitHub API 응답
샘플을 보관한다.

## 파일

- `ceph-tags.json`: tags API에서 가져온 최근 repository tag 목록.
- `ceph-tag-ref-v21.3.0.json`: `v21.3.0` tag의 git ref 응답.
- `ceph-tag-object-v21.3.0.json`: `v21.3.0`의 annotated tag object 응답.
- `ceph-tag-ref-v20.2.2.json`: `v20.2.2` tag의 git ref 응답.
- `ceph-tag-object-v20.2.2.json`: `v20.2.2`의 annotated tag object 응답.

## 생성 명령

tag 목록은 다음 명령으로 만들었다.

```bash
curl -L --fail --retry 2 \
  -H 'Accept: application/vnd.github+json' \
  -H 'X-GitHub-Api-Version: 2022-11-28' \
  -o ceph-tags.json \
  'https://api.github.com/repos/ceph/ceph/tags?per_page=30'
```

tag ref 응답은 다음 명령으로 만들었다.

```bash
curl -L --fail --retry 2 \
  -H 'Accept: application/vnd.github+json' \
  -H 'X-GitHub-Api-Version: 2022-11-28' \
  -o ceph-tag-ref-v21.3.0.json \
  'https://api.github.com/repos/ceph/ceph/git/ref/tags/v21.3.0'
```

```bash
curl -L --fail --retry 2 \
  -H 'Accept: application/vnd.github+json' \
  -H 'X-GitHub-Api-Version: 2022-11-28' \
  -o ceph-tag-ref-v20.2.2.json \
  'https://api.github.com/repos/ceph/ceph/git/ref/tags/v20.2.2'
```

annotated tag object SHA는 tag ref 응답에서 다음 명령으로 추출했다.

```bash
jq -r '.object.sha' ceph-tag-ref-v21.3.0.json
jq -r '.object.sha' ceph-tag-ref-v20.2.2.json
```

수집 시점의 SHA는 다음과 같았다.

```text
v21.3.0: b44498fd6d18e36065848c80e00edb64acae1adf
v20.2.2: cdc68ca21135a7b5fa913f216c3903628f3b61bb
```

annotated tag object 응답은 다음 명령으로 만들었다.

```bash
curl -L --fail --retry 2 \
  -H 'Accept: application/vnd.github+json' \
  -H 'X-GitHub-Api-Version: 2022-11-28' \
  -o ceph-tag-object-v21.3.0.json \
  'https://api.github.com/repos/ceph/ceph/git/tags/b44498fd6d18e36065848c80e00edb64acae1adf'
```

```bash
curl -L --fail --retry 2 \
  -H 'Accept: application/vnd.github+json' \
  -H 'X-GitHub-Api-Version: 2022-11-28' \
  -o ceph-tag-object-v20.2.2.json \
  'https://api.github.com/repos/ceph/ceph/git/tags/cdc68ca21135a7b5fa913f216c3903628f3b61bb'
```

## Verification 참고

repository tags API 응답에는 GitHub의 signature verification 상태가 포함되지
않는다. annotated tag의 verification 상태를 확인하려면 git ref 응답에서 tag
object SHA를 확인한 뒤, 해당 tag object 응답을 다시 조회해야 한다.

수집한 verification 결과는 다음과 같다.

```text
v21.3.0: verified=false, reason=unknown_key
v20.2.2: verified=false, reason=unsigned
```
