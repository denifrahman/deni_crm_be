## Installation
### Cloning the Repository
```bash
git clone https://github.com/denifrahman/deni_crm_be.git
cd deni_crm_be
```

1. Install dependencies:
    ```bash
    go mod tidy
    ```


2. Start the development server:
    ```bash
    Run: go run cmd/server/main.go
    ```


3. Example env : 
    ``` bash
    DATABASE_URL=string_url_postgresql
    JWT_SECRET=my_very_secret_key
    PORT=8080
    ```
    
### Deployment 

Deployed on google cloud

>url api : https://service-crm-231482427087.asia-southeast2.run.app

