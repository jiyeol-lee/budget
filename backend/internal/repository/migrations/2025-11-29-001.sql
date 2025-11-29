-- Migration: 2025-11-29-001
-- Description: Initial schema migration - consolidated from existing migrations
-- This file combines all table definitions from the original migrations.go

-- ============================================================================
-- Budget Limits Table
-- Stores monthly budget limits and notification thresholds
-- ============================================================================
CREATE TABLE IF NOT EXISTS budget_limits (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    month INTEGER NOT NULL,
    year INTEGER NOT NULL,
    amount REAL NOT NULL,
    notification_threshold REAL DEFAULT 0.8,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(month, year)
);

-- ============================================================================
-- Expected Expenses Table
-- Stores recurring expected expenses (weekly/monthly)
-- ============================================================================
CREATE TABLE IF NOT EXISTS expected_expenses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    item_name TEXT NOT NULL,
    source TEXT NOT NULL,
    expected_amount REAL NOT NULL,
    expense_type TEXT NOT NULL CHECK(expense_type IN ('weekly', 'monthly')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- Actual Expenses Table
-- Stores actual expense records with optional link to expected expenses
-- ============================================================================
CREATE TABLE IF NOT EXISTS actual_expenses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    item_name TEXT NOT NULL,
    source TEXT NOT NULL,
    actual_amount REAL NOT NULL,
    expense_type TEXT NOT NULL CHECK(expense_type IN ('weekly', 'monthly', 'misc', 'tax')),
    item_code TEXT,
    expected_expense_id INTEGER,
    receipt_date DATE DEFAULT (DATE('now')),
    month INTEGER NOT NULL,
    year INTEGER NOT NULL,
    receipt_number INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (expected_expense_id) REFERENCES expected_expenses(id) ON DELETE SET NULL
);

-- Indexes for actual_expenses table
CREATE INDEX IF NOT EXISTS idx_actual_expenses_month_year ON actual_expenses(year, month);
CREATE INDEX IF NOT EXISTS idx_actual_expenses_expected ON actual_expenses(expected_expense_id);
CREATE INDEX IF NOT EXISTS idx_actual_expenses_receipt_number ON actual_expenses(receipt_number);
