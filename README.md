# Transaction Processor API
This application is a transaction processor API built with Go. It uses a PostgreSQL database and communicates with the Mailtrap API for sending emails.

## Prerequisites
- Docker
- Docker Compose

## Getting Started
The application is containerized using Docker and orchestrated with Docker Compose.

To start the application, navigate to the project directory where the `docker-compose.yaml` file is located and run the following command:

```bash
docker compose up --build
```

This command will start the golang_app and postgres_db services defined in the docker-compose.yaml file.  

The golang_app service is the main application, and it communicates with the postgres_db service, which is a PostgreSQL database.  

The application will be accessible at http://localhost:8000 by using the following endpoints:

```curl
curl --location --request POST 'http://localhost:8000/transactions/run-daily-report'
```

### Environment Variables
The application uses several environment variables, which are defined in the docker-compose.yaml file:  
- DATABASE_DSN: The data source name (DSN) for the PostgreSQL database.
- MAILTRAP_HOST: The API host for Mailtrap.
- MAILTRAP_TOKEN: The API token for Mailtrap. You can obtain this from your [Mailtrap](https://mailtrap.io/) account.
- MAILTRAP_FROM_EMAIL: The email address to send from when using Mailtrap.
- TRANSACTIONS_FILE_PATH: The file path for the transactions CSV file.
- EMAIL_LOGO_URL: The URL for the email logo.

Please replace the placeholders in the docker-compose.yaml file with your actual values before starting the application. 