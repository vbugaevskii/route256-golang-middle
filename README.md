# Домашнее задание

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