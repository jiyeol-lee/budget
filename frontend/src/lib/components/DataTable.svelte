<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Column {
		key: string;
		header: string;
		headerClass?: string;
		cellClass?: string;
		hideOnMobile?: boolean;
		skeletonWidth?: string;
	}

	interface Group<T> {
		key: string | number;
		items: T[];
	}

	interface Props<T> {
		// Data (use one)
		data?: T[];
		groups?: Group<T>[];

		// Config
		columns: Column[];
		keyField: keyof T;

		// Loading state
		loading?: boolean;
		skeletonRows?: number;

		// Empty state
		emptyIcon?: Snippet;
		emptyTitle: string;
		emptyDescription?: string;

		// Customization
		rowClass?: (item: T) => string;

		// Content snippets - ALL rendering delegated to parent
		cell: Snippet<[{ item: T; column: Column }]>;
		mobileCard: Snippet<[{ item: T }]>;
		groupHeader?: Snippet<[{ group: Group<T> }]>;
		groupFooter?: Snippet<[{ group: Group<T> }]>;
		tableFooter?: Snippet;
	}

	let {
		data,
		groups,
		columns,
		keyField,
		loading = false,
		skeletonRows = 5,
		emptyIcon,
		emptyTitle,
		emptyDescription,
		rowClass,
		cell,
		mobileCard,
		groupHeader,
		groupFooter,
		tableFooter
	}: Props<Record<string, unknown>> = $props();

	/**
	 * Convert skeletonWidth prop to Tailwind class
	 */
	function getSkeletonWidthClass(width?: string): string {
		if (!width) return 'w-3/4';
		// Handle fractional widths
		if (width.includes('/')) return `w-${width}`;
		// Handle numeric widths (like '16', '20', '24')
		if (/^\d+$/.test(width)) return `w-${width}`;
		// Handle 'full'
		if (width === 'full') return 'w-full';
		return `w-${width}`;
	}

	/**
	 * Check if the table has any data to display
	 */
	let isEmpty = $derived.by(() => {
		if (groups) {
			return groups.length === 0 || groups.every((g) => g.items.length === 0);
		}
		return !data || data.length === 0;
	});

	/**
	 * Get the key value for an item
	 */
	function getKey(item: Record<string, unknown>): string | number {
		return item[keyField as string] as string | number;
	}

	/**
	 * Get the row class for an item
	 */
	function getRowClass(item: Record<string, unknown>): string {
		const baseClass = 'hover:bg-surface-light transition-colors';
		const customClass = rowClass ? rowClass(item) : '';
		return `${baseClass} ${customClass}`.trim();
	}

	/**
	 * Get the number of visible columns (for colspan calculations)
	 */
	let columnCount = $derived(columns.length);
</script>

