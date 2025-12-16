# 🕒 Time Logger

A lightweight **time tracking web app** built with **Go + React + PostgreSQL + Docker**.  
The app allows you to log work by ticket and day, visualize it in a timesheet grid, and export later (Google Sheets and invoices planned).

---

## Running locally

- HTTP is exposed on `localhost:8085` (SSL on `localhost:8443`) to avoid clashing with other projects that may already use port 80/443.
- Bring everything up with `docker compose up --build` from the repo root.
- Postgres is reachable from your host at `localhost:5433` (user: `timelogger`, password: `password`, db: `timelogger`) if you want to run the Go server or psql locally; containers still talk to it via the internal `db:5432` address.
