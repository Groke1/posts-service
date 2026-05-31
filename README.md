# Posts-service

## Описание

Система для добавления и чтения постов и комментариев с использованием `GraphQL`.

Реализовано `GraphQL Subscriptions`: есть возможность получения комментариев в асинхронном формате.

### Характеристики системы постов:
1. Можно просмотреть список постов.
2. Можно просмотреть пост и комментарии под ним.
3. Пользователь, написавший пост, может запретить оставление комментариев к своему посту.

### Характеристики системы комментариев к постам:
1. Комментарии организованы иерархически, позволяя вложенность без ограничений.
2. Длина текста комментария ограничена до 2000 символов.
3. Система пагинации для получения списка комментариев.

Для запуска надо сгенерировать файлы `make generate` и запустить `docker-compose`.

## Примеры запросов

### Создание поста
``` graphql
mutation {
  createPost(
    input: {
      title: "Название поста"
      text: "Текст поста"
      authorID: 1
      commentsEnabled: true
    }
  ) {
    id
    title
    text
    authorID
    commentsEnabled
    createdAt
  }
}
```

### Получение данных поста

```graphql
query {
  post(id: 1) {
    id,
    createdAt,
    text,
    title,
    authorID,
    commentsEnabled
  }
}
```

### Получение списка постов
Получаем последние 10 постов для автора:
```graphql
query {
  posts(
    authorID: 1
    limit: 10
    offset: 0
  ) {
    id
    text
  }
}
```
Если не указать `authorID`, получим 10 последних постов всех авторов.

### Изменение возможности комментирования поста
```graphql
mutation {
  setPostCommentsEnabled(
    postID: 1
    authorID: 1
    enabled: false
  ) {
    id
    updatedAt
  }
}
```
Если пост с `postID` написал не `authorID`, вернется ошибка.

### Добавление комментария к посту
```graphql
mutation {
  createComment(input: {
    authorID: 5
    postID: 1
    text: "Комментарий к посту"
  }) {
    id
    postID
    createdAt
    text
  }
}
```
### Добавление ответа к комментарию
```graphql
mutation {
  createComment(input: {
    authorID: 5
    postID: 1
    parentID: 1
    text: "Ответ к комментарию с id 1"
  }) {
    id
    postID
    createdAt
    text
  }
}
```
Если поста с `postID` не существует или среди всех комментариев к посту нет комментария с `parentID`, вернется ошибка.

###  Получение комментария
```graphql
query {
  comment(id: 1) {
    text
  }
}
```

### Получение списка комментариев
Получаем первые 10 комментариев к посту:
```graphql
query {
  comments(
    postID: 1
    limit: 10
    offset: 0
  ) {
    id
    text
  }
}
```

### Получение списка ответов на комментарий
Получаем первые 10 ответов на комментарий с `id = 1` 
```graphql
query {
  comments(
    postID: 1
    parentID: 1
    limit: 10
    offset: 0
  ) {
    id
    text
  }
}
```
Если комментарий с `parentID` написан не под постом `postID`, вернется ошибка.

### Подписка на комментарии
```graphql
subscription {
  commentAdded(postID: 1) {
    id
    text
    authorID
  }
}
```

## Скрипты

Запустить генерацию кода:
```bash
make generate
```

Запустить unit-тесты:
```bash
make test
```

Запустить сервер с `in-memory` хранилищем:
```bash
make run-inmemory
```

Запустить сервер с `postgres` хранилищем:
```bash
make run-postgres
```

Остановить сервер:
```bash
make stop
```
