# Лаба 1. Первый сервис: HTTP CRUD в памяти

**Цель:** написать HTTP-сервис для главной сущности своей темы. Сервис умеет создавать, читать, изменять и удалять записи (это и есть CRUD: Create, Read, Update, Delete) и отвечает по правилам REST. Данные пока хранятся в памяти программы, в обычном списке или словаре. База данных появится в лабе 2.

Язык и фреймворк любые.

Срок: PR из `lab-1` в `main` до 23:59 накануне встречи 2.

## Что понадобится

| Что | Зачем |
| --- | --- |
| Git и аккаунт на GitHub | Репозиторий проекта и сдача через PR |
| Язык и фреймворк | Например, Python 3.10+ и FastAPI, Node.js и Express, Go 1.22+ |
| Менеджер пакетов своего языка | pip, npm, Go modules, Maven или Gradle, NuGet, Composer. Через него ставятся фреймворк и библиотеки, а список зависимостей хранится в файле проекта |
| VS Code с расширением REST Client | Отправлять запросы из файла `requests.http` прямо из редактора. В JetBrains IDE это встроено |
| Python 3.10+ и pipx | Чтобы поставить чекер: `pipx install git+https://github.com/bozhedom/backcheck@v1` |

## Что получится в итоге

Пример для темы «Библиотека», где главная сущность это книга.

```http
POST /api/books
Content-Type: application/json

{"title": "Мастер и Маргарита", "author": "Булгаков", "copies": 3}
```

```http
HTTP/1.1 201 Created
Content-Type: application/json

{"id": 1, "title": "Мастер и Маргарита", "author": "Булгаков", "copies": 3}
```

| Запрос | Что делает | Ответ |
| --- | --- | --- |
| `GET /health` | Проверка, что сервис жив | `200` и `{"status": "ok"}` |
| `POST /api/books` | Создать книгу | `201` и книга с `id`; `400`, если тело неправильное |
| `GET /api/books` | Список книг | `200` и массив |
| `GET /api/books/1` | Одна книга | `200`; `404`, если такой нет |
| `PATCH /api/books/1` | Изменить часть полей | `200` и обновлённая книга; `400`; `404` |
| `DELETE /api/books/1` | Удалить | `204` без тела; `404` |

## Что сделать по шагам

1. **Создать ветку** `lab-1` от `main`: `git switch -c lab-1`.
2. **Придумать поля сущности.** 3-5 полей, хотя бы одно обязательное строковое (например, `title`) и хотя бы одно число или статус. `id` ставит сервер, присылать его в теле не нужно.
3. **Поднять сервис на порту 8000** с эндпоинтом `GET /health`, который отвечает `{"status": "ok"}`.
4. **Сделать CRUD** по пути `/api/<сущности во множественном числе>`: `/api/books`, `/api/tours`, `/api/rooms`. Данные хранятся в словаре `{id: запись}`, id выдаёт счётчик, который только растёт.
5. **Обработать ошибки.** Неизвестный id даёт 404. Сломанный JSON, пустое или отсутствующее обязательное поле дают 400. Сервер при этом не падает и не отвечает 500.
6. **Записать зависимости в файл**: `requirements.txt`, `package.json`, `go.mod`, `pom.xml` или аналог для своего языка. По этому файлу проект ставится на любой машине одной командой, например `pip install -r requirements.txt` или `npm install`.
7. **Заполнить `README.md`**: тема, как установить зависимости и запустить, таблица эндпоинтов.
8. **Заполнить `requests.http`** всеми запросами. По этому файлу идёт демонстрация на паре, поэтому каждый запрос нужно хоть раз запустить.
9. **Заполнить `contract.json`** под свою тему, по нему работает чекер. Как заполнять, описано ниже.
10. **Проверить себя чекером**: `backcheck --lab 1 --stars`.
11. **Сдать**: `backcheck --lab 1 --stars --report`, коммит, push, PR из `lab-1` в `main`, в описании PR отметить чек-лист и вставить содержимое `backcheck-report.md`.

## Чек-лист базы (6 баллов) и как его проверяет чекер

