# GoCallAPI
API that you can call to create calls and stops them.
The main functionality is to record the history of past and current calls. 

## Setup
### Create an .env file
Create a file called `.env` in the route of the project and add those variables with your own values.
- APP_DB_USERNAME=postgres
- APP_DB_PASSWORD=password
- APP_DB_NAME=postgres

### Create a Postgres Docker container
Run the following command to create a container to host the database: `docker run -it -p 5432:5432 -e POSTGRES_PASSWORD=password -d postgres`

## Routes
### Implemented
- POST `/call`: Call other User
- GET `/call/[id]`: cet information of a call by id
- PUT `/call/[id]`: Update call by id
- DELETE `/call/[id]`: Delete call by id

- GET `/calls`: Browse through all the calls

- GET `/stop/[id]`: End call with User

- GET `/health`: To get server health check

### To Be Implemented
`/auth`: Authenticate the User
