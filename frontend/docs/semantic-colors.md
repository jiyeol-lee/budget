# Semantic Colors Documentation

This document describes the semantic color system used in the Budget Tracker frontend application.

## Overview

We use semantic color naming instead of primitive Tailwind colors (like `bg-blue-500`). This approach provides:

- **Consistency**: Colors have meaning throughout the app
- **Maintainability**: Change colors in one place
- **Accessibility**: Easier to ensure proper contrast
- **Theming**: Simpler to implement dark/light themes
- **Tailwind Alignment**: Perfect consistency with Tailwind's default color palette

## Tailwind CSS v4 Primitive Variables

Our semantic color system is built on top of **Tailwind CSS v4 primitive color variables**. Instead of defining raw hex values, we reference Tailwind's built-in color primitives using the `var(--color-xxx)` syntax.

### Benefits of This Approach

1. **Perfect Consistency**: Colors are guaranteed to match Tailwind's default palette
2. **Easier Maintenance**: No need to look up or copy hex values manually
3. **Future-Proof**: Automatically benefits from any Tailwind color updates
4. **Reduced Errors**: Eliminates typos in hex color codes

### How It Works

Instead of defining colors with raw hex values:

```css
/* Old approach - direct hex values */
--color-primary: #3b82f6;
--color-primary-hover: #2563eb;
```

We now reference Tailwind's primitive variables:

```css
/* New approach - Tailwind CSS v4 primitives */
--color-primary: var(--color-blue-500);
--color-primary-hover: var(--color-blue-600);
```

> **Note**: The semantic variable names (like `--color-primary`) remain unchanged in your components. Only the underlying value definitions have changed to use Tailwind primitives. This means **no changes are required in your component code**.

## Color Categories

### Text Colors

We use 4 variants of text colors for visual hierarchy:

| Class                  | CSS Variable              | Usage                          |
| ---------------------- | ------------------------- | ------------------------------ |
| `text-text-primary`    | `--color-text-primary`    | Main content, headings         |
| `text-text-secondary`  | `--color-text-secondary`  | Supporting text, descriptions  |
| `text-text-tertiary`   | `--color-text-tertiary`   | Less important text, hints     |
| `text-text-quaternary` | `--color-text-quaternary` | Least prominent text, disabled |

**Text on colored backgrounds:**

| Class                          | CSS Variable                      | Usage                               |
| ------------------------------ | --------------------------------- | ----------------------------------- |
| `text-text-on-primary`         | `--color-text-on-primary`         | Text on primary color background    |
| `text-text-on-success`         | `--color-text-on-success`         | Text on success color background    |
| `text-text-on-warning`         | `--color-text-on-warning`         | Text on warning color background    |
| `text-text-on-danger`          | `--color-text-on-danger`          | Text on danger color background     |
| `text-text-on-info`            | `--color-text-on-info`            | Text on info color background       |
| `text-text-inverted`           | `--color-text-inverted`           | Text on light backgrounds           |
| `text-text-inverted-secondary` | `--color-text-inverted-secondary` | Secondary text on light backgrounds |

### Brand / Primary Colors

| Class                     | CSS Variable                 | Usage                          |
| ------------------------- | ---------------------------- | ------------------------------ |
| `bg-primary`              | `--color-primary`            | Primary buttons, active states |
| `bg-primary-hover`        | `--color-primary-hover`      | Hover state for primary        |
| `bg-primary-active`       | `--color-primary-active`     | Active/pressed state           |
| `bg-primary-light`        | `--color-primary-light`      | Light variant for backgrounds  |
| `bg-primary-dark`         | `--color-primary-dark`       | Dark variant                   |
| `text-primary`            | `--color-primary`            | Primary colored text           |
| `text-primary-foreground` | `--color-primary-foreground` | Text on primary background     |
| `border-primary`          | `--color-primary`            | Primary colored borders        |

### Status Colors

#### Success

