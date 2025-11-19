Проект написан на Go, использует Gin для HTTP API и PostgreSQL для хранения данных.

```
.
├── cmd/
│   └── main.go            # Точка входа
├── internal/
│   ├── app/http/           # HTTP сервер и маршруты
│   ├── domain/            # Сущности, репозитории, сервисы
│   └── infrastructure/db/ # PostgreSQL репозитории
├── migrations/            # SQL миграции для таблиц
├── Dockerfile
├── docker-compose.yml
└── README.md
```

Требования

- Docker >= 20.x
- Docker Compose >= 1.29.x
- Go >= 1.23.x (для локального запуска без Docker)

Запуск через Docker Compose

Склонируйте репозиторий:
```
git clone https://github.com/f4ke-n0name/avito.git
cd avito
```

Проверьте, что в корне есть папка migrations/ с SQL-файлами создания таблиц:
```
migrations/
├── 001_init.sql
```

Запустите контейнеры:

```
docker-compose up --build
```

> ⚠️ **Важно!**
>
> Файлы shell-скриптов (`*.sh`), созданные в Windows, могут автоматически сохраняться  
> с окончаниями строк **CRLF**, из-за чего Alpine Linux выдаёт ошибку:
>
> ```
> avito-app: ./migrate.sh: not found
> ```
>
> При этом файл существует и имеет права на выполнение — проблема именно в формате строк.
>
> **Перед использованием обязательно убедитесь, что файл сохранён в формате LF.**  
> В VS Code это меняется справа внизу (**CRLF → LF**).

PostgreSQL будет доступен на порту 5432.
Go-приложение стартует на http://localhost:8080.

В Docker Compose прокидывается:
```
DATABASE_URL=postgres://user:pass@db:5432/yourdb
```

Примеры API
1. Добавить команду
```
POST /team/add

{
  "team_name": "TeamAlpha",
  "members": [
    {"user_id": "u1", "username": "Alice", "is_active": true},
    {"user_id": "u2", "username": "Bob", "is_active": true},
    {"user_id": "u3", "username": "Charlie", "is_active": false}
  ]
}
```


Ответ:
```
{
  "team": {
    "team_name": "TeamAlpha",
    "members": [
      {"user_id": "u1", "username": "Alice", "is_active": true, "team_name": "TeamAlpha"},
      {"user_id": "u2", "username": "Bob", "is_active": true, "team_name": "TeamAlpha"},
      {"user_id": "u3", "username": "Charlie", "is_active": false, "team_name": "TeamAlpha"}
    ]
  }
}
```
2. Получить команду
```
GET /team/get?team_name=TeamAlpha

Ответ:

{
  "team_name": "TeamAlpha",
  "members": [
    {"user_id": "u1", "username": "Alice", "is_active": true, "team_name": "TeamAlpha"},
    {"user_id": "u2", "username": "Bob", "is_active": true, "team_name": "TeamAlpha"},
    {"user_id": "u3", "username": "Charlie", "is_active": false, "team_name": "TeamAlpha"}
  ]
}
```
3. Создать Pull Request
```
POST /pullRequest/create

{
  "pull_request_id": "pr1",
  "pull_request_name": "Add new feature",
  "author_id": "u1"
}


Ответ:

{
  "pr": {
    "pr_id": "pr1",
    "name": "Add new feature",
    "author_id": "u1",
    "status": "OPEN",
    "created_at": "2025-11-16T12:00:00Z",
    "merged_at": null,
    "reviewers": []
  }
}
```
4. Получить список PR для ревьюера
```
GET /users/getReview?user_id=u2

Ответ:

{
  "user_id": "u2",
  "pull_requests": [
    {
      "pr_id": "pr1",
      "name": "Add new feature",
      "author_id": "u1",
      "status": "OPEN",
      "created_at": "2025-11-16T12:00:00Z",
      "merged_at": null,
      "reviewers": ["u2"]
    }
  ]
}
```
5. Merge PR
```
POST /pullRequest/merge

{
  "pull_request_id": "pr1"
}
```

Ответ:
```
{
  "pr": {
    "pr_id": "pr1",
    "name": "Add new feature",
    "author_id": "u1",
    "status": "MERGED",
    "created_at": "2025-11-16T12:00:00Z",
    "merged_at": "2025-11-16T12:30:00Z",
    "reviewers": ["u2"]
  }
}
```
6. Поменять ревьюера
```
POST /pullRequest/reassign

{
  "pull_request_id": "pr1",
  "old_user_id": "u2"
}
```

Ответ:
```
{
  "pr": { ... },
  "replaced_by": "u3"
}
```
