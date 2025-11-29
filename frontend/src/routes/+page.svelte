<script lang="ts">
	import { onMount } from 'svelte';
	import { get } from '$lib/utils/api';
	import {
		formatCurrency,
		formatMonthYear,
		formatPercentageValue,
		safeNumber
	} from '$lib/utils/format';
	import { getMonths } from '$lib/utils';
	import { expectedExpensesStore } from '$lib/stores/expectedExpenses.svelte';
	import { budgetStore } from '$lib/stores/budget.svelte';
	import { actualExpensesStore } from '$lib/stores/actualExpenses.svelte';
	import Skeleton from '$lib/components/Skeleton.svelte';
	import { ExpenseFilterTypeEnum } from '$lib/types/enums';
	import { Button, YearSelector } from '$lib';
	import {
		AlertCircleIcon,
		DollarSignIcon,
		CalendarIcon,
		CalculatorIcon,
		BarChart3Icon,
		PlusIcon,
		FileTextIcon,
		InfoIcon
	} from 'lucide-svelte';
	import * as m from '$lib/paraglide/messages';

	// Budget Status API Response Type
	interface BudgetStatus {
		current_budget: { amount: number } | null;
		total_spent: number;
		expected_total: number;
		percentage_used: number;
		status: 'safe' | 'warning' | 'danger' | 'over';
		message: string;
	}

	// Date Selection State
	let selectedMonth = $state(new Date().getMonth() + 1);
	let selectedYear = $state(new Date().getFullYear());

	// State
	let budgetStatus = $state<BudgetStatus | null>(null);
	let budgetStatusLoading = $state(true);
	let budgetStatusError = $state<string | null>(null);

	// Derived values from stores
	let expensesLoading = $derived(expectedExpensesStore.loading);
	let budgetsLoading = $derived(budgetStore.loading);
	let actualExpensesLoading = $derived(actualExpensesStore.loading);

	// Computed expense summaries
	let totalWeekly = $derived(safeNumber(expectedExpensesStore.getTotalWeekly()));
	let totalMonthly = $derived(safeNumber(expectedExpensesStore.getTotalMonthly()));

	// Use budgetStatus.total_spent if available to ensure consistency with the status card and budget page
	// Fallback to actualExpensesStore summary if budgetStatus is not yet loaded
	let actualSpent = $derived(
		budgetStatus
			? safeNumber(budgetStatus.total_spent)
			: safeNumber(actualExpensesStore.summary?.total_actual || 0)
	);

	// Check if current month has a budget
	let currentMonthBudget = $derived(budgetStore.getBudgetByMonthYear(selectedMonth, selectedYear));
	let hasBudgetForCurrentMonth = $derived(currentMonthBudget !== undefined);
	let budgetAmount = $derived(safeNumber(currentMonthBudget?.amount || 0));
	let remainingBudget = $derived(budgetAmount - actualSpent);

	// Overall loading state
	let isLoading = $derived(
		budgetStatusLoading || expensesLoading || budgetsLoading || actualExpensesLoading
	);

	/**
	 * Get status color classes
	 */
	function getStatusColors(status: string): {
		bg: string;
		border: string;
		text: string;
		progress: string;
		badge: string;
	} {
		switch (status) {
			case 'over':
				return {
					bg: 'bg-budget-over-bg',
					border: 'border-budget-over',
					text: 'text-budget-over',
					progress: 'bg-budget-over',
					badge: 'bg-danger-dark text-danger-light'
				};
			case 'danger':
				return {
					bg: 'bg-budget-danger-bg',
					border: 'border-budget-danger',
					text: 'text-budget-danger',
					progress: 'bg-budget-danger',
					badge: 'bg-danger-dark text-danger-light'
				};
			case 'warning':
				return {
					bg: 'bg-budget-warning-bg',
					border: 'border-budget-warning',
					text: 'text-budget-warning',
					progress: 'bg-budget-warning',
					badge: 'bg-warning-dark text-warning-light'
				};
			default:
				return {
					bg: 'bg-budget-safe-bg',
					border: 'border-budget-safe',
					text: 'text-budget-safe',
					progress: 'bg-budget-safe',
					badge: 'bg-success-dark text-success-light'
				};
		}
	}

	/**
	 * Get status label
	 */
	function getStatusLabel(status: string): string {
		switch (status) {
			case 'over':
				return m.budget_status_over();
			case 'danger':
				return m.budget_status_critical();
			case 'warning':
				return m.budget_status_warning();
			default:
				return m.budget_status_on_track();
		}
	}

	/**
	 * Fetch budget status from API
	 */
	async function fetchBudgetStatus(
		month: number = selectedMonth,
		year: number = selectedYear
	): Promise<void> {
		budgetStatusLoading = true;
		budgetStatusError = null;

		try {
			const data = await get<BudgetStatus>(
				`/notifications/budget-status?month=${month}&year=${year}`
			);
			budgetStatus = data;
		} catch (err) {
			// If 404, no budget for current month - this is expected
			if (err instanceof Error && err.message.includes('404')) {
				budgetStatus = null;
			} else {
				budgetStatusError = err instanceof Error ? err.message : 'Failed to fetch budget status';
			}
		} finally {
			budgetStatusLoading = false;
		}
	}

	/**
	 * Handle date change
	 */
	async function handleDateChange(month: number, year: number) {
		selectedMonth = month;
		selectedYear = year;

		// Sync with actual expenses store so other pages stay in sync
		actualExpensesStore.setMonthYear(month, year);

		await Promise.all([
			fetchBudgetStatus(month, year),
			actualExpensesStore.fetchSummary(month, year)
		]);
	}

	// Fetch all data on mount
	onMount(() => {
		// Initialize store with current selection if needed
		actualExpensesStore.setMonthYear(selectedMonth, selectedYear);

		fetchBudgetStatus();
		expectedExpensesStore.fetchExpenses(ExpenseFilterTypeEnum.ALL);
		budgetStore.fetchBudgets();
		actualExpensesStore.fetchSummary(selectedMonth, selectedYear);
	});
