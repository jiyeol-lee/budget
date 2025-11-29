<script lang="ts">
	type ValueColor = 'primary' | 'info' | 'purple' | 'warning';

	interface Props {
		title: string;
		value: string;
		subtitle?: string;
		valueColor?: ValueColor;
		class?: string;
	}

	let { title, value, subtitle, valueColor = 'primary', class: className = '' }: Props = $props();

	// Color mapping for value text
	const colorStyles: Record<ValueColor, string> = {
		primary: 'text-text-primary',
		info: 'text-info',
		purple: 'text-purple-400',
		warning: 'text-warning'
	};

	let valueColorClass = $derived(colorStyles[valueColor]);
</script>

<div class="bg-surface rounded-lg p-4 shadow sm:p-6 {className}">
	<h3 class="text-text-secondary text-sm font-medium">{title}</h3>
	<p class="mt-2 text-2xl font-bold {valueColorClass}">{value}</p>
	{#if subtitle}
		<p class="text-text-secondary mt-1 text-sm">{subtitle}</p>
	{/if}
</div>