| Class                     | CSS Variable                 | Usage                         |
| ------------------------- | ---------------------------- | ----------------------------- |
| `bg-success`              | `--color-success`            | Success states, confirmations |
| `bg-success-hover`        | `--color-success-hover`      | Success hover state           |
| `bg-success-active`       | `--color-success-active`     | Success active state          |
| `bg-success-light`        | `--color-success-light`      | Light success backgrounds     |
| `bg-success-dark`         | `--color-success-dark`       | Dark success variant          |
| `text-success`            | `--color-success`            | Success text                  |
| `text-success-foreground` | `--color-success-foreground` | Text on success background    |
| `border-success`          | `--color-border-success`     | Success borders               |

#### Warning

| Class                     | CSS Variable                 | Usage                      |
| ------------------------- | ---------------------------- | -------------------------- |
| `bg-warning`              | `--color-warning`            | Warning states, cautions   |
| `bg-warning-hover`        | `--color-warning-hover`      | Warning hover state        |
| `bg-warning-active`       | `--color-warning-active`     | Warning active state       |
| `bg-warning-light`        | `--color-warning-light`      | Light warning backgrounds  |
| `bg-warning-dark`         | `--color-warning-dark`       | Dark warning variant       |
| `text-warning`            | `--color-warning`            | Warning text               |
| `text-warning-foreground` | `--color-warning-foreground` | Text on warning background |
| `border-warning`          | `--color-warning`            | Warning borders            |

#### Danger / Error

| Class                    | CSS Variable                | Usage                             |
| ------------------------ | --------------------------- | --------------------------------- |
| `bg-danger`              | `--color-danger`            | Error states, destructive actions |
| `bg-danger-hover`        | `--color-danger-hover`      | Danger hover state                |
| `bg-danger-active`       | `--color-danger-active`     | Danger active state               |
| `bg-danger-light`        | `--color-danger-light`      | Light danger backgrounds          |
| `bg-danger-dark`         | `--color-danger-dark`       | Dark danger variant               |
| `text-danger`            | `--color-danger`            | Error text                        |
| `text-danger-foreground` | `--color-danger-foreground` | Text on danger background         |
| `border-danger`          | `--color-border-error`      | Error borders                     |

#### Info

| Class                  | CSS Variable              | Usage                   |
| ---------------------- | ------------------------- | ----------------------- |
| `bg-info`              | `--color-info`            | Informational states    |
| `bg-info-hover`        | `--color-info-hover`      | Info hover state        |
| `bg-info-active`       | `--color-info-active`     | Info active state       |
| `bg-info-light`        | `--color-info-light`      | Light info backgrounds  |
| `bg-info-dark`         | `--color-info-dark`       | Dark info variant       |
| `text-info`            | `--color-info`            | Info text               |
| `text-info-foreground` | `--color-info-foreground` | Text on info background |
| `border-info`          | `--color-info`            | Info borders            |

### Surface / Background Colors

| Class                 | CSS Variable               | Usage                                 |
| --------------------- | -------------------------- | ------------------------------------- |
| `bg-surface`          | `--color-surface`          | Card backgrounds, sections            |
| `bg-surface-light`    | `--color-surface-light`    | Lighter surface variant               |
| `bg-surface-lighter`  | `--color-surface-lighter`  | Even lighter surface                  |
| `bg-surface-dark`     | `--color-surface-dark`     | Darker surface variant                |
| `bg-surface-darker`   | `--color-surface-darker`   | Darkest surface                       |
| `bg-surface-elevated` | `--color-surface-elevated` | Elevated elements (modals, dropdowns) |
| `bg-surface-overlay`  | `--color-surface-overlay`  | Modal/dialog overlays                 |
| `bg-background`       | `--color-background`       | Main app background                   |
| `bg-background-alt`   | `--color-background-alt`   | Alternative background                |

### Border Colors

| Class                   | CSS Variable             | Usage                 |
| ----------------------- | ------------------------ | --------------------- |
| `border-border`         | `--color-border`         | Default borders       |
| `border-border-light`   | `--color-border-light`   | Lighter borders       |
| `border-border-dark`    | `--color-border-dark`    | Darker borders        |
| `border-border-focus`   | `--color-border-focus`   | Focus state borders   |
| `border-border-error`   | `--color-border-error`   | Error state borders   |
| `border-border-success` | `--color-border-success` | Success state borders |

