# DataTable Component Documentation

This document describes the reusable DataTable component used in the Budget Tracker frontend application.

## Overview

The `DataTable` component is a flexible, responsive data table that supports both flat and grouped data displays. It provides automatic responsive layouts with a desktop table view and mobile card view. Built with Svelte 5 using runes and snippets, it delegates all cell rendering to the parent component for maximum customization.

**Key Features:**

- **Dual Data Modes**: Flat data array or grouped data with headers/footers
- **Responsive Design**: Desktop table view with mobile card fallback
- **Loading Skeleton State**: Animated skeleton placeholders during data fetching
- **Empty State**: Configurable empty state with icon, title, and description
- **Custom Row Styling**: Function-based row class customization
- **Snippet-Based Rendering**: Full control over cell, mobile card, and group rendering
- **Semantic Colors**: Uses the project's semantic color system

## Interfaces

### Column

Defines the structure of a table column.

```typescript
interface Column {
	key: string; // Unique identifier for the column
	header: string; // Display text for column header
	headerClass?: string; // Additional CSS classes for header cell
	cellClass?: string; // Additional CSS classes for data cells
	hideOnMobile?: boolean; // Whether to hide column on mobile (unused in card view)
	skeletonWidth?: string; // Width class for skeleton placeholder (e.g., '3/4', '1/2', '16')
}
```

### Group

Defines a group of items for grouped data display.

```typescript
interface Group<T> {
	key: string | number; // Unique identifier for the group
	items: T[]; // Array of items in the group
}
```

## Props

| Prop               | Type                                     | Required | Default     | Description                                                       |
| ------------------ | ---------------------------------------- | -------- | ----------- | ----------------------------------------------------------------- |
| `data`             | `T[]`                                    | No\*     | `undefined` | Flat array of data items (use either `data` or `groups`)          |
| `groups`           | `Group<T>[]`                             | No\*     | `undefined` | Grouped data with key and items (use either `data` or `groups`)   |
| `columns`          | `Column[]`                               | Yes      | -           | Column definitions for the table                                  |
| `keyField`         | `keyof T`                                | Yes      | -           | Property name to use as unique key for each row                   |
| `loading`          | `boolean`                                | No       | `false`     | Whether to show skeleton loading state                            |
| `skeletonRows`     | `number`                                 | No       | `5`         | Number of skeleton rows to display when loading                   |
| `emptyTitle`       | `string`                                 | Yes      | -           | Title text shown when table is empty                              |
| `emptyDescription` | `string`                                 | No       | `undefined` | Description text shown below empty title                          |
| `emptyIcon`        | `Snippet`                                | No       | `undefined` | Icon snippet rendered in empty state                              |
| `rowClass`         | `(item: T) => string`                    | No       | `undefined` | Function returning additional CSS classes for a row               |
| `cell`             | `Snippet<[{ item: T; column: Column }]>` | Yes      | -           | Snippet for rendering table cells                                 |
| `mobileCard`       | `Snippet<[{ item: T }]>`                 | Yes      | -           | Snippet for rendering mobile card view                            |
| `groupHeader`      | `Snippet<[{ group: Group<T> }]>`         | No       | `undefined` | Snippet for rendering group header (grouped mode only)            |
| `groupFooter`      | `Snippet<[{ group: Group<T> }]>`         | No       | `undefined` | Snippet for rendering group footer (grouped mode only)            |
| `tableFooter`      | `Snippet`                                | No       | `undefined` | Snippet for rendering content in `<tfoot>` element (desktop only) |

\*Note: You must provide either `data` or `groups`, but not both.

## Loading Skeleton State

The DataTable component supports a loading skeleton state that displays placeholder content while data is being fetched. This provides visual feedback to users and maintains layout stability during async operations.

### Overview

When `loading={true}`, the component renders animated skeleton placeholders instead of actual data:

- **Desktop**: Displays a skeleton table with header and row placeholders
- **Mobile**: Displays skeleton cards that match the mobile card layout

### Desktop Skeleton Table

The desktop skeleton view renders:

- Column headers with animated skeleton bars
- Configurable number of rows (via `skeletonRows` prop)
- Per-column skeleton widths (via `skeletonWidth` in column definition)

### Mobile Skeleton Cards

The mobile skeleton view renders card-style placeholders with:

- Title area placeholder (3/4 width)
- Status badge placeholder
- Two description line placeholders (1/2 and 1/3 width)
- Action button placeholders

### Skeleton Width Configuration