| № | Пункт | Как проверяется |
| --- | --- | --- |
| 1 | README (тема, как запустить, таблица эндпоинтов), `requests.http` со всеми запросами, файл зависимостей | Чекер находит файлы, таблицу с путём ресурса и файл зависимостей. Раздел «как запустить» проверяется вручную |
| 2 | `GET /health` отвечает 200 и `{"status": "ok"}` | Код 200, заголовок `Content-Type: application/json`, в теле `"status": "ok"` |
| 3 | `POST /api/<items>` отвечает 201 и возвращает объект с `id` | Тело из `contract.create` даёт 201, в ответе есть `id`, все отправленные поля вернулись с теми же значениями |
| 4 | `GET /api/<items>` и `GET /api/<items>/{id}` работают, неизвестный id даёт 404 | Созданная запись есть в списке, по id приходит та же запись, `missing_id` даёт 404 |
| 5 | `PUT` или `PATCH` обновляет, `DELETE` отвечает 204, после удаления `GET` даёт 404 | PATCH (или PUT, если PATCH нет), потом GET видит изменение, DELETE отвечает ровно 204, потом GET даёт 404 |
| 6 | Сломанный JSON или пустое обязательное поле дают 400, а не 500 и не падение | Битый JSON, пустое и отсутствующее `required_field` дают 400 (422 тоже засчитывается), после этого сервер отвечает |

## ★ (2 балла): валидация и пагинация

- Проверяются **все** поля, ошибки приходят в одном формате: `{"error": "текст", "field": "имя поля"}` с кодом **400**.
- Список умеет фильтр и пагинацию: `?limit=10&offset=0&<поле>=<значение>`. `limit` это сколько записей отдать, `offset` это сколько пропустить с начала.

Чекер отправляет каждый случай из `contract.errors` (нужен хотя бы один случай на каждое поле из `create`) и ждёт 400 с правильным `field`. Фильтр должен вернуть только подходящие записи. `limit=2&offset=0` и `limit=2&offset=2` должны дать две разные страницы.

## ★★ (2 балла): документация и логи

- Документация OpenAPI (Swagger) по адресу `/docs`.
- Каждый запрос пишет в терминал строку лога: метод, путь, статус, время ответа в миллисекундах. Например: `GET /api/books/7 404 1.3ms`.

Чекер проверяет, что `/docs` отдаёт HTML (редирект тоже подходит), находит OpenAPI-спецификацию и ищет в ней пути ресурса. Спецификацию он ищет сам на `/openapi.json`, `/swagger.json`, `/v3/api-docs`, `/swagger/doc.json`, `/api-docs`. Если она лежит по другому адресу, путь пишется в `contract.json` в поле `openapi_path`. Для лога чекер отправляет запрос с меткой `backcheck-log-XXXX`, эту строку нужно показать в терминале сервиса.

### Что такое OpenAPI и Swagger

**OpenAPI** это стандартное описание HTTP API в файле YAML или JSON: какие есть пути и методы, что отправлять в теле, какие ответы и коды приходят. **Swagger UI** это страница, которая показывает это описание и позволяет отправлять запросы прямо из браузера (кнопка Try it out). Описание либо строится фреймворком из кода, либо пишется файлом `openapi.yaml` руками.

Пример описания одного эндпоинта:

```yaml
openapi: 3.0.3
info:
  title: Библиотека
  version: 1.0.0
paths:
  /api/books/{id}:
    get:
      summary: Получить книгу
      parameters:
        - name: id
          in: path
          required: true
          schema: { type: integer }
      responses:
        "200": { description: Книга }
        "404": { description: Книги нет }
```

### Как подключить на своём языке

