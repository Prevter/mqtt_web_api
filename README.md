### Setup
1. Build the frontend
    ```bash
    cd frontend
    npm install
    npm run build
    ```
2. Build the backend
    ```bash
    go build .
    ```
3. Create a config file and edit it if needed
    ```bash
    cp .env.example .env
    ```
4. Run the server
    ```bash
    ./mqtt_web_api
    ```
   
### Details
The server is written in Go and uses a PostgreSQL database.
The frontend is written in Vue.js and uses Tailwind CSS.
Frontend is built with Vite and is served by the backend. 
Backend in addition to serving the frontend also provides a REST API for communication with the database.