Control the width of skeleton placeholders per column using the `skeletonWidth` property in the column definition. This allows you to approximate the expected content width for a more realistic loading appearance.

**Available width values:**

| Value    | Tailwind Class | Description         |
| -------- | -------------- | ------------------- |
| `'full'` | `w-full`       | 100% width          |
| `'3/4'`  | `w-3/4`        | 75% width (default) |
| `'1/2'`  | `w-1/2`        | 50% width           |
| `'1/4'`  | `w-1/4`        | 25% width           |
| `'16'`   | `w-16`         | Fixed 4rem (64px)   |
| `'20'`   | `w-20`         | Fixed 5rem (80px)   |
| `'24'`   | `w-24`         | Fixed 6rem (96px)   |

### Example Usage

```svelte
<script lang="ts">
	import { DataTable } from '$lib';

	let isLoading = $state(true);

	const columns = [
		{ key: 'name', header: 'Name', skeletonWidth: '3/4' },
		{ key: 'status', header: 'Status', skeletonWidth: '16' },
		{ key: 'amount', header: 'Amount', skeletonWidth: '20' },
		{ key: 'actions', header: 'Actions', skeletonWidth: '24', headerClass: 'text-right' }
	];

	// Simulate data fetching
	async function loadData() {
		isLoading = true;
		const response = await fetch('/api/expenses');
		expenses = await response.json();
		isLoading = false;
	}
</script>

<DataTable
	loading={isLoading}
	skeletonRows={5}
	data={expenses}
	{columns}
	keyField="id"
	emptyTitle="No expenses found"
>
	{#snippet cell({ item, column })}
		<!-- Cell rendering -->
	{/snippet}

	{#snippet mobileCard({ item })}
		<!-- Mobile card rendering -->
	{/snippet}
</DataTable>
```

### State Priority

The DataTable renders states in the following priority order:

1. **Loading** (`loading={true}`) - Shows skeleton placeholders
2. **Empty** (no data) - Shows empty state message
3. **Data** - Shows actual table content

This means that while loading, the empty state will not be shown even if there is no data yet.

## Snippets

### cell (required)

Renders the content of each table cell on desktop. Receives the item and column context.

```svelte
{#snippet cell({ item, column })}
	{#if column.key === 'name'}
		<span class="font-medium">{item.name}</span>
	{:else if column.key === 'amount'}
		<span>{formatCurrency(item.amount)}</span>
	{/if}
{/snippet}
```

### mobileCard (required)

Renders the mobile card view for each item. Full control over card layout.

```svelte
{#snippet mobileCard({ item })}
	<div class="bg-surface border-border rounded-lg border p-4">
		<h3 class="font-medium">{item.name}</h3>
		<p class="text-lg">{formatCurrency(item.amount)}</p>
	</div>
{/snippet}
```

### emptyIcon (optional)

Renders an icon in the empty state.

```svelte
{#snippet emptyIcon()}
	<ClipboardListIcon class="text-text-secondary h-12 w-12" />
{/snippet}
```

### groupHeader (optional, grouped mode only)

Renders a header row for each group.

```svelte
{#snippet groupHeader({ group })}
	<div class="bg-surface-dark px-4 py-3 font-medium">
		Group: {group.key}
	</div>
{/snippet}
```

### groupFooter (optional, grouped mode only)

Renders a footer row for each group.

```svelte
{#snippet groupFooter({ group })}
	<div class="bg-surface-light px-4 py-2 text-right">
		Total: {calculateGroupTotal(group.items)}
	</div>
{/snippet}
```

### tableFooter (optional)

Renders content inside a `<tfoot>` element at the end of the table. Use this for table-level footer rows like totals or special items (e.g., tax rows). Only renders on desktop table view.

```svelte
{#snippet tableFooter()}
	<tr class="bg-amber-950/30">
		<td colspan="3">Total</td>
		<td>$100.00</td>
	</tr>
{/snippet}
```

## Usage Examples

### Basic Flat Data Table

