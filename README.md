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

1. Для работы требуется unsealed Vault хранилище с secret engine, в котором хранятся логины и пароли к терминалам. 
2. Если Vault запущен, перейти к шагу 6.
3. Если Vault не запущен, то добавить терминалы в deploy/vault/terminals.json
4. Выполнить `make run.vault`
5. Скопировать адрес и root key из консоли.
6. Добавить адрес в VAULT_ADDR и root key в VAULT_TOKEN.
7. 