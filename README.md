#### Local Setup
1.> Install golang and set GOPATH and GOROOT accordingly.

2.> Clone the github repo in the "src" folder of golang installation. (https://github.com/pphowakande/accountservice.git)


3.> Build compiles the accountservice named by the import paths, along with their dependencies by running below command

    $ go build

4.> Now the binary file created by "go build" command as

    $ ./accountservice

5.> With above command , you would get this app running at - http://localhost:8080/


## REST Endpoints

1.> Endpoint to Generate Authentication Token
----
This endpoint returns JWT token.

* **URL**

  /authenticate

* **Method:**

  `POST`

*  **URL Params**

    None

* **Data Params**

    `{"username":"poonam","password":"poonam"}`

* **Success Response:**

  * **Code:** 201 <br />
    **Content:** `{"token":"12345"}`

* **Error Response:**

  * **Code:** `<Error code number>` <br />
    **Content:** `{ errors : "<Respective error message>" }`


* **Sample Call**

  curl --header "Content-Type: application/json" --request POST --data '{"username": "poonam", "password":"poonam"}' http://localhost:8080/authenticate


2.> Endpoint to Create Account
----
This endpoint creates account.

* **URL**

  /account

* **Method:**

  `POST`

*  **URL Params**

    None

* **Data Params**

    `{"id":"1","name":"poonam"}`

* **Success Response:**

  * **Code:** 201 <br />
    **Content:** `{"id":"1","name":"poonam"}`

* **Error Response:**

  * **Code:** `<Error code number>` <br />
    **Content:** `{ errors : "<Respective error message>" }`


* **Sample Call**

  curl -H "authorization: JWT {token}" --request POST --data '{"id": "2","name": "person_2"}' http://localhost:8080/account


3.> Endpoint to get account details
----
This endpoint return account details using accoun id.

* **URL**

  /account/{id}

* **Method:**

  `GET`

*  **URL Params**

  None

* **Data Params**

  `{"id": "1"}`

* **Success Response:**

* **Code:** 200 <br />
  **Content:** `{"id":"1","name":"poonam"}`

* **Error Response:**

* **Code:** `<Error code number>` <br />
  **Content:** `{ errors : "<Respective error message>" }`


* **Sample Call**

  curl -H "authorization: JWT {pass token here}" -X GET http://localhost:8080/account/1


#### Run Test cases

  go test ./...
