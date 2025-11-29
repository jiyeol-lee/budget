<script lang="ts">
	import ReceiptPreview from './ReceiptPreview.svelte';
	import ExtractedItemsTable from './ExtractedItemsTable.svelte';
	import { Button } from '$lib';
	import type { ExtractedItem } from '$lib/stores/receipt.svelte';
	import { expectedExpensesStore } from '$lib/stores/expectedExpenses.svelte';
	import { onMount } from 'svelte';
	import { EyeIcon, EyeOffIcon } from 'lucide-svelte';

	interface Props {
		imageUrl: string | null;
		items: ExtractedItem[];
		onUpdateItem: (
			index: number,
			field: keyof ExtractedItem,
			value: string | number | boolean
		) => void;
		onRemoveItem: (index: number) => void;
		onAddItem: () => void;
		onToggleSelection: (index: number) => void;
		onToggleSelectAll: (selectAll: boolean) => void;
	}

	let {
		imageUrl,
		items,
		onUpdateItem,
		onRemoveItem,
		onAddItem,
		onToggleSelection,
		onToggleSelectAll
	}: Props = $props();

	onMount(async () => {
		if (expectedExpensesStore.expenses.length === 0) {
			await expectedExpensesStore.fetchExpenses();
		}
	});

	$effect(() => {
		if (expectedExpensesStore.expenses.length > 0 && items.length > 0) {
			matchItems();
		}
	});

	function matchItems() {
		const expected = expectedExpensesStore.expenses;

		items.forEach((item, index) => {
			// Skip if already matched
			if (item.expected_expense_id) return;

			// Try match by Item Name (case insensitive)
			if (item.item_name) {
				const match = expected.find(
					(e) => e.item_name.toLowerCase() === item.item_name.toLowerCase()
				);
				if (match) {
					onUpdateItem(index, 'expected_expense_id', match.id);
					return;
				}
			}
		});
	}

	let showReceiptImage = $state(false);

	function toggleReceiptImage() {
		showReceiptImage = !showReceiptImage;
	}
</script>

<div class="grid grid-cols-1 gap-6">
	<!-- Extracted Items Table (Full Width) -->
	<div class="bg-surface overflow-hidden rounded-lg p-4 shadow">
		<div class="mb-4 flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
			<h3 class="text-text-primary text-lg font-semibold">Extracted Items</h3>

			<Button variant="secondary" onclick={toggleReceiptImage}>
				{#snippet leading()}
					{#if showReceiptImage}
						<EyeOffIcon class="h-4 w-4" />
					{:else}
						<EyeIcon class="h-4 w-4" />
					{/if}
				{/snippet}
				{showReceiptImage ? 'Hide Receipt Image' : 'Show Receipt Image'}
			</Button>
		</div>

		{#if showReceiptImage}
			<div class="bg-surface-dark border-border mb-6 overflow-hidden rounded-lg border shadow">
				<div class="h-[500px]">
					<ReceiptPreview {imageUrl} alt="Uploaded receipt" />
				</div>
			</div>
		{/if}

		<ExtractedItemsTable
			{items}
			{onUpdateItem}
			{onRemoveItem}
			{onAddItem}
			{onToggleSelection}
			{onToggleSelectAll}
		/>
	</div>
</div>
