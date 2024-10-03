# Тестовое задание для компании MEDODS на позицию Junior Backend Developer

## Описание
Проект реализует два REST маршрута:

- Первый маршрут выдает пару Access и Refresh токенов для пользователя с идентификатором (GUID), указанным в параметре запроса.
- Второй маршрут выполняет операцию Refresh для пары Access и Refresh токенов.

## Установка
1. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/Pechenkoff/medods.git
   ```
2. Перейдите в директорию проекта:
   ```bash
   cd medods
   ```
3. Запустите Docker Compose, который создаст нужные образы и подгрузит необходимые зависимости:
   ```bash
   docker compose up
   ```
## Использование
1. **Создание токенов**
* **POST** `/access`
* **Описание:** Генерирует пару Access и Resfresh токенов по GUID и email.
* **Пример запроса:**
```bash
curl -X POST http://localhost:8080/access -H "Content-Type: application/json" -d '{"user_id": "123e4567-e89b-12d3-a456-426614174000", "email": "user@example.com"}'
```
2. **Обновление токенов**
* **POST** `/refresh`
* **Описание:** Получает пару Access и Refresh токенов, расшифровыет Access токен, проверяет Refresh токен и возвращает новую пару Access и Refresh токенов.

## Тестирование
Тесты будут позже.

## Конфигурация
Конфигурационные, docker compose и dockerfile файлы находяться в репозитории, можете изменить под себя.

## Документация API
Swagger документация доступна после запуска сервиса по адресу:
http://localhost:8080/swagger/index.html.

## Дополнительные улучшения и допущения
* Рассмотреть возможность добавления мехинизмов мониторинга для отслеживания производительности API
* Было допущено, что есть отдельный сервис нотификаций, который общается с другими сервисами через брокер сообщений Kafka, который и отправит нотификацию пользователю о смене IP
* Было допущено, что пользователь отправляет сам свою почту, было бы лучше если у нас была бы отдельная база данных с данными о пользователе, где через GUID мы могли бы получать его email
* Было допущено, что без реализации отдельной базы данных используется метод POST для получения GUID и email'а, в случае если база есть можно в хендлерах переделать логику получения GUID, например через `/access/{user_guid}` и в слоях сервисов и репозитория реализовать логику получения email'а из общей базы данных о пользователе
* Было допущено, что генерируется новый refresh токен для пользователя, но приэтом можно реализовать логику "протухания" этого токена и отчищать базу данных при истечении срока жизни токена

## Контакты
Для обратной связи: kryukovskoi.vv@phystech.edu