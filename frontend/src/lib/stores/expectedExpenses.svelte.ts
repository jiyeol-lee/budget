/**
 * Expected Expenses Store for Budget Tracker
 * Manages planned recurring expense data with support for WEEKLY/MONTHLY filtering
 */

import { api } from '$lib/utils/api';
import { ExpenseTypeEnum, ExpenseFilterTypeEnum } from '$lib/types/enums';

/**
 * Type for expected expense types (excludes MISC)
 */
export type ExpectedExpenseType = ExpenseTypeEnum.WEEKLY | ExpenseTypeEnum.MONTHLY;

/**
 * Expected Expense type interface matching backend model
 */
export interface ExpectedExpense {
	id: number;
	item_name: string;
	source: string;
	expected_amount: number;
	expense_type: ExpectedExpenseType;
	created_at: string;
	updated_at: string;
}

/**
 * Expected Expense input type for create/update operations
 */
export interface ExpectedExpenseInput {
	item_name: string;
	source: string;
	expected_amount: number;
	expense_type: ExpectedExpenseType;
}

/**
 * Filter type for expected expense queries
 */
export type ExpectedExpenseFilterType =
	| ExpenseFilterTypeEnum.ALL
	| ExpenseFilterTypeEnum.WEEKLY
	| ExpenseFilterTypeEnum.MONTHLY;

/**
 * Store state interface
 */
export interface ExpectedExpensesState {
	expenses: ExpectedExpense[];
	loading: boolean;
	error: string | null;
	filter: ExpectedExpenseFilterType;
}

/**
 * Backend response format for expected expense list
 */
interface ExpectedExpenseListResponse {
	expenses: ExpectedExpense[];
	filter: string;
	count: number;
}

/**
 * Create a reactive expected expenses store
 */
function createExpectedExpensesStore() {
	let expenses = $state<ExpectedExpense[]>([]);
	let loading = $state(false);
	let error = $state<string | null>(null);
	let filter = $state<ExpectedExpenseFilterType>(ExpenseFilterTypeEnum.ALL);
	let sourceFilter = $state<string | null>(null);

	/**
	 * Fetch expected expenses from API with optional type filter
	 */
	async function fetchExpenses(
		type: ExpectedExpenseFilterType = ExpenseFilterTypeEnum.ALL
	): Promise<void> {
		loading = true;
		error = null;
		filter = type;

		try {
			const endpoint =
				type === ExpenseFilterTypeEnum.ALL
					? '/expected-expenses'
					: `/expected-expenses?type=${type}`;
			const response = await api.get<ExpectedExpenseListResponse>(endpoint);
			// Backend returns { expenses: [], filter: "ALL", count: 0 }
			expenses = response?.expenses ?? [];
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to fetch expected expenses';
			expenses = [];
		} finally {
			loading = false;
		}
	}

	/**
	 * Create a new expected expense
	 */
	async function createExpense(input: ExpectedExpenseInput): Promise<ExpectedExpense | null> {
		loading = true;
		error = null;

		try {
			const newExpense = await api.post<ExpectedExpense, ExpectedExpenseInput>(
				'/expected-expenses',
				input
			);
			// Add to local state if it matches current filter
			if (filter === ExpenseFilterTypeEnum.ALL || filter === (input.expense_type as string)) {
				expenses = [...expenses, newExpense];
			}
			return newExpense;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to create expected expense';
			return null;
		} finally {
			loading = false;
		}
	}

	/**
	 * Update an existing expected expense
	 */
	async function updateExpense(
		id: number,
		input: ExpectedExpenseInput
	): Promise<ExpectedExpense | null> {
		loading = true;
		error = null;

		try {
			const updatedExpense = await api.put<ExpectedExpense, ExpectedExpenseInput>(
				`/expected-expenses/${id}`,
				input
			);
			// Update local state
			expenses = expenses.map((exp) => (exp.id === id ? updatedExpense : exp));
			// If filter changed and doesn't match current filter, remove from local state
			if (
				filter !== ExpenseFilterTypeEnum.ALL &&
				filter !== (updatedExpense.expense_type as string)
			) {
				expenses = expenses.filter((exp) => exp.id !== id);
			}
			return updatedExpense;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to update expected expense';
			return null;
		} finally {
			loading = false;
		}
	}

	/**
	 * Delete an expected expense
	 */
	async function deleteExpense(id: number): Promise<boolean> {
		loading = true;
		error = null;

		try {
			await api.delete(`/expected-expenses/${id}`);
			// Remove from local state
			expenses = expenses.filter((exp) => exp.id !== id);
			return true;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to delete expected expense';
			return false;
		} finally {
			loading = false;
		}
	}

	/**
	 * Set the filter type and refetch if needed
	 */
	async function setFilter(type: ExpectedExpenseFilterType): Promise<void> {
		if (filter !== type) {
			await fetchExpenses(type);
		}
	}

	/**
	 * Set the source filter
	 */
	function setSourceFilter(source: string | null): void {
		sourceFilter = source;
	}

	/**
	 * Clear error state
	 */
	function clearError(): void {
		error = null;
	}

	/**
	 * Get filtered expenses based on current filter and source filter
	 */
	function getFilteredExpenses(): ExpectedExpense[] {
		let result = expenses ?? [];

		// Apply expense type filter
		if (filter !== ExpenseFilterTypeEnum.ALL) {
			result = result.filter((exp) => (exp.expense_type as string) === filter);
		}

		// Apply source filter
		if (sourceFilter !== null) {
			result = result.filter((exp) => exp.source === sourceFilter);
		}

		return result;
	}

	/**
	 * Get distinct sources from current expenses, sorted alphabetically
	 */
	function getDistinctSources(): string[] {
		const sources = new Set((expenses ?? []).map((exp) => exp.source));
		return Array.from(sources).sort((a, b) => a.localeCompare(b));
	}

	/**
	 * Calculate total weekly expenses
	 */
	function getTotalWeekly(): number {
		return (expenses ?? [])
			.filter((exp) => exp.expense_type === ExpenseTypeEnum.WEEKLY)
			.reduce((sum, exp) => sum + exp.expected_amount, 0);
	}

	/**
	 * Calculate total monthly expenses
	 */
	function getTotalMonthly(): number {
		return (expenses ?? [])
			.filter((exp) => exp.expense_type === ExpenseTypeEnum.MONTHLY)
			.reduce((sum, exp) => sum + exp.expected_amount, 0);
	}

	/**
	 * Calculate estimated monthly total (weekly * 4 + monthly)
	 */
	function getEstimatedMonthlyTotal(): number {
		return getTotalWeekly() * 4 + getTotalMonthly();
	}

	return {
		// State getters (using $state requires getter functions or direct access)
		get expenses() {
			return expenses ?? [];
		},
		get loading() {
			return loading;
		},
		get error() {
			return error;
		},
		get filter() {
			return filter;
		},
		get sourceFilter() {
			return sourceFilter;
		},
		// Actions
		fetchExpenses,
		createExpense,
		updateExpense,
		deleteExpense,
		setFilter,
		setSourceFilter,
		clearError,
		// Computed values
		getFilteredExpenses,
		getDistinctSources,
		getTotalWeekly,
		getTotalMonthly,
		getEstimatedMonthlyTotal
	};
}

// Export a singleton store instance
export const expectedExpensesStore = createExpectedExpensesStore();
