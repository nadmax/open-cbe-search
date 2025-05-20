# open-cbe-search

**open-cbe-search** is a Belgian enterprise search engine powered by **Elasticsearch**, written in **Go**, and designed to ingest and index large-scale datasets from CSV files.  
It provides fast, scalable, and full-text search capabilities over structured Belgian company data.

## 📦 Requirements
- [Docker](https://docs.docker.com/engine/)
- [Docker Compose](https://docs.docker.com/compose/)
- A `data` folder at the root of the project containing all relevant `.csv` files
- A `.env` file in `backend/engine` folder, based on the [`.env.example`](https://github.com/nadmax/open-cbe-search/blob/engine/.env.example) file

## 🛠️ Getting Started
### 1. Clone the Repository
```bash
git clone https://github.com/nadmax/open-cbe-search.git
cd open-cbe-search
```

### 2. Prepare the Required Files
- Place all .csv files in the `data` directory (e.g., `enterprise.csv`, `branch.csv`, etc.)
- Create `.env` inside `backend/engine`, based on the [`.env.example`](https://github.com/nadmax/open-cbe-search/blob/engine/.env.example) file

### 3. Run with Docker Compose
```bash
docker compose up -d
```