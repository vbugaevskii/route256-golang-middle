# Домашнее задание 1

- [x] Создать скелеты трёх сервисов по описанию АПИ из файла contracts.md
- [x] Структуру проекта сделать с учетом разбиения на слои, бизнес-логику писать отвязанной от реализаций клиентов и хендлеров
- [x] Все хендлеры отвечают просто заглушками
- [x] Сделать удобный враппер для сервера по тому принципу, по которому делали на воркшопе
- [x] Придумать самостоятельно удобный враппер для клиента
- [x] Все межсервисные вызовы выполняются. Если хендлер по описанию из contracts.md должен ходить в другой сервис, он должен у вас это успешно делать в коде.
- [x] Общение сервисов по http-json-rpc
- [x] должны успешно проходить make precommit и make run-all в корневой папке
- [x] Наладить общение с product-service (в хендлере Checkout.listCart). Токен для общения с product-service получить, написав в личку @badger_za

*Дедлайн: 27 мая, 23:59 (сдача) / 30 мая, 23:59 (проверка)*

# Домашнее задание 2

Перевести всё взаимодействие c сервисами на протокол gRPC.

Для этого:

- [x] Использовать разделение на слои, созданное ранее, заменив слой HTTP на GRPC.
- [x] Взаимодействие по HTTP полностью удалить и оставить только gRPC.
- [x] В каждом проекте нужно добавить в Makefile команды для генерации кода из proto файла и установки нужных зависимостей.

Дополнительное задание на алмазик: добавить HTTP-gateway и proto-валидацию.

*Дедлайн: 3 июня, 23:59 (сдача) / 6 июня, 23:59 (проверка)*

# Домашнее задание №3

1. [ ] Для каждого сервиса(где необходимо что-то сохранять/брать) поднять отдельную БД в __docker-compose__.
2. [ ] Сделать миграции в каждом сервисе (достаточно папки миграций и скрипта).
3. [ ] Создать необходимые таблицы.
4. [ ] Реализовать логику репозитория в каждом сервисе.
5. [ ] В качестве query builder-а можно использовать любую библиотеку (согласовать индивидуально с тьютором). Рекомендуется https://github.com/Masterminds/squirrel.
6. [ ] Драйвер для работы с postgresql: только __pgx__ (pool).
7. [ ] В одном из сервисов сделать транзакционность запросов (как на воркшопе).

Задание на алмазик:
1. Для каждой БД полнять свой балансировщик (pgbouncer или odyssey, можно и то и то). Сервисы ходят не на прямую в БД, а через балансировщик

*Дедлайн: 10 июня, 23:59 (сдача) / 13 июня, 23:59 (проверка)*

# Полезные команды
```bash
make build  # собрать бинарь

# послать GET запрос сервису LOMS
curl -i localhost:8081/createOrder -d '{"user": 1, "items": [{"sku": 12, "count": 23}]}'
curl -i localhost:8081/listOrder -d '{"orderID": 42}'
curl -i localhost:8081/orderPayed -d '{"orderID": 42}'
curl -i localhost:8081/cancelOrder -d '{"orderID": 42}'
curl -i localhost:8081/stocks -d '{"sku": 12}'

# послать GET запрос сервису Checkout
curl -i localhost:8080/addToCart -d '{"user": 1, "sku": 12, "count": 23}'
curl -i localhost:8080/deleteFromCart -d '{"user": 1, "sku": 12, "count": 23}'
curl -i localhost:8080/listCart -d '{"user": 1}'
curl -i localhost:8080/purchase -d '{"user": 1}'
```

Клиент для grpc – [`grpcurl`](https://github.com/fullstorydev/grpcurl).

```bash
# посмотреть, какие методы есть у сервиса
grpcurl -plaintext route256.pavl.uk:8082 list

# послать GRPC запрос сервису LOMS
grpcurl -plaintext -d '{"user": 1, "items": [{"sku": 12, "count": 23}]}' localhost:8081 loms.Loms/CreateOrder
grpcurl -plaintext -d '{"orderID": 42}' localhost:8081 loms.Loms/ListOrder
grpcurl -plaintext -d '{"orderID": 42}' localhost:8081 loms.Loms/OrderPayed
grpcurl -plaintext -d '{"orderID": 42}' localhost:8081 loms.Loms/CancelOrder
grpcurl -plaintext -d '{"sku": 12}' localhost:8081 loms.Loms/Stocks

# послать GRPC запрос сервису Checkout
grpcurl -plaintext -d '{"user": 1, "sku": 12, "count": 23}' localhost:8080 checkout.Checkout/AddToCart
grpcurl -plaintext -d '{"user": 1, "sku": 12, "count": 23}' localhost:8080 checkout.Checkout/DeleteFromCart
grpcurl -plaintext -d '{"user": 1}' localhost:8080 checkout.Checkout/ListCart
grpcurl -plaintext -d '{"user": 1}' localhost:8080 checkout.Checkout/Purchase
```

Клиент для накатки миграций - [`goose`](https://github.com/pressly/goose)

```bash
# создать миграцию init
goose create init sql
```

Ручное тестирование:

```bash
set -x

grpcurl -plaintext -d '{"user": 1, "sku": 773587830, "count": 5}' localhost:8080 checkout.Checkout/AddToCart # OK
grpcurl -plaintext -d '{"user": 1}' localhost:8080 checkout.Checkout/ListCart # OK
grpcurl -plaintext -d '{"user": 1, "sku": 773587830, "count": 5}' localhost:8080 checkout.Checkout/AddToCart # ERROR

grpcurl -plaintext -d '{"user": 1, "sku": 773587830, "count": 1}' localhost:8080 checkout.Checkout/DeleteFromCart # OK
grpcurl -plaintext -d '{"user": 1}' localhost:8080 checkout.Checkout/ListCart # OK

grpcurl -plaintext -d '{"user": 1, "sku": 773596051, "count": 3}' localhost:8080 checkout.Checkout/AddToCart # OK
grpcurl -plaintext -d '{"user": 1}' localhost:8080 checkout.Checkout/ListCart # OK

grpcurl -plaintext -d '{"user": 1}' localhost:8080 checkout.Checkout/Purchase # OK -> orderId=2
grpcurl -plaintext -d '{"user": 1}' localhost:8080 checkout.Checkout/ListCart # ERROR
grpcurl -plaintext -d '{"orderID": 2}' localhost:8081 loms.Loms/ListOrder # OK

grpcurl -plaintext -d '{"user": 2, "sku": 773587830, "count": 5}' localhost:8080 checkout.Checkout/AddToCart # OK
grpcurl -plaintext -d '{"user": 2}' localhost:8080 checkout.Checkout/ListCart # OK

grpcurl -plaintext -d '{"orderID": 2}' localhost:8081 loms.Loms/OrderPayed # OK
grpcurl -plaintext -d '{"orderID": 42}' localhost:8081 loms.Loms/ListOrder # OK

grpcurl -plaintext -d '{"user": 2}' localhost:8080 checkout.Checkout/Purchase # ERROR -> orderId=3
grpcurl -plaintext -d '{"orderID": 3}' localhost:8081 loms.Loms/ListOrder # OK
```