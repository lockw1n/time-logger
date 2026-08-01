# `tl` — Time Logger CLI

A terminal client for the Time Logger backend: log time (manually or with a live
timer), view your timesheet, fix or delete entries, and pull invoices — all
without opening the web UI. It talks only to the REST API (through nginx) — no
direct database access.

## Install

From the `backend/` directory:

```bash
go build -o tl ./cmd/tl        # produces ./tl
# or install onto your PATH:
go install ./cmd/tl            # produces $GOBIN/tl (or $GOPATH/bin/tl)
```

`tl version` prints the build version (`dev` unless set via
`-ldflags "-X main.version=..."`).

## First-run setup

Config and session state live in `~/.config/tl/` (override with `TL_CONFIG_DIR`):
`config.json` (settings), `token.json` (JWT, mode `0600`), and `timer.json`
(the running timer, when one is active).

1. **Point at the backend.** In the current local stack the backend is reachable
   only via nginx over HTTPS with a self-signed cert, so use HTTPS + `insecure`:

   ```bash
   tl config set base_url https://localhost:8443
   tl config set insecure true          # skip TLS verify for the self-signed cert
   # (alternatively: tl config set ca_cert /path/to/ca.pem)
   ```

2. **Log in.** The JWT lasts 30 minutes; re-run `login` when commands start
   reporting `not logged in or session expired`.

   ```bash
   tl login                             # prompts for email + password (no echo)
   tl logout                            # deletes the stored token
   ```

3. **Set your defaults** so you don't repeat `--company`/`--activity`:

   ```bash
   tl companies                         # find your company id
   tl config set default_company 1
   tl config set default_activity "Spryker Housekeeping"
   tl config list                       # show all settings
   ```

## Commands

### `tl add` — log time

```
tl add <ticket-code> <duration> [comment...]
  -d, --date <date>       today (default) | yesterday | mon..sun | YYYY-MM-DD
  -a, --activity <name>   activity name (case-insensitive); required unless
                          default_activity is set
  -c, --company <id>      overrides default_company
  -y, --yes               skip confirmations (merge + future-date warning)
```

**Durations:** `90m`, `90`, `1h`, `1h30m`, `1h30`, `1.5h`, `0.25h` (1m–24h).

```bash
tl add APP-123 1h30m "code review"
tl add APP-200 45m -d yesterday
tl add APP-123 2h -a "myblum" -d fri     # most recent Friday
```

**Merge behavior:** logging again to the *same ticket + activity + day* never
creates a duplicate — it offers to add the minutes to the existing entry
(joining comments with `; `):

```
$ tl add APP-123 45m "more"
APP-123/Spryker Housekeeping already has 1h30m on 2026-07-30 — add 45m for a total of 2h15m? [Y/n]
```

Pass `-y` to auto-confirm. Logging to a **future date** asks for confirmation
(also skipped by `-y`). A date not covered by an active contract fails with
`no active contract covers <date>`.

### `tl start` / `tl status` / `tl stop` / `tl cancel` — live timer

Instead of logging a fixed duration, time your work as you go. There is one
timer at a time (state in `~/.config/tl/timer.json`, no background daemon).

```
tl start [ticket-code]
  -a, --activity <name>   activity (required unless default_activity is set)
      --note <text>       note attached to the entry when the timer stops
      --at <HH:MM>        backdate today's start (must be in the past)

tl status                 show the running timer + elapsed (local only, no network)
tl stop [comment...]      log the elapsed time and clear the timer
      --duration <dur>    log this amount instead of the measured elapsed time
  -y, --yes               skip confirmations
tl cancel                 discard the running timer without logging it
  -y, --yes
```

```bash
$ tl start                       # on branch feature/APP-123-widget
inferred ticket APP-123 from branch 'feature/APP-123-widget'
⏱ started APP-123 (Spryker Feature Development) at 09:15
$ tl status
⏱ APP-123 (Spryker Feature Development) — 1h20m elapsed (started 09:15)
$ tl stop "implemented widget"
```

- **Ticket inference:** with no `ticket-code`, `tl start` reads the current git
  branch and extracts a ticket like `APP-123` (e.g. from `feature/APP-123-widget`),
  printing what it inferred. Outside a repo, or when no code is found, pass one
  explicitly.
