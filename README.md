# 🧠 Архитектура микросервисного Twitter-клона

## 📌 Общая идея

Создание микросервисного клона Twitter с фокусом на реалистичную архитектуру backend-а:

* Высокая масштабируемость
* Расширяемость и независимость сервисов
* Реалистичная нагрузка через ботов и нагрузочное тестирование
* Продвинутая система мониторинга и логирования
* Использование real-time коммуникации через WebSocket
* Работа с Kafka, gRPC, Redis, PostgreSQL и др.

## 🔧 Общие технологии

| Компонент       | Технологии                                |
| --------------- | ----------------------------------------- |
| Язык            | Go                                        |
| Контейнеризация | Docker + Docker Compose                   |
| API             | REST + gRPC                               |
| Сообщения       | Kafka, Redis Streams                      |
| СУБД            | PostgreSQL, Redis, ClickHouse (аналитика) |
| Кэш             | Redis                                     |
| CI/CD           | GitHub Actions                            |
| Мониторинг      | Prometheus, Grafana                       |
| Логирование     | Loki + Grafana, zap/logrus                |
| Реалтайм        | WebSocket-шлюз                            |
| Инфраструктура  | 1–2 VPS (Ubuntu), TLS через Let's Encrypt |

---

# 📦 Сервисы

## 1. Gateway API

**Назначение:** точка входа в систему, маршрутизатор запросов

**Функции:**

* REST-обработка внешних HTTP-запросов
* Подключение к gRPC сервисам
* Авторизация (JWT)
* Поддержка WebSocket

**Стек:**

* Go + Gin/Chi
* gRPC-клиенты ко всем внутренним сервисам
* Redis (сессии, ratelimiter)
* WebSocket сервер
* Middleware: Auth, RequestID, Logger

---

## 2. Auth Service

**Назначение:** аутентификация, авторизация, refresh-токены

**Функции:**

* Регистрация и логин
* Генерация JWT и refresh-токенов
* Валидация токенов

**Стек:**

* PostgreSQL (таблица пользователей)
* Redis (сессии, временные токены)
* gRPC Server
* Kafka producer (UserCreated)
* bcrypt, jwt-go

---

## 3. User Profile Service

**Назначение:** профиль пользователя

**Функции:**

* Редактирование профиля
* Получение информации о пользователе
* Отображение аватарки, био и др.

**Стек:**

* PostgreSQL (таблица user\_profiles)
* gRPC Server
* Kafka (ProfileUpdated)
* Минималистичная REST для SSR

---

## 4. Tweet Service

**Назначение:** управление твитами

**Функции:**

* Создание, удаление, редактирование твитов
* Хранение текста, ссылок, изображений (если потребуется)

**Стек:**

* PostgreSQL
* Redis (кэш популярных твитов)
* Kafka (TweetCreated, TweetDeleted)
* gRPC Server

---

## 5. Feed Service

**Назначение:** генерация фидов пользователей

**Функции:**

* Fan-out on write или on read
* Обработка новых твитов

**Стек:**

* Redis (user feed cache)
* Kafka Consumer (TweetCreated)
* gRPC Server

---

## 6. Follow Service

**Назначение:** подписки/отписки

**Функции:**

* Подписка на пользователей
* Хранение связи many-to-many

**Стек:**

* PostgreSQL (таблица follows)
* Kafka (FollowCreated)
* gRPC Server

---

## 7. Like Service

**Назначение:** лайки твитов

**Функции:**

* Ставить и удалять лайк
* Подсчёт лайков

**Стек:**

* Redis (реалтайм)
* PostgreSQL (история)
* Kafka (LikeCreated)
* gRPC Server

---

## 8. Notification Service

**Назначение:** уведомления пользователей

**Функции:**

* Push-уведомления (WebSocket)
* Email/SMS (внешние)

**Стек:**

* PostgreSQL (уведомления)
* Redis PubSub → WebSocket Gateway
* Kafka Consumer (TweetLiked, Followed)

---

## 9. Search Service

**Назначение:** полнотекстовый поиск

**Функции:**

* Поиск твитов, пользователей

**Стек:**

* PostgreSQL FTS или Bleve/Meilisearch
* gRPC Server
* Kafka Consumer

---

## 10. Analytics Service

**Назначение:** аналитика активности

**Функции:**

* Подсчёт твитов, лайков, фолловеров
* Хранение статистики

**Стек:**

* ClickHouse
* Kafka Consumers (TweetCreated и др.)
* REST/gRPC для запроса отчётов

---

## 11. Media Service (опционально)

**Назначение:** загрузка и хранение медиа-файлов

**Функции:**

* Загрузка аватарок
* Раздача CDN ссылок

**Стек:**

* MinIO или файловая система
* REST API
* gRPC Server

---

## 12. Admin Service

**Назначение:** административная панель

**Функции:**

* Бан пользователей
* Просмотр логов и метрик

**Стек:**

* PostgreSQL
* gRPC + REST
* RBAC Middleware

---

## 13. WebSocket Gateway

**Назначение:** прокси-сервер для уведомлений и real-time взаимодействия

**Функции:**

* Подключение пользователей
* Подписка на события

**Стек:**

* Redis PubSub
* WebSocket сервер
* gRPC клиент → Notification Service

---

## 14. Email/SMS Service

**Назначение:** отправка внешних уведомлений

**Функции:**

* Email-уведомления
* Интеграция со сторонними API (или заглушки)

**Стек:**

* Kafka Consumer
* SMTP/Mailgun API (mock)
* gRPC Server

---

## 15. Logger / Log Collector

**Назначение:** централизованный сбор логов

**Функции:**

* Структурированные логи
* Отправка в Loki

**Стек:**

* Zap или Logrus
* Loki
* Grafana Dashboards

---

## 16. Rate Limiter Service

**Назначение:** ограничение количества запросов

**Функции:**

* RPS/объём по IP/токену
* Используется через SDK другими сервисами

**Стек:**

* Redis
* gRPC Server

---

# 🤔 Прочее

**Мониторинг:**

* Prometheus + Grafana
* AlertManager с отправкой в Telegram

**CI/CD:**

* GitHub Actions (build/test/lint)

**TLS:**

* Let's Encrypt

**Real-Time:**

* WebSocket Gateway вынесен в отдельный сервис

**Backup:**

* Пока не реализуется, фокус на фичах

**Service Mesh:**

* Не используется на MVP-этапе

**ML:**

* Не используется, real-time scoring не нужен

---

Если хочешь, можно начать с реализации `auth-service` и подготовить шаблоны всех сервисов с Makefile, Dockerfile и базовой gRPC-инфраструктурой.
