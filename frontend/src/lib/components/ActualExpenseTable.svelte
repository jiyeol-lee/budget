<script lang="ts">
	import type { ActualExpense } from '$lib/stores/actualExpenses.svelte';
	import { formatCurrency, formatDate } from '$lib/utils/format';
	import { ExpenseTypeEnum, ExpenseFilterTypeEnum } from '$lib/types/enums';
	import { ClipboardListIcon, PencilIcon, TrashIcon, ReceiptIcon } from 'lucide-svelte';
	import * as m from '$lib/paraglide/messages';
	import { Button, DataTable } from '$lib';
	import { getTypeBadgeClass } from '$lib/utils/expense';
	import { SvelteMap } from 'svelte/reactivity';

	interface Props {
		expenses: ActualExpense[];
		onEdit: (expense: ActualExpense) => void;
		onDelete: (expense: ActualExpense) => void;
		filterType?: ExpenseFilterTypeEnum;
		loading?: boolean;
	}

	let {
		expenses,
		onEdit,
		onDelete,
		filterType = ExpenseFilterTypeEnum.ALL,
		loading = false
	}: Props = $props();

	/**
	 * Receipt group with tax item separated
	 */
	interface ReceiptGroup {
		receiptNumber: number;
		items: ActualExpense[];
		taxItem: ActualExpense | null;
		total: number;
	}

	/**
	 * DataTable group structure
	 */
	interface DataTableGroup {
		key: string | number;
		items: ActualExpense[];
		meta: ReceiptGroup;
	}

	/**
	 * Check if we're in grouped mode
	 */
	let isGroupedMode = $derived(filterType === ExpenseFilterTypeEnum.ALL);

	/**
	 * Group expenses by receipt_number for grouped view
	 */
	let groupedData = $derived.by((): DataTableGroup[] | null => {
		if (!isGroupedMode) {
			return null;
		}

		const groups = new SvelteMap<number, ReceiptGroup>();

		for (const expense of expenses) {
			const receiptNum = expense.receipt_number || 0;
			if (!groups.has(receiptNum)) {
				groups.set(receiptNum, {
					receiptNumber: receiptNum,
					items: [],
					taxItem: null,
					total: 0
				});
			}

			const group = groups.get(receiptNum);
			if (group) {
				if (expense.expense_type === ExpenseTypeEnum.TAX) {
					group.taxItem = expense;
				} else {
					group.items.push(expense);
				}
				group.total += expense.actual_amount;
			}
		}

		// Sort groups by receipt number (descending - newest first)
		// Convert to DataTable group structure
		return Array.from(groups.values())
			.sort((a, b) => b.receiptNumber - a.receiptNumber)
			.map((group) => ({
				key: group.receiptNumber,
				items: group.items,
				meta: group
			}));
	});

	/**
	 * Column definitions for grouped mode
	 */
	const groupedColumns = [
		{ key: 'receipt_date', header: m.actual_expenses_month(), skeletonWidth: '20' },
		{ key: 'item_name', header: m.expense_form_item_name(), skeletonWidth: '3/4' },
		{ key: 'source', header: m.expense_form_source(), skeletonWidth: '1/2' },
		{ key: 'expense_type', header: m.expenses_type(), skeletonWidth: '16' },
		{ key: 'actual_amount', header: m.expense_form_amount(), skeletonWidth: '20' },
		{ key: 'actions', header: m.common_actions(), headerClass: 'text-right', skeletonWidth: '24' }
	];

	/**
	 * Column definitions for flat mode (includes receipt_number)
	 */
	const flatColumns = [
		{ key: 'receipt_number', header: m.receipt_number_label(), skeletonWidth: '12' },
		{ key: 'receipt_date', header: m.actual_expenses_month(), skeletonWidth: '20' },
		{ key: 'item_name', header: m.expense_form_item_name(), skeletonWidth: '3/4' },
		{ key: 'source', header: m.expense_form_source(), skeletonWidth: '1/2' },
		{ key: 'expense_type', header: m.expenses_type(), skeletonWidth: '16' },
		{ key: 'actual_amount', header: m.expense_form_amount(), skeletonWidth: '20' },
		{ key: 'actions', header: m.common_actions(), headerClass: 'text-right', skeletonWidth: '24' }
	];

	/**
	 * Check if expense is a tax item
	 */
	function isTaxItem(expense: ActualExpense): boolean {
		return expense.expense_type === ExpenseTypeEnum.TAX;
	}

	/**
	 * Get row class for tax item highlighting
	 */
	function getRowClass(item: Record<string, unknown>): string {
		const expense = item as unknown as ActualExpense;
		return isTaxItem(expense) ? 'bg-amber-950/30' : '';
	}

	/**
	 * Get expense type display text
	 */
	function getExpenseTypeText(type: ExpenseTypeEnum): string {
		return type === ExpenseTypeEnum.TAX ? m.expenses_type_tax() : type;
	}

	/**
	 * Cast Record to ActualExpense for type safety in snippets
	 */
	function asExpense(item: Record<string, unknown>): ActualExpense {
		return item as unknown as ActualExpense;
	}
