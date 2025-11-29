# i18n Conventions Documentation

This document describes the internationalization (i18n) conventions used in the Budget Tracker frontend application.

## Overview

We use **Paraglide JS** for internationalization. It's the officially recommended i18n solution by Svelte, offering:

- Compiler-based optimization (tiny ~300B runtime)
- Full TypeScript support with autocomplete
- Svelte 5 compatibility
- Tree-shakable messages

## Setup

Translation files are located in:

- `messages/en.json` - English (source)
- `messages/ko.json` - Korean

Generated code is output to `src/lib/paraglide/` (gitignored).

## Key Naming Conventions

### Pattern: `{scope}_{context}_{element}`

We use **flat keys with snake_case** prefixes. This provides:

- Better tree-shaking
- IDE autocomplete support
- Clear organization

### Naming Categories

| Prefix          | Usage                    | Examples                                             |
| --------------- | ------------------------ | ---------------------------------------------------- |
| `common_`       | Shared/reusable strings  | `common_save`, `common_cancel`, `common_loading`     |
| `nav_`          | Navigation items         | `nav_dashboard`, `nav_budget`, `nav_receipts`        |
| `{page}_`       | Page-specific content    | `dashboard_title`, `budget_description`              |
| `{page}_form_`  | Form labels/placeholders | `budget_form_month`, `expense_form_amount`           |
| `action_`       | Action buttons/links     | `action_add_expense`, `action_process_receipt`       |
| `stats_`        | Statistics/metrics       | `stats_weekly_expected`, `stats_remaining_budget`    |
| `error_`        | Error messages           | `error_network`, `error_not_found`                   |
| `validation_`   | Form validation          | `validation_required`, `validation_invalid_amount`   |
| `toast_`        | Toast notifications      | `toast_success`, `toast_error`                       |
| `delete_`       | Delete confirmations     | `delete_confirm_title`, `delete_confirm_description` |
| `notification_` | System notifications     | `notification_budget_exceeded`                       |
| `expenses_`     | Expense-related terms    | `expenses_weekly`, `expenses_monthly`                |

### Key Naming Rules

1. **Use snake_case** for all keys

   - ✅ `budget_form_amount`
   - ❌ `budgetFormAmount`
   - ❌ `budget-form-amount`

2. **Start with scope/context**

   - ✅ `expected_expenses_title`
   - ❌ `title_expected_expenses`

3. **Be descriptive but concise**

   - ✅ `dashboard_no_budget_description`
   - ❌ `dashboard_the_text_shown_when_no_budget_is_set`

4. **Use consistent suffixes**

   - `_title` - Page/section titles
   - `_description` - Descriptions/subtitles
   - `_label` - Form labels
   - `_placeholder` - Input placeholders
   - `_button` - Button text (optional, often just action name)
   - `_empty` - Empty state titles
   - `_empty_description` - Empty state descriptions
   - `_add` - Add action text
   - `_edit` - Edit action text
   - `_total` - Total/sum labels

5. **Group related keys**
   ```json
   {
   	"budget_title": "Budget Settings",
   	"budget_description": "Configure your monthly budgets",
   	"budget_add": "Add Budget",
   	"budget_edit": "Edit Budget",
   	"budget_form_month": "Month",
   	"budget_form_year": "Year",
   	"budget_form_amount": "Budget Amount"
   }
   ```

## Usage in Components

### Importing Messages

```typescript
import * as m from '$lib/paraglide/messages';
```

### Using in Templates

```svelte
<script lang="ts">
	import * as m from '$lib/paraglide/messages';
</script>

<h1>{m.dashboard_title()}</h1>
<p>{m.dashboard_description()}</p>
<button>{m.common_save()}</button>
```

### With Parameters

Define in JSON:

```json
{
	"stats_of_budget": "of {amount} budget",
	"receipt_confirm_add_description": "This will add {count} item(s) to your actual expenses for the current month.",
	"stats_used": "{percentage}% used",
	"stats_remaining": "{amount} remaining",
	"stats_over_budget": "{amount} over budget"
}
```

Use in template:

