# Network File Management System - Run Guide

## Prerequisites
- **Go**: 1.18+
- **Node.js**: 16+
- **PostgreSQL**: Running on localhost:5432 (User: postgres, Pass: password, DB: netfilessys)
- **Redis**: Running on localhost:6379

## Backend Setup

1.  **Database Setup**:
    Ensure PostgreSQL is running and create the database:
    ```sql
    CREATE DATABASE netfilessys;
    ```

2.  **Configuration**:
    Check `internal/config/config.yaml`. Adjust database/redis credentials if needed.

3.  **Run Migration**:
    Initialize the database schema:
    ```bash
    go run cmd/migrate/main.go
    ```

4.  **Start Server**:
    ```bash
    go run cmd/server/main.go
    ```
    The backend will start on `http://localhost:8080`.

## Frontend Setup

1.  **Install Dependencies**:
    ```bash
    cd frontend
    npm install
    ```

2.  **Start Dev Server**:
    ```bash
    npm run dev
    ```
    The frontend will start on `http://localhost:5173` (usually).

## Usage

1.  Open Browser at `http://localhost:5173`.
2.  **Register** a new user.
3.  **Login**.
4.  **Upload** files (supports chunked upload).
5.  **Share** files by selecting them and clicking "Share".