{#if loading}
	<!-- Loading Skeleton State -->

	<!-- Mobile Skeleton Cards (visible on small screens) -->
	<div class="space-y-3 p-4 sm:hidden">
		{#each Array(skeletonRows) as _, i (i)}
			<div class="bg-surface border-border animate-pulse rounded-lg border p-4">
				<div class="mb-3 flex items-start justify-between">
					<div class="bg-skeleton h-5 w-3/4 rounded"></div>
					<div class="bg-skeleton h-6 w-16 rounded-full"></div>
				</div>
				<div class="space-y-2">
					<div class="bg-skeleton h-4 w-1/2 rounded"></div>
					<div class="bg-skeleton h-4 w-1/3 rounded"></div>
				</div>
				<div class="border-border mt-3 flex justify-end gap-2 border-t pt-3">
					<div class="bg-skeleton h-8 w-16 rounded"></div>
					<div class="bg-skeleton h-8 w-16 rounded"></div>
				</div>
			</div>
		{/each}
	</div>

	<!-- Desktop Table Skeleton (hidden on small screens) -->
	<div class="-mx-4 hidden overflow-x-auto sm:mx-0 sm:block">
		<div class="inline-block min-w-full align-middle">
			<table class="divide-border min-w-full divide-y">
				<thead class="bg-surface-light">
					<tr>
						{#each columns as column (column.key)}
							<th
								scope="col"
								class="text-text-secondary px-4 py-3 text-left text-xs font-medium tracking-wider uppercase sm:px-6 {column.headerClass ||
									''}"
							>
								<div class="bg-skeleton h-3 w-16 animate-pulse rounded"></div>
							</th>
						{/each}
					</tr>
				</thead>
				<tbody class="bg-surface divide-border divide-y">
					{#each Array(skeletonRows) as _, rowIndex (rowIndex)}
						<tr class="animate-pulse">
							{#each columns as column (column.key)}
								<td class="px-4 py-4 whitespace-nowrap sm:px-6 {column.cellClass || ''}">
									<div
										class="bg-skeleton h-4 rounded {getSkeletonWidthClass(column.skeletonWidth)}"
									></div>
								</td>
							{/each}
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	</div>
{:else if isEmpty}
	<!-- Empty State -->
	<div class="text-text-secondary py-12 text-center">
		{#if emptyIcon}
			<div class="mx-auto h-12 w-12">
				{@render emptyIcon()}
			</div>
		{/if}
		<p class="mt-4 text-lg font-medium">{emptyTitle}</p>
		{#if emptyDescription}
			<p class="mt-2 text-sm">{emptyDescription}</p>
		{/if}
	</div>
{:else if groups}
	<!-- Grouped View -->

	<!-- Mobile Card View (visible on small screens) -->
	<div class="space-y-6 sm:hidden">
		{#each groups as group (group.key)}
			<div class="space-y-4">
				<!-- Group Header (Mobile) -->
				{#if groupHeader}
					{@render groupHeader({ group })}
				{/if}

				<!-- Mobile Cards -->
				<div class="space-y-3">
					{#each group.items as item (getKey(item))}
						{@render mobileCard({ item })}
					{/each}
				</div>

				<!-- Group Footer (Mobile) -->
				{#if groupFooter}
					{@render groupFooter({ group })}
				{/if}
			</div>
		{/each}
	</div>

	<!-- Desktop Table View (hidden on small screens) -->
	<div class="-mx-4 hidden overflow-x-auto sm:mx-0 sm:block">
		<div class="inline-block min-w-full align-middle">
			<table class="divide-border min-w-full divide-y">
				<thead class="bg-surface-light">
					<tr>
						{#each columns as column (column.key)}
							<th
								scope="col"
								class="text-text-secondary px-4 py-3 text-left text-xs font-medium tracking-wider uppercase sm:px-6 {column.headerClass ||
									''}"
							>
								{column.header}
							</th>
						{/each}
					</tr>
				</thead>
				<tbody class="bg-surface divide-border divide-y">
					{#each groups as group (group.key)}
						<!-- Group Header Row -->
						{#if groupHeader}
							<tr>
								<td colspan={columnCount} class="p-0">
									{@render groupHeader({ group })}
								</td>
							</tr>
						{/if}

						<!-- Data Rows -->
						{#each group.items as item (getKey(item))}
							<tr class={getRowClass(item)}>
								{#each columns as column (column.key)}
									<td class="px-4 py-4 whitespace-nowrap sm:px-6 {column.cellClass || ''}">
										{@render cell({ item, column })}
									</td>
								{/each}
							</tr>
						{/each}

						<!-- Group Footer Row -->
						{#if groupFooter}
							<tr>
								<td colspan={columnCount} class="p-0">
									{@render groupFooter({ group })}
								</td>
							</tr>
						{/if}
					{/each}
				</tbody>
				{#if tableFooter}
					<tfoot class="bg-surface">
						{@render tableFooter()}
					</tfoot>
				{/if}
			</table>
		</div>
	</div>
{:else if data}
	<!-- Flat View (non-grouped) -->

	<!-- Mobile Card View (visible on small screens) -->
	<div class="space-y-4 sm:hidden">
		{#each data as item (getKey(item))}
			{@render mobileCard({ item })}
		{/each}
	</div>

	<!-- Desktop Table View (hidden on small screens) -->
	<div class="-mx-4 hidden overflow-x-auto sm:mx-0 sm:block">
		<div class="inline-block min-w-full align-middle">
			<table class="divide-border min-w-full divide-y">
				<thead class="bg-surface-light">
					<tr>
						{#each columns as column (column.key)}
							<th
								scope="col"
								class="text-text-secondary px-4 py-3 text-left text-xs font-medium tracking-wider uppercase sm:px-6 {column.headerClass ||
									''}"
							>
								{column.header}
							</th>
						{/each}
					</tr>
				</thead>
				<tbody class="bg-surface divide-border divide-y">
					{#each data as item (getKey(item))}
						<tr class={getRowClass(item)}>
							{#each columns as column (column.key)}
								<td class="px-4 py-4 whitespace-nowrap sm:px-6 {column.cellClass || ''}">
									{@render cell({ item, column })}
								</td>
							{/each}
						</tr>
					{/each}
				</tbody>
				{#if tableFooter}
					<tfoot class="bg-surface">
						{@render tableFooter()}
					</tfoot>
				{/if}
			</table>
		</div>
	</div>
{/if}
