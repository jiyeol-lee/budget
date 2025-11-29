<script lang="ts">
	import { onMount } from 'svelte';
	import {
		expectedExpensesStore,
		type ExpectedExpense,
		type ExpectedExpenseInput,
		type ExpectedExpenseFilterType
	} from '$lib/stores/expectedExpenses.svelte';
	import { toastStore } from '$lib/stores/toast.svelte';
	import ExpenseTable from '$lib/components/ExpenseTable.svelte';
	import ExpenseForm from '$lib/components/ExpenseForm.svelte';
	import ExpenseTabs from '$lib/components/ExpenseTabs.svelte';
	import Skeleton from '$lib/components/Skeleton.svelte';
	import { ExpenseFilterTypeEnum } from '$lib/types/enums';
	import { PlusIcon, AlertTriangleIcon } from 'lucide-svelte';
	import { Button, Dialog, SourceFilter, SummaryCard } from '$lib';
	import * as m from '$lib/paraglide/messages';

	// UI State
	let showForm = $state(false);
	let editingExpense = $state<ExpectedExpense | null>(null);
	let deleteConfirmExpense = $state<ExpectedExpense | null>(null);
	let deleting = $state(false);

	// Derived values from store
	let expenses = $derived(expectedExpensesStore.getFilteredExpenses());
	let loading = $derived(expectedExpensesStore.loading);
	// Error is handled via toasts now - accessing store.error in handlers
	let activeTab = $derived(expectedExpensesStore.filter);
	let sourceFilter = $derived(expectedExpensesStore.sourceFilter);
	let distinctSources = $derived(expectedExpensesStore.getDistinctSources());

	// Computed summary values
	let totalWeekly = $derived(expectedExpensesStore.getTotalWeekly());
	let totalMonthly = $derived(expectedExpensesStore.getTotalMonthly());
	let estimatedMonthlyTotal = $derived(expectedExpensesStore.getEstimatedMonthlyTotal());
	let expenseCount = $derived(expenses.length);

	/**
	 * Format amount as currency
	 */
	function formatCurrency(amount: number): string {
		return new Intl.NumberFormat('en-US', {
			style: 'currency',
			currency: 'USD'
		}).format(amount);
	}

	/**
	 * Handle tab change
	 */
	async function handleTabChange(tab: `${ExpenseFilterTypeEnum}`) {
		await expectedExpensesStore.setFilter(tab as ExpectedExpenseFilterType);
	}

	/**
	 * Open form for creating new expense
	 */
	function handleAddExpense() {
		editingExpense = null;
		showForm = true;
	}

	/**
	 * Open form for editing expense
	 */
	function handleEditExpense(expense: ExpectedExpense) {
		editingExpense = expense;
		showForm = true;
	}

	/**
	 * Show delete confirmation
	 */
	function handleDeleteClick(expense: ExpectedExpense) {
		deleteConfirmExpense = expense;
	}

	/**
	 * Confirm and execute delete
	 */
	async function handleConfirmDelete() {
		if (deleteConfirmExpense) {
			deleting = true;
			const success = await expectedExpensesStore.deleteExpense(deleteConfirmExpense.id);
			if (success) {
				toastStore.success(`"${deleteConfirmExpense.item_name}" deleted successfully`);
			} else if (expectedExpensesStore.error) {
				toastStore.error(expectedExpensesStore.error);
			}
			deleteConfirmExpense = null;
			deleting = false;
		}
	}

	/**
	 * Cancel delete
	 */
	function handleCancelDelete() {
		deleteConfirmExpense = null;
	}

	/**
	 * Handle form submission
	 */
	async function handleFormSubmit(data: ExpectedExpenseInput) {
		if (editingExpense) {
			// Update existing expense
			const result = await expectedExpensesStore.updateExpense(editingExpense.id, data);
			if (result) {
				toastStore.success(`"${data.item_name}" updated successfully`);
				showForm = false;
				editingExpense = null;
			} else if (expectedExpensesStore.error) {
				toastStore.error(expectedExpensesStore.error);
			}
		} else {
			// Create new expense
			const result = await expectedExpensesStore.createExpense(data);
			if (result) {
				toastStore.success(`"${data.item_name}" added successfully`);
				showForm = false;
			} else if (expectedExpensesStore.error) {
				toastStore.error(expectedExpensesStore.error);
			}
		}
	}

	/**
	 * Handle form cancel
	 */
	function handleFormCancel() {
		showForm = false;
		editingExpense = null;
	}

	// Fetch expenses on mount
	onMount(() => {
		expectedExpensesStore.fetchExpenses(ExpenseFilterTypeEnum.ALL);
	});
