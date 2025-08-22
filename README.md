# Order Notification Service

Сервис для работы с заказами: сохранение в PostgreSQL, кеширование в памяти и получение по UUID.  
Поднимается в Docker вместе с Postgres, Kafka и Zookeeper.

---

## 🚀 Стек технологий
- **Go**
- **PostgreSQL**
- **Kafka + Zookeeper**
- **Docker, Docker Compose**

---

## ⚙️ Запуск проекта

### 1. Сборка и запуск

`docker-compose up --build -d`

### 2. Проверка контейнеров

`docker ps`


Ожидаемые контейнеры:

app — Go-сервис

postgres — база данных

zookeeper — сервис ZooKeeper

kafka — брокер Kafka

### 3. Логи приложения

`docker logs -f app`

📡 API
Создать заказ
`curl -X POST http://localhost:8080/new_order \
  -H "Content-Type: application/json" \
  -d @example_order.json`

Получить заказ по UUID
curl -X GET "http://localhost:8080/order/<uuid>"

📂 Структура проекта

/cmd            — точка входа

/internal

/handler      — http-обработчики

/service      — бизнес-логика

/repo         — работа с БД и кэшем

/entity       — модели данных

/pkg/postgres   — подключение к БД

/config         — конфиги

/web            — статические файлы


🛠️ Makefile

В проекте есть Makefile для удобной работы:

make build      # сборка контейнеров

make up         # запуск docker-compose

make down       # остановка и удаление

make logs       # логи приложения

make restart    # рестарт

make ps         # список контейнеров

make curl-order # тестовый POST заказа

### 📦 Тестовый заказ

Файл example_order.json в корне проекта:

{
"order_uid": "550e8400-e29b-41d4-a716-446655440000",
"track_number": "WBILMTESTTRACK",
"entry": "WBIL",
"delivery": {
"name": "Test Testov",
"phone": "+9720000000",
"zip": "2639809",
"city": "Kiryat Mozkin",
"address": "Ploshad Mira 15",
"region": "Kraiot",
"email": "test@gmail.com"
},
"payment": {
"transaction": "550e8400-e29b-41d4-a716-446655440000",
"request_id": "",
"currency": "USD",
"provider": "wbpay",
"amount": 1817,
"payment_dt": 1637907727,
"bank": "alpha",
"delivery_cost": 1500,
"goods_total": 317,
"custom_fee": 0
},
"items": [
{
"chrt_id": 9934930,
"track_number": "WBILMTESTTRACK",
"price": 453,
"rid": "ab4219087a764ae0b123456789",
"name": "Mascaras",
"sale": 30,
"size": "0",
"total_price": 317,
"nm_id": 2389212,
"brand": "Vivienne Sabo",
"status": 202
}
],
"locale": "en",
"internal_signature": "",
"customer_id": "test",
"delivery_service": "meest",
"shardkey": "9",
"sm_id": 99,
"date_created": "2021-11-26T06:22:19Z",
"oof_shard": "1"
}