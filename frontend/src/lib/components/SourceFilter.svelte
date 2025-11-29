<script lang="ts">
	import * as m from '$lib/paraglide/messages';

	interface Props {
		sources: string[];
		activeSource: string | null;
		onFilterChange: (source: string | null) => void;
	}

	let { sources, activeSource, onFilterChange }: Props = $props();

	function handleChange(event: Event) {
		const select = event.target as HTMLSelectElement;
		const value = select.value;
		onFilterChange(value === '' ? null : value);
	}
</script>

<div class="flex items-center gap-2">
	<label for="source-filter" class="text-text-secondary text-sm font-medium">
		{m.filter_source()}
	</label>
	<select
		id="source-filter"
		value={activeSource ?? ''}
		onchange={handleChange}
		class="border-border focus:ring-primary focus:border-border-focus bg-surface-light text-text-primary rounded-md border px-3 py-1.5 text-sm shadow-sm transition-colors"
	>
		<option value="">{m.filter_source_all()}</option>
		{#each sources as source (source)}
			<option value={source}>{source}</option>
		{/each}
	</select>
</div>
