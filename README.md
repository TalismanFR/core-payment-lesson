core payment service
===

Тимур application layer. 
- Работа с бд: BeginTransaction
- dto входа и выхода
- Когда нужно искать Pay в бд
- 

Денис domain payment. Проведение платежа используя sdk bepaid
на основе данных с сущности pay. pay.


### Как запускать
1. Создать `.env` с содержимым:
```dotenv
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_USER=payservice
POSTGRES_PASSWORD=payservice

VAULT_ADDR=http://vault:8200
VAULT_TOKEN=hvs.9X0GyCKE4qBhODNOcDpev8eF
```
2. Выполнить `make vault` 
3. Скопировать значение `VAULT_TOKEN` в `.env`. 
Если нужно, то скопировать любой uuid как `TerminalId` в tests/grpc-client/main.go
4. Выполнить `make run`
5. Если нужно проверить соединение, то выполнить `tests/grpc-client/main.go`
6. Выполнить `make stop`, чтобы остановить все контейнеры
7. Если нужно ,то выполнить `make clear`
