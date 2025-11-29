<script lang="ts">
	import * as m from '$lib/paraglide/messages';

	interface Props {
		value: number;
		onYearChange: (year: number) => void;
		/**
		 * Start year of the range (inclusive)
		 * @default 2025
		 */
		startYear?: number;
		/**
		 * End year of the range (inclusive)
		 * @default 2030
		 */
		endYear?: number;
		/**
		 * Optional id for the select element
		 */
		id?: string;
		/**
		 * Whether to show the label
		 * @default false
		 */
		showLabel?: boolean;
	}

	let {
		value,
		onYearChange,
		startYear = 2025,
		endYear = 2030,
		id = 'year-selector',
		showLabel = false
	}: Props = $props();

	// Generate years array from startYear to endYear
	let years = $derived(Array.from({ length: endYear - startYear + 1 }, (_, i) => startYear + i));

	function handleChange(event: Event) {
		const select = event.target as HTMLSelectElement;
		const selectedYear = parseInt(select.value, 10);
		onYearChange(selectedYear);
	}
</script>

{#if showLabel}
	<div class="flex items-center gap-2">
		<label for={id} class="text-text-secondary text-sm font-medium">
			{m.year_selector_label()}
		</label>
		<select
			{id}
			{value}
			onchange={handleChange}
			class="border-input-border bg-input-bg text-text-primary focus:border-input-focus focus:ring-input-focus block w-full rounded-md shadow-sm sm:text-sm"
		>
			{#each years as year (year)}
				<option value={year}>{year}</option>
			{/each}
		</select>
	</div>
{:else}
	<select
		{id}
		{value}
		onchange={handleChange}
		class="border-input-border bg-input-bg text-text-primary focus:border-input-focus focus:ring-input-focus block w-full rounded-md shadow-sm sm:text-sm"
	>
		{#each years as year (year)}
			<option value={year}>{year}</option>
		{/each}
	</select>
{/if}
