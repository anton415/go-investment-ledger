# Go Investment Ledger

Учебный проект для последовательного восстановления Go и развития production-like backend на примере приложения для учёта инвестиций.

Investment Ledger начинается как modular monolith для учёта инвестиционных операций, затем становится многопользовательским web-продуктом с Keycloak и котировками MOEX, а после этого развивается в распределённую систему с Kafka, Kubernetes и service mesh в Yandex Cloud.

План не привязан к датам. Переход между этапами определяется работающей контрольной вертикалью, тестами и способностью объяснить решение без открытого редактора.

## Путь проекта

1. **Core backend:** Go, ошибки, доменная модель, интерфейсы, тесты, HTTP, PostgreSQL, идемпотентность и конкурентность.
2. **Secure connected product:** регистрация пользователей, Keycloak/OIDC, полноценная ролевая модель, ownership и интеграция с MOEX ISS.
3. **Web demo:** production-ready backend разворачивается в Yandex Cloud; отдельный простой frontend обращается к нему по HTTPS.
4. **Distributed system:** выделяются market data и reporting services, а события доставляются через Kafka с outbox/inbox.
5. **Cloud platform:** Docker, CI/CD, Kubernetes, Helm, service mesh, Terraform, наблюдаемость и финальный deploy в Yandex Cloud.
6. **Подготовка к трудоустройству:** System Design, Go Core, goroutines, microservices interviews и упаковка проекта.

## Итоговый функциональный scope

- пользователи, регистрация, login/logout и восстановление доступа через Keycloak;
- глобальные и объектные роли, изоляция данных владельцев и audit действий;
- портфели, счета и финансовые инструменты;
- append-only пополнения, выводы, покупки, продажи, дивиденды, комиссии и compensating operations;
- cash balance, позиции, average cost, market/total value и P&L;
- ручные котировки и delayed market data из MOEX ISS с явной давностью и fallback;
- идемпотентный HTTP API, PostgreSQL, миграции и integration tests;
- отдельный frontend с Authorization Code + PKCE;
- ledger, market-data и reporting services с версионированными контрактами;
- Kafka, transactional outbox, inbox/idempotent consumers, replay и DLQ;
- Docker/CI, Kubernetes/Helm, service mesh, metrics/logs/traces и IaC;
- публичный frontend и backend в Yandex Cloud.

Реальные сделки, брокерские ключи, realtime market data без соответствующего договора, налоги и мультивалютность сознательно не входят в учебный продукт.

## Прогресс обучения

README не хранит текущий этап или выполненные пункты. Актуальный прогресс ведётся в [GitHub milestones](https://github.com/anton415/go-investment-ledger/milestones) и [issues](https://github.com/anton415/go-investment-ledger/issues): один issue соответствует одному этапу и закрывается после воспроизводимой контрольной работы.

Запуск текущего примера:

```bash
go run .
```

Полный бездатный roadmap, контрольные работы и cloud architecture находятся в [LEARNING_PLAN.md](LEARNING_PLAN.md).

Промежуточные результаты этапа S02:

- [Product goal и границы Investment Ledger](docs/domain/product-goal.md).
- [Жизненные циклы команды и Operation](docs/domain/operation-lifecycle.md).
- [Process Map L0 — текстовый черновик](docs/processes/investment-ledger-l0.md).
- [Record operation L1 — текстовый черновик](docs/processes/record-operation-l1.md).