```svelte
<script lang="ts">
	import { DataTable } from '$lib';
	import { ClipboardListIcon } from 'lucide-svelte';

	interface Expense {
		id: number;
		name: string;
		amount: number;
		category: string;
	}

	let expenses: Expense[] = $state([
		{ id: 1, name: 'Groceries', amount: 150.0, category: 'Food' },
		{ id: 2, name: 'Gas', amount: 45.0, category: 'Transport' }
	]);

	const columns = [
		{ key: 'name', header: 'Name' },
		{ key: 'amount', header: 'Amount' },
		{ key: 'category', header: 'Category' }
	];

	function formatCurrency(amount: number): string {
		return new Intl.NumberFormat('en-US', {
			style: 'currency',
			currency: 'USD'
		}).format(amount);
	}
</script>

<DataTable
	data={expenses as unknown as Record<string, unknown>[]}
	{columns}
	keyField="id"
	emptyTitle="No expenses found"
	emptyDescription="Add your first expense to get started"
>
	{#snippet emptyIcon()}
		<ClipboardListIcon class="text-text-secondary h-12 w-12" />
	{/snippet}

	{#snippet cell({ item, column })}
		{@const expense = item as unknown as Expense}
		{#if column.key === 'name'}
			<span class="text-text-primary font-medium">{expense.name}</span>
		{:else if column.key === 'amount'}
			<span class="text-text-primary">{formatCurrency(expense.amount)}</span>
		{:else if column.key === 'category'}
			<span class="text-text-secondary">{expense.category}</span>
		{/if}
	{/snippet}

	{#snippet mobileCard({ item })}
		{@const expense = item as unknown as Expense}
		<div class="bg-surface border-border rounded-lg border p-4">
			<div class="flex items-start justify-between">
				<div>
					<h3 class="text-text-primary font-medium">{expense.name}</h3>
					<p class="text-text-secondary text-sm">{expense.category}</p>
				</div>
				<span class="text-text-primary font-semibold">
					{formatCurrency(expense.amount)}
				</span>
			</div>
		</div>
	{/snippet}
</DataTable>
```

### Grouped Data Table

```svelte
<script lang="ts">
	import { DataTable } from '$lib';
	import { ReceiptIcon } from 'lucide-svelte';

	interface ExpenseItem {
		id: number;
		name: string;
		amount: number;
	}

	interface ExpenseGroup {
		key: string | number;
		items: ExpenseItem[];
	}

	let groups: ExpenseGroup[] = $state([
		{
			key: 'Receipt #1',
			items: [
				{ id: 1, name: 'Milk', amount: 4.99 },
				{ id: 2, name: 'Bread', amount: 3.49 }
			]
		},
		{
			key: 'Receipt #2',
			items: [{ id: 3, name: 'Coffee', amount: 12.99 }]
		}
	]);

	const columns = [
		{ key: 'name', header: 'Item' },
		{ key: 'amount', header: 'Amount' }
	];

	function calculateTotal(items: ExpenseItem[]): number {
		return items.reduce((sum, item) => sum + item.amount, 0);
	}
</script>

<DataTable
	groups={groups as unknown as { key: string | number; items: Record<string, unknown>[] }[]}
	{columns}
	keyField="id"
	emptyTitle="No receipts found"
>
	{#snippet groupHeader({ group })}
		<div class="bg-surface-dark flex items-center gap-2 px-4 py-3">
			<ReceiptIcon class="h-5 w-5" />
			<span class="font-medium">{group.key}</span>
		</div>
	{/snippet}

	{#snippet groupFooter({ group })}
		{@const items = group.items as unknown as ExpenseItem[]}
		<div class="bg-surface-light px-4 py-2 text-right font-semibold">
			Total: ${calculateTotal(items).toFixed(2)}
		</div>
	{/snippet}

	{#snippet cell({ item, column })}
		{@const expense = item as unknown as ExpenseItem}
		{#if column.key === 'name'}
			<span>{expense.name}</span>
		{:else if column.key === 'amount'}
			<span>${expense.amount.toFixed(2)}</span>
		{/if}
	{/snippet}

	{#snippet mobileCard({ item })}
		{@const expense = item as unknown as ExpenseItem}
		<div class="bg-surface-light rounded-lg p-3">
			<div class="flex justify-between">
				<span>{expense.name}</span>
				<span class="font-medium">${expense.amount.toFixed(2)}</span>
			</div>
		</div>
	{/snippet}
</DataTable>
```

### Custom Row Styling

Use the `rowClass` prop to conditionally style rows based on item data.

```svelte
<script lang="ts">
	import { DataTable } from '$lib';

	interface Transaction {
		id: number;
		description: string;
		amount: number;
		type: 'income' | 'expense' | 'tax';
	}

	let transactions: Transaction[] = $state([...]);

	const columns = [
		{ key: 'description', header: 'Description' },
		{ key: 'amount', header: 'Amount' },
		{ key: 'type', header: 'Type' }
	];

	/**
	 * Apply custom styling for tax rows
	 */
	function getRowClass(item: Record<string, unknown>): string {
		const transaction = item as unknown as Transaction;
		if (transaction.type === 'tax') {
			return 'bg-amber-950/30 border-amber-800';
		}
		if (transaction.type === 'income') {
			return 'bg-green-950/20';
		}
		return '';
	}
</script>

<DataTable
	data={transactions as unknown as Record<string, unknown>[]}
	{columns}
	keyField="id"
	rowClass={getRowClass}
	emptyTitle="No transactions"
>
	{#snippet cell({ item, column })}
		<!-- Cell rendering -->
	{/snippet}

	{#snippet mobileCard({ item })}
		<!-- Mobile card rendering -->
	{/snippet}
</DataTable>
```

