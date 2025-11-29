// place files you want to import through the `$lib` alias in this folder.
export * from './utils';
export * from './types/enums';

// Stores (using .svelte.ts extension for Svelte 5 runes)
export {
	budgetStore,
	type Budget,
	type CreateBudgetRequest,
	type UpdateBudgetRequest
} from './stores/budget.svelte';
export {
	expectedExpensesStore,
	type ExpectedExpense,
	type ExpectedExpenseInput,
	type ExpectedExpenseFilterType,
	type ExpectedExpensesState
} from './stores/expectedExpenses.svelte';
export {
	actualExpensesStore,
	type ActualExpense,
	type ActualExpenseInput,
	type ActualExpenseSummary,
	type ActualExpenseFilterType
} from './stores/actualExpenses.svelte';
export {
	receiptStore,
	getReceiptStore,
	setReceiptStore,
	createReceiptStore,
	type ExtractedItem,
	type ReceiptState,
	type ReceiptStore
} from './stores/receipt.svelte';
export { toastStore, type Toast, type ToastType } from './stores/toast.svelte';

// Components
export { default as Button } from './components/Button.svelte';
export { default as DataTable } from './components/DataTable.svelte';
export { default as Dialog } from './components/Dialog.svelte';
export { default as SourceFilter } from './components/SourceFilter.svelte';
export { default as SummaryCard } from './components/SummaryCard.svelte';
export { default as YearSelector } from './components/YearSelector.svelte';
