# План обучения и развития Investment Ledger

## Цель и устройство плана

План не привязан к датам. Этапы идут в порядке технических зависимостей: следующий начинается после воспроизводимой контрольной работы предыдущего, а не после наступления календарной даты.

У плана три равноправные цели:

1. Последовательно восстановить и углубить Go: от языка, ошибок и тестов до конкурентности, распределённых систем и эксплуатации.
2. Развить Investment Ledger от консольного примера до многопользовательского web-продукта и затем до распределённой системы в Yandex Cloud.
3. Подготовиться к трудоустройству Go-разработчиком через воспроизводимые продуктовые артефакты, системный дизайн и практику собеседований.

Все продуктовые этапы обязательны для итогового portfolio release. Пометка **lab** означает, что тему нужно понять и проверить на небольшом примере, но не нужно искусственно встраивать её в продукт, если для неё нет реальной задачи.

## Итоговый продукт

Investment Ledger должен пройти четыре демонстрационных состояния:

| Checkpoint | Результат |
|---|---|
| Core backend | Домен, HTTP API, PostgreSQL, идемпотентность, конкурентные инварианты и тесты работают как modular monolith. |
| Secure connected product | Есть регистрация, Keycloak/OIDC, роли и ownership, котировки MOEX, публичный backend в Yandex Cloud и отдельный простой frontend. |
| Distributed system | Market data и reporting выделены в сервисы, обмен событиями идёт через Kafka с outbox/inbox и наблюдаемой доставкой. |
| Cloud portfolio release | Система развёрнута в Managed Kubernetes, использует управляемые хранилища, service mesh, IaC, CI/CD и проходит полный E2E-сценарий. |

### Функциональное ядро

- `UserProfile` связывает доменную модель с неизменяемым Keycloak `sub`.
- `Portfolio` хранит базовую валюту и принадлежит пользователю.
- `PortfolioMembership` задаёт объектные роли `owner`, `editor`, `viewer`.
- Глобальная роль `admin` предназначена для эксплуатации и справочников; доступ администратора к пользовательским данным журналируется.
- `Account` описывает брокерский или денежный счёт.
- `Instrument` хранит внутренний ID, `SECID`, board, тип и валюту инструмента.
- `Operation` является append-only записью: пополнение, вывод, покупка, продажа, дивиденд, комиссия и компенсирующая операция.
- `Quote` хранит цену, источник, `quoted_at`, `received_at` и признак устаревания.
- `Position` и `PortfolioSummary` вычисляются из операций и котировок.

### Инварианты

- Операции являются source of truth; производные проекции можно полностью перестроить.
- Покупка или вывод не делает cash balance отрицательным.
- Продажа не делает позицию отрицательной.
- Для денег запрещён `float64`; масштаб и округление фиксируются в доменной модели.
- Количество, цена, UTC-время и порядок событий имеют документированную точность и семантику.
- В первой продуктовой вертикали используется одна базовая валюта без FX-конвертации.
- Повтор одного idempotency key с тем же payload возвращает сохранённый результат, а с другим payload — `409 Conflict`.
- Любой пользовательский запрос фильтруется по `sub`, membership и роли; существование чужого ресурса не раскрывается. Машинный запрос использует отдельную service identity и минимальный scope.
- Рыночные данные не считаются realtime: клиент видит источник и давность котировки.
- Ни один сервис не читает таблицы другого сервиса.

### Осознанно не входит

- размещение реальных заявок и хранение брокерских ключей;
- обещание realtime MOEX без отдельного договора на данные;
- налоговый учёт, tax lots и корпоративные действия;
- мультивалютные портфели и FX;
- сложный UI, мобильное приложение и дизайн-система;
- внедрение Redis, MongoDB, Cgo или другого инструмента только ради упоминания в резюме.

## Как проходит каждый урок

1. Воспроизвести предыдущую идею по памяти.
2. Сформулировать одну новую идею простыми словами.
3. Получить небольшое условие и критерии приёмки.
4. Сделать самостоятельную попытку.
5. Добавить тест или другую воспроизводимую проверку.
6. Выполнить `gofmt`, `go vet`, `go test`; для конкурентного кода — `go test -race`.
7. Сделать self-review и один небольшой осмысленный commit.
8. Обновить checklist и при необходимости оставить поясняющий комментарий в GitHub issue текущего этапа.