</script>

<svelte:head>
	<title>Expected Expenses | Budget Tracker</title>
</svelte:head>

<div class="space-y-4 sm:space-y-6">
	<!-- Page Header -->
	<div class="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
		<div>
			<h1 class="text-text-primary text-xl font-bold sm:text-2xl">{m.expected_expenses_title()}</h1>
			<p class="text-text-secondary mt-1 text-sm">
				{m.expected_expenses_description()}
			</p>
		</div>
		<Button onclick={handleAddExpense}>
			{#snippet leading()}<PlusIcon class="h-5 w-5" />{/snippet}
			{m.expected_expenses_add()}
		</Button>
	</div>

	<!-- Expenses Summary -->
	<div class="grid grid-cols-1 gap-4 sm:grid-cols-2 sm:gap-6 lg:grid-cols-3">
		{#if loading && expenses.length === 0}
			<Skeleton variant="card" />
			<Skeleton variant="card" />
			<Skeleton variant="card" />
		{:else}
			<SummaryCard
				title={m.expected_expenses_total_weekly()}
				value={formatCurrency(totalWeekly)}
				subtitle={m.expenses_weekly()}
				valueColor="info"
			/>
			<SummaryCard
				title={m.expected_expenses_total_monthly()}
				value={formatCurrency(totalMonthly)}
				subtitle={m.expenses_monthly()}
				valueColor="purple"
			/>
			<SummaryCard
				title={m.stats_monthly_expected()}
				value={formatCurrency(estimatedMonthlyTotal)}
				subtitle={`(${m.expenses_weekly()} × 4) + ${m.expenses_monthly()}`}
				valueColor="primary"
				class="sm:col-span-2 lg:col-span-1"
			/>
		{/if}
	</div>

	<!-- Form Modal -->
	<Dialog
		open={showForm}
		onClose={handleFormCancel}
		title={editingExpense ? m.expected_expenses_edit() : m.expected_expenses_add()}
		size="md"
	>
		<ExpenseForm expense={editingExpense} onSubmit={handleFormSubmit} onCancel={handleFormCancel} />
	</Dialog>

	<!-- Delete Confirmation Modal -->
	<Dialog
		open={deleteConfirmExpense !== null}
		onClose={handleCancelDelete}
		title={m.delete_confirm_title()}
		size="sm"
	>
		<div class="mb-4 flex items-center">
			<div class="bg-danger-dark flex h-10 w-10 shrink-0 items-center justify-center rounded-full">
				<AlertTriangleIcon class="text-danger h-6 w-6" />
			</div>
			<p class="text-text-tertiary ml-4 text-sm">
				{m.delete_confirm_description()}
			</p>
		</div>
		<div class="flex flex-col-reverse gap-3 sm:flex-row sm:justify-end">
			<Button variant="secondary" onclick={handleCancelDelete} disabled={deleting}>
				{m.common_cancel()}
			</Button>
			<Button variant="danger" onclick={handleConfirmDelete} loading={deleting}>
				{deleting ? m.common_processing() : m.common_delete()}
			</Button>
		</div>
	</Dialog>

	<!-- Expenses List Section -->
	<div class="bg-surface rounded-lg shadow">
		<div
			class="border-border flex flex-col gap-4 border-b px-4 py-4 sm:flex-row sm:items-center sm:justify-between sm:px-6"
		>
			<div class="flex items-center space-x-4">
				<h2 class="text-text-primary text-lg font-semibold">
					{activeTab === ExpenseFilterTypeEnum.ALL
						? 'All'
						: activeTab.charAt(0).toUpperCase() + activeTab.slice(1)} Expenses
				</h2>
				<span
					class="bg-surface-light text-text-tertiary inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium"
				>
					{expenseCount}
					{expenseCount === 1 ? 'item' : 'items'}
				</span>
			</div>
			<div class="flex flex-col gap-3 sm:flex-row sm:items-center">
				<SourceFilter
					sources={distinctSources}
					activeSource={sourceFilter}
					onFilterChange={(source) => expectedExpensesStore.setSourceFilter(source)}
				/>
				<ExpenseTabs {activeTab} onTabChange={handleTabChange} showMisc={false} />
			</div>
		</div>

		<div class="p-4 sm:p-6">
			<ExpenseTable
				{expenses}
				loading={loading && expenses.length === 0}
				onEdit={handleEditExpense}
				onDelete={handleDeleteClick}
			/>
		</div>
	</div>
</div>
