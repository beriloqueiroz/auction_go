# Auction Bid

## expiração do Auction

- A todo momento uma função é executada verificando quais actions estão expirados e os mesmos são marcados como "Completed"

## variáveis de ambiente

- no arquivo .env há uma variável AUCTION_DURATION_HOUR que pode ser definida em unidade de time.Duration , exemplo: 5m -> 5 minutos, 3h -> 3 horas. Ao incluir um action o mesmo fica com status Active = 0 inicialmente e automaticamente ao passar o tempo definido pela  AUCTION_DURATION_HOUR o mesmo ficará com o valor Completed = 1 (pode se verificar esse cenário pela listagem de actions)

## para rodar

- aplicação na porta 8080, e mongodb

  ```bash
    docker compose up -d
  ```

## Requests

- criar auction
  
  ```bash
    curl  -X POST \
      '<http://localhost:8080/auction>' \
      --header 'Accept: */*' \
      --header 'User-Agent: Thunder Client (<https://www.thunderclient.com>)' \
      --header 'Content-Type: application/json' \
      --data-raw '{
      "product_name":"teste 1",
      "category": "category test",
      "description": "description test",
      "condition": 1
    }'
  ```

- listar auctions
  
  ```bash
    curl  -X GET \
    'http://localhost:8080/auction?status=0' \
    --header 'Accept: */*' \
    --header 'User-Agent: Thunder Client (https://www.thunderclient.com)'
  ```

- criar bid (o user_id pode ser qualquer uuid, o auction_id deve ser preenchido com o resultado do item anterior)

  ```bash
  curl  -X POST \
    'http://localhost:8080/bid' \
    --header 'Accept: */*' \
    --header 'User-Agent: Thunder Client (https://www.thunderclient.com)' \
    --header 'Content-Type: application/json' \
    --data-raw '{
    "user_id": "151ae271-d04f-41d0-a3c1-956d78380d8b",
    "auction_id": "f3ab2e11-efce-4a79-aa1b-b908d2d38437",
    "amount": 152
  }'
  ```