- The activity is resolved **up front**, so `stop` can't fail hours later on a typo.
- `tl stop` logs to the day it runs, through the same merge flow as `tl add`, and
  only clears the timer **after** the entry is safely saved. A comment given to
  `stop` wins over the `--note` from `start`. If your session has expired when you
  stop, the timer is kept (exit 2) — run `tl login`, then `tl stop` again.
- **Forgotten-timer guards:** a measured run over 24h is refused (use `--duration`
  to log a specific amount); over 12h it asks for confirmation first.
- `tl status` is a pure local read — it works offline and is cheap enough for a
  shell prompt (`--json` supported; prints `{"running":false}` and exits 1 when idle).

### `tl today` / `tl week` — view the timesheet

```bash
tl today                 # single day; rows + entry #ids + total
tl today -d yesterday
tl week                  # Mon–Sun grid for the current week
tl week  -d 2026-07-15   # the week containing that date
```

The entry `#id` shown by `tl today` is what `tl edit`/`tl delete` operate on.

### `tl ui` — interactive weekly timesheet

```bash
tl ui                    # full-screen grid for the current week
tl ui -d 2026-07-15      # open at the week containing that date
tl ui -c 2               # open for a specific company
```

A full-screen, editable version of `tl week`: the same Mon–Sun grid, reviewed and
fixed up in place. It is a view/edit layer over the same API the one-shot commands
use — adds go through the identical create-or-merge path, so it never creates a
duplicate entry.

| Key | Action |
|---|---|
| arrows / `hjkl` | move the cell cursor |
| `enter` | edit the cell — empty cell adds, a single entry edits inline, 2+ entries open a picker |
| `a` | add an entry for the cursor's day (ticket, activity `←/→`, duration, comment) |
| `d` | delete the cell's entry (with confirm; a multi-entry cell opens the picker) |
| `[` / `]` | previous / next week · `t` jumps to this week |
| `c` | switch company · `r` refresh |
| `q` / `Ctrl-C` | quit |

In the inline duration editor, submitting an **empty** value deletes the entry
(after confirmation). Editing a cell writes to the server and the totals row
recalculates; API errors appear in the footer without exiting. The header shows
the running timer and pending-outbox count when present. A terminal narrower than
100 columns shows a resize notice instead of a broken grid. `tl ui` needs a real
terminal (it will not run with its input piped).

### `tl edit` — fix an entry's duration or comment

```
tl edit <entry-id> [--duration <dur>] [--comment <text>]
  --duration <dur>   new duration (e.g. 2h, 90m, 1h30m)
  --comment <text>   replace the comment
  -y, --yes          skip the confirmation
```

Only the **duration** and **comment** can be edited. To change the ticket,
activity or date, `tl delete` the entry and `tl add` it again. At least one of
`--duration`/`--comment` must be given.

Before mutating, `edit` fetches the entry and shows a one-line before → after
summary, then confirms (skip with `-y`):

```bash
$ tl edit 421 --duration 2h
edit #421 APP-123 (2026-07-13): 45m → 2h? [y/N]
```

**Clearing a comment is not possible** through the current API: the backend
drops empty comments, so both a `null` and an empty string are ignored. To
remove a comment, `tl delete` the entry and `tl add` it again without one.
(That is why there is no `--clear-comment` flag; passing an empty `--comment` is
rejected with this guidance.)

### `tl delete` — remove an entry

```
tl delete <entry-id>
  -y, --yes          skip the confirmation
```

Fetches the entry, shows it, and confirms before deleting (skip with `-y`):

```bash
$ tl delete 421
delete 45m on APP-123 (2026-07-13) "fixed mapper"? [y/N]
deleted entry #421
```

### `tl invoice` — download an invoice document

```
tl invoice [--month YYYY-MM | --start YYYY-MM-DD --end YYYY-MM-DD]
  --month <YYYY-MM>   whole month (default: the previous calendar month)
  --start / --end     explicit range (mutually exclusive with --month)
  --format <fmt>      pdf (default) | excel
  -o, --output <path> output file or directory (default: current directory)
  --force             overwrite an existing file
  -c, --company <id>  overrides default_company
```

Streams the generated document to disk and prints the absolute path written. The
filename comes from the server's `Content-Disposition`, falling back to
`invoice-<start>-<end>.<ext>`. `-o` may point at a directory (the server
filename is used inside it) or be a full target path. An existing file is not
overwritten without `--force`.