### Table Footer (e.g., Tax Row)

Use `tableFooter` to add footer rows that align with table columns.

```svelte
<DataTable data={items} {columns} keyField="id" emptyTitle="No items">
	{#snippet cell({ item, column })}
		<!-- cell content -->
	{/snippet}

	{#snippet mobileCard({ item })}
		<!-- mobile card content -->
	{/snippet}

	{#snippet tableFooter()}
		{#if taxItem}
			<tr class="border-t border-amber-800 bg-amber-950/30">
				<td class="px-4 py-4">Tax</td>
				<td class="px-4 py-4">-</td>
				<td class="px-4 py-4">${taxItem.amount}</td>
				<td class="px-4 py-4 text-right">
					<button>Edit</button>
				</td>
			</tr>
		{/if}
	{/snippet}
</DataTable>
```

### Column Visibility Control

Use `headerClass` and `cellClass` to control column visibility at different breakpoints.

```svelte
<script lang="ts">
	const columns = [
		{ key: 'name', header: 'Name' },
		{ key: 'email', header: 'Email' },
		// Hide on screens smaller than lg
		{
			key: 'phone',
			header: 'Phone',
			headerClass: 'hidden lg:table-cell',
			cellClass: 'hidden lg:table-cell'
		},
		// Hide on screens smaller than xl
		{
			key: 'address',
			header: 'Address',
			headerClass: 'hidden xl:table-cell',
			cellClass: 'hidden xl:table-cell'
		},
		// Right-align actions column
		{ key: 'actions', header: 'Actions', headerClass: 'text-right' }
	];
</script>
```

## Styling Details

### Base Styles

- Desktop table hidden on mobile (`hidden sm:block`)
- Mobile card view visible only on small screens (`sm:hidden`)
- Table uses full width with horizontal scroll on overflow
- Header row with `bg-surface-light` background
- Data rows with `hover:bg-surface-light` transition

### Row Styles

| State   | Class                                      |
| ------- | ------------------------------------------ |
| Default | `hover:bg-surface-light transition-colors` |
| Custom  | Base + value from `rowClass` function      |

### Cell Styles

| Element     | Class                                                                        |
| ----------- | ---------------------------------------------------------------------------- |
| Header cell | `text-text-secondary text-xs font-medium tracking-wider uppercase px-4 py-3` |
| Data cell   | `px-4 py-4 whitespace-nowrap`                                                |

### Empty State Styles

- Centered layout with `py-12 text-center`
- Icon container with `h-12 w-12`
- Title with `text-lg font-medium`
- Description with `text-sm`

## Best Practices

1. **Always provide both snippets** - `cell` and `mobileCard` are required for complete responsive behavior
2. **Type cast items** - Use `item as unknown as YourType` inside snippets for type safety
3. **Use descriptive keyField** - Choose a unique identifier field (e.g., `id`) for proper row tracking
4. **i18n for headers** - Use message functions for column headers to support internationalization
5. **Responsive column hiding** - Use `headerClass` and `cellClass` with Tailwind breakpoints for responsive tables
6. **Consistent card design** - Mirror desktop column data in mobile cards for a cohesive experience
7. **Empty states matter** - Provide helpful empty titles and descriptions with actionable guidance
8. **Group semantically** - When using grouped mode, ensure groups represent meaningful data relationships

## Accessibility

- Uses semantic `<table>`, `<thead>`, `<tbody>`, `<tr>`, `<th>`, and `<td>` elements
- Header cells use `scope="col"` for proper screen reader association
- Table has horizontal scroll with `overflow-x-auto` for small screens
- Focus and hover states are properly styled for keyboard navigation
- Mobile cards should include interactive elements with proper `aria-label` attributes
- Empty state messaging is visible and descriptive

## Related Documentation

- [Button Component](./button-component.md) - Often used for action columns
- [Semantic Colors](./semantic-colors.md) - Color system used by table styling
- [i18n Conventions](./i18n-conventions.md) - For column header translations
