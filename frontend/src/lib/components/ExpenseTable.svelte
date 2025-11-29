<script lang="ts">
	import type { ExpectedExpense } from '$lib/stores/expectedExpenses.svelte';
	import { Button, DataTable } from '$lib';
	import { ClipboardListIcon, PencilIcon, TrashIcon } from 'lucide-svelte';
	import * as m from '$lib/paraglide/messages';
	import { formatCurrency } from '$lib/utils/format';
	import { getTypeBadgeClass } from '$lib/utils/expense';

	interface Props {
		expenses: ExpectedExpense[];
		onEdit: (expense: ExpectedExpense) => void;
		onDelete: (expense: ExpectedExpense) => void;
		loading?: boolean;
	}

	let { expenses, onEdit, onDelete, loading = false }: Props = $props();

	/**
	 * Column definitions with i18n headers
	 */
	const columns = $derived([
		{ key: 'item_name', header: m.expense_form_item_name(), skeletonWidth: '3/4' },
		{ key: 'source', header: m.expense_form_source(), skeletonWidth: '1/2' },
		{ key: 'expected_amount', header: m.expense_form_amount(), skeletonWidth: '20' },
		{ key: 'expense_type', header: m.expenses_type(), skeletonWidth: '16' },
		{ key: 'actions', header: m.common_actions(), headerClass: 'text-right', skeletonWidth: '24' }
	]);
</script>

<DataTable
	{loading}
	data={expenses as unknown as Record<string, unknown>[]}
	{columns}
	keyField="id"
	emptyTitle={m.expected_expenses_empty()}
	emptyDescription={m.expected_expenses_empty_description()}
>
	{#snippet emptyIcon()}
		<ClipboardListIcon class="text-text-secondary h-12 w-12" />
	{/snippet}

	{#snippet cell({ item, column })}
		{@const expense = item as unknown as ExpectedExpense}
		{#if column.key === 'item_name'}
			<div class="text-text-primary text-sm font-medium">
				{expense.item_name}
			</div>
		{:else if column.key === 'source'}
			<div class="text-text-secondary text-sm">
				{expense.source}
			</div>
		{:else if column.key === 'expected_amount'}
			<div class="text-text-primary text-sm font-medium">
				{formatCurrency(expense.expected_amount)}
			</div>
		{:else if column.key === 'expense_type'}
			<span
				class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium uppercase {getTypeBadgeClass(
					expense.expense_type
				)}"
			>
				{expense.expense_type}
			</span>
		{:else if column.key === 'actions'}
			<div class="text-right">
				<Button variant="link" onclick={() => onEdit(expense)} class="mr-3">
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
		{@const expense = item as unknown as ExpectedExpense}
		<div class="bg-surface border-border rounded-lg border p-4 shadow-sm">
			<div class="mb-3 flex items-start justify-between">
				<div class="min-w-0 flex-1">
					<h3 class="text-text-primary truncate text-sm font-medium">{expense.item_name}</h3>
					<p class="text-text-secondary text-sm">{expense.source}</p>
				</div>
				<span
					class="ml-2 inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium uppercase {getTypeBadgeClass(
						expense.expense_type
					)}"
				>
					{expense.expense_type}
				</span>
			</div>
			<div class="flex items-center justify-between">
				<div>
					<p class="text-text-primary text-lg font-semibold">
						{formatCurrency(expense.expected_amount)}
					</p>
				</div>
				<div class="flex space-x-3">
					<Button
						variant="ghost"
						onclick={() => onEdit(expense)}
						class="text-primary hover:bg-primary-dark/30 p-2"
						aria-label={m.common_edit()}
					>
						<PencilIcon class="h-5 w-5" />
					</Button>
					<Button
						variant="ghost"
						onclick={() => onDelete(expense)}
						class="!text-danger hover:!text-danger-light hover:bg-danger-dark/30 p-2"
						aria-label={m.common_delete()}
					>
						<TrashIcon class="h-5 w-5" />
					</Button>
				</div>
			</div>
		</div>
	{/snippet}
</DataTable>