</script>

<svelte:head>
	<title>Dashboard | Budget Tracker</title>
</svelte:head>

<div class="space-y-4 sm:space-y-6">
	<!-- Page Header -->
	<div class="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
		<div>
			<h1 class="text-text-primary text-xl font-bold sm:text-2xl">{m.dashboard_title()}</h1>
			<p class="text-text-secondary mt-1 text-sm">
				{m.dashboard_description()}
			</p>
		</div>

		<!-- Month/Year Selector -->
		<div class="flex space-x-2">
			<select
				value={selectedMonth}
				onchange={(e) => handleDateChange(parseInt(e.currentTarget.value), selectedYear)}
				class="border-input-border bg-input-bg text-text-primary focus:border-input-focus focus:ring-input-focus block w-full rounded-md shadow-sm sm:text-sm"
			>
				{#each getMonths() as mo (mo.value)}
					<option value={mo.value}>{mo.label}</option>
				{/each}
			</select>
			<YearSelector
				value={selectedYear}
				onYearChange={(year) => handleDateChange(selectedMonth, year)}
				id="dashboard-year-selector"
			/>
		</div>
	</div>

	<!-- Budget Status Card -->
	<div class="bg-surface rounded-lg shadow">
		<div class="border-border flex items-center justify-between border-b px-4 py-4 sm:px-6">
			<h2 class="text-text-primary text-lg font-semibold">{m.dashboard_budget_status()}</h2>
			<a href="/budget" class="text-primary hover:text-primary-hover text-sm font-medium">
				{m.dashboard_manage_budget()}
			</a>
		</div>
		<div class="p-4 sm:p-6">
			{#if budgetStatusLoading}
				<!-- Loading Skeleton -->
				<div class="animate-pulse space-y-4">
					<div class="flex items-center justify-between">
						<div class="bg-skeleton h-4 w-24 rounded"></div>
						<div class="bg-skeleton h-6 w-20 rounded"></div>
					</div>
					<div class="bg-skeleton h-3 w-full rounded"></div>
					<div class="flex justify-between">
						<div class="bg-skeleton h-4 w-32 rounded"></div>
						<div class="bg-skeleton h-4 w-32 rounded"></div>
					</div>
					<div class="bg-skeleton h-4 w-3/4 rounded"></div>
				</div>
			{:else if budgetStatusError}
				<!-- Error State -->
				<div class="py-4 text-center">
					<AlertCircleIcon class="text-danger mx-auto h-8 w-8" />
					<p class="text-danger mt-2 text-sm">{budgetStatusError}</p>
					<Button variant="link" onclick={() => fetchBudgetStatus()} class="mt-2">
						{m.common_retry()}
					</Button>
				</div>
			{:else if !budgetStatus}
				<!-- No Budget Set -->
				<div class="py-6 text-center">
					<DollarSignIcon class="text-text-tertiary mx-auto h-12 w-12" />
					<h3 class="text-text-primary mt-2 text-sm font-medium">
						{m.dashboard_no_budget_title()}
					</h3>
					<p class="text-text-secondary mt-1 text-sm">
						{m.dashboard_no_budget_description()}
					</p>
					<a
						href="/budget"
						class="bg-primary text-text-on-primary hover:bg-primary-hover mt-4 inline-flex items-center rounded-md px-4 py-2 text-sm font-medium transition-colors"
					>
						{m.dashboard_set_budget()}
					</a>
				</div>
			{:else}
				<!-- Budget Status Display -->
				{@const colors = getStatusColors(budgetStatus.status)}
				{@const spent = safeNumber(budgetStatus.total_spent)}
				{@const budgetAmount = safeNumber(budgetStatus.current_budget?.amount || 0)}
				{@const remaining = budgetAmount - spent}
				{@const percentageUsed = safeNumber(budgetStatus.percentage_used)}
				<div class="space-y-4">
					<!-- Status Header -->
					<div class="flex items-center justify-between">
						<span class="text-text-secondary text-sm">
							{formatMonthYear(selectedMonth, selectedYear)}
						</span>
						<span
							class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium {colors.badge}"
						>
							{getStatusLabel(budgetStatus.status)}
						</span>
					</div>

					<!-- Progress Bar -->
					<div>
						<div class="mb-1 flex justify-between text-sm">
							<span class="text-text-tertiary font-medium">
								{formatCurrency(spent)}
								{m.amount_spent()}
							</span>
							<span class="text-text-secondary">
								{m.stats_of_budget({ amount: formatCurrency(budgetAmount) })}
							</span>
						</div>
						<div class="bg-surface-light h-3 w-full overflow-hidden rounded-full">
							<div
								class="{colors.progress} h-3 rounded-full transition-all duration-500"
								style="width: {Math.min(percentageUsed, 100)}%"
							></div>
						</div>
					</div>

					<!-- Summary Stats -->
					<div class="flex justify-between text-sm">
						{#if remaining >= 0}
							<span class="text-text-tertiary">
								<span class="text-success font-medium">{formatCurrency(remaining)}</span>
								{m.amount_remaining()}
							</span>
						{:else}
							<span class="text-text-tertiary">
								<span class="text-danger font-medium">{formatCurrency(Math.abs(remaining))}</span>
								{m.amount_over_budget()}
							</span>
						{/if}
						<span class="text-text-secondary">
							{m.stats_used({ percentage: formatPercentageValue(percentageUsed) })}
						</span>
					</div>

					<!-- Message -->
					{#if budgetStatus.message}
						<div class="{colors.bg} {colors.border} rounded-lg border p-3">
							<p class="text-sm {colors.text}">{budgetStatus.message}</p>
						</div>
					{/if}
				</div>
			{/if}
		</div>
	</div>

	<!-- Quick Stats Grid -->
	<div class="grid grid-cols-1 gap-3 sm:grid-cols-2 sm:gap-4 lg:grid-cols-4">
		{#if isLoading}
			<!-- Loading Skeletons -->
			<Skeleton variant="card" count={4} />
		{:else}
			<!-- Weekly Expenses -->
			<div class="bg-surface rounded-lg p-4 shadow sm:p-6">
				<div class="flex items-center">
					<div class="bg-info-dark rounded-lg p-2">
						<CalendarIcon class="text-info h-5 w-5 sm:h-6 sm:w-6" />
					</div>
					<h3 class="text-text-secondary ml-3 text-sm font-medium">
						{m.stats_weekly_expected()}
					</h3>
				</div>
				<p class="text-text-primary mt-3 text-xl font-bold sm:mt-4 sm:text-2xl">
					{formatCurrency(totalWeekly)}
				</p>
				<p class="text-text-secondary mt-1 text-xs">{m.expenses_weekly()}</p>
			</div>

			<!-- Monthly Expenses -->
			<div class="bg-surface rounded-lg p-4 shadow sm:p-6">
				<div class="flex items-center">
					<div class="rounded-lg bg-purple-900/30 p-2">
						<CalculatorIcon class="h-5 w-5 text-purple-400 sm:h-6 sm:w-6" />
					</div>
					<h3 class="text-text-secondary ml-3 text-sm font-medium">
						{m.stats_monthly_expected()}
					</h3>
				</div>
				<p class="text-text-primary mt-3 text-xl font-bold sm:mt-4 sm:text-2xl">
					{formatCurrency(totalMonthly)}
				</p>
				<p class="text-text-secondary mt-1 text-xs">{m.expenses_monthly()}</p>
			</div>

			<!-- Actual Spent -->
			<div class="bg-surface rounded-lg p-4 shadow sm:p-6">
				<div class="flex items-center">
					<div class="bg-warning-dark rounded-lg p-2">
						<DollarSignIcon class="text-warning h-5 w-5 sm:h-6 sm:w-6" />
					</div>
					<h3 class="text-text-secondary ml-3 text-sm font-medium">{m.stats_actual_spent()}</h3>
				</div>
				<p class="text-text-primary mt-3 text-xl font-bold sm:mt-4 sm:text-2xl">
					{formatCurrency(actualSpent)}
				</p>
				<p class="text-text-secondary mt-1 text-xs">{m.expenses_monthly()}</p>
			</div>

			<!-- Remaining Budget -->
			<div class="bg-surface rounded-lg p-4 shadow sm:p-6">
				<div class="flex items-center">
					<div class="bg-success-dark rounded-lg p-2">
						<BarChart3Icon class="text-success h-5 w-5 sm:h-6 sm:w-6" />
					</div>
					<h3 class="text-text-secondary ml-3 text-sm font-medium">
						{m.stats_remaining_budget()}
					</h3>
				</div>
				<p
					class="mt-3 text-xl font-bold sm:mt-4 sm:text-2xl {remainingBudget < 0
						? 'text-danger'
						: 'text-text-primary'}"
				>
					{formatCurrency(remainingBudget)}
				</p>
				<p class="text-text-secondary mt-1 text-xs">{m.nav_budget()}</p>
			</div>
		{/if}
	</div>

	<!-- Quick Actions -->
	<div class="bg-surface rounded-lg p-4 shadow sm:p-6">
		<h2 class="text-text-primary mb-4 text-lg font-semibold">{m.dashboard_quick_actions()}</h2>
		<div class="grid grid-cols-1 gap-3 sm:grid-cols-2 sm:gap-4 lg:grid-cols-3">
			<!-- Add Expense -->
			<a
				href="/expected-expenses"
				class="border-border hover:bg-hover hover:border-border-light group flex items-center rounded-lg border p-4 transition-colors"
			>
				<div
					class="bg-primary-dark group-hover:bg-primary-dark/70 rounded-lg p-3 transition-colors"
				>
					<PlusIcon class="text-primary h-6 w-6" />
				</div>
				<div class="ml-4">
					<h3 class="text-text-primary text-sm font-medium">
						{m.action_add_expected_expense()}
					</h3>
					<p class="text-text-secondary text-xs">{m.expected_expenses_description()}</p>
				</div>
			</a>

			<!-- Process Receipt -->
			<a
				href="/receipts/process"
				class="border-border hover:bg-hover hover:border-border-light group flex items-center rounded-lg border p-4 transition-colors"
			>
				<div class="rounded-lg bg-purple-900/30 p-3 transition-colors group-hover:bg-purple-900/50">
					<FileTextIcon class="h-6 w-6 text-purple-400" />
				</div>
				<div class="ml-4">
					<h3 class="text-text-primary text-sm font-medium">{m.action_process_receipt()}</h3>
					<p class="text-text-secondary text-xs">{m.receipt_description()}</p>
				</div>
			</a>

			<!-- Set Budget (conditional) -->
			{#if !hasBudgetForCurrentMonth}
				<a
					href="/budget"
					class="border-warning hover:bg-warning-dark/20 hover:border-warning-hover group flex items-center rounded-lg border p-4 transition-colors"
				>
					<div
						class="bg-warning-dark group-hover:bg-warning-dark/70 rounded-lg p-3 transition-colors"
					>
						<DollarSignIcon class="text-warning h-6 w-6" />
					</div>
					<div class="ml-4">
						<h3 class="text-warning-light text-sm font-medium">
							{m.action_set_monthly_budget()}
						</h3>
						<p class="text-warning text-xs">{m.dashboard_no_budget_title()}</p>
					</div>
				</a>
			{:else}
				<a
					href="/budget"
					class="border-border hover:bg-hover hover:border-border-light group flex items-center rounded-lg border p-4 transition-colors"
				>
					<div
						class="bg-success-dark group-hover:bg-success-dark/70 rounded-lg p-3 transition-colors"
					>
						<DollarSignIcon class="text-success h-6 w-6" />
					</div>
					<div class="ml-4">
						<h3 class="text-text-primary text-sm font-medium">{m.budget_title()}</h3>
						<p class="text-text-secondary text-xs">{m.budget_description()}</p>
					</div>
				</a>
			{/if}
		</div>
	</div>

	<!-- Tips Section -->
	<div class="bg-info-bg border-info rounded-lg border p-4">
		<div class="flex">
			<div class="shrink-0">
				<InfoIcon class="text-info h-5 w-5" />
			</div>
			<div class="ml-3">
				<h3 class="text-info text-sm font-medium">{m.dashboard_tips_title()}</h3>
				<div class="text-info mt-2 text-sm">
					<p>{m.dashboard_tips_description()}</p>
				</div>
			</div>
		</div>
	</div>
</div>