### Form / Input Colors

| Class                                | CSS Variable                | Usage                     |
| ------------------------------------ | --------------------------- | ------------------------- |
| `bg-input-bg`                        | `--color-input-bg`          | Input field background    |
| `border-input-border`                | `--color-input-border`      | Input field border        |
| `focus:border-input-focus`           | `--color-input-focus`       | Input focus border        |
| `placeholder:text-input-placeholder` | `--color-input-placeholder` | Placeholder text          |
| `bg-input-disabled`                  | `--color-input-disabled`    | Disabled input background |

### Budget Status Colors

| Class                  | CSS Variable                | Usage                     |
| ---------------------- | --------------------------- | ------------------------- |
| `bg-budget-safe`       | `--color-budget-safe`       | Under budget indicator    |
| `bg-budget-safe-bg`    | `--color-budget-safe-bg`    | Safe status background    |
| `bg-budget-warning`    | `--color-budget-warning`    | Near limit indicator      |
| `bg-budget-warning-bg` | `--color-budget-warning-bg` | Warning status background |
| `bg-budget-danger`     | `--color-budget-danger`     | Over budget indicator     |
| `bg-budget-danger-bg`  | `--color-budget-danger-bg`  | Danger status background  |
| `bg-budget-over`       | `--color-budget-over`       | Exceeded budget indicator |
| `bg-budget-over-bg`    | `--color-budget-over-bg`    | Over status background    |
| `text-budget-safe`     | `--color-budget-safe`       | Safe status text          |
| `text-budget-warning`  | `--color-budget-warning`    | Warning status text       |
| `text-budget-danger`   | `--color-budget-danger`     | Danger status text        |
| `text-budget-over`     | `--color-budget-over`       | Over budget text          |

### Table Colors

| Class                      | CSS Variable              | Usage                      |
| -------------------------- | ------------------------- | -------------------------- |
| `bg-table-header`          | `--color-table-header`    | Table header background    |
| `bg-table-row`             | `--color-table-row`       | Table row background       |
| `bg-table-row-alt`         | `--color-table-row-alt`   | Alternating row background |
| `hover:bg-table-row-hover` | `--color-table-row-hover` | Row hover state            |
| `border-table-border`      | `--color-table-border`    | Table borders              |

### Navigation Colors

| Class                       | CSS Variable                 | Usage                      |
| --------------------------- | ---------------------------- | -------------------------- |
| `bg-nav-bg`                 | `--color-nav-bg`             | Navigation background      |
| `text-nav-item`             | `--color-nav-item`           | Navigation item text       |
| `hover:text-nav-item-hover` | `--color-nav-item-hover`     | Nav item hover text        |
| `text-nav-item-active`      | `--color-nav-item-active`    | Active nav item text       |
| `bg-nav-item-active-bg`     | `--color-nav-item-active-bg` | Active nav item background |

### Interactive States

| Class                   | CSS Variable            | Usage                |
| ----------------------- | ----------------------- | -------------------- |
| `hover:bg-hover`        | `--color-hover`         | Generic hover state  |
| `active:bg-active`      | `--color-active`        | Active/pressed state |
| `focus:ring-focus-ring` | `--color-focus-ring`    | Focus ring           |
| `bg-disabled`           | `--color-disabled`      | Disabled background  |
| `text-disabled-text`    | `--color-disabled-text` | Disabled text        |

### Loading / Skeleton

| Class               | CSS Variable             | Usage                        |
| ------------------- | ------------------------ | ---------------------------- |
| `bg-skeleton`       | `--color-skeleton`       | Skeleton loading background  |
| `bg-skeleton-shine` | `--color-skeleton-shine` | Skeleton animation highlight |

## Usage Examples

### Button Examples

```svelte
<!-- Primary Button -->
<button class="bg-primary hover:bg-primary-hover text-text-on-primary rounded px-4 py-2">
	Save
</button>

<!-- Danger Button -->
<button class="bg-danger hover:bg-danger-hover text-text-on-danger rounded px-4 py-2">
	Delete
</button>

<!-- Secondary/Ghost Button -->
<button class="border-border text-text-secondary hover:bg-hover rounded border px-4 py-2">
	Cancel
</button>

<!-- Success Button -->
<button class="bg-success hover:bg-success-hover text-text-on-success rounded px-4 py-2">
	Confirm
</button>
```

