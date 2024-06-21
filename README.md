# API/CRUD - GOLANG (Estudos)

API desenvolvida em Go para realizar operações CRUD básicas, criada como parte dos estudos para a **Imersão 18 Full Cycle**.

## Requisitos

A aplicação deve ler um arquivo JSON ao iniciar (`go run main.go`) e manipular os dados exclusivamente na memória, sem salvar alterações no arquivo JSON.

### Endpoints Disponíveis:

- **GET /events**
  - Lista todos os eventos.

- **GET /events/:eventId**
  - Retorna os dados de um evento específico.

- **GET /events/:eventId/spots**
  - Lista os lugares disponíveis para um evento.

- **POST /event/:eventId/reserve**
  - Reserva um lugar para um evento específico.
  - **Dados:** `spots` é um array como por exemplo `[A1, B2]`.
  - Não é possível reservar o mesmo spot duas vezes. Retorna erro 400 se isso ocorrer.

## Testes no Postman

### GET /events
![Listar todos os eventos](https://github.com/JessanyKaline/api-golang-crud/blob/main/postman-test/GET-events.png)

### GET /events/:eventId
![Listar os dados de um evento](https://github.com/JessanyKaline/api-golang-crud/blob/main/postman-test/GET-eventById.png)

### GET /events/:eventId/spots
![Listar os lugares de um evento](https://github.com/JessanyKaline/api-golang-crud/blob/main/postman-test/GET-spots.png)

### POST /event/:eventId/reserve - Sucesso
![Reservar um lugar - Sucesso](https://github.com/JessanyKaline/api-golang-crud/blob/main/postman-test/POST-sucessfulyReserve.png)

### POST /event/:eventId/reserve - Erro
![Reservar um lugar - Erro](https://github.com/JessanyKaline/api-golang-crud/blob/main/postman-test/POST-%20alreadyReserved.png)

**OBESERVAÇÃO**: o código está todo comentado, pois visa servir de referência para estudos.
