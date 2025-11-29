<script lang="ts">
	import { onMount } from 'svelte';
	import {
		actualExpensesStore,
		type ActualExpense,
		type ActualExpenseFilterType
	} from '$lib/stores/actualExpenses.svelte';
	import { toastStore } from '$lib/stores/toast.svelte';
	import ActualExpenseTable from '$lib/components/ActualExpenseTable.svelte';
	import ExpenseTabs from '$lib/components/ExpenseTabs.svelte';
	import Skeleton from '$lib/components/Skeleton.svelte';
	import { Button, SummaryCard, SourceFilter } from '$lib';
	import { ExpenseFilterTypeEnum } from '$lib/types/enums';
	import { AlertTriangleIcon } from 'lucide-svelte';
	import * as m from '$lib/paraglide/messages';

	// UI State
	let deleteConfirmExpense = $state<ActualExpense | null>(null);
	let deleting = $state(false);

	// Derived values from store
	let expenses = $derived(actualExpensesStore.getFilteredExpenses());
	let loading = $derived(actualExpensesStore.loading);
	let activeTab = $derived(actualExpensesStore.filterType);
	let summary = $derived(actualExpensesStore.summary);
	let currentMonth = $derived(actualExpensesStore.currentMonth);
	let currentYear = $derived(actualExpensesStore.currentYear);
	let sourceFilter = $derived(actualExpensesStore.sourceFilter);
	let distinctSources = $derived(actualExpensesStore.getDistinctSources());

	let expenseCount = $derived(expenses.length);

	const months = [
		'January',
		'February',
		'March',
		'April',
		'May',
		'June',
		'July',
		'August',
		'September',
		'October',
		'November',
		'December'
	];

	// Generate years (current year - 5 to current year + 1)
	const currentYearVal = new Date().getFullYear();
	const years = Array.from({ length: 7 }, (_, i) => currentYearVal - 5 + i);

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
	function handleTabChange(tab: ActualExpenseFilterType) {
		actualExpensesStore.setFilter(tab);
		actualExpensesStore.fetchExpenses();
	}

	/**
	 * Handle month/year change
	 */
	async function handleDateChange(month: number, year: number) {
		actualExpensesStore.setMonthYear(month, year);
		await Promise.all([actualExpensesStore.fetchExpenses(), actualExpensesStore.fetchSummary()]);
	}

	/**
	 * Placeholder for edit (since actual expenses come from receipts mainly, but manual add/edit might be supported)
	 */
	function handleEditExpense(_expense: ActualExpense) {
		// TODO: Implement manual edit if needed, or link to receipt view
		toastStore.info(
			'Editing actual expenses is not yet implemented directly. Please edit via receipts.'
		);
	}

	/**
	 * Show delete confirmation
	 */
	function handleDeleteClick(expense: ActualExpense) {
		deleteConfirmExpense = expense;
	}

	/**
	 * Confirm and execute delete
	 */
	async function handleConfirmDelete() {
		if (deleteConfirmExpense) {
			deleting = true;
			const success = await actualExpensesStore.deleteExpense(deleteConfirmExpense.id);
			if (success) {
				toastStore.success(`"${deleteConfirmExpense.item_name}" deleted successfully`);
			} else if (actualExpensesStore.error) {
				toastStore.error(actualExpensesStore.error);
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
	 * Handle click outside modal to close
	 */
	function handleModalBackdropClick(event: MouseEvent) {
		if (event.target === event.currentTarget) {
			if (deleteConfirmExpense) {
				handleCancelDelete();
			}
		}
	}

	// Fetch expenses on mount
	onMount(() => {
		Promise.all([actualExpensesStore.fetchExpenses(), actualExpensesStore.fetchSummary()]);
	});
</script>

<svelte:head>
	<title>Actual Expenses | Budget Tracker</title>
</svelte:head>

<div class="space-y-4 sm:space-y-6">
	<!-- Page Header -->
	<div class="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
		<div>
			<h1 class="text-text-primary text-xl font-bold sm:text-2xl">{m.actual_expenses_title()}</h1>
			<p class="text-text-secondary mt-1 text-sm">
				{m.actual_expenses_description()}
			</p>
		</div>

		<!-- Month/Year Selector -->
		<div class="flex space-x-2">
			<select
				value={currentMonth}
				onchange={(e) => handleDateChange(parseInt(e.currentTarget.value), currentYear)}
				class="border-input-border bg-input-bg text-text-primary focus:border-input-focus focus:ring-input-focus block w-full rounded-md shadow-sm sm:text-sm"
			>
				{#each months as month, i (i)}
					<option value={i + 1}>{month}</option>
				{/each}
			</select>
			<select
				value={currentYear}
				onchange={(e) => handleDateChange(currentMonth, parseInt(e.currentTarget.value))}
				class="border-input-border bg-input-bg text-text-primary focus:border-input-focus focus:ring-input-focus block w-full rounded-md shadow-sm sm:text-sm"
			>
				{#each years as year (year)}
					<option value={year}>{year}</option>
				{/each}
			</select>
		</div>
	</div>

	<!-- Expenses Summary -->
	<div class="grid grid-cols-1 gap-4 sm:grid-cols-2 sm:gap-6 lg:grid-cols-4">
		{#if loading && !summary}
			<Skeleton variant="card" />
			<Skeleton variant="card" />
			<Skeleton variant="card" />
			<Skeleton variant="card" />
		{:else if summary}
			<SummaryCard
				title={m.actual_expenses_total()}
				value={formatCurrency(summary.total_actual)}
				valueColor="primary"
			/>
			<SummaryCard
				title={m.expected_expenses_total_weekly()}
				value={formatCurrency(summary.total_weekly)}
				valueColor="info"
			/>
			<SummaryCard
				title={m.expected_expenses_total_monthly()}
				value={formatCurrency(summary.total_monthly)}
				valueColor="purple"
			/>
			<SummaryCard
				title={m.expected_expenses_total_misc()}
				value={formatCurrency(summary.total_misc)}
				valueColor="warning"
			/>
		{/if}
	</div>

<!-- Delete Confirmation Modal -->
{#if deleteConfirmExpense}
	<div
		class="bg-surface-overlay fixed top-0 left-0 z-50 flex h-screen w-screen items-center justify-center p-4"
		onclick={handleModalBackdropClick}
		onkeydown={(e) => e.key === 'Escape' && handleCancelDelete()}
		role="dialog"
		aria-modal="true"
		aria-labelledby="delete-modal-title"
		tabindex="-1"
	>
			<div class="bg-surface w-full max-w-md rounded-lg p-6 shadow-xl">
				<div class="mb-4 flex items-center">
					<div
						class="bg-danger-dark flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-full"
					>
						<AlertTriangleIcon class="text-danger h-6 w-6" />
					</div>
					<h3 id="delete-modal-title" class="text-text-primary ml-4 text-lg font-medium">
						{m.delete_confirm_title()}
					</h3>
				</div>
				<p class="text-text-tertiary mb-6 text-sm">
					{m.delete_confirm_description()}
				</p>
				<div class="flex flex-col-reverse gap-3 sm:flex-row sm:justify-end">
					<Button
						onclick={handleCancelDelete}
						disabled={deleting}
						variant="secondary"
					>
						{m.common_cancel()}
					</Button>
					<Button
						onclick={handleConfirmDelete}
						disabled={deleting}
						loading={deleting}
						variant="danger"
					>
						{m.common_delete()}
					</Button>
				</div>
			</div>
		</div>
	{/if}

	<!-- Expenses List Section -->
	<div class="bg-surface rounded-lg shadow">
		<div
			class="border-border flex flex-col gap-4 border-b px-4 py-4 sm:flex-row sm:items-center sm:justify-between sm:px-6"
		>
			<div class="flex items-center space-x-4">
				<h2 class="text-text-primary text-lg font-semibold">
					{activeTab === 'all' ? 'All' : activeTab.charAt(0).toUpperCase() + activeTab.slice(1)} Expenses
				</h2>
				<span
					class="bg-surface-light text-text-tertiary inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium"
				>
					{expenseCount}
					{expenseCount === 1 ? 'item' : 'items'}
				</span>
			</div>
			<div class="flex flex-col items-start gap-3 sm:flex-row sm:items-center">
				<SourceFilter
					sources={distinctSources}
					activeSource={sourceFilter}
					onFilterChange={(source) => actualExpensesStore.setSourceFilter(source)}
				/>
				<ExpenseTabs {activeTab} onTabChange={handleTabChange} showMisc={true} />
			</div>
		</div>

		<div class="p-4 sm:p-6">
			<ActualExpenseTable
				{expenses}
				loading={loading && expenses.length === 0}
				onEdit={handleEditExpense}
				onDelete={handleDeleteClick}
				filterType={activeTab as ExpenseFilterTypeEnum}
			/>
		</div>
	</div>
</div>