### Card Example

```svelte
<div class="bg-surface border-border rounded-lg border p-4">
	<h2 class="text-text-primary text-lg font-semibold">Card Title</h2>
	<p class="text-text-secondary mt-2">Card description goes here.</p>
	<span class="text-text-tertiary text-sm">Additional hint text</span>
</div>
```

### Status Message Example

```svelte
<!-- Success Message -->
<div class="bg-success-light border-border-success text-success rounded border p-4">
	Operation completed successfully!
</div>

<!-- Error Message -->
<div class="bg-danger-light border-border-error text-danger rounded border p-4">
	An error occurred. Please try again.
</div>

<!-- Warning Message -->
<div class="bg-warning-light border-warning text-warning rounded border p-4">
	Please review your input before proceeding.
</div>

<!-- Info Message -->
<div class="bg-info-light border-info text-info rounded border p-4">
	Here's some helpful information.
</div>
```

### Form Input Example

```svelte
<label class="block">
	<span class="text-text-secondary text-sm">Email</span>
	<input
		type="email"
		class="bg-input-bg border-input-border text-text-primary placeholder:text-input-placeholder focus:border-input-focus focus:ring-focus-ring w-full
           rounded border
           px-3 py-2 focus:ring-1"
		placeholder="Enter your email..."
	/>
</label>
```

### Budget Status Example

```svelte
<!-- Safe Status -->
<div class="bg-budget-safe-bg border-budget-safe border-l-4 p-4">
	<span class="text-budget-safe font-semibold">On Track</span>
	<p class="text-text-secondary">You're within your budget.</p>
</div>

<!-- Warning Status -->
<div class="bg-budget-warning-bg border-budget-warning border-l-4 p-4">
	<span class="text-budget-warning font-semibold">Warning</span>
	<p class="text-text-secondary">Approaching your budget limit.</p>
</div>

<!-- Danger Status -->
<div class="bg-budget-danger-bg border-budget-danger border-l-4 p-4">
	<span class="text-budget-danger font-semibold">Critical</span>
	<p class="text-text-secondary">Very close to budget limit.</p>
</div>

<!-- Over Budget -->
<div class="bg-budget-over-bg border-budget-over border-l-4 p-4">
	<span class="text-budget-over font-semibold">Over Budget</span>
	<p class="text-text-secondary">You've exceeded your budget.</p>
</div>
```

### Table Example

```svelte
<table class="w-full">
	<thead class="bg-table-header">
		<tr>
			<th class="text-text-primary border-table-border border-b p-3 text-left">Name</th>
			<th class="text-text-primary border-table-border border-b p-3 text-left">Amount</th>
		</tr>
	</thead>
	<tbody>
		<tr class="bg-table-row hover:bg-table-row-hover">
			<td class="text-text-primary border-table-border border-b p-3">Groceries</td>
			<td class="text-text-secondary border-table-border border-b p-3">$150.00</td>
		</tr>
		<tr class="bg-table-row-alt hover:bg-table-row-hover">
			<td class="text-text-primary border-table-border border-b p-3">Utilities</td>
			<td class="text-text-secondary border-table-border border-b p-3">$85.00</td>
		</tr>
	</tbody>
</table>
```

### Skeleton Loading Example

```svelte
<div class="animate-skeleton">
	<div class="bg-skeleton mb-2 h-4 w-3/4 rounded"></div>
	<div class="bg-skeleton h-4 w-1/2 rounded"></div>
</div>
```

### Navigation Example

```svelte
<nav class="bg-nav-bg p-4">
	<a href="/" class="text-nav-item hover:text-nav-item-hover">Dashboard</a>
	<a href="/budget" class="text-nav-item-active bg-nav-item-active-bg rounded px-3 py-2">
		Budget
	</a>
</nav>
```

## Adding New Colors

To add new semantic colors, update `src/app.css` using Tailwind CSS v4 primitive variables:

