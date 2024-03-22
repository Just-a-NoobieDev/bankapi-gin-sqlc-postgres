# Golang Bank API

This is a simple bank api that I made using Go programming language and other technologies

## Technologies I used

- Framework Gin (gin-gonic)
- SQLc
- golang-migrate
- PostgresSQL
- golang-viper (config)
- testify (unit testing)

## Endpoints (so far)

- Version 1 `/api/v1`

  - accounts

    - `GET` all accounts paginated

      - endpoint `/accounts?page=?&size=?`
      - Query Params
        - `page` `required` page number
        - `size` `required` size of data per page

    - `GET` account

      - endpoint `/accounts/:id`
      - Params -`:id` specific account id

    - `POST` create account

      - endpoint `/accounts`
      - Body
        - `name` full name of account
        - `currency` currency supported currently (USD EUR CAD)

    - `POST` deposit

      - endpoint `/accounts/deposit`
      - Body
        - `id` id of the account
        - `amount` number of money to be deposit (currently the data type of it is integer will change later)

    - `DELETE` account

      - endpoint `/accounts/:id`
      - Params -`:id` specific account id

  - transfers

    - `GET` all transfers by account paginated

      - endpoint `/transfers?id=?&page=?&size=?`
      - Query Params
        - `id` `required` id of the account
        - `page` `required` page number
        - `size` `required` size of data per page

    - `GET` transfer

      - endpoint `/transfers/:id`
      - Params -`:id` specific transfer id

    - `POST` create transfer

      - endpoint `/transfers`
      - Body
        - `from_account_id` id of the sender
        - `to_account_id` id of the receiver
        - `amount` amount to be transfer
        - `currency` currency supported currently (USD EUR CAD)

  - entry

    - `GET` all entry by account paginated

      - endpoint `/entry?id=?&page=?&size=?`
      - Query Params
        - `id` `required` id of the account
        - `page` `required` page number
        - `size` `required` size of data per page

  - users

    - `POST` create / register user

      - endpoint `/users/register`
      - Body `all fields required`
        - `username` primary key, foreign key to the accounts model
        - `full_name` name of the user
        - `email` `unique` email of the user
        - `password` `min=6` 
        - `password_again`
