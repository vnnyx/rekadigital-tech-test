# GOLANG-DOT-API

## Tech stack

- Golang v1.20
- Echo V4
- Redis for caching
- Golang ozzo-validator
- Gomock for mock generator
- Docker and kubernetes for deployment

## Design Pattern

**Entity > Repository > Usecase > Delivery (Echo REST API)**

Desain pattern yang digunakan pada aplikasi adalah menggunakan clean architecture, dimana pada desain ini memiliki 4 layer utama, yaitu entity, repository, usecase, dan juga delivery dengan menggunakan http.

## Hal yang perlu diperhatikan

1. Bagaimana jika terdapat ribuan transaksi pada database?

   Untuk menangani ini, hal yang pertama saya lakukan adalah menggunakan pagination untuk endpoint `GET` transaction sehingga data yang diambil menjadi lebih sedikit dan dapat disesuaikan dengan kebutuhan. Selanjutanya untuk menagani ini juga bisa melakukan indexing pada skema databae, indexing ini dilakukan untuk kolom yang sering dilakukan query.

2. Bagaimana jika terdapat banyak user yang mengakses API tersebut secara bersamaan?

   Dalam hal ini apabila banyak user yang mengakses API secara bersamaan maka server dapat mengalami down. Dengan demikian maka diperlukan caching pada data yang sering di request, seperti get list transaction. Hal ini juga dapat diatasi dengan melakukan locking untuk setiap request yang masuk, sehingga request akan masuk secara bergantian dan hal ini juga dapat mengatasi server down.

## Setup

### Local

1. run `export ENV=local`
2. run `go mod tidy`
3. run main app with `go run main.go server`

### Docker

1. run `docker-compose up -d`
2. try to hit `http://localhost:3000/rekadigital-api`

## Live Demo

I have deployed the application that can be tested via:

`https://cloud.vnnyx.my.id/dot-api/{ENDPOINT}`

## List Endpoint

```
BASEURL: https://cloud.vnnyx.my.id/rekadigital-api

POST /transaction
GET /transaction
```

## Testing (Integration and Unit Testing)

How to run:

1. Customizing the configuration file to perform integration tests located in `/test/integration/configs/test.yaml`.
2. After customizing the test.yaml file, run `make test` (if you have the make command), you can also copy the test command from the Makefile to perform testing.
3. You can also view the test coverage using `make cover` or run the cover command in the Makefile.

![testing](https://imgur.com/NDjrexA.png)

## Deployment

![deployment](https://imgur.com/ePe9oT3.png)

![deployment](https://imgur.com/O8FkDYu.png)

## API Docs

For the full documentation, it can be accessed through the following link.

[API Docs](https://documenter.getpostman.com/view/24450154/2s93XyUPEu)