```bash
tl invoice                       # previous month, PDF, into the current dir
tl invoice --month 2026-06 -o ~/invoices/
tl invoice --start 2026-06-01 --end 2026-06-15 --format excel
```

Invoice generation round-trips a PDF renderer and can be slow, so `tl invoice`
prints `generating…` to stderr and allows up to 60s for this request. A period
with no entries (or any validation problem) surfaces the server's error message
and writes no file.

### `tl companies` / `tl activities` — reference data

```bash
tl companies             # id  name  name_short
tl activities            # id  name  billable  priority (for the resolved company)
tl activities -c 1
```

### `tl ping` — connectivity check

```bash
tl ping                  # GET /health — proves nginx+TLS reachability
```

### `tl sync` — offline queue

When the backend is unreachable, `tl add` and `tl stop` don't fail — they queue
the entry to `~/.config/tl/outbox/` (exit code `3`) and print e.g.
`api unreachable — queued 45m on APP-123 for sync (1 pending)`. The measured time
is captured in the queued payload, so `tl stop` also drops the timer: its job is
done. Every subsequent API command replays the queue in the background before its
own work, printing `synced N queued entries` when it finishes.

```bash
tl sync                  # force a foreground replay now; report and exit
tl sync --list           # show pending + failed ops without replaying
tl sync --discard <file> # drop a permanently-failed op (from --list) after confirm
```

Details worth knowing:

- **Only entry creates are queued.** Edits, deletes and invoices fail fast when
  offline — deferring a mutation against server state you can't see is a
  corruption risk, not a convenience.
- **Only unreachable writes queue.** A write is queued only when the request
  provably never reached the backend (connection refused, DNS failure). A request
  that was cancelled (Ctrl-C) or timed out *after being sent* is ambiguous — the
  server may already have committed it — so it is reported as interrupted and
  **not** queued, to avoid replaying a create the server already applied.
- **Merge on replay.** Two queued entries for the same ticket/activity/day
  collapse into a single entry with summed minutes via the same 409-merge the
  interactive path uses — non-interactively, no prompt.
- **At-least-once (during replay).** If a replay op is sent but not confirmed
  (the sync is cancelled or the process dies), it stays queued and is retried; the
  server's unique index + merge keeps it to one row. Because merge *sums*, a
  redelivered create can over-count minutes in that rare window — replay is
  intentionally at-least-once, so prefer `tl sync` to let it finish cleanly.
- **Permanent failures** (a 400, a date no contract covers) and **unreadable op
  files** are moved to `outbox/failed/` — the reason recorded for the former, the
  original bytes preserved for the latter — never silently dropped. Inspect with
  `tl sync --list`, clear with `tl sync --discard`.
- Offline `tl add` resolves the activity name from a local cache of the last
  activities fetch, so run `tl activities` (or any add/start/stop) online once per
  company first; otherwise offline add fails with guidance instead of queueing.

## Scripting with `--json`

Every read command (`companies`, `activities`, `today`, `week`) can print the raw
API response for piping into `jq`, and `tl status` emits a JSON view of the
running timer:

```bash
tl week --json | jq '.totals.overall_minutes'
tl activities --json | jq '.[] | {id, name}'
tl today --json | jq '.rows[].entries[].id'
tl status --json | jq 'select(.running) | .elapsed_minutes'
```

## Global flags

`-c, --company <id>` · `--json` · `--base-url <url>` · `--insecure` ·
`--ca-cert <pem>` — all override the corresponding config value for a single run.

## Exit codes

| Code | Meaning |
|---|---|
| `0` | success |
| `1` | API / validation / other error |
| `2` | not logged in or session expired — run `tl login` |
| `3` | operation queued offline — accepted locally but not yet confirmed by the backend (see `tl sync`) |

## Notes

- Durations always display as `7h30m` (never decimal hours); dates as `YYYY-MM-DD`.
- The backend returns response dates as `DD.MM.YYYY` but accepts `YYYY-MM-DD` in
  requests. Human-readable output normalizes dates to `YYYY-MM-DD`, so you always
  see and type that form. `--json`, by contrast, prints the backend's response
  **verbatim** (dates as `DD.MM.YYYY`) so scripts see the true payload.
