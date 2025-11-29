<script lang="ts">
	import type { Snippet } from 'svelte';
	import { LoaderCircle } from 'lucide-svelte';

	type ButtonVariant = 'primary' | 'secondary' | 'danger' | 'ghost' | 'link';
	type ButtonType = 'button' | 'submit' | 'reset';

	interface Props {
		variant?: ButtonVariant;
		type?: ButtonType;
		disabled?: boolean;
		loading?: boolean;
		fullWidth?: boolean;
		class?: string;
		'aria-label'?: string;
		onclick?: (event: MouseEvent) => void;
		children?: Snippet;
		leading?: Snippet;
		trailing?: Snippet;
	}

	let {
		variant = 'primary',
		type = 'button',
		disabled = false,
		loading = false,
		fullWidth = false,
		class: className = '',
		'aria-label': ariaLabel,
		onclick,
		children,
		leading,
		trailing
	}: Props = $props();

	// Computed: button is effectively disabled when loading or explicitly disabled
	let isDisabled = $derived(disabled || loading);

	// Variant-specific styles
	const variantStyles: Record<ButtonVariant, string> = {
		primary:
			'bg-primary hover:bg-primary-hover text-text-on-primary border border-transparent shadow-sm',
		secondary: 'bg-surface-light hover:bg-surface-lighter text-text-secondary border border-border',
		danger:
			'bg-danger hover:bg-danger-hover text-text-on-danger border border-transparent shadow-sm',
		ghost: 'bg-transparent hover:bg-surface-light text-text-secondary border border-transparent',
		link: 'bg-transparent text-primary hover:text-primary-light border-none'
	};

	// Base styles (excluding padding for link variant)
	let baseStyles = $derived(
		variant === 'link'
			? 'inline-flex items-center justify-center font-medium text-sm cursor-pointer transition-colors focus:outline-none focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2'
			: 'inline-flex items-center justify-center rounded-md px-4 py-2 text-sm font-medium cursor-pointer transition-colors focus:outline-none focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2'
	);

	// Disabled styles
	const disabledStyles = 'disabled:cursor-not-allowed disabled:opacity-50';

	// Full width style
	let widthStyle = $derived(fullWidth ? 'w-full' : '');

	// Combined classes
	let buttonClasses = $derived(
		[baseStyles, variantStyles[variant], disabledStyles, widthStyle, className]
			.filter(Boolean)
			.join(' ')
	);
</script>

<button {type} disabled={isDisabled} class={buttonClasses} aria-label={ariaLabel} {onclick}>
	{#if loading}
		<LoaderCircle class="mr-2 -ml-1 h-4 w-4 animate-spin" />
	{:else if leading}
		<span class="mr-2 -ml-1">
			{@render leading()}
		</span>
	{/if}

	{#if children}
		{@render children()}
	{/if}

	{#if trailing && !loading}
		<span class="-mr-1 ml-2">
			{@render trailing()}
		</span>
	{/if}
</button>
