# Case task Software Engineering School
  
## Description

The Exchange API provides exchange rate information and allows users to subscribe and send emails related to exchange rates.

## Installation

1. Clone the repository:

```bash
git clone https://github.com/DenPedosenko/BTC_Exchange_Rate.git
```

2. Navigate to the project directory:

```bash 
cd BTC_Exchange_Rate
```
  
3. Build the Docker image:
```bash 
docker-compose build
```

## Usage

1. Run the Docker container:

```bash 
docker-compose run --rm --service-ports exchangeapi
```

2. The API will be accessible at `http://localhost:8080`.

## Endpoints

The following endpoints are available:

### [GET] /api/rate

- Description: Retrieve the latest exchange rate.

- Response: Returns the latest exchange rate.

Example curl command for testing:

```bash
curl -X 'GET' \
  'http://localhost:8080/api/rate' \
  -H 'accept: application/json'
```

Response:
```
1041856.44656
```

### [POST] api/subscribe

- Description: Subscribe to receive exchange rate notifications via email.

- Request: Requires a URL encoded value with the email address.

- Response: Returns a 200 code in the header if success and 409 if the email already exists in storage.


Example curl command for testing:

```bash
curl -X 'POST' \
  'http://localhost:8080/api/subscribe' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -d 'email=test@gmail.com'
```
  
### [POST] /sendEmails

- Description: Send emails to subscribed users with exchange rate updates.

- Response: Returns a 200 code in the header if the emails are sent successfully.
  
Example curl command for testing:

```bash
curl -X 'POST' \
  'http://localhost:8080/api/sendEmails' \
  -H 'accept: application/json' \
  -d ''
```
