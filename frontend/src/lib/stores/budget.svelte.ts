/**
 * Budget Store for Budget Tracker
 * Manages budget data with Svelte 5 runes
 */

import { get, post, put, del } from '$lib/utils/api';

/**
 * Budget type interface matching backend model
 */
export interface Budget {
	id: number;
	month: number;
	year: number;
	amount: number;
	notification_threshold: number;
	created_at: string;
	updated_at: string;
}

/**
 * Create budget request payload
 */
export interface CreateBudgetRequest {
	month: number;
	year: number;
	amount: number;
	notification_threshold: number;
}

/**
 * Update budget request payload
 */
export interface UpdateBudgetRequest {
	month?: number;
	year?: number;
	amount?: number;
	notification_threshold?: number;
}

/**
 * Budget store state
 */
interface BudgetStoreState {
	budgets: Budget[];
	loading: boolean;
	error: string | null;
}

/**
 * Create a reactive budget store using Svelte 5 patterns
 */
function createBudgetStore() {
	let state = $state<BudgetStoreState>({
		budgets: [],
		loading: false,
		error: null
	});

	return {
		// Getters
		get budgets() {
			return state.budgets;
		},
		get loading() {
			return state.loading;
		},
		get error() {
			return state.error;
		},

		/**
		 * Fetch all budgets from the API
		 */
		async fetchBudgets(): Promise<void> {
			state.loading = true;
			state.error = null;
			try {
				const budgets = await get<Budget[] | null>('/budgets');
				// Handle null response - backend may return null for empty list
				state.budgets = budgets ?? [];
			} catch (err) {
				state.error = err instanceof Error ? err.message : 'Failed to fetch budgets';
				state.budgets = [];
				console.error('Error fetching budgets:', err);
			} finally {
				state.loading = false;
			}
		},

		/**
		 * Create a new budget
		 */
		async createBudget(data: CreateBudgetRequest): Promise<Budget | null> {
			state.loading = true;
			state.error = null;
			try {
				const newBudget = await post<Budget, CreateBudgetRequest>('/budgets', data);
				state.budgets = [...state.budgets, newBudget];
				return newBudget;
			} catch (err) {
				state.error = err instanceof Error ? err.message : 'Failed to create budget';
				console.error('Error creating budget:', err);
				return null;
			} finally {
				state.loading = false;
			}
		},

		/**
		 * Update an existing budget
		 */
		async updateBudget(id: number, data: UpdateBudgetRequest): Promise<Budget | null> {
			state.loading = true;
			state.error = null;
			try {
				const updatedBudget = await put<Budget, UpdateBudgetRequest>(`/budgets/${id}`, data);
				state.budgets = state.budgets.map((b) => (b.id === id ? updatedBudget : b));
				return updatedBudget;
			} catch (err) {
				state.error = err instanceof Error ? err.message : 'Failed to update budget';
				console.error('Error updating budget:', err);
				return null;
			} finally {
				state.loading = false;
			}
		},

		/**
		 * Delete a budget
		 */
		async deleteBudget(id: number): Promise<boolean> {
			state.loading = true;
			state.error = null;
			try {
				await del(`/budgets/${id}`);
				state.budgets = state.budgets.filter((b) => b.id !== id);
				return true;
			} catch (err) {
				state.error = err instanceof Error ? err.message : 'Failed to delete budget';
				console.error('Error deleting budget:', err);
				return false;
			} finally {
				state.loading = false;
			}
		},

		/**
		 * Get budget for a specific month/year
		 */
		getBudgetByMonthYear(month: number, year: number): Budget | undefined {
			return state.budgets?.find((b) => b.month === month && b.year === year);
		},

		/**
		 * Get current month's budget
		 */
		getCurrentMonthBudget(): Budget | undefined {
			const now = new Date();
			return this.getBudgetByMonthYear(now.getMonth() + 1, now.getFullYear());
		},

		/**
		 * Clear any errors
		 */
		clearError(): void {
			state.error = null;
		}
	};
}

// Export singleton instance
export const budgetStore = createBudgetStore();
