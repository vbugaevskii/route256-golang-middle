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

Дополнительное задание на алмазик:
- [x] Добавить HTTP-gateway и proto-валидацию.

*Дедлайн: 3 июня, 23:59 (сдача) / 6 июня, 23:59 (проверка)*

# Домашнее задание №3

1. [x] Для каждого сервиса(где необходимо что-то сохранять/брать) поднять отдельную БД в __docker-compose__.
2. [x] Сделать миграции в каждом сервисе (достаточно папки миграций и скрипта).
3. [x] Создать необходимые таблицы.
4. [x] Реализовать логику репозитория в каждом сервисе.
5. [x] В качестве query builder-а можно использовать любую библиотеку (согласовать индивидуально с тьютором). Рекомендуется https://github.com/Masterminds/squirrel.
6. [x] Драйвер для работы с postgresql: только __pgx__ (pool).
7. [x] В одном из сервисов сделать транзакционность запросов (как на воркшопе).

Задание на алмазик:
1. [x] Для каждой БД полнять свой балансировщик (pgbouncer или odyssey, можно и то и то). Сервисы ходят не на прямую в БД, а через балансировщик

*Дедлайн: 10 июня, 23:59 (сдача) / 13 июня, 23:59 (проверка)*

# Домашнее задание №4

1. [x] Уменьшить время ответа Checkout.listCart при помощи worker pool. Запрашивать не более 5 SKU одновременно. Worker pool нужно написать самостоятельно. Обязательное требование - читаемость кода и покрытие комментариями.
2. [x] При общении с Product Service необходимо использовать лимит 10 RPS на клиентской стороне. Допускается использование библиотечных рейт-лимитеров. В случае собственного читаемый код и комментарии обязательны.
3. [x] Во всех слоях сервиса необходимо использовать контекст для возможности отмены вызова.

Задания на алмазик:

1. [ ] Синхронизировать рейт-лимитер при помощи БД.
2. [x] Аннулировать заказы старше 10 минут в фоне. Позаботиться о том, чтобы реплики сервиса не штурмовали БД все вместе.

*Дедлайн: 17 июня, 23:59 (сдача) / 20 июня, 23:59 (проверка)*

# Домашнее задание №5

- [x] Необходимо обеспечить полное покрытие бизнес-логики ручек ListCart или Purchase модульными тестами (go test -coverprofile).
- [x] Если вдруг ваши слои до сих пор не изолированны друг от друга через интерфейсы, необходимо это сделать.
- [x] В качестве генератора моков можете использовать, что душе угодно: **mockery**, minimoc, gomock, ...

Задание на алмазик:
- [x] добавить интеграционные тесты для проверки слоя взаимодействия с базой данных.

*Дедлайн: 24 июня 23:59 (сдача) / 27 июня, 23:59 (проверка)*

# Домашнее задание №6

1. [x] LOMS пишет в Кафку изменения статусов заказов
2. [x] Сервис нотификаций должен их вычитывать и отправлять нотификации об изменениях статуса заказа (писать в телегу)
3. [x] Нотификация должна быть доставлена гарантированно и ровно один раз
4. [x] Нотификации должны доставляться в правильном порядке

Задания на алмазик:

[ ] Весь новый функционал покрыт юнит тестами плюс сам код написан таким образом, что конфигурацию для нотификаций можно задавать в конфиге и прокидывать в основную логику

*Дедлайн: 01.07.2023 (сдача) / 04.07.2023 (проверка)*

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

grpcurl -plaintext -d '{"sku": 773587830}' localhost:8081 loms.Loms/Stocks # OK

grpcurl -plaintext -d '{"user": 1, "sku": 773587830, "count": 5}' localhost:8080 checkout.Checkout/AddToCart # OK
grpcurl -plaintext -d '{"user": 1}' localhost:8080 checkout.Checkout/ListCart # OK
grpcurl -plaintext -d '{"user": 1, "sku": 773587830, "count": 5}' localhost:8080 checkout.Checkout/AddToCart # ERROR

grpcurl -plaintext -d '{"user": 1, "sku": 773587830, "count": 1}' localhost:8080 checkout.Checkout/DeleteFromCart # OK
grpcurl -plaintext -d '{"user": 1}' localhost:8080 checkout.Checkout/ListCart # OK

grpcurl -plaintext -d '{"user": 1, "sku": 773596051, "count": 3}' localhost:8080 checkout.Checkout/AddToCart # OK
grpcurl -plaintext -d '{"user": 1}' localhost:8080 checkout.Checkout/ListCart # OK

grpcurl -plaintext -d '{"user": 1}' localhost:8080 checkout.Checkout/Purchase # OK -> orderId=1
grpcurl -plaintext -d '{"user": 1}' localhost:8080 checkout.Checkout/ListCart # OK
grpcurl -plaintext -d '{"orderID": 1}' localhost:8081 loms.Loms/ListOrder # OK

grpcurl -plaintext -d '{"user": 2, "sku": 773587830, "count": 5}' localhost:8080 checkout.Checkout/AddToCart # OK
grpcurl -plaintext -d '{"user": 2}' localhost:8080 checkout.Checkout/ListCart # OK

grpcurl -plaintext -d '{"orderID": 1}' localhost:8081 loms.Loms/OrderPayed # OK
grpcurl -plaintext -d '{"orderID": 42}' localhost:8081 loms.Loms/ListOrder # ERROR

grpcurl -plaintext -d '{"user": 2}' localhost:8080 checkout.Checkout/Purchase # ERROR -> orderId=2
grpcurl -plaintext -d '{"orderID": 2}' localhost:8081 loms.Loms/ListOrder # OK

grpcurl -plaintext -d '{"sku": 773587830}' localhost:8081 loms.Loms/Stocks # OK
```

Протестировать rps можно так:
```bash
#!/usr/bin/env sh

mkdir -p answers

for i in {1..5}; do
    grpcurl -plaintext -d '{"user": 1}' localhost:8080 checkout.Checkout/ListCart \
        >"answers/$i.ans" 2>&1 &
done

wait
```

Клиент для просмотра топиков Kafka – [Offset Explorer](https://www.kafkatool.com/download.html)

Настройка Telegram бота:
1. получить [токен](https://core.telegram.org/bots/tutorial#getting-ready) для бота;
2. разрешить общение с ботом (нужно найти в Телеграме созданного бота и нажать кнопку `Start`);
3. узнать свой идентификатор через другого [бота](https://t.me/getmyid_bot);
4. настроить конфиг для сервиса `notifications`.

Смотрелка Jaeger: http://localhost:16686
Проверить, что шлются метрики:
- http://localhost:8070/metrics
- http://localhost:9090/targets