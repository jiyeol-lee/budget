/**
 * Actual Expenses Store for Budget Tracker
 * Manages real spending data from receipts with monthly filtering
 */

import { api } from '$lib/utils/api';
import { ExpenseTypeEnum, ExpenseFilterTypeEnum } from '$lib/types/enums';

export interface ActualExpense {
	id: number;
	item_name: string;
	source: string;
	actual_amount: number;
	expense_type: ExpenseTypeEnum;
	item_code?: string;
	expected_expense_id?: number;
	receipt_date: string;
	receipt_number: number;
	month: number;
	year: number;
	created_at: string;
	updated_at: string;
}

export interface ActualExpenseInput {
	item_name: string;
	source: string;
	actual_amount: number;
	expense_type: ExpenseTypeEnum;
	item_code?: string;
	expected_expense_id?: number;
	receipt_date?: string;
	receipt_number?: number;
}

export interface ActualExpenseSummary {
	month: number;
	year: number;
	total_weekly: number;
	total_monthly: number;
	total_misc: number;
	total_actual: number;
}

export type ActualExpenseFilterType = `${ExpenseFilterTypeEnum}`;

interface ActualExpenseListResponse {
	expenses: ActualExpense[];
	total: number;
}

function createActualExpensesStore() {
	let expenses = $state<ActualExpense[]>([]);
	let loading = $state(false);
	let error = $state<string | null>(null);
	let currentMonth = $state(new Date().getMonth() + 1);
	let currentYear = $state(new Date().getFullYear());
	let summary = $state<ActualExpenseSummary | null>(null);
	let filterType = $state<ActualExpenseFilterType>(ExpenseFilterTypeEnum.ALL);
	let sourceFilter = $state<string | null>(null);

	async function fetchExpenses(
		month?: number,
		year?: number,
		type?: ActualExpenseFilterType
	): Promise<void> {
		loading = true;
		error = null;
		try {
			const m = month ?? currentMonth;
			const y = year ?? currentYear;
			const t = type ?? filterType;

			let url = `/actual-expenses?month=${m}&year=${y}`;
			if (t && t !== ExpenseFilterTypeEnum.ALL) {
				url += `&type=${t}`;
			}

			const response = await api.get<ActualExpenseListResponse>(url);
			expenses = response.expenses || [];
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to fetch actual expenses';
			expenses = [];
		} finally {
			loading = false;
		}
	}

	async function fetchSummary(month?: number, year?: number): Promise<void> {
		try {
			const m = month ?? currentMonth;
			const y = year ?? currentYear;
			summary = await api.get<ActualExpenseSummary>(
				`/actual-expenses/summary?month=${m}&year=${y}`
			);
		} catch (e) {
			console.error('Failed to fetch summary:', e);
			summary = null;
		}
	}

	async function fetchNextReceiptNumber(): Promise<number> {
		try {
			const response = await api.get<{ next_receipt_number: number }>(
				'/actual-expenses/next-receipt-number'
			);
			return response.next_receipt_number;
		} catch (e) {
			console.error('Failed to fetch next receipt number:', e);
			// Return 1 as fallback if the API fails
			return 1;
		}
	}

	async function createExpense(input: ActualExpenseInput): Promise<ActualExpense | null> {
		loading = true;
		error = null;
		try {
			const newExpense = await api.post<ActualExpense>('/actual-expenses', input);
			expenses = [newExpense, ...expenses];
			// Refresh summary
			await fetchSummary();
			return newExpense;
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to create actual expense';
			return null;
		} finally {
			loading = false;
		}
	}

	async function createBatch(inputs: ActualExpenseInput[]): Promise<ActualExpense[]> {
		loading = true;
		error = null;
		const created: ActualExpense[] = [];
		try {
			for (const input of inputs) {
				const expense = await api.post<ActualExpense>('/actual-expenses', input);
				created.push(expense);
			}
			// Refresh the list and summary
			await fetchExpenses();
			await fetchSummary();
			return created;
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to create actual expenses';
			return created;
		} finally {
			loading = false;
		}
	}

	async function updateExpense(
		id: number,
		input: Partial<ActualExpenseInput>
	): Promise<ActualExpense | null> {
		loading = true;
		error = null;
		try {
			const updated = await api.put<ActualExpense>(`/actual-expenses/${id}`, input);
			expenses = expenses.map((e) => (e.id === id ? updated : e));
			await fetchSummary();
			return updated;
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to update actual expense';
			return null;
		} finally {
			loading = false;
		}
	}

	async function deleteExpense(id: number): Promise<boolean> {
		loading = true;
		error = null;
		try {
			await api.delete(`/actual-expenses/${id}`);
			expenses = expenses.filter((e) => e.id !== id);
			await fetchSummary();
			return true;
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to delete actual expense';
			return false;
		} finally {
			loading = false;
		}
	}

	function setMonthYear(month: number, year: number): void {
		currentMonth = month;
		currentYear = year;
		sourceFilter = null; // Reset source filter when month/year changes
	}

	function setFilter(type: ActualExpenseFilterType): void {
		filterType = type;
	}

	function setSourceFilter(source: string | null): void {
		sourceFilter = source;
	}

	// Computed getters
	function getTotalByType(type: ExpenseTypeEnum): number {
		return expenses
			.filter((e) => e.expense_type === type)
			.reduce((sum, e) => sum + e.actual_amount, 0);
	}

	function getMonthlyTotal(): number {
		return expenses.reduce((sum, e) => sum + e.actual_amount, 0);
	}

	function getFilteredExpenses(): ActualExpense[] {
		let filtered = expenses;

		// Filter by expense type
		if (filterType !== ExpenseFilterTypeEnum.ALL) {
			filtered = filtered.filter((e) => e.expense_type === filterType);
		}

		// Filter by source
		if (sourceFilter !== null) {
			filtered = filtered.filter((e) => e.source === sourceFilter);
		}

		return filtered;
	}

	function getDistinctSources(): string[] {
		const sources = new Set(expenses.map((exp) => exp.source).filter(Boolean));
		return Array.from(sources).sort();
	}

	return {
		get expenses() {
			return expenses;
		},
		get loading() {
			return loading;
		},
		get error() {
			return error;
		},
		get currentMonth() {
			return currentMonth;
		},
		get currentYear() {
			return currentYear;
		},
		get summary() {
			return summary;
		},
		get filterType() {
			return filterType;
		},
		get sourceFilter() {
			return sourceFilter;
		},
		// Actions
		fetchExpenses,
		fetchSummary,
		fetchNextReceiptNumber,
		createExpense,
		createBatch,
		updateExpense,
		deleteExpense,
		setMonthYear,
		setFilter,
		setSourceFilter,
		// Computed
		getTotalByType,
		getMonthlyTotal,
		getFilteredExpenses,
		getDistinctSources
	};
}

export const actualExpensesStore = createActualExpensesStore();
