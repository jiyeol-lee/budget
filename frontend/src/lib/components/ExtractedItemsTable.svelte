<script lang="ts">
	import type { ExtractedItem } from '$lib/stores/receipt.svelte';
	import { ClipboardListIcon, TrashIcon, PlusIcon } from 'lucide-svelte';
	import * as m from '$lib/paraglide/messages';
	import { ExpenseTypeEnum } from '$lib/types/enums';
	import { Button } from '$lib';

	interface Props {
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
		items,
		onUpdateItem,
		onRemoveItem,
		onAddItem,
		onToggleSelection,
		onToggleSelectAll
	}: Props = $props();

	// Computed: check if all items are selected
	let allSelected = $derived(items.length > 0 && items.every((item) => item.selected));
	let someSelected = $derived(items.some((item) => item.selected) && !allSelected);

	/**
	 * Format currency for display
	 */
	function formatCurrency(amount: number): string {
		return new Intl.NumberFormat('en-US', {
			style: 'currency',
			currency: 'USD'
		}).format(amount);
	}

	/**
	 * Handle price input change with validation
	 */
	function handlePriceChange(index: number, value: string): void {
		const numValue = parseFloat(value);
		if (!isNaN(numValue)) {
			onUpdateItem(index, 'item_price', numValue);
		}
	}

	/**
	 * Handle select all checkbox change
	 */
	function handleSelectAllChange(event: Event): void {
		const checkbox = event.target as HTMLInputElement;
		onToggleSelectAll(checkbox.checked);
	}
</script>

