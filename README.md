# WB Tech: level # 0 (Golang)
Это пример простого приложения для управления заказами. Приложение разбито на модули для удобства чтения и поддержки.
## Конфигурация Docker Compose

В файле `docker-compose.yml` содержатся следующие сервисы:
1. **postgres**: Этот сервис использует образ Docker `postgres:15.3`. Он настраивает контейнер базы данных PostgreSQL и связывает каталог `./volumes/pgdata` на хост-машине с папкой `/var/lib/postgresql/data` внутри контейнера. Данные базы данных будут храниться в указанном каталоге на хост-машине, обеспечивая сохранность данных. Сервис также использует файл окружения `.env` для загрузки переменных окружения, если он присутствует. Сервис доступен на порту `5555` на хост-машине, который отображается на порт `5432` внутри контейнера.
2. **nats**: Этот сервис использует образ Docker `nats-streaming:latest`. Он настраивает контейнер сервера NATS Streaming и открывает порты `4222` и `8222` на хост-машине, которые отображаются на соответствующие порты внутри контейнера.

## script/producer.go
Скрипт посылает json в nats-streaming канал, выполняются следующие шаги:
1. Устанавливается соединение с NATS Streaming сервером, работающим на адресе 127.0.0.1:4222.
2. Создается пример заказа (order) с использованием функции getExampleOrder().Эта функция используется для загрузки примера заказа из файла example.json. Она открывает файл, декодирует его содержимое в структуру Order и возвращает эту структуру.
3. Затем запускается бесконечный цикл, в котором заказу присваивается новый UID, и этот заказ отправляется через NATS с использованием сериализации в JSON. В консоли выводится информация о том, какой заказ был отправлен. После каждой отправки происходит пауза в 2 секунды.


## Основное приложение. Описание модулей

### 1. `app`

Модуль `app` содержит код, связанный с инициализацией и конфигурацией приложения.

- `app.go`: Описывает структуру приложения и его методы.

### 2. `handler`

Модуль `handler` отвечает за обработку HTTP-запросов.

- `handler.go`: Содержит обработчики запросов для работы с заказами.

### 3. `infrastructure`

Модуль `infrastructure` включает код для работы с внешними ресурсами, такими как базы данных и очереди сообщений.

- `database.go`: Инициализация и настройка соединения с базой данных PostgreSQL.

- `nats.go`: Инициализация соединения с NATS для обработки сообщений.

- `cache.go`: Кэш для хранения заказов в памяти.

### 4. `model`

Модуль `model` содержит структуры данных, представляющих заказы и их компоненты.

- `model.go`: Описывает структуры `Order`, `Delivery`, `Payment`, `Item`.

### 5. `repository`

Модуль `repository` отвечает за взаимодействие с базой данных и кэшем заказов.

- `repository.go`: Реализует методы для загрузки данных в кэш и сохранения заказов в базу данных.


## Запуск приложения

Для запуска приложения необходимо выполнить следующие шаги:

1. Убедитесь, что у вас установлены Go, Docker compose.

2. Создайте `.env` в корневой папке с параметрами подключения к бд, пример можете посмотреть в `.env.example`

3. Запустите приложение с помощью команды:
```bash
docker compose up
go run main.go
```
для запуска продюсера:
```bash
go run script/producer.go
```