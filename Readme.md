


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
