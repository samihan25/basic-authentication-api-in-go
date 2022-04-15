# basic-authentication-api-in-go
Basic authentication APIs in Go

Used https://github.com/gorilla/mux for implementation


## Setup
git clone git@github.com:samihan25/basic-authentication-api-in-go.git

cd basic-authentication-api-in-go



## How to use
-  cd src
-  go run .


## Available API

1. signup
   - This API will accept unique usernames
   - Password and confirm_password should match
   - Each field should be non-empty
   - URL: POST http://localhost:3333/signup
   - BODY: {
            "username": "samihan25",
            "password": "abcd",
            "confirm_password": "abcd",
            "fullname": "Samihan Deshmukh"
        }
2. login
   - Username should match with existing username
   - Username and Password combination should be present. Signup should have een performed earlier.
   - URL: POST http://localhost:3333/login
   - BODY: {
            "username": "samihan25",
            "password": "abcd"
        }
3. profile
   - Username and OTP should be provided. OTP can be collected by executing login request.
   - URL: POST http://localhost:3333/profile
   - BODY: {
            "username": "samihan25",
            "otp": 123456
        }
4. logout
   - OTP generated will be purged. You cannot use old OTP, you will need to generate new OTP by Login.
   - URL: POST http://localhost:3333/logout
   - BODY: {
            "username": "samihan25"
        }
