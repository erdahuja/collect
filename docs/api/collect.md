# Project: collect
# 📁 Collection: 1. users 


## End-point: create user
### Method: POST
>```
>http://localhost:3000/v1/users
>```
### Body (**raw**)

```json
{
    "name": "rishabh",
    "email": "rishabh@example.com",
    "roles": [
        "USER"
    ],
    "password": "rishabh",
    "passwordConfirm": "rishabh"
}
```

### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJtZXJnZWR1cCIsInN1YiI6IjEiLCJleHAiOjE2ODA4ODY2MjEsImlhdCI6MTY4MDg4MzAyMSwicm9sZXMiOlsiQURNSU4iLCJVU0VSIl19.wqUWPldniYveFk8hdjETi99fAxHnrYdUGnOMt0yj3uw|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: get token
### Method: GET
>```
>http://localhost:3000/v1/users/token/1
>```
### 🔑 Authentication basic

|Param|value|Type|
|---|---|---|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: view users
### Method: GET
>```
>http://localhost:3000/v1/users
>```
### Headers

|Content-Type|Value|
|---|---|
|Authorization|Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJtZXJnZWR1cCIsInN1YiI6IjEiLCJleHAiOjE2ODA4NjIyMjUsImlhdCI6MTY4MDg1ODYyNSwicm9sZXMiOlsiQURNSU4iLCJVU0VSIl19.PV4JRkWmkQ-z3lyWL2LHV40Bglykokf3xyBGtk2hKs4|


### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJtZXJnZWR1cCIsInN1YiI6IjEiLCJleHAiOjE2ODA4OTEwNDAsImlhdCI6MTY4MDg4NzQ0MCwicm9sZXMiOlsiQURNSU4iLCJVU0VSIl19.zRptzPqhHKr0gXx8Er0KP7HKXzgqI2CtWLAEb6iEBro|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃
# 📁 Collection: 2. forms 


## End-point: Get form
### Method: GET
>```
>undefined
>```

⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: add form
### Method: POST
>```
>http://localhost:3000/v1/forms
>```
### Body (**raw**)

```json
{
          "form_title": "Test survey",
        "form_description": "Did you do a good job deepak?"
}
```

### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjb2xsZWN0Iiwic3ViIjoiMSIsImV4cCI6MTY4NDM4ODk4NiwiaWF0IjoxNjg0Mzg1Mzg2LCJyb2xlcyI6WyJBRE1JTiIsIlVTRVIiXX0.F4IWdxxZSwsoYmMJ9PPHpJgkNOtG3ia4Nk-2K4at660|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: get questions
### Method: GET
>```
>undefined
>```

⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃
# 📁 Collection: 3. questions 


## End-point: Create questions
### Method: POST
>```
>http://localhost:3000/v1/items
>```
### Body (**raw**)

```json
{
    "name": "test product",
    "quantity": 10,
    "cost": 100
}
```

### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJtZXJnZWR1cCIsInN1YiI6IjEiLCJleHAiOjE2ODA4NjM4MTAsImlhdCI6MTY4MDg2MDIxMCwicm9sZXMiOlsiQURNSU4iLCJVU0VSIl19.AGRscERM3QHYJRAsOZDCgg5Ayw8UGLtE_00ZvlJIsIc|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃
# 📁 Collection: 4. response 


## End-point: Create Response
### Method: POST
>```
>http://localhost:3000/v1/response
>```
### Body (**raw**)

```json
{
  "form_id": 2,
  "respondent_id": "deepak"
}
```

### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjb2xsZWN0Iiwic3ViIjoiMCIsImV4cCI6MTY4NDMwNDU0OCwiaWF0IjoxNjg0Mjk4NTQ4LCJyb2xlcyI6WyJVU0VSIiwiQ09MTEVDVE9SIl19.wkn8HVOd-wHk7M0m5Ffbb3mc1i85IZHId71aYZNx94k|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: get responses
### Method: GET
>```
>undefined
>```

⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃
# 📁 Collection: 5. answer 


## End-point: Create answer
### Method: POST
>```
>http://localhost:3000/v1/response/33/answer
>```
### Body (**raw**)

```json
{
  "question_id": 1,
  "answer_text": "This is the answer text."
}

```

### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjb2xsZWN0Iiwic3ViIjoiMCIsImV4cCI6MTY4NDMwNDU0OCwiaWF0IjoxNjg0Mjk4NTQ4LCJyb2xlcyI6WyJVU0VSIiwiQ09MTEVDVE9SIl19.wkn8HVOd-wHk7M0m5Ffbb3mc1i85IZHId71aYZNx94k|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃
_________________________________________________
Powered By: [postman-to-markdown](https://github.com/bautistaj/postman-to-markdown/)
