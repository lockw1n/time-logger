# 🕒 Time Logger

A lightweight **time tracking web app** built with **Go + React + PostgreSQL + Docker**.  
The app allows you to log work by ticket and day, visualize it in a timesheet grid, and export later (Google Sheets and invoices planned).

---

## Running locally

- HTTP is exposed on `localhost:8085` (SSL on `localhost:8443`) to avoid clashing with other projects that may already use port 80/443.
- Bring everything up with `docker compose up --build` from the repo root.
- Postgres is reachable from your host at `localhost:5433` (user: `timelogger`, password: `password`, db: `timelogger`) if you want to run the Go server or psql locally; containers still talk to it via the internal `db:5432` address.

## Company profile for invoices

The invoice PDF now uses company details stored in the database (table `companies`). Set or update the single company record via `PUT /api/company`:

```bash
curl -X PUT http://localhost:8085/api/company \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Acme Inc",
    "uid": "AT U1234567",
    "address_line1": "Main Street 1",
    "zip": "1010",
    "city": "Vienna",
    "country": "Austria"
  }'
```

Fetch the current values with `GET /api/company`. Invoice generation will fail with `400 company not configured` until a record exists.

Invoice numbers are auto-generated based on the report month: they are always the first day of the following month in `YYYYMMDD` format (e.g., report for `2025-11` → invoice number `20251201`).
