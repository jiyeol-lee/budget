# Budget Tracker

A modern budget tracking application with AI-powered receipt processing. Manage your monthly budgets, track fixed expenses, and automatically extract expense data from receipt PDFs.

![Budget Tracker](https://img.shields.io/badge/Status-In%20Development-yellow)
![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)
![Node Version](https://img.shields.io/badge/Node.js-18+-339933?logo=node.js)
![License](https://img.shields.io/badge/License-AGPL--3.0-blue)

## Overview

Budget Tracker helps you stay on top of your finances by:

- **Setting Monthly Budget Limits** - Define spending limits and get notified when approaching your threshold
- **Tracking Expected Expenses** - Manage recurring weekly and monthly expenses
- **AI Receipt Processing** - Upload receipt PDFs and let AI extract expense details automatically

## Tech Stack

| Layer        | Technology                                         |
| ------------ | -------------------------------------------------- |
| **Frontend** | SvelteKit (Svelte 5) + TailwindCSS v4 + TypeScript |
| **Backend**  | Pure Go (net/http) + SQLite                        |
| **AI**       | Claude Sonnet 4.5 via Anthropic API                |

## Project Structure

```
budget-tracker/
├── frontend/                    # SvelteKit application
│   ├── src/
│   │   ├── lib/
│   │   │   ├── components/      # Reusable UI components
│   │   │   ├── stores/          # Svelte stores (state management)
│   │   │   └── utils/           # API client, helpers
│   │   └── routes/              # Pages (SvelteKit file-based routing)
│   ├── static/                  # Static assets
│   ├── package.json
│   └── svelte.config.js
├── backend/                     # Go application
│   ├── cmd/server/              # Entry point (main.go)
│   └── internal/
│       ├── api/                 # HTTP handlers, router, middleware
│       ├── models/              # Data structures
│       ├── repository/          # Database operations (SQLite)
│       └── services/            # AI client (OpenAI integration)
└── README.md
```

## Getting Started

### Prerequisites

- **Go** 1.22 or higher
- **Node.js** 18 or higher
- **OpenAI API Key** (required for AI receipt processing)

### Environment Variables

| Variable             | Required    | Description                                                                                                |
| -------------------- | ----------- | ---------------------------------------------------------------------------------------------------------- |
| `ANTHROPIC_API_KEY`  | Yes         | API key for Claude Sonnet 4.5 receipt processing                                                           |
| `TURSO_MODE`         | No          | Database connection mode: `local` (default, file-based SQLite) or `remote` (Turso cloud)                   |
| `TURSO_LOCAL_PATH`   | No          | File path for local SQLite database (default: `./data/budget.db`). Used when `TURSO_MODE=local` or not set |
| `TURSO_DATABASE_URL` | Conditional | Turso database URL. Required when `TURSO_MODE=remote`                                                      |
| `TURSO_AUTH_TOKEN`   | Conditional | Turso authentication token. Required when `TURSO_MODE=remote`                                              |

### Running the Backend

```bash
# Navigate to backend directory
cd backend

# Download dependencies
go mod download

# Start the server with your API key (local SQLite database)
ANTHROPIC_API_KEY=your-api-key-here go run ./cmd/server

# Server starts on http://localhost:8080

# For production with Turso cloud database (optional)
TURSO_MODE=remote \
  TURSO_DATABASE_URL=libsql://your-database.turso.io \
  TURSO_AUTH_TOKEN=your-auth-token \
  ANTHROPIC_API_KEY=your-api-key \
  go run ./cmd/server
```

### Running the Frontend

```bash
# Navigate to frontend directory
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev

# App available at http://localhost:5173
```

## Features

### Monthly Budget Limits

Set a maximum spending limit for each month to keep your finances in check.

- **Configurable Threshold**: Get notified when reaching your budget threshold (default: 80%)
- **Visual Progress**: Track your spending progress with visual indicators
- **Real-time Updates**: See your remaining budget update as you add expenses

### Expected Expenses Tracking

Organize and track your recurring expenses by frequency.

#### Weekly Expenses

For regular weekly purchases like groceries:

| Item              | Source | Expected Amount |
| :---------------- | :----- | :-------------- |
| Organic Spaghetti | Costco | $12.59          |
| Sweet Onions      | Costco | $4.39           |

#### Monthly Expenses

For monthly recurring costs like subscriptions:

| Item   | Source | Expected Amount |
| :----- | :----- | :-------------- |
| Rice   | H-Mart | $30.00          |
| Kimchi | H-Mart | $30.00          |

The app automatically calculates your estimated monthly total based on weekly (x4) and monthly expenses.

### AI Receipt Processing (Core Feature)

Transform paper receipts into structured expense data with AI.

**How it works:**

1. **Upload PDF** - Select a receipt PDF file from your device
2. **AI Extraction** - Claude Sonnet 4.5 analyzes the document and extracts item details
3. **Side-by-Side Comparison** - Review extracted data alongside the original receipt
4. **Edit & Confirm** - Correct any discrepancies before saving

**Supported Format:** PDF files only (max 10MB)

**Extracted Data Format:**

| Source | Type    | Item Code | Price     | Item Name (AI Extracted) |
| :----- | :------ | :-------- | :-------- | :----------------------- |
| Costco | WEEKLY  | ORG SPA   | $12.99    | Organic Spaghetti        |
| H-Mart | MONTHLY | SWON      | $4.49     | Sweet Onions             |
| IRS    | TAX     | FED TAX   | $1,500.00 | Federal Income Tax       |

> **Note**: Items from the same receipt are automatically grouped together with the same receipt number for easy tracking.

> **Note**: Receipt processing is stateless - receipt files are not stored after processing.

## API Endpoints

### Budgets

| Method   | Endpoint            | Description         |
| -------- | ------------------- | ------------------- |
| `GET`    | `/api/budgets`      | List all budgets    |
| `POST`   | `/api/budgets`      | Create a new budget |
| `GET`    | `/api/budgets/{id}` | Get budget by ID    |
| `PUT`    | `/api/budgets/{id}` | Update budget       |
| `DELETE` | `/api/budgets/{id}` | Delete budget       |

### Expected Expenses

| Method   | Endpoint                      | Description                                                         |
| -------- | ----------------------------- | ------------------------------------------------------------------- |
| `GET`    | `/api/expected-expenses`      | List expected expenses (supports `?type=WEEKLY` or `?type=MONTHLY`) |
| `POST`   | `/api/expected-expenses`      | Create a new expected expense                                       |
| `GET`    | `/api/expected-expenses/{id}` | Get expected expense by ID                                          |
| `PUT`    | `/api/expected-expenses/{id}` | Update expected expense                                             |
| `DELETE` | `/api/expected-expenses/{id}` | Delete expected expense                                             |

### Actual Expenses

| Method   | Endpoint                                   | Description                       |
| -------- | ------------------------------------------ | --------------------------------- |
| `GET`    | `/api/actual-expenses`                     | List actual expenses              |
| `POST`   | `/api/actual-expenses`                     | Create new actual expense         |
| `GET`    | `/api/actual-expenses/next-receipt-number` | Get next available receipt number |
| `GET`    | `/api/actual-expenses/summary`             | Get monthly expense summary       |
| `GET`    | `/api/actual-expenses/{id}`                | Get actual expense by ID          |
| `PUT`    | `/api/actual-expenses/{id}`                | Update actual expense             |
| `DELETE` | `/api/actual-expenses/{id}`                | Delete actual expense             |

### Receipt Processing

| Method | Endpoint                | Description                 |
| ------ | ----------------------- | --------------------------- |
| `POST` | `/api/receipts/process` | Process receipt PDF with AI |

**Request Format:**

- Content-Type: `multipart/form-data`
- Form field: `document` (the PDF file)
- Max file size: 10MB
- Supported format: **PDF only** (JPEG, PNG not supported)

### Notifications

| Method | Endpoint                           | Description                          |
| ------ | ---------------------------------- | ------------------------------------ |
| `GET`  | `/api/notifications/budget-status` | Get current budget status and alerts |

## Database Schema

The application uses SQLite with three main tables:

### `budget_limits`

Stores monthly budget configurations.

| Column                 | Type     | Description                                   |
| ---------------------- | -------- | --------------------------------------------- |
| id                     | INTEGER  | Primary key                                   |
| month                  | INTEGER  | Month (1-12)                                  |
| year                   | INTEGER  | Year                                          |
| amount                 | REAL     | Budget limit amount                           |
| notification_threshold | REAL     | Notification threshold (0.0-1.0), default 0.8 |
| created_at             | DATETIME | Record creation timestamp                     |
| updated_at             | DATETIME | Last update timestamp                         |

> **Note**: A unique constraint exists on `(month, year)` to ensure only one budget per month.

### `expected_expenses`

Stores planned recurring expense items.

| Column          | Type     | Description                |
| --------------- | -------- | -------------------------- |
| id              | INTEGER  | Primary key                |
| item_name       | TEXT     | Item name                  |
| source          | TEXT     | Store/vendor name          |
| expected_amount | REAL     | Expected amount            |
| expense_type    | TEXT     | Frequency (WEEKLY/MONTHLY) |
| created_at      | DATETIME | Record creation timestamp  |
| updated_at      | DATETIME | Last update timestamp      |

### `actual_expenses`

Stores actual expense records from receipts.

| Column              | Type     | Description                                 |
| ------------------- | -------- | ------------------------------------------- |
| id                  | INTEGER  | Primary key                                 |
| item_name           | TEXT     | Item name                                   |
| source              | TEXT     | Store/vendor name                           |
| actual_amount       | REAL     | Actual amount paid                          |
| expense_type        | TEXT     | Category (WEEKLY/MONTHLY/MISC/TAX)          |
| item_code           | TEXT     | Optional short code                         |
| expected_expense_id | INTEGER  | Foreign key to expected_expenses (nullable) |
| receipt_date        | DATE     | Date on receipt                             |
| receipt_number      | INTEGER  | Receipt grouping number                     |
| month               | INTEGER  | Month (1-12)                                |
| year                | INTEGER  | Year                                        |
| created_at          | DATETIME | Record creation timestamp                   |
| updated_at          | DATETIME | Last update timestamp                       |

## Development

### Backend Commands

```bash
cd backend

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Build binary
go build -o budget-tracker ./cmd/server

# Run with race detector (development)
go run -race ./cmd/server
```

### Frontend Commands

```bash
cd frontend

# Type checking
npm run check

# Linting
npm run lint

# Build for production
npm run build

# Preview production build
npm run preview
```

### Code Quality

- **Backend**: Follow standard Go conventions, use `go fmt` and `go vet`
- **Frontend**: TypeScript strict mode enabled, ESLint + Prettier configured

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0) - see the [LICENSE](LICENSE) file for details.

---

**Built with Go and SvelteKit**