```css
@theme {
	/* Add your new semantic color using Tailwind primitives */
	--color-new-semantic: var(--color-indigo-500);
	--color-new-semantic-hover: var(--color-indigo-600);
	--color-new-semantic-light: var(--color-indigo-100);
}
```

### Available Tailwind Primitives

You can reference any color from [Tailwind's default color palette](https://tailwindcss.com/docs/colors):

```css
/* Example primitives available */
var(--color-slate-500)
var(--color-gray-500)
var(--color-zinc-500)
var(--color-red-500)
var(--color-orange-500)
var(--color-amber-500)
var(--color-yellow-500)
var(--color-lime-500)
var(--color-green-500)
var(--color-emerald-500)
var(--color-teal-500)
var(--color-cyan-500)
var(--color-sky-500)
var(--color-blue-500)
var(--color-indigo-500)
var(--color-violet-500)
var(--color-purple-500)
var(--color-fuchsia-500)
var(--color-pink-500)
var(--color-rose-500)
/* Each color has shades from 50 to 950 */
```

Then use in components (same as before):

```svelte
<div class="bg-new-semantic hover:bg-new-semantic-hover">Content</div>

<span class="text-new-semantic">Colored text</span>

<div class="border-new-semantic border">Bordered element</div>
```

## Color Values Reference

Below are the semantic colors defined in our theme. The **Tailwind Primitive** column shows which Tailwind CSS v4 variable is used, and the **Hex Value** is shown for reference only (these map to Tailwind's default palette).

### Primary Colors

| Variable                 | Tailwind Primitive      | Hex Value |
| ------------------------ | ----------------------- | --------- |
| `--color-primary`        | `var(--color-blue-500)` | `#3b82f6` |
| `--color-primary-hover`  | `var(--color-blue-600)` | `#2563eb` |
| `--color-primary-active` | `var(--color-blue-700)` | `#1d4ed8` |
| `--color-primary-light`  | `var(--color-blue-100)` | `#dbeafe` |
| `--color-primary-dark`   | `var(--color-blue-800)` | `#1e40af` |

### Status Colors

| Variable          | Tailwind Primitive       | Hex Value |
| ----------------- | ------------------------ | --------- |
| `--color-success` | `var(--color-green-500)` | `#22c55e` |
| `--color-warning` | `var(--color-amber-500)` | `#f59e0b` |
| `--color-danger`  | `var(--color-red-500)`   | `#ef4444` |
| `--color-info`    | `var(--color-blue-500)`  | `#3b82f6` |

### Surface Colors

| Variable                | Tailwind Primitive      | Hex Value |
| ----------------------- | ----------------------- | --------- |
| `--color-surface`       | `var(--color-gray-800)` | `#1f2937` |
| `--color-surface-light` | `var(--color-gray-700)` | `#374151` |
| `--color-background`    | `var(--color-gray-900)` | `#111827` |

### Text Colors

| Variable                  | Tailwind Primitive      | Hex Value |
| ------------------------- | ----------------------- | --------- |
| `--color-text-primary`    | `var(--color-white)`    | `#ffffff` |
| `--color-text-secondary`  | `var(--color-gray-400)` | `#9ca3af` |
| `--color-text-tertiary`   | `var(--color-gray-500)` | `#6b7280` |
| `--color-text-quaternary` | `var(--color-gray-600)` | `#4b5563` |

## Best Practices

1. **Never use primitive colors** like `bg-blue-500` - always use semantic names
2. **Use text hierarchy** - primary for main content, secondary for supporting, tertiary for hints
3. **Match text to background** - use `text-on-*` variants when on colored backgrounds
4. **Use status colors consistently** - success for positive, warning for caution, danger for errors
5. **Maintain contrast** - ensure text is readable against backgrounds
6. **Use appropriate hover states** - always provide visual feedback for interactive elements
7. **Consider accessibility** - test color combinations for sufficient contrast ratios
8. **Use budget-specific colors** - for budget-related UI, use the `budget-*` color variants
9. **Keep surfaces consistent** - use surface colors for cards and elevated content
10. **Test in different contexts** - ensure colors work well in both light and dark themes