Подсказки выдаются ступенчато: наводящий вопрос → схема или псевдокод → ревью попытки → готовый фрагмент только последним шагом.

## Этапы и отслеживание прогресса

Этот файл фиксирует содержание и порядок обучения, но не текущий статус. Прогресс ведётся в [GitHub milestones](https://github.com/anton415/go-investment-ledger/milestones) и связанных [issues](https://github.com/anton415/go-investment-ledger/issues): один issue соответствует одному этапу и содержит его checklist и контрольную работу.

| Этап Investment Ledger | Проверяемый результат |
|---|---|
| [1. Базовый Go](https://github.com/anton415/go-investment-ledger/issues/1) | Консольный ledger и небольшие упражнения по языку. |
| [2. Процессы, ошибки, домен](https://github.com/anton415/go-investment-ledger/issues/2) | Process Map Investment Ledger, финансовые инварианты, error flow и tests. |
| [3. Пакеты и тесты](https://github.com/anton415/go-investment-ledger/issues/3) | Границы `domain/application/ports`, DI, unit и contract tests. |
| [4. HTTP](https://github.com/anton415/go-investment-ledger/issues/4) | OpenAPI, handlers, middleware и `httptest`. |
| [5. PostgreSQL](https://github.com/anton415/go-investment-ledger/issues/5) | PGX/Goose, транзакции, Testcontainers и анализ индексов. |
| [6. Конкурентность](https://github.com/anton415/go-investment-ledger/issues/6) | Race-safe операции, cancellation, worker-pool exercises и `go test -race`. |
| [7. Keycloak/RBAC](https://github.com/anton415/go-investment-ledger/issues/7) | Регистрация, OIDC/JWT, роли, ownership и auth integration tests. |
| [8. MOEX](https://github.com/anton415/go-investment-ledger/issues/8) | MOEX ISS adapter, contract tests и degradation policy. |
| [9. Production readiness](https://github.com/anton415/go-investment-ledger/issues/9) | Logs/metrics/traces, profiling, Docker и CI basics. |
| [10. Первый cloud deploy](https://github.com/anton415/go-investment-ledger/issues/10) | Terraform, Yandex Cloud, secrets, rollout и recovery runbooks. |
| [11. Frontend](https://github.com/anton415/go-investment-ledger/issues/11) | Отдельный web-client, PKCE, security headers и browser E2E. |
| [12. Микросервисы](https://github.com/anton415/go-investment-ledger/issues/12) | Service boundaries, template, sync contracts, gRPC и resilience tests. |
| [13. Kafka](https://github.com/anton415/go-investment-ledger/issues/13) | Версионированные события, outbox/inbox, replay и DLQ. |
| [14. Kubernetes](https://github.com/anton415/go-investment-ledger/issues/14) | Helm, probes, NetworkPolicy и local-cluster E2E. |
| [15. CI/CD и mesh](https://github.com/anton415/go-investment-ledger/issues/15) | Promotion/rollback, GitLab Runner lab, mTLS и traffic policy. |
| [16. Финальный cloud deploy](https://github.com/anton415/go-investment-ledger/issues/16) | Managed Kubernetes/Kafka/PostgreSQL, IaC, monitoring и disaster-recovery tests. |
| [17. Собеседования](https://github.com/anton415/go-investment-ledger/issues/17) | HLD, mock interviews, решения задач, резюме и demo проекта. |

## Этап 1. Базовый Go

В текущем Investment Ledger реализован консольный срез. Package boundaries, интерфейсы и инфраструктурный каркас появятся в соответствующих последующих этапах и не считаются уже готовым кодом этого репозитория.

1. Модули, пакеты, `internal`, импорты и точка входа.
2. Переменные, константы, zero values и именованные типы.
3. Функции и несколько возвращаемых значений.
4. `if`, `switch`, `for`, `range` и closures.
5. Arrays, slices, maps и строки/runes.
6. Structs, методы, receivers и pointers.
7. Инкапсуляция и защита внутренних slices/maps от утечки.
8. Интерфейсы и структурная типизация.

Контрольная работа: консольный Investment Ledger объясняется без редактора и воспроизводит изученные языковые конструкции. PGX/Goose и HTTP-фреймворк осознанно появляются только на этапах 4–5.

## Этап 2. Процессинговое мышление, ошибки и финансовый домен

1. Product goal, роли и интересы участников.
2. Ubiquitous language, границы домена, Process Map L0/L1 и жизненные циклы.
3. Бизнес-правила, политики, доменные события и негативные сценарии.
4. `error` как значение, sentinel errors, `%w`, `errors.Is` и `errors.As`.
5. Валидация идентификатора, `SECID`, количества, цены и времени.
6. `defer`, освобождение ресурсов и составные ошибки cleanup.
7. `time.Time`, UTC, `as_of` и стабильный порядок операций.
8. Money/Quantity: minor units или fixed scale, округление и запрет `float64`.
9. Average cost, cash balance, total value, P&L, отсутствие и устаревание котировки.
10. Семантика compensating operation и idempotency key фиксируется в ADR.

Контрольная работа: доменная модель и Process Map согласованы; тесты покрывают oversell, отрицательный cash, duplicate ID, округление, одинаковое время и отсутствие котировки.

## Этап 3. Пакеты, интерфейсы и тестирование

1. Пакеты `domain`, `application`, `ports`, `memory` и направление зависимостей.
2. Композиция, маленькие интерфейсы рядом с потребителем и constructor injection.
3. In-memory repositories без глобального состояния.
4. `io.Reader`/`io.Writer`, JSON и границы serialization.
5. Unit tests, table tests, subtests, helpers и deterministic clock/ID generator.
6. Fakes, stubs, mocks и contract tests без привязки домена к библиотеке.
7. Coverage, `gofmt`, `go vet`, линтер, Makefile и повторяемые команды.
8. KISS, DRY, YAGNI, SOLID и GRASP с поправкой на идиомы Go.
9. **Lab:** generics/comparable на маленькой задаче; reflection — только на JSON/struct tags.

Контрольная работа: in-memory use cases создают пользователя, портфель, счёт, инструмент и операции, вычисляют summary и изолируют данные двух владельцев без HTTP и БД.

## Этап 4. HTTP, API и сети

1. `net/http`, lifecycle запроса и HTTP semantics.
2. REST resources, JSON DTO и преобразование в доменную модель.
3. Status codes, единый problem format и различие validation/not found/conflict.
4. `context.Context`, deadlines, timeouts и graceful shutdown.
5. Middleware: request ID, recovery, access log, body limit и rate limit.
6. Пагинация, фильтры, versioning, OpenAPI и примеры `curl`.
7. `httptest`, handler tests и E2E in-memory flow.
8. TCP/IP, DNS, TLS, CORS и reverse proxy на пути одного запроса.
9. **Lab:** реализовать один endpoint на Fiber и сравнить со стандартной библиотекой в ADR.

Контрольная работа: стабильный in-memory API проводит полный ledger flow, а тесты покрывают happy path, validation, not found, conflict, timeout и отмену запроса.

## Этап 5. PostgreSQL, транзакции и идемпотентность

1. Реляционная модель, ключи, ограничения и нормализация.
2. DDL/DML, `JOIN`, `GROUP BY`, CTE и оконные функции.
3. PGX pool, Goose migrations и точное SQL-отображение денежных типов.
4. Repository adapters и отображение ошибок БД в ошибки приложения.
5. Integration tests с PostgreSQL через Testcontainers.
6. Транзакции, isolation levels, row locks, deadlocks и retry policy.
7. Атомарная проверка cash/position и append operation без сохранения позиции как source of truth.
8. Idempotency fingerprint, операция и сохранённый HTTP-response в одной транзакции.
9. Индексы и `EXPLAIN (ANALYZE, BUFFERS)` на реальных запросах.
10. Separate database user, migrations user, backup/restore smoke test и безопасные seed data.
11. **Lab:** сравнить PGX и `sqlx` на одном repository method, не смешивая их в production code.

Контрольная работа: миграции поднимают пустую БД; операции переживают рестарт; повтор key/payload воспроизводит ответ; конфликтующий payload получает `409`; конкурентные записи не нарушают инварианты.

## Этап 6. Конкурентность Go

1. Ядра, процессы, goroutines и планировщик Go.
2. Channels, ownership, close и `select`.
3. `WaitGroup`, `Mutex`, `RWMutex`, `Once`, `Pool`, atomics и `sync.Map`.
4. Memory model, data race и race detector.
5. Deadlock, goroutine leak и cancellation через context.
6. Pipeline, FanOut/FanIn и backpressure.
7. Worker pool, `errgroup`, semaphore и rate limiting.
8. Поиск с отменой и продвинутый worker pool как отдельные упражнения.
9. Шардированный cache и измерение contention.
10. Конкурентные integration tests для продажи и HTTP-idempotency.

Контрольная работа: `go test -race ./...` проходит; параллельные продажи, повторные запросы, cancellation и shutdown проверены воспроизводимо.

## Этап 7. Регистрация, Keycloak и ролевая модель

1. OAuth 2.0, OpenID Connect, JWT, access/refresh token, `issuer`, `audience`, `exp` и JWKS.
2. Keycloak realm, backend client, frontend public client, client roles и service accounts.
3. Самостоятельная регистрация, login, logout, refresh, verified email и сброс пароля. Локально используется mail catcher, а публично — настроенный SMTP provider, generic responses, rate limit, audit и abuse controls.
4. Глобальные роли `user`, `admin` и объектные роли `owner`, `editor`, `viewer`.
5. `Portfolio.owner_subject = sub`, membership и запрет cross-user access.
6. Backend принимает `Authorization: Bearer <access_token>` и локально проверяет подпись и claims через JWKS.
7. `401` означает отсутствие/недействительность identity, `403` — недостаточные права; `404` скрывает чужой ресурс.
8. Конфигурация realm/client экспортируется и версионируется без паролей и client secrets.
9. Для service-to-service взаимодействия готовится отдельный client credentials flow.
10. Тесты: два пользователя, все роли, истёкший token, неверный issuer/audience, rotation ключа и недоступный Keycloak.

В продуктовом контуре refresh token не передаётся backend на каждом запросе: браузер использует Authorization Code + PKCE, а API валидирует access token.

Контрольная работа: пользователь регистрируется в Keycloak, управляет только своими портфелями, viewer не пишет, editor не меняет membership, admin-действия аудируются, а матрица разрешений покрыта integration tests.

## Этап 8. Интеграция с MOEX ISS

1. Разобрать ISS endpoints, engines, markets, boards, `SECID`, securities, marketdata, history и candles.
2. Отделить внутренний `InstrumentID` от MOEX identifiers и выбрать правило board.
3. Ввести порт `MarketDataProvider` и MOEX adapter, не протекая ISS DTO в домен.
4. Хранить `source`, `quoted_at`, `received_at`, freshness и исходную валюту.
5. Добавить HTTP timeout, retry с backoff/jitter, rate limit, cache и circuit breaker.
6. Определить поведение при закрытой бирже, выходных, пустом ответе, stale quote и неизвестном `SECID`.
7. Сохранить ручную котировку как явный fallback с audit trail.
8. Использовать fixtures/contract tests; live smoke test запускается отдельно и не делает CI зависимым от MOEX.
9. Проверить правила использования и отображения рыночных данных перед публичным demo.

Контрольная работа: сервис получает delayed quote из ISS, корректно маркирует давность, переживает timeout/429/5xx, использует manual fallback и не обещает realtime.

## Этап 9. Production-ready modular monolith и основы CI/CD

1. Проследить путь изменения от commit до локального release artifact и определить владельца каждого gate.
2. Dockerfile, Docker Compose, non-root runtime и versioned images.
3. Compose поднимает backend, PostgreSQL, Keycloak и observability dependencies одной командой.
4. Structured logs, correlation ID, Prometheus metrics, Grafana dashboard, OpenTelemetry/Jaeger и health/readiness.
5. Benchmarks, allocations, GC, `pprof`, escape analysis и сохранённый baseline.
6. CI architecture: format, lint, unit/integration/E2E, race, vulnerability scan, image build и SBOM.
7. Immutable image tags, artifact retention, локальный registry lab и правила promotion между окружениями.
8. Secrets не лежат в Git/image и не выводятся в pipeline logs; approvals и least privilege описаны до cloud deploy.
9. Миграции выполняются отдельным release step и проверяются на чистой БД.

Контрольная работа: clean checkout одной командой поднимает modular monolith, CI выпускает immutable image и артефакты, а logs/metrics/traces проводят один пользовательский запрос end-to-end.

## Этап 10. Первый backend deploy в Yandex Cloud

1. Terraform baseline: IAM/service account, VPC, security groups и Container Registry.
2. Compute Cloud VM запускает versioned backend, Keycloak и reverse proxy через Docker Compose.
3. Managed PostgreSQL доступна только по приватной сети; у ledger и Keycloak отдельные databases/users, а migrations user отделён от runtime user.
4. Keycloak использует realm import/export, persistent database, backup, proxy/hostname settings; admin console закрыта allowlist/VPN.
5. Lockbox, DNS, Certificate Manager и TLS убирают secrets и открытый HTTP из deploy flow.
6. CI публикует image, выполняет migration step, обновляет окружение и проверяет smoke scenario.
7. Настроить backup/restore, rollback и runbook недоступности БД/Keycloak/MOEX.
8. Добавить budget alerts, quotas note и безопасную teardown инструкцию для всех платных ресурсов.

Контрольная работа: backend имеет стабильный HTTPS URL в Yandex Cloud, login и MOEX работают извне, deploy воспроизводится из Terraform/CI, а rollback, восстановление БД и удаление ресурсов проверены.

## Этап 11. Отдельный простой frontend

Frontend создаётся как отдельный компонент в каталоге `web/` этого monorepo. Базовый вариант — React + TypeScript + Vite; более сложный framework не нужен без конкретной причины.

1. Независимые build, lint, tests, image/artifact и CI.
2. Keycloak Authorization Code + PKCE, login/logout/register и восстановление сессии без client secret.
3. Страницы: портфели, счета, инструменты, добавление операции, список операций и portfolio summary.
4. Отображение loading/empty/error, `401`/`403`, stale quote и недоступности backend/MOEX.
5. Typed API client из OpenAPI или небольшой ручной adapter; base URL задаётся окружением.
6. CORS allowlist, CSP, XSS/CSRF-модель и безопасное хранение token state.
7. Component tests и Playwright E2E для login → operation → summary.
8. Deployment статических файлов в Yandex Object Storage + Cloud CDN, DNS и HTTPS.
9. Immutable frontend artifact фиксируется вместе с backend/service image digests в `deploy/release-manifest.yaml`.

Контрольная работа: публичный frontend проходит браузерный E2E и запрашивает данные только у уже работающего backend в Yandex Cloud; в bundle и репозитории нет secrets.

## Этап 12. Границы микросервисов и синхронные вызовы

1. Объяснить, почему modular monolith стал недостаточен, и зафиксировать decision record разделения.
2. Выделить `market-data-service`; `ledger-service` сохраняет Portfolio/Account/Operation и транзакционные инварианты.
3. Зафиксировать будущую границу reporting projection, но не выделять `reporting-service` до появления событий следующего этапа.
4. Создать единый service template: config, logs, metrics, traces, health, shutdown, migrations и CI.
5. Каждый существующий сервис владеет своей БД, пользователем, миграциями и schema; cross-service SQL запрещён.
6. Версионировать REST/gRPC/Protobuf contracts и проверять backward compatibility.
7. User-initiated flow сохраняет end-user identity: используется propagation access token или token exchange с корректной audience, а downstream повторно проверяет claims и business permissions.
8. Machine/background flow использует отдельный Keycloak client credentials token с минимальным service scope.
9. Настроить timeout budget, retry только безопасных/idempotent операций, circuit breaker и bulkhead.
10. Добавить service discovery/API gateway и contract/failure tests.

Контрольная работа: ledger и market data независимо собираются и разворачиваются, не разделяют таблицы, а ledger продолжает работать при недоступности market data в documented degraded mode.

## Этап 13. Kafka, события и распределённая согласованность

1. Kafka broker, topics, partitions, keys, offsets, consumer groups и ordering guarantees.
2. События `operation.recorded.v1` и `quote.updated.v1`; envelope содержит `event_id`, `schema_version`, occurred/published time, producer и trace ID.
3. Partition key выбирается по агрегату: portfolio для операций, instrument для котировок.
4. Transactional outbox публикует событие после локальной транзакции без dual write.
5. Inbox/idempotent consumer защищает от повторной доставки; delivery semantics явно называются at-least-once.
6. Выделить `reporting-service`: он строит собственную projection из событий и умеет rebuild/replay.
7. Retry topics или backoff, DLQ, poison message policy и ручное восстановление.
8. Schema evolution, compatibility tests и запрет ломать старых consumers.
9. Метрики producer errors, consumer lag, retries, duplicates и DLQ; trace проходит через message headers.
10. Проверить duplicate, reorder в допустимых границах, consumer restart, broker outage и replay.
11. Изучить saga и compensation; в ADR объяснить, нужен ли saga выбранному процессу. Если нет — выполнить изолированный failure-injection lab, не добавляя искусственную business saga.
12. Зафиксировать eventual-consistency contract: отдельные watermarks/versions для operation и quote streams, допустимый projection lag и отсутствие общего порядка между разными partition keys.
13. Frontend показывает `processing`/`stale`, не обещает мгновенный read-after-write и умеет дождаться нужной projection version.

Контрольная работа: повторная доставка или рестарт reporting consumer не создаёт дублей, projection перестраивается с нуля, а потерянные/сломанные сообщения диагностируются через metrics, logs и DLQ.

## Этап 14. Kubernetes как runtime

1. Pods, Deployments, Services, Jobs/CronJobs, namespaces и labels.
2. ConfigMaps, Secrets и различие cloud IAM с Kubernetes service accounts/RBAC.
3. Startup/liveness/readiness probes и graceful pod termination.
4. Requests/limits, HPA, PDB и placement basics.
5. Migration Job, Ingress/Gateway, TLS и service discovery.
6. NetworkPolicy и минимальные разрешения между сервисами.
7. Helm charts с values для local/staging/prod без копирования manifest trees.
8. Local cluster через kind или k3d, локальный registry и repeatable bootstrap.
9. Rolling update, rollback, pod-kill и node-failure exercises.

Контрольная работа: fresh local cluster поднимается одной документированной последовательностью, проходит E2E, rolling update и pod-kill test без потери принятых операций.

## Этап 15. Продвинутый CI/CD, локальные окружения и service mesh

1. Расширить путь изменения от готовых release artifacts до локальных dev/staging namespaces и определить владельцев promotion gates.
2. Выполнить изолированные labs с локальными GitLab и GitLab Runner; проектовый pipeline может остаться в GitHub Actions.
3. Pipeline использует по одному проверенному immutable image на компонент и единый release manifest: deploy, verification, promotion и rollback идут без повторной сборки.
4. Semantic versions, artifact retention, approvals и environment protection.
5. Secrets не выводятся в logs и выдаются по принципу least privilege.
6. Ошибки pipeline имеют runbook, владельца и критерий остановки promotion.
7. Выбрать Istio или Linkerd в ADR и объяснить эксплуатационную цену mesh.
8. Включить strict mTLS между сервисами, service identity и authorization policy.
9. Проверить traffic split/canary, timeout и telemetry; mesh retry запрещён для неидемпотентного POST без idempotency key.
10. Сравнить app-level resilience и mesh policy, не дублируя бесконтрольно retries.

Контрольная работа: pipeline продвигает один release manifest через локальные dev/staging окружения, rollback воспроизводим, mTLS проверен тестом, а failure scenario виден в distributed trace. Применение этого flow к cloud production выполняется на следующем этапе.

## Этап 16. Финальный distributed deploy в Yandex Cloud

1. Terraform создаёт VPC/subnets/security groups, IAM/service accounts и Container Registry; remote state шифруется, блокируется и не содержит secret values.
2. Managed Kubernetes, node groups и cluster access поднимаются без ручных скрытых шагов.
3. Managed PostgreSQL использует раздельные databases/users, backups и private connectivity.
4. Managed Kafka использует private connectivity, TLS/SASL, минимальные ACL, topic retention и monitoring; события не содержат лишние PII.
5. Lockbox secrets доставляются в pods без записи в image/state, имеют rotation procedure; Certificate Manager, DNS, Application Load Balancer и TLS настраиваются как код.
6. Keycloak запускается как отдельный Kubernetes workload с dedicated Managed PostgreSQL database/user, realm import/migration, backup/restore, proxy/hostname settings и закрытой admin console.
7. Helm разворачивает ledger, market data, reporting, Keycloak и observability, после чего cloud pipeline продвигает release manifest через staging в production.
8. Frontend остаётся в Object Storage + CDN и обращается к публичному API по HTTPS.
9. Yandex Monitoring/alerts, Prometheus/Grafana и distributed tracing покрывают пользовательский сценарий.
10. Backup/restore, Keycloak migration, rolling update, broker outage, pod failure и disaster recovery smoke tests.
11. Budget alerts, quotas, cost note и безопасный teardown всех платных ресурсов.

Контрольная работа: чистое окружение разворачивается из IaC и проходит сценарий Keycloak registration/login → operation → PostgreSQL → MOEX quote → Kafka → reporting projection → portfolio summary во frontend. Отдельный сценарий доказывает отсутствие дублей после retry/restart.

## Этап 17. System Design, собеседования и упаковка опыта

1. Для Investment Ledger подготовить requirements, capacity estimate, HLD, storage choice и границы сервисов.
2. Объяснить sync/async, consistency, idempotency, reliability, cache, replicas и load balancing.
3. Разобрать failures, retry/timeout/DLQ, observability, security, rate limiting и external integrations.
4. Провести mock system design и записать ошибки объяснения.
5. Повторить Go Core и goroutines по interview checklist, решая задачи без редактора.
6. Повторить SQL, indexes, Kafka, Redis, MongoDB, gRPC, Docker, Kubernetes, CAP, saga/outbox, discovery/gateway/resilience.
7. **Labs:** Redis и MongoDB на уменьшенных примерах; Cgo — устройство и риски на минимальном упражнении. В основной продукт они входят только при появлении задачи.
8. Подготовить README, OpenAPI, ERD/C4, ADR, runbooks, demo script, резюме и историю решений проекта.
9. Пройти HR, technical screening, full interview, system design и TeamFit rehearsals.

Контрольная работа: автор проводит demo и mock interview без редактора, связывает ответы с работающей системой и не выдаёт учебные labs за production experience.

## Definition of Done итогового portfolio release

- Из одного чистого checkout команда читает `deploy/release-manifest.yaml` и поднимает закреплённые image digests backend, frontend и сервисов; frontend также собирается из каталога `web/`.
- Пользователь регистрируется и входит через Keycloak; роли и ownership проверены тестами двух пользователей.
- Ledger invariants, idempotency и конкурентные сценарии покрыты unit/integration/E2E tests.
- MOEX adapter имеет fixtures, live smoke test, stale-data policy и manual fallback.
- Отдельный frontend работает с backend в Yandex Cloud и не содержит secrets.
- Сервисы имеют раздельное владение данными и версионированные HTTP/gRPC/event contracts.
- Kafka delivery выдерживает duplicates, restart и replay благодаря outbox/inbox.
- Local Kubernetes и Yandex Managed Kubernetes разворачиваются воспроизводимо.
- Service mesh даёт mTLS, telemetry и одну осмысленную traffic policy.
- CI/CD, migrations, backup/restore, rollback, monitoring, alerts, budget и teardown документированы.
- Опубликованы OpenAPI, ERD/C4, ADR, runbooks, безопасные demo data и сценарий демонстрации.
- Все 17 roadmap issues закрыты только после выполнения checklist и контрольной работы; GitHub milestones показывают актуальный прогресс без правок этого файла.

## Основные источники

- [Официальная документация Keycloak OIDC](https://www.keycloak.org/securing-apps/oidc-layers) и [JavaScript adapter с PKCE](https://www.keycloak.org/securing-apps/javascript-adapter).
- [MOEX ISS reference](https://iss.moex.com/iss/reference/) и [MOEX Market Data Policy](https://www.moex.com/en/datapolicy/).
- [Yandex Managed Kubernetes](https://yandex.cloud/en/docs/managed-kubernetes/), [Managed PostgreSQL](https://yandex.cloud/en/docs/managed-postgresql/), [Managed Kafka](https://yandex.cloud/en/docs/managed-kafka/), [Container Registry](https://yandex.cloud/en/docs/container-registry/) и [Object Storage + CDN](https://yandex.cloud/en/docs/storage/tutorials/cdn-hosting/console).
- [Kubernetes workloads](https://kubernetes.io/docs/concepts/workloads/), [Services](https://kubernetes.io/docs/concepts/services-networking/service/) и [probes](https://kubernetes.io/docs/concepts/workloads/pods/probes/).
- [Istio в Yandex Managed Kubernetes](https://yandex.cloud/en/docs/marketplace/tutorials/istio).