<div class="overflow-x-auto">
	{#if items.length === 0}
		<!-- Empty State -->
		<div class="text-text-secondary py-12 text-center">
			<ClipboardListIcon class="text-text-secondary mx-auto h-12 w-12" />
			<p class="mt-4 text-lg font-medium">{m.receipt_no_items()}</p>
			<p class="mt-2 text-sm">{m.receipt_description()}</p>
		</div>
	{:else}
		<table class="divide-border min-w-full divide-y">
			<thead class="bg-surface-dark">
				<tr>
					<th scope="col" class="w-12 px-3 py-3">
						<input
							type="checkbox"
							checked={allSelected}
							indeterminate={someSelected}
							onchange={handleSelectAllChange}
							class="border-border text-primary focus:ring-primary h-4 w-4 rounded"
						/>
					</th>
					<th
						scope="col"
						class="text-text-secondary px-3 py-3 text-left text-xs font-medium tracking-wider uppercase"
					>
						{m.expense_form_source()}
					</th>
					<th
						scope="col"
						class="text-text-secondary px-3 py-3 text-left text-xs font-medium tracking-wider uppercase"
					>
						{m.expenses_type()}
					</th>
					<th
						scope="col"
						class="text-text-secondary px-3 py-3 text-left text-xs font-medium tracking-wider uppercase"
					>
						{m.expense_form_item_code()}
					</th>
					<th
						scope="col"
						class="text-text-secondary px-3 py-3 text-left text-xs font-medium tracking-wider uppercase"
					>
						{m.expense_form_amount()}
					</th>
					<th
						scope="col"
						class="text-text-secondary px-3 py-3 text-left text-xs font-medium tracking-wider uppercase"
					>
						{m.expense_form_item_name()}
					</th>
					<th
						scope="col"
						class="text-text-secondary w-16 px-3 py-3 text-right text-xs font-medium tracking-wider uppercase"
					>
						{m.common_actions()}
					</th>
				</tr>
			</thead>
			<tbody class="bg-surface divide-border divide-y">
				{#each items as item, index (index)}
					<tr
						class="hover:bg-surface-dark transition-colors {item.selected
							? 'bg-primary-dark/30'
							: ''}"
					>
						<!-- Select Checkbox -->
						<td class="px-3 py-2">
							<input
								type="checkbox"
								checked={item.selected}
								onchange={() => onToggleSelection(index)}
								class="border-border text-primary focus:ring-primary h-4 w-4 rounded"
							/>
						</td>

						<!-- Source -->
						<td class="px-3 py-2">
							<input
								type="text"
								value={item.source}
								onchange={(e) =>
									onUpdateItem(index, 'source', (e.target as HTMLInputElement).value)}
								class="border-border focus:border-border-focus focus:ring-primary bg-surface-light text-text-primary w-full rounded border px-2 py-1 text-sm focus:ring-1"
								placeholder="Store name"
							/>
						</td>

						<!-- Type -->
						<td class="px-3 py-2">
							<select
								value={item.type}
								onchange={(e) => onUpdateItem(index, 'type', (e.target as HTMLSelectElement).value)}
								class="border-border focus:border-border-focus focus:ring-primary bg-surface-light text-text-primary w-full rounded border px-2 py-1 text-sm focus:ring-1"
							>
								<option value={ExpenseTypeEnum.WEEKLY}>WEEKLY</option>
								<option value={ExpenseTypeEnum.MONTHLY}>MONTHLY</option>
								<option value={ExpenseTypeEnum.MISC}>MISC</option>
								<option value={ExpenseTypeEnum.TAX}>{m.expenses_type_tax()}</option>
							</select>
						</td>

						<!-- Item Code -->
						<td class="px-3 py-2">
							<input
								type="text"
								value={item.item_code}
								onchange={(e) =>
									onUpdateItem(index, 'item_code', (e.target as HTMLInputElement).value)}
								class="border-border focus:border-border-focus focus:ring-primary bg-surface-light text-text-primary w-full rounded border px-2 py-1 text-sm focus:ring-1"
								placeholder="Code"
							/>
						</td>

						<!-- Price -->
						<td class="px-3 py-2">
							<div class="relative">
								<span
									class="text-text-secondary absolute inset-y-0 left-0 flex items-center pl-2 text-sm"
									>$</span
								>
								<input
									type="number"
									value={item.item_price}
									onchange={(e) => handlePriceChange(index, (e.target as HTMLInputElement).value)}
									step="0.01"
									class="border-border focus:border-border-focus focus:ring-primary bg-surface-light text-text-primary w-full rounded border py-1 pr-2 pl-6 text-sm focus:ring-1"
									placeholder="0.00"
								/>
							</div>
						</td>

						<!-- Item Name -->
						<td class="px-3 py-2">
							<div class="flex items-center gap-2">
								<input
									type="text"
									value={item.item_name}
									onchange={(e) =>
										onUpdateItem(index, 'item_name', (e.target as HTMLInputElement).value)}
									class="border-border focus:border-border-focus focus:ring-primary bg-surface-light text-text-primary w-full rounded border px-2 py-1 text-sm focus:ring-1"
									placeholder="Item name"
								/>
							</div>
						</td>

						<!-- Actions -->
						<td class="px-3 py-2 text-right">
							<Button
								variant="ghost"
								onclick={() => onRemoveItem(index)}
								class="!text-danger hover:!text-danger-light p-1"
							>
								<TrashIcon class="h-4 w-4" />
							</Button>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>

		<!-- Add Row Button -->
		<div class="mt-4">
			<Button variant="ghost" onclick={onAddItem}>
				{#snippet leading()}<PlusIcon class="h-4 w-4" />{/snippet}
				{m.common_add()}
			</Button>
		</div>
	{/if}
</div>

<!-- Summary -->
{#if items.length > 0}
	<div class="border-border mt-4 border-t pt-4">
		<div class="flex items-center justify-between text-sm">
			<span class="text-text-tertiary">
				{items.filter((i) => i.selected).length} of {items.length} items selected
			</span>
			<span class="text-text-primary font-medium">
				Total: {formatCurrency(items.reduce((sum, item) => sum + item.item_price, 0))}
			</span>
		</div>
		{#if items.filter((i) => i.selected).length > 0}
			<div class="text-text-tertiary mt-2 text-sm">
				Selected total: {formatCurrency(
					items.filter((i) => i.selected).reduce((sum, item) => sum + item.item_price, 0)
				)}
			</div>
		{/if}
	</div>
{/if}
