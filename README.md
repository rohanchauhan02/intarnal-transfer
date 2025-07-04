# Internal Transfers Application

A robust internal transfers application developed using Golang and PostgreSQL.

## ✅ Features

- Account creation with configurable initial balances
- Real-time account balance queries
- Secure internal fund transfers

## 🛠 Technology Stack

- Go (Echo, Viper, Gorm)
- PostgreSQL
- HTTP RESTful APIs

## 🚀 Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/yourname/internal-transfers.git
cd internal-transfers
```

### 2. Prerequisites

Ensure `make` is installed on your system.

This project leverages a `Makefile` to streamline setup and development workflows.

### 3. Setup

Initialize Go modules and configuration files:

```bash
make setup
```

### 4. Running the Application

Start the application:

```bash
make start
```

### 5. Running Tests

Execute tests with coverage reporting:

```bash
make test
```

### 6. Working with Mocks

Install `mockgen` for generating mocks:

```bash
make install-mockgen
```

Generate mocks:

```bash
make mock
```
