веб-сервер с двумя обработчиками.
первый обработчик пополняет баланс указанного пользователя на указанную сумму. то есть получает на вход 2 значения.

второй переводит указанную сумму со счета первого пользователя на счет другого. в минус уходить нельзя. принимает на вход 3 значения:
1. ид юзера (или баланса), с которого списание
2. ид юзера, кому на счет поступают средства
3. сумма средств для перевода

в проекте нужна миграция для создания таблицы. если получится по времени, то можно написать тесты.
предусмотреть останов веб-сервера без потери обрабатываемых запросов

субд postgre

go build .\cmd\main.go

Можно запустить через докер

`docker build . --tag=finstar_api`

`docker run -e PG_CONN="postgres://postgres:postgres@localhost/postgres?sslmode=disable" finstar_api`


Пополнить баланс

`curl --location --request POST 'http://localhost:8000/api/v1/user/deposit' \
--header 'Content-Type: application/json' \
--data-raw '{
"to" : 2,
"total" : 0.01
}'`

Сделать перевод

`curl --location --request POST 'http://localhost:8000/api/v1/user/transfer' \
--header 'Content-Type: application/json' \
--data-raw '{
"from": 13,
"to": 12,
"total": 55.01
}'`
