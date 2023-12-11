# RESTful APIs

@author: Deming Chen (陈德铭)

[TOC]

## Sign in

Path: `\signin`

Methods: POST

Data:

```json
{
    "email": "example@example.com",
    "password": "example"
}
```

Header: `{"Content-Type": 'application/json'}`

Content-Type: `"application/json" `

Response:

- 200 (OK) `{"token": tokenString, "username": user.Username, "id": user.ID}`
- 400 (StatusBadRequest): `{"error": [error message]}`
- 401 (Unauthorized): `{"status": "unauthorized"}`
- 500 (InternalServerError): `{"error": [error message]}`

## Sign up

Path: `\signup`

Methods: POST

Data:

```json
{
    "email": "example@example.com",
    "password": "example",
    "username": "Example"
}
```

Header: `{ "Content-Type": 'application/json' }`

Response:

- 200 (OK) `"success"`

- 400 (StatusBadRequest): 

  ```json
  {
      "error": [error message],
      "status": "Invalid username or password or email" | "Email is already taken"	
  }
  ```

- 500 (InternalServerError): `{"error": [error message]}`

## Sellers Add Items

Path: `/additem`

Methods: POST

Data (`formData`):

```
"price": float
"tag": int
"description": string
"title": string
```

Header: `{ Authorization: "Bearer [token]"},`

Response: besides [authentication errors](#authErr), `additem` may return one of the following result:

- 200 (OK)

- 400 (StatusBadRequest): `{"error": ""invalid [price|tag]"}`
- 500 (StatusInternalServerError): `{"error": "cannot open file [filename]: [error]"}`
- 500 (StatusInternalServerError): `{"error": "server cannot add item: [error]"}`

## Seller Query List

Path: `/qlist`

Methods: GET

Data: no data needed

Header: `{ Authorization: "Bearer [token]"}`

Response: besides [authentication errors](#authErr), `qlist` may return one of the following result:

- 200 (OK)
- 500 (StatusInternalServerError): `{"error": "[error]"}`

## Query Item's Detail

Path: `/qitem`

Methods: GET

Data:

```json
{
    "item_id": "example_id"
}
```

Header: nothing

Response:

- 200 (OK)

  ```json
  {
  	"seller_id":   uint,
  	"title":       string,
  	"price":       float,
  	"tag":         int,
  	"description": string,
  	"status":      int,
  	"image_urls":  {
          "[internal_name]": [url] string,
          ...
      }
  }
  ```

  

- 400 (StatusBadRequest): `{"status": "invalid item ID: [item_id]"}`

- 500 (StatusInternalServerError): `{"error": unable to read item: [error]}`


## Buyer Add an Order

Path: `/addorder`

Methods: POST

Data: `{item_id: [id]}`

Header: `{ Authorization: "Bearer [token]"}`

Response: besides [authentication errors](#authErr), `qorder` may return one of the following result

- 500 (StatusInternalServerError): `{"error": [error]}`
- 200 (OK): `{[item1, item2, ...]}`

## Buyer Query Orders

Path: `/qorder`

Methods: GET

Data: no data needed

Header: `{ Authorization: "Bearer [token]"}`

Response: besides [authentication errors](#authErr), `qorder` may return one of the following result

- 500 (StatusInternalServerError): `{"error": [error]}`
- 200 (OK):

  ```json
  {
  	Id1: remainTime1,
      Id2: remainTime2,
      ...
  }
  ```

  

## Buyer Cancel an Order

Path: `/qorder`

Methods: Post

Data: `{item_id: [id]}`

Header: `{ Authorization: "Bearer [token]"}`

Response: besides [authentication errors](#authErr), `qorder` may return one of the following result

- 500 (StatusInternalServerError): `{"error": [error]}`
- 400 (StatusBadRequest): `{"status": "Order does not exist!"}`
- 200 (OK)

## Buyer Check Out an Order

Path: `/checkout`

Methods: Post

Data: `{item_id: [id]}`

Header: `{ Authorization: "Bearer [token]"}`

Response: besides [authentication errors](#authErr), `qorder` may return one of the following result

- 500 (StatusInternalServerError): `{"error": [error]}`
- 400 (StatusBadRequest): `{"status": "order does not exist!"}`
- 200 (OK)

## Search Items

Path: `/search`

Methods: GET

Data:

```json
{
    "tag": "0",
    "keywords": "some keywords"
}
```

Header: nothing

Response:

- 200 (OK): `[item_id1, item_id2, item_id3, ...]`

- 400 (StatusBadRequest): `{"status": "invalid tag: [tag]"}`

- 500 (StatusInternalServerError): `{"error": [error]}`


## Delete an Item

Path: `/ditem`

Method: DELETE

Data: `{"item_id": "example_id"}`

Header: `{ Authorization: "Bearer [token]"}`

Response:

- 200 (OK)
- 400 (StatusBadRequest): `{"error": [Item [ID] does not exist."}`
- 400 (StatusBadRequest): `{"status": "cannot delete item whose status is [Status]"}`
- 401 (StatusUnauthorized): `{"status": "unauthorized"}`
- 500 (StatusInternalServerError): `{"error": [error]}`

## <a name="authErr">Authentication Errors</a> 

- 401 (StatusUnauthorized): `{"status": "Invalid token"}`
- 401 (StatusUnauthorized): `{"status": "Token expired"}`
- 401 (StatusUnauthorized): `{"error": "Unable to parse [token|id|expire time]: [error]"}`
- 500 (StatusInternalServerError): `{"error": [error]}`