```svelte
<p>{m.stats_of_budget({ amount: '$1,000' })}</p>
<p>{m.receipt_confirm_add_description({ count: 5 })}</p>
<p>{m.stats_used({ percentage: '75' })}</p>
<p>{m.stats_remaining({ amount: '$250' })}</p>
```

### With Svelte 5 Runes

```svelte
<script lang="ts">
	import * as m from '$lib/paraglide/messages';

	let amount = $state(1000);
	let message = $derived(m.stats_of_budget({ amount: `$${amount}` }));
</script>

<p>{message}</p>
```

## Adding New Translations

### 1. Add Key to Source File (en.json)

```json
{
	"new_feature_title": "New Feature",
	"new_feature_description": "This is a new feature."
}
```

### 2. Add to All Language Files

Add the same key to all language files (e.g., ko.json):

```json
{
	"new_feature_title": "새로운 기능",
	"new_feature_description": "새로운 기능입니다."
}
```

### 3. Run Build

```bash
npm run build
```

This regenerates `src/lib/paraglide/messages.js` with the new function.

### 4. Use in Component

```svelte
<script lang="ts">
	import * as m from '$lib/paraglide/messages';
</script>

<h1>{m.new_feature_title()}</h1>
```

## Complete Key Reference

### Common Keys

| Key                 | English       | Description      |
| ------------------- | ------------- | ---------------- |
| `common_loading`    | Loading...    | Loading state    |
| `common_processing` | Processing... | Processing state |
| `common_save`       | Save          | Save button      |
| `common_cancel`     | Cancel        | Cancel button    |
| `common_delete`     | Delete        | Delete button    |
| `common_edit`       | Edit          | Edit button      |
| `common_confirm`    | Confirm       | Confirm button   |
| `common_close`      | Close         | Close button     |
| `common_add`        | Add           | Add button       |
| `common_retry`      | Try again     | Retry/try again  |
| `common_clear`      | Clear         | Clear button     |
| `common_submit`     | Submit        | Submit button    |
| `common_back`       | Back          | Back button      |
| `common_next`       | Next          | Next button      |
| `common_yes`        | Yes           | Affirmative      |
| `common_no`         | No            | Negative         |
| `common_or`         | or            | Conjunction      |
| `common_and`        | and           | Conjunction      |

### Navigation Keys

| Key                     | English           | Description            |
| ----------------------- | ----------------- | ---------------------- |
| `nav_dashboard`         | Dashboard         | Dashboard link         |
| `nav_budget`            | Budget            | Budget link            |
| `nav_expected_expenses` | Expected Expenses | Expected Expenses link |
| `nav_actual_expenses`   | Actual Expenses   | Actual Expenses link   |
| `nav_receipts`          | Receipts          | Receipts link          |

### Page Title/Description Pattern

| Pattern                    | Example Keys                                       |
| -------------------------- | -------------------------------------------------- |
| `{page}_title`             | `dashboard_title`, `budget_title`, `receipt_title` |
| `{page}_description`       | `dashboard_description`, `budget_description`      |
| `{page}_empty`             | `expected_expenses_empty`, `actual_expenses_empty` |
| `{page}_empty_description` | `expected_expenses_empty_description`              |
| `{page}_add`               | `budget_add`, `expected_expenses_add`              |
| `{page}_edit`              | `budget_edit`, `expected_expenses_edit`            |

### Form Label Pattern

| Pattern                              | Example Keys                            |
| ------------------------------------ | --------------------------------------- |
| `{context}_form_{field}`             | `budget_form_month`, `budget_form_year` |
| `{context}_form_{field}_placeholder` | `expense_form_item_name_placeholder`    |
| `{context}_form_{field}_description` | `budget_form_threshold_description`     |

### Status/State Pattern

| Pattern                    | Example Keys                                      |
| -------------------------- | ------------------------------------------------- |
| `{context}_status_{state}` | `budget_status_on_track`, `budget_status_warning` |
| `{context}_{state}`        | `expenses_weekly`, `expenses_monthly`             |

### Action Pattern

| Pattern                | Example Keys                                            |
| ---------------------- | ------------------------------------------------------- |
| `action_{verb}_{noun}` | `action_add_expected_expense`, `action_process_receipt` |
| `action_{verb}`        | `action_set_monthly_budget`                             |

### Error/Validation Pattern