</script>

{#if isGroupedMode && groupedData}
	<!-- Grouped View -->
	<div class="space-y-6">
		{#each groupedData as group (group.key)}
			{@const meta = group.meta}
			<!-- Receipt Group -->
			<div class="bg-surface border-border overflow-hidden rounded-lg border shadow-sm">
				<!-- Group Header -->
				<div
					class="bg-surface-dark border-border flex items-center justify-between border-b px-4 py-3"
				>
					<div class="flex items-center gap-2">
						<ReceiptIcon class="text-text-secondary h-5 w-5" />
						<span class="text-text-primary font-medium">
							{m.receipt_number_label()}{meta.receiptNumber}
						</span>
					</div>
					<span class="text-text-primary font-semibold">
						{m.receipt_group_total()}: {formatCurrency(meta.total)}
					</span>
				</div>

				<DataTable
					{loading}
					data={group.items as unknown as Record<string, unknown>[]}
					columns={groupedColumns}
					keyField="id"
					emptyTitle={m.actual_expenses_empty()}
				>
					{#snippet emptyIcon()}
						<ClipboardListIcon class="text-text-secondary h-12 w-12" />
					{/snippet}

					{#snippet cell({ item, column })}
						{@const expense = asExpense(item)}
						{#if column.key === 'receipt_date'}
							<div class="text-text-primary text-sm">
								{formatDate(expense.receipt_date)}
							</div>
						{:else if column.key === 'item_name'}
							<div class="text-text-primary text-sm font-medium">
								{expense.item_name}
								{#if expense.item_code}
									<span class="text-text-secondary ml-1 text-xs">({expense.item_code})</span>
								{/if}
							</div>
						{:else if column.key === 'source'}
							<div class="text-text-tertiary text-sm">
								{expense.source}
							</div>
						{:else if column.key === 'expense_type'}
							<span
								class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium uppercase {getTypeBadgeClass(
									expense.expense_type
								)}"
							>
								{expense.expense_type}
							</span>
						{:else if column.key === 'actual_amount'}
							<div class="text-text-primary text-sm font-medium">
								{formatCurrency(expense.actual_amount)}
							</div>
						{:else if column.key === 'actions'}
							<div class="text-right">
								<Button variant="link" onclick={() => onEdit(expense)} class="text-primary mr-3">
									{m.common_edit()}
								</Button>
								<Button
									variant="link"
									onclick={() => onDelete(expense)}
									class="!text-danger hover:!text-danger-light"
								>
									{m.common_delete()}
								</Button>
							</div>
						{/if}
					{/snippet}

					{#snippet mobileCard({ item })}
						{@const expense = asExpense(item)}
						<div
							class="bg-surface-light border-border rounded-lg border p-3 {isTaxItem(expense)
								? 'border-amber-800 bg-amber-950/30'
								: ''}"
						>
							<div class="mb-2 flex items-start justify-between">
								<div class="min-w-0 flex-1">
									<div class="flex items-center space-x-2">
										<span class="text-text-secondary text-xs"
											>{formatDate(expense.receipt_date)}</span
										>
										<span
											class="inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium uppercase {getTypeBadgeClass(
												expense.expense_type
											)}"
										>
											{getExpenseTypeText(expense.expense_type)}
										</span>
									</div>
									<h3 class="text-text-primary mt-1 truncate text-sm font-medium">
										{expense.item_name}
									</h3>
									<p class="text-text-secondary text-sm">{expense.source}</p>
								</div>
							</div>
							<div class="flex items-center justify-between">
								<div>
									<p class="text-text-primary text-lg font-semibold">
										{formatCurrency(expense.actual_amount)}
									</p>
									{#if expense.item_code}
										<p class="text-text-secondary text-xs">
											{m.expense_item_code_label()}: {expense.item_code}
										</p>
									{/if}
								</div>
								<div class="flex space-x-3">
									<Button
										variant="ghost"
										onclick={() => onEdit(expense)}
										class="text-primary p-2"
										aria-label={m.common_edit()}
									>
										<PencilIcon class="h-5 w-5" />
									</Button>
									<Button
										variant="ghost"
										onclick={() => onDelete(expense)}
										class="!text-danger hover:!text-danger-light p-2"
										aria-label={m.common_delete()}
									>
										<TrashIcon class="h-5 w-5" />
									</Button>
								</div>
							</div>
						</div>
					{/snippet}

					{#snippet tableFooter()}
						{#if meta.taxItem}
							<tr class="border-t border-amber-800 bg-amber-950/30">
								<td class="px-4 py-4 whitespace-nowrap sm:px-6">
									<div class="text-sm text-amber-200">
										{formatDate(meta.taxItem.receipt_date)}
									</div>
								</td>
								<td class="px-4 py-4 whitespace-nowrap sm:px-6">
									<div class="text-sm font-medium text-amber-200">
										{m.receipt_group_tax()}
									</div>
								</td>
								<td class="px-4 py-4 whitespace-nowrap sm:px-6">
									<div class="text-sm text-amber-300">-</div>
								</td>
								<td class="px-4 py-4 whitespace-nowrap sm:px-6">
									<span
										class="inline-flex items-center rounded-full bg-amber-900 px-2.5 py-0.5 text-xs font-medium text-amber-200 uppercase"
									>
										{m.expenses_type_tax()}
									</span>
								</td>
								<td class="px-4 py-4 whitespace-nowrap sm:px-6">
									<div class="text-sm font-semibold text-amber-200">
										{formatCurrency(meta.taxItem.actual_amount)}
									</div>
								</td>
								<td class="px-4 py-4 text-right text-sm font-medium whitespace-nowrap sm:px-6">
									<Button
										variant="link"
										onclick={() => meta.taxItem && onEdit(meta.taxItem)}
										class="text-primary mr-3"
									>
										{m.common_edit()}
									</Button>
									<Button
										variant="link"
										onclick={() => meta.taxItem && onDelete(meta.taxItem)}
										class="!text-danger hover:!text-danger-light"
									>
										{m.common_delete()}
									</Button>
								</td>
							</tr>
						{/if}
					{/snippet}
				</DataTable>

				<!-- Tax Row (Mobile) -->
				{#if meta.taxItem}
					<div class="p-4 pt-0 sm:hidden">
						<div class="rounded-lg border border-amber-800 bg-amber-950/30 p-3">
							<div class="flex items-center justify-between">
								<div>
									<span class="text-sm font-medium text-amber-200">{m.receipt_group_tax()}</span>
								</div>
								<div class="flex items-center gap-3">
									<span class="font-semibold text-amber-200"
										>{formatCurrency(meta.taxItem.actual_amount)}</span
									>
									<div class="flex space-x-2">
										<Button
											variant="ghost"
											onclick={() => meta.taxItem && onEdit(meta.taxItem)}
											class="text-primary p-1"
											aria-label={m.common_edit()}
										>
											<PencilIcon class="h-4 w-4" />
										</Button>
										<Button
											variant="ghost"
											onclick={() => meta.taxItem && onDelete(meta.taxItem)}
											class="!text-danger hover:!text-danger-light p-1"
											aria-label={m.common_delete()}
										>
											<TrashIcon class="h-4 w-4" />
										</Button>
									</div>
								</div>
							</div>
						</div>
					</div>
				{/if}
			</div>
		{/each}
	</div>

	<!-- Empty state for grouped mode -->
	{#if groupedData.length === 0}
		<div class="text-text-secondary py-12 text-center">
			<ClipboardListIcon class="text-text-secondary mx-auto h-12 w-12" />
			<p class="mt-4 text-lg font-medium">{m.actual_expenses_empty()}</p>
			<p class="mt-2 text-sm">{m.actual_expenses_empty_description()}</p>
		</div>
	{/if}
{:else}
	<!-- Flat View (Filtered) -->
	<DataTable
		{loading}
		data={expenses as unknown as Record<string, unknown>[]}
		columns={flatColumns}
		keyField="id"
		rowClass={getRowClass}
		emptyTitle={m.actual_expenses_empty()}
		emptyDescription={m.actual_expenses_empty_description()}
	>
		{#snippet emptyIcon()}
			<ClipboardListIcon class="text-text-secondary h-12 w-12" />
		{/snippet}

		{#snippet cell({ item, column })}
			{@const expense = asExpense(item)}
			{@const tax = isTaxItem(expense)}
			{#if column.key === 'receipt_number'}
				<div class="{tax ? 'text-amber-200' : 'text-text-primary'} text-sm font-medium">
					{#if expense.receipt_number}
						#{expense.receipt_number}
					{:else}
						-
					{/if}
				</div>
			{:else if column.key === 'receipt_date'}
				<div class="{tax ? 'text-amber-200' : 'text-text-primary'} text-sm">
					{formatDate(expense.receipt_date)}
				</div>
			{:else if column.key === 'item_name'}
				<div class="{tax ? 'text-amber-200' : 'text-text-primary'} text-sm font-medium">
					{expense.item_name}
					{#if expense.item_code}
						<span class="text-text-secondary ml-1 text-xs">({expense.item_code})</span>
					{/if}
				</div>
			{:else if column.key === 'source'}
				<div class="{tax ? 'text-amber-300' : 'text-text-tertiary'} text-sm">
					{expense.source}
				</div>
			{:else if column.key === 'expense_type'}
				<span
					class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium uppercase {getTypeBadgeClass(
						expense.expense_type
					)}"
				>
					{getExpenseTypeText(expense.expense_type)}
				</span>
			{:else if column.key === 'actual_amount'}
				<div class="{tax ? 'text-amber-200' : 'text-text-primary'} text-sm font-medium">
					{formatCurrency(expense.actual_amount)}
				</div>
			{:else if column.key === 'actions'}
				<div class="text-right">
					<Button variant="link" onclick={() => onEdit(expense)} class="text-primary mr-3">
						{m.common_edit()}
					</Button>
					<Button
						variant="link"
						onclick={() => onDelete(expense)}
						class="!text-danger hover:!text-danger-light"
					>
						{m.common_delete()}
					</Button>
				</div>
			{/if}
		{/snippet}

		{#snippet mobileCard({ item })}
			{@const expense = asExpense(item)}
			{@const tax = isTaxItem(expense)}
			<div
				class="bg-surface border-border rounded-lg border p-4 shadow-sm {tax
					? 'border-amber-800 bg-amber-950/30'
					: ''}"
			>
				<div class="mb-3 flex items-start justify-between">
					<div class="min-w-0 flex-1">
						<div class="flex items-center space-x-2">
							{#if expense.receipt_number}
								<span class="text-text-secondary text-xs font-medium"
									>#{expense.receipt_number}</span
								>
							{/if}
							<span class="text-text-secondary text-xs">{formatDate(expense.receipt_date)}</span>
							<span
								class="inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium uppercase {getTypeBadgeClass(
									expense.expense_type
								)}"
							>
								{getExpenseTypeText(expense.expense_type)}
							</span>
						</div>
						<h3 class="text-text-primary mt-1 truncate text-sm font-medium">{expense.item_name}</h3>
						<p class="text-text-secondary text-sm">{expense.source}</p>
					</div>
				</div>
				<div class="flex items-center justify-between">
					<div>
						<p class="{tax ? 'text-amber-200' : 'text-text-primary'} text-lg font-semibold">
							{formatCurrency(expense.actual_amount)}
						</p>
						{#if expense.item_code}
							<p class="text-text-secondary text-xs">
								{m.expense_item_code_label()}: {expense.item_code}
							</p>
						{/if}
					</div>
					<div class="flex space-x-3">
						<Button
							variant="ghost"
							onclick={() => onEdit(expense)}
							class="text-primary p-2"
							aria-label={m.common_edit()}
						>
							<PencilIcon class="h-5 w-5" />
						</Button>
						<Button
							variant="ghost"
							onclick={() => onDelete(expense)}
							class="!text-danger hover:!text-danger-light p-2"
							aria-label={m.common_delete()}
						>
							<TrashIcon class="h-5 w-5" />
						</Button>
					</div>
				</div>
			</div>
		{/snippet}
	</DataTable>
{/if}