| Стек | Как подключить | Как сделать страницу на `/docs` |
| --- | --- | --- |
| Python, FastAPI | Ничего не нужно, строится по маршрутам и моделям | Уже на `/docs`, описание на `/openapi.json` |
| Python, Django REST Framework | `drf-spectacular` | В `urls.py`: `SpectacularAPIView` на `openapi.json` и `SpectacularSwaggerView` на `docs/` |
| Node.js, Express | `swagger-ui-express` и свой `openapi.yaml` (пакет `yaml` для чтения файла) | `app.use("/docs", swaggerUi.serve, swaggerUi.setup(spec))`, описание отдать на `/openapi.json` |
| Node.js, NestJS | `@nestjs/swagger` | `SwaggerModule.setup("docs", app, document)` |
| Java, Spring Boot | `springdoc-openapi-starter-webmvc-ui` | В `application.properties`: `springdoc.swagger-ui.path=/docs`; описание на `/v3/api-docs` |
| C#, ASP.NET Core | Swashbuckle (`AddSwaggerGen`, `UseSwagger`, `UseSwaggerUI`) | `RoutePrefix = "docs"` в настройках `UseSwaggerUI`; описание на `/swagger/v1/swagger.json`, путь указать в `openapi_path` |
| Go | `swaggo/swag` (комментарии над обработчиками, `swag init`) и `http-swagger`, либо свой `openapi.yaml` и HTML-страница со Swagger UI с CDN | Маршрут `GET /docs/` в своём роутере; описание swaggo отдаёт на `/swagger/doc.json` |
| PHP, Laravel | L5-Swagger или Scramble | Путь страницы задаётся в конфиге пакета |

Если описание пишется руками, оно должно совпадать с настоящими эндпоинтами: на приёме могут попросить отправить запрос прямо со страницы `/docs`.

## Как заполнить contract.json

Чекер один на всех, а темы разные. В `contract.json` описывается своя тема: как называется ресурс, какое тело правильное, какое поле обязательное.

```json
{
  "base_url": "http://localhost:8000",
  "resource": "books",
  "id_field": "id",
  "missing_id": "999999",
  "list_key": null,
  "create": { "title": "Мастер и Маргарита", "author": "Булгаков", "copies": 3 },
  "update": { "copies": 5 },
  "required_field": "title",
  "errors": [
    { "body": { "title": "" }, "field": "title" },
    { "body": { "author": "" }, "field": "author" },
    { "body": { "copies": -1 }, "field": "copies" }
  ],
  "filter": { "param": "author", "value": "Булгаков", "other": "Пушкин" }
}
```

| Поле | Что писать |
| --- | --- |
| `base_url` | Адрес сервиса, обычно `http://localhost:8000` |
| `resource` | Имя ресурса из пути: для `/api/books` это `books` |
| `id_field` | Как называется id в ответе, обычно `id` |
| `missing_id` | id, которого точно нет. Для UUID подойдёт любой несуществующий UUID |
| `list_key` | `null`, если `GET /api/books` отдаёт массив. `"items"`, если отдаёт `{"items": [...]}` |
| `create` | Правильное тело для создания |
| `update` | Какие поля поменять при обновлении |
| `required_field` | Обязательное поле, чекер отправит его пустым и без него |
| `errors` | Только для ★: неправильные значения полей и какое поле должно попасть в ошибку |
| `filter` | Только для ★: параметр фильтра, подходящее и неподходящее значение |
| `openapi_path` | Только для ★★ и только если OpenAPI-описание лежит не на стандартном адресе, например `"/swagger/v1/swagger.json"` |

## Как проверить себя

```bash
# сервис запущен в соседнем терминале
backcheck --lab 1              # база
backcheck --lab 1 --stars      # база и звёздочки
backcheck --lab 1 -v --only 6  # подробно: какие запросы уходят и что приходит в пункте 6
```

Если чекер пишет «сервис не отвечает», значит сервис не запущен, запущен на другом порту или `base_url` в `contract.json` указан неверно.

## Сценарий демонстрации на паре

health, создать, список, получить по id, обновить, удалить, снова получить (404), создать со сломанным JSON (400). Всё это есть в `requests.http`.

## С чего начать на своём языке

Минимальный сервис на Python (FastAPI). Запуск: `pip install fastapi uvicorn`, затем `uvicorn main:app --reload --port 8000`.

