# Button Component Documentation

This document describes the reusable Button component used in the Budget Tracker frontend application.

## Overview

The `Button` component is a flexible, accessible button element that supports multiple visual variants, loading states, and icon placement. It's built with Svelte 5 using runes and snippets for optimal performance and composability.

**Key Features:**

- **5 Visual Variants**: Primary, Secondary, Danger, Ghost, and Link
- **Loading State**: Built-in spinner with automatic disabled behavior
- **Icon Support**: Leading and trailing icon slots via Svelte snippets
- **Full Width Option**: Easily expand to container width
- **Accessible**: Focus-visible ring and proper disabled states
- **Semantic Colors**: Uses the project's semantic color system

## Props

| Prop         | Type                                                        | Default     | Description                                                                  |
| ------------ | ----------------------------------------------------------- | ----------- | ---------------------------------------------------------------------------- |
| `variant`    | `'primary' \| 'secondary' \| 'danger' \| 'ghost' \| 'link'` | `'primary'` | Visual style of the button                                                   |
| `type`       | `'button' \| 'submit' \| 'reset'`                           | `'button'`  | HTML button type attribute                                                   |
| `disabled`   | `boolean`                                                   | `false`     | Whether the button is disabled                                               |
| `loading`    | `boolean`                                                   | `false`     | Shows loading spinner and disables the button                                |
| `fullWidth`  | `boolean`                                                   | `false`     | Whether button should expand to full container width                         |
| `class`      | `string`                                                    | `''`        | Additional CSS classes to apply                                              |
| `aria-label` | `string`                                                    | `undefined` | Accessible label for screen readers, especially useful for icon-only buttons |
| `onclick`    | `(event: MouseEvent) => void`                               | `undefined` | Click event handler                                                          |
| `children`   | `Snippet`                                                   | `undefined` | Button content (text, elements)                                              |
| `leading`    | `Snippet`                                                   | `undefined` | Content rendered before the main content (e.g., icon)                        |
| `trailing`   | `Snippet`                                                   | `undefined` | Content rendered after the main content (e.g., icon)                         |

## Variants

### Primary (default)

Use for the main call-to-action on a page. Should be used sparingly—typically one per section.

```svelte
<Button>Save Budget</Button>
<Button variant="primary">Confirm</Button>
```

### Secondary

Use for secondary actions that complement the primary action.

```svelte
<Button variant="secondary">Cancel</Button>
<Button variant="secondary">View Details</Button>
```

### Danger

Use for destructive actions like delete or remove. Draws attention to potentially irreversible actions.

```svelte
<Button variant="danger">Delete Expense</Button>
<Button variant="danger">Remove Item</Button>
```

### Ghost

Use for tertiary actions or in toolbars where you need subtle buttons.

```svelte
<Button variant="ghost">Edit</Button>
<Button variant="ghost">More Options</Button>
```

### Link

Use when a button should appear as a text link but maintain button semantics.

```svelte
<Button variant="link">Learn more</Button>
<Button variant="link">View all expenses</Button>
```

## Usage Examples

### Basic Usage

```svelte
<script lang="ts">
	import { Button } from '$lib';
</script>

<Button onclick={() => console.log('Clicked!')}>Click Me</Button>
```

### With Loading State

The button automatically shows a spinner and becomes disabled when `loading` is true.

```svelte
<script lang="ts">
	import { Button } from '$lib';

	let isSubmitting = $state(false);

	async function handleSubmit() {
		isSubmitting = true;
		try {
			await saveData();
		} finally {
			isSubmitting = false;
		}
	}
</script>

<Button loading={isSubmitting} onclick={handleSubmit}>Save Changes</Button>
```

### With Icons

Use the `leading` and `trailing` snippets to add icons.

```svelte
<script lang="ts">
	import { Button } from '$lib';
	import { Plus, ArrowRight, Trash2 } from 'lucide-svelte';
</script>

<!-- Leading icon -->
<Button>
	{#snippet leading()}
		<Plus class="h-4 w-4" />
	{/snippet}
	Add Expense
</Button>

<!-- Trailing icon -->
<Button variant="secondary">
	Next Step
	{#snippet trailing()}
		<ArrowRight class="h-4 w-4" />
	{/snippet}
</Button>

<!-- Danger with icon -->
<Button variant="danger">
	{#snippet leading()}
		<Trash2 class="h-4 w-4" />
	{/snippet}
	Delete
</Button>
```

### Icon-Only Buttons

For buttons that contain only an icon without text, always provide an `aria-label` for accessibility.