| Pattern             | Example Keys                                       |
| ------------------- | -------------------------------------------------- |
| `error_{type}`      | `error_network`, `error_not_found`, `error_server` |
| `validation_{rule}` | `validation_required`, `validation_invalid_amount` |

### Stats Pattern

| Pattern                     | Example Keys                                  |
| --------------------------- | --------------------------------------------- |
| `stats_{metric}`            | `stats_weekly_expected`, `stats_actual_spent` |
| `stats_{metric}_{modifier}` | `stats_remaining_budget`, `stats_over_budget` |

### Toast Pattern

| Pattern        | Example Keys                                                  |
| -------------- | ------------------------------------------------------------- |
| `toast_{type}` | `toast_success`, `toast_error`, `toast_warning`, `toast_info` |

### Delete Confirmation Pattern

| Pattern                    | Example Keys                                         |
| -------------------------- | ---------------------------------------------------- |
| `delete_confirm_{element}` | `delete_confirm_title`, `delete_confirm_description` |

## Existing Keys by Category

### Dashboard Page

```json
{
	"dashboard_title": "Dashboard",
	"dashboard_description": "Welcome to your budget tracker...",
	"dashboard_budget_status": "Budget Status",
	"dashboard_manage_budget": "Manage Budget",
	"dashboard_no_budget_title": "No budget set for this month",
	"dashboard_no_budget_description": "Set up a monthly budget...",
	"dashboard_set_budget": "Set Budget",
	"dashboard_quick_stats": "Quick Stats",
	"dashboard_quick_actions": "Quick Actions",
	"dashboard_tips_title": "Tips",
	"dashboard_tips_description": "Set a monthly budget..."
}
```

### Budget Page

```json
{
	"budget_title": "Budget Settings",
	"budget_description": "Configure your monthly budgets...",
	"budget_add": "Add Budget",
	"budget_edit": "Edit Budget",
	"budget_form_month": "Month",
	"budget_form_year": "Year",
	"budget_form_amount": "Budget Amount",
	"budget_form_threshold": "Notification Threshold (%)",
	"budget_form_threshold_description": "You'll be notified when...",
	"budget_status_on_track": "On Track",
	"budget_status_warning": "Warning",
	"budget_status_critical": "Critical",
	"budget_status_over": "Over Budget",
	"budget_no_budgets": "No budgets configured",
	"budget_no_budgets_description": "Add a budget to start...",
	"budget_tips_title": "Tips",
	"budget_tips_set_threshold": "Set a notification threshold..."
}
```

### Expenses (Shared)

```json
{
	"expenses_weekly": "Weekly",
	"expenses_monthly": "Monthly",
	"expenses_misc": "Misc",
	"expenses_all": "All",
	"expenses_type": "Type"
}
```

### Expected Expenses Page

```json
{
	"expected_expenses_title": "Expected Expenses",
	"expected_expenses_description": "Manage your recurring expenses...",
	"expected_expenses_add": "Add Expected Expense",
	"expected_expenses_edit": "Edit Expected Expense",
	"expected_expenses_empty": "No expected expenses",
	"expected_expenses_empty_description": "Add your recurring expenses...",
	"expected_expenses_total_weekly": "Total Weekly",
	"expected_expenses_total_monthly": "Total Monthly",
	"expected_expenses_total_misc": "Total Misc"
}
```

### Actual Expenses Page

```json
{
	"actual_expenses_title": "Actual Expenses",
	"actual_expenses_description": "Track your real spending...",
	"actual_expenses_add": "Add Actual Expense",
	"actual_expenses_edit": "Edit Actual Expense",
	"actual_expenses_empty": "No actual expenses",
	"actual_expenses_empty_description": "Start recording your actual spending...",
	"actual_expenses_total": "Total Spent",
	"actual_expenses_month": "Month",
	"actual_expenses_year": "Year"
}
```

### Expense Form Fields

```json
{
	"expense_form_item_name": "Item Name",
	"expense_form_item_name_placeholder": "e.g., Groceries, Rent...",
	"expense_form_source": "Source",
	"expense_form_source_placeholder": "e.g., Costco, Amazon...",
	"expense_form_amount": "Amount",
	"expense_form_amount_placeholder": "0.00",
	"expense_form_expense_type": "Expense Type",
	"expense_form_item_code": "Item Code",
	"expense_form_item_code_placeholder": "Optional identifier"
}
```