```python
from fastapi import FastAPI, HTTPException, Response
from pydantic import BaseModel, Field

app = FastAPI()
books: dict[int, dict] = {}
next_id = 1


class BookIn(BaseModel):
    title: str = Field(min_length=1)
    author: str = Field(min_length=1)
    copies: int = Field(default=1, ge=0)


@app.get("/health")
def health():
    return {"status": "ok"}


@app.post("/api/books", status_code=201)
def create_book(data: BookIn):
    global next_id
    book = {"id": next_id, **data.model_dump()}
    books[next_id] = book
    next_id += 1
    return book


@app.get("/api/books/{book_id}")
def get_book(book_id: int):
    if book_id not in books:
        raise HTTPException(404, "книга не найдена")
    return books[book_id]


@app.delete("/api/books/{book_id}", status_code=204)
def delete_book(book_id: int):
    if books.pop(book_id, None) is None:
        raise HTTPException(404, "книга не найдена")
    return Response(status_code=204)
```

Список и PATCH дописываются так же. FastAPI на неправильное тело отвечает 422; в базе это засчитывается, а для ★ нужен свой обработчик `RequestValidationError`, который отвечает 400.

| Язык | Минимум | Подсказка |
| --- | --- | --- |
| Python | FastAPI и uvicorn | `@app.post("/api/items", status_code=201)`; 422 превращается в 400 через `@app.exception_handler(RequestValidationError)` |
| JS/TS | Express | `app.use(express.json())`; ошибку разбора JSON ловит middleware с четырьмя аргументами `(err, req, res, next)` |
| Go | `net/http` (Go 1.22+) | `mux.HandleFunc("GET /api/items/{id}", ...)`, id достаётся через `r.PathValue("id")` |
| Java | Spring Boot (Web) | `@RestController`, `ResponseEntity.status(201)`; `@ControllerAdvice` для 400 |
| C# | ASP.NET Core Minimal API | `app.MapPost(...)`, `Results.Created(...)`, `Results.NoContent()` |
| PHP | Slim или Laravel | `$response->withStatus(201)` |

## Частые ошибки

- `DELETE` отвечает `200` и `{"ok": true}` вместо `204` без тела.
- Неизвестный id даёт 500 (например, `KeyError` в Python) или 200 с `null` вместо 404.
- Сервер падает на сломанном JSON. В Express не хватает обработчика ошибок, в Go не проверена ошибка от `Decode`.
- id считается как `len(items) + 1`. После удаления id начинают повторяться.
- `requests.http` ни разу не запускали, и пути в нём не совпадают с настоящими.

## Вопросы, которые могут задать на приёме

1. Что происходит по шагам от ввода адреса в браузере до показа страницы: DNS, TCP, TLS, HTTP?
2. Чем отличаются IP-адрес и порт? Что значит, что сервис слушает `0.0.0.0:8000`, а не `127.0.0.1:8000`?
3. Из каких частей состоят HTTP-запрос и HTTP-ответ? Где в них метод, путь, заголовки, тело, статус?
4. Чем PUT отличается от PATCH, а POST от PUT?
5. Что такое идемпотентность? Какие методы HTTP идемпотентны и почему POST нет?
6. Почему создание отвечает 201, а удаление 204? Что будет в теле этих ответов?
7. Чем 400, 404 и 500 отличаются друг от друга? На чьей стороне ошибка в каждом случае?
8. Зачем заголовок `Content-Type` и что будет, если сервер отдаст JSON с `text/html`?
9. Какие правила REST соблюдены в проекте? Почему в пути существительное во множественном числе, а не `/createItem`?
10. Что такое stateless в REST и почему сервер не должен помнить прошлые запросы клиента?
11. Что будет с данными после перезапуска сервиса и почему? Как устроено хранение сейчас?
12. Где в коде проверяется обязательное поле и что произойдёт без этой проверки?
13. Что такое OpenAPI и чем он отличается от Swagger UI? Зачем такая документация, если есть README?
14. Что такое фреймворк и менеджер пакетов? Какие зависимости у проекта и где они записаны?
15. Чем HTTP отличается от HTTPS? Что именно шифруется и что видно постороннему в сети?

Возможные живые правки: добавить поле в сущность; поменять код ответа; добавить `GET /api/<items>/count`; ограничить длину строки.
