# PR Reviewer Assignment Service

Сервис для автоматического назначения ревьюверов на Pull Request'ы.

## Что делает

- Создаёт команды разработчиков
- Автоматически назначает до 2 ревьюверов на PR из команды автора
- Переназначает ревьюверов если нужно
- Блокирует изменения после merge PR

## Технологии

- Go 1.25.2
- PostgreSQL 15
- Echo (HTTP framework)
- SQLC (генерация SQL кода)
- Docker & Docker Compose

## Быстрый старт

Для запуска проекта нужно выполнить команду `docker-compose up --build`.
После этого сервис будет доступен на порту `:8080`

## API

### Создание команды

**Endpoint:** `POST /team/add`

**Request:**
```http
POST http://localhost:8080/team/add
Content-Type: application/json
```
```json
{
  "team_name": "backend",
  "members": [
    {
      "user_id": "u1",
      "username": "Alice", 
      "is_active": true
    },
    {
      "user_id": "u2",
      "username": "Bob",
      "is_active": true
    },
    {
      "user_id": "u3", 
      "username": "Charlie",
      "is_active": true
    }
  ]
}
```

### Получение команды

**Endpoint:** `GET /team/get`

**Request:**
```http
GET http://localhost:8080/team/get
Content-Type: application/json
```

### Установка актвиности для пользователя

**Endpoint:** `POST /users/setIsActive`

**Request:**
```http
POST http://localhost:8080/users/setIsActive
Content-Type: application/json
```
```json
{
  "user_id": "u2",
  "is_active": false
}
```

### Переназначить ревьюера

**Endpoint:** `POST /pullRequest/reassign`

**Request:**
```http
POST http://localhost:8080/pullRequest/reassign
Content-Type: application/json
```
```json
{
  "pull_request_id": "pr-1001",
  "old_reviewer_id": "u2"
}
```

### Смержить PR (идемпотентно)

**Endpoint:** `POST /pullRequest/merge`

**Request:**
```http
POST http://localhost:8080/pullRequest/merge
Content-Type: application/json
```
```json
{
  "pull_request_id": "pr-1001"
}
```

### Посмотреть PR ревьювера

**Endpoint:** `GET /users/getReview`

**Request:**
```http
GET http://localhost:8080/users/getReview?user_id=u3
Content-Type: application/json
```

## Дополнительно
### Описал конфигурацию линтера
Описана в файле `.golangci.yml`
### Добавил простой эндпоинт статистики
#### Статистика по пользователям

**Endpoint:** `GET /stats/users`

**Request:**
```http
GET http://localhost:8080/stats/users
Content-Type: application/json
```
**Response:**
```json
{
  "users": [
    {"user_id": "u1", "username": "Alice", "team_name": "backend", "assignments_count": 3},
    {"user_id": "u2", "username": "Bob", "team_name": "backend", "assignments_count": 2}
  ]
}
```

#### Статистика по PR

**Endpoint:** `GET /stats/prs`

**Request:**
```http
GET http://localhost:8080/stats/prs
Content-Type: application/json
```
**Response:**
```json
{
  "total_prs": 5,
  "open_prs": 2,
  "merged_prs": 3
}

```
#### Статистика по нагрузке ревьюеров

**Endpoint:** `GET /stats/workload`

**Request:**
```http
GET http://localhost:8080/stats/workload
Content-Type: application/json
```

**Response:**
```json
{
  "reviewers": [
    {"user_id": "u1", "username": "Alice", "team_name": "backend", "open_prs_count": 2},
    {"user_id": "u2", "username": "Bob", "team_name": "backend", "open_prs_count": 1}
  ]
}

```