### Receipt Processing Page

```json
{
	"receipt_title": "Process Receipt",
	"receipt_description": "Upload a receipt PDF...",
	"receipt_upload": "Upload Receipt",
	"receipt_upload_description": "Drag and drop or click to upload",
	"receipt_upload_formats": "PDF files only",
	"receipt_process": "Process Receipt",
	"receipt_processing": "Processing receipt...",
	"receipt_review": "Review Extracted Items",
	"receipt_add_to_expenses": "Add to Expenses",
	"receipt_clear_start_over": "Clear & Start Over",
	"receipt_instructions_title": "Instructions",
	"receipt_instructions_1": "Upload a PDF file of your receipt",
	"receipt_instructions_2": "AI will extract items and prices",
	"receipt_instructions_3": "Review and edit extracted items",
	"receipt_instructions_4": "Add items to your actual expenses",
	"receipt_success": "Receipt processed successfully!",
	"receipt_error": "Failed to process receipt",
	"receipt_no_items": "No items extracted",
	"receipt_confirm_add": "Confirm Add to Expenses",
	"receipt_confirm_add_description": "This will add {count} item(s)..."
}
```

### Statistics

```json
{
	"stats_weekly_expected": "Weekly Expected",
	"stats_monthly_expected": "Monthly Expected",
	"stats_actual_spent": "Actual Spent",
	"stats_remaining_budget": "Remaining Budget",
	"stats_of_budget": "of {amount} budget",
	"stats_used": "{percentage}% used",
	"stats_remaining": "{amount} remaining",
	"stats_over_budget": "{amount} over budget"
}
```

### Notifications

```json
{
	"notification_budget_exceeded": "Budget exceeded!",
	"notification_budget_warning": "Approaching budget limit",
	"notification_budget_critical": "Critical: Very close to budget limit"
}
```

### Validation Messages

```json
{
	"validation_required": "This field is required",
	"validation_min_length": "Must be at least {min} characters",
	"validation_invalid_amount": "Please enter a valid amount",
	"validation_invalid_year": "Please enter a valid year",
	"validation_invalid_month": "Please select a valid month",
	"validation_positive_number": "Must be a positive number"
}
```

### Error Messages

```json
{
	"error_generic": "Something went wrong. Please try again.",
	"error_network": "Network error. Please check your connection.",
	"error_not_found": "Item not found",
	"error_server": "Server error. Please try again later."
}
```

## Best Practices

1. **Always add keys to source (en.json) first**
2. **Keep keys alphabetically sorted** within their scope
3. **Use parameters for dynamic content** instead of string concatenation
4. **Reuse common keys** - don't duplicate "Save" as `budget_save` and `expense_save`
5. **Write for translators** - provide context in the key name
6. **Avoid technical jargon** - use user-friendly language
7. **Test with long translations** - some languages expand significantly
8. **Run build after changes** - to regenerate message functions
9. **Keep translations in sync** - ensure all language files have the same keys
10. **Use consistent terminology** - "Budget" not "Monthly Limit" in one place and "Spending Cap" in another

## Troubleshooting

### Key not found

- Ensure the key exists in the JSON file
- Run `npm run build` to regenerate messages
- Check for typos in the key name
- Verify the import: `import * as m from '$lib/paraglide/messages'`

### TypeScript error

- The generated `messages.js` includes types
- Make sure `src/lib/paraglide` is in `.gitignore`
- Run build to generate fresh files
- Check that the function is being called: `m.key()` not `m.key`

### Parameters not working

- Ensure parameter name in JSON matches usage
- Parameters use `{paramName}` syntax in JSON
- Pass object with matching keys: `m.key({ paramName: value })`
- Parameter names are case-sensitive

### Missing translations

- Add keys to all language files (`en.json`, `ko.json`)
- Use the same key names across all files
- Check for JSON syntax errors (trailing commas, missing quotes)

### Changes not appearing

- Run `npm run build` after modifying message files
- Clear browser cache if needed
- Restart dev server: `npm run dev`