```svelte
<script lang="ts">
	import { Button } from '$lib';
	import { Settings, Menu, X, Edit, Trash2 } from 'lucide-svelte';
</script>

<!-- Settings button -->
<Button variant="ghost" aria-label="Open settings">
	{#snippet leading()}
		<Settings class="h-4 w-4" />
	{/snippet}
</Button>

<!-- Menu toggle -->
<Button variant="ghost" aria-label="Toggle menu">
	{#snippet leading()}
		<Menu class="h-4 w-4" />
	{/snippet}
</Button>

<!-- Close button -->
<Button variant="ghost" aria-label="Close dialog">
	{#snippet leading()}
		<X class="h-4 w-4" />
	{/snippet}
</Button>

<!-- Action buttons in a table row -->
<div class="flex gap-1">
	<Button variant="ghost" aria-label="Edit expense">
		{#snippet leading()}
			<Edit class="h-4 w-4" />
		{/snippet}
	</Button>
	<Button variant="danger" aria-label="Delete expense">
		{#snippet leading()}
			<Trash2 class="h-4 w-4" />
		{/snippet}
	</Button>
</div>
```

### Full Width Button

Use `fullWidth` for buttons that should span their container, common in forms or mobile layouts.

```svelte
<Button fullWidth>Submit Form</Button>

<Button variant="secondary" fullWidth>Cancel</Button>
```

### Form Submit Button

Set `type="submit"` for buttons within forms.

```svelte
<form onsubmit={handleSubmit}>
	<input type="text" name="amount" />

	<Button type="submit" loading={isSubmitting}>Save Budget</Button>
</form>
```

### Disabled State

```svelte
<Button disabled>Cannot Click</Button>

<!-- Disabled with reason -->
<Button disabled={!isFormValid}>Submit</Button>
```

### Custom Classes

Add additional Tailwind classes for specific styling needs.

```svelte
<Button class="mt-4 sm:w-auto">Responsive Button</Button>
```

### Combined Example

A typical form footer with multiple button variants:

```svelte
<script lang="ts">
	import { Button } from '$lib';
	import { Save, X } from 'lucide-svelte';

	let isSaving = $state(false);
</script>

<div class="flex justify-end gap-3">
	<Button variant="ghost" onclick={handleCancel}>
		{#snippet leading()}
			<X class="h-4 w-4" />
		{/snippet}
		Cancel
	</Button>

	<Button loading={isSaving} onclick={handleSave}>
		{#snippet leading()}
			<Save class="h-4 w-4" />
		{/snippet}
		Save Changes
	</Button>
</div>
```

## Styling Details

### Base Styles

- `inline-flex items-center justify-center` - Flexbox for icon alignment
- `rounded-md px-4 py-2` - Standard padding and border radius (except `link` variant)
- `text-sm font-medium` - Typography
- `transition-colors` - Smooth color transitions
- `focus-visible:ring-2 focus-visible:ring-primary` - Accessible focus indicator

### Variant Styles

| Variant     | Background                                | Text                   | Border          |
| ----------- | ----------------------------------------- | ---------------------- | --------------- |
| `primary`   | `bg-primary` / `bg-primary-hover`         | `text-text-on-primary` | Transparent     |
| `secondary` | `bg-surface-light` / `bg-surface-lighter` | `text-text-secondary`  | `border-border` |
| `danger`    | `bg-danger` / `bg-danger-hover`           | `text-text-on-danger`  | Transparent     |
| `ghost`     | Transparent / `bg-surface-light`          | `text-text-secondary`  | Transparent     |
| `link`      | Transparent                               | `text-primary`         | None            |

### Disabled State

All variants receive `opacity-50` and `cursor-not-allowed` when disabled.

## Accessibility

- Uses native `<button>` element for proper semantics
- `disabled` attribute is properly applied when loading or explicitly disabled
- Focus-visible ring provides clear keyboard focus indication
- Loading spinner includes visual indication of processing state
- `aria-label` prop for providing accessible names to icon-only buttons

## Best Practices

1. **One primary per section** - Use primary buttons sparingly for main actions
2. **Clear labels** - Use action-oriented text like "Save Changes" not just "Submit"
3. **Loading feedback** - Always show loading state during async operations
4. **Icon consistency** - Use 16x16 (h-4 w-4) icons with the component
5. **Dangerous actions** - Use danger variant for destructive actions
6. **Keyboard accessibility** - Test with keyboard navigation
7. **Don't disable without reason** - If disabled, consider showing why
8. **Icon-only buttons need aria-label** - Always provide `aria-label` for buttons without visible text

## Related Documentation

- [Semantic Colors](./semantic-colors.md) - Color system used by button variants
- [i18n Conventions](./i18n-conventions.md) - For button label translations
