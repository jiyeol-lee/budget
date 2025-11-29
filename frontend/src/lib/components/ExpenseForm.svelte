<script lang="ts">
	import type { ExpectedExpense, ExpectedExpenseInput } from '$lib/stores/expectedExpenses.svelte';
	import { ExpenseTypeEnum } from '$lib/types/enums';
	import { Button } from '$lib';
	import { AlertCircleIcon } from 'lucide-svelte';
	import * as m from '$lib/paraglide/messages';

	interface Props {
		expense?: ExpectedExpense | null;
		onSubmit: (data: ExpectedExpenseInput) => void;
		onCancel: () => void;
	}

	let { expense = null, onSubmit, onCancel }: Props = $props();

	// Form state
	let item_name = $state(expense?.item_name ?? '');
	let source = $state(expense?.source ?? '');
	let expected_amount = $state(expense?.expected_amount?.toString() ?? '');
	let expense_type = $state<ExpenseTypeEnum.WEEKLY | ExpenseTypeEnum.MONTHLY>(
		expense?.expense_type ?? ExpenseTypeEnum.WEEKLY
	);

	// Validation state
	let errors = $state<Record<string, string>>({});
	let touched = $state<Record<string, boolean>>({});

	// Check if form is valid (for disabling submit button)
	let isFormValid = $derived(() => {
		const nameValid = item_name.trim().length >= 2;
		const sourceValid = source.trim().length > 0;
		const amountValid =
			expected_amount.trim().length > 0 &&
			!isNaN(parseFloat(expected_amount)) &&
			parseFloat(expected_amount) > 0;
		return nameValid && sourceValid && amountValid;
	});

	/**
	 * Validate a single field
	 */
	function validateField(field: string, value: string): string {
		switch (field) {
			case 'item_name':
				if (!value.trim()) return 'Item name is required';
				if (value.trim().length < 2) return 'Item name must be at least 2 characters';
				break;
			case 'source':
				if (!value.trim()) return 'Source is required';
				break;
			case 'expected_amount': {
				if (!value.trim()) return 'Expected amount is required';
				const amount = parseFloat(value);
				if (isNaN(amount)) return 'Please enter a valid number';
				if (amount <= 0) return 'Amount must be greater than 0';
				break;
			}
		}
		return '';
	}

	/**
	 * Validate all fields
	 */
	function validateForm(): boolean {
		const newErrors: Record<string, string> = {};

		newErrors.item_name = validateField('item_name', item_name);
		newErrors.source = validateField('source', source);
		newErrors.expected_amount = validateField('expected_amount', expected_amount);

		errors = newErrors;

		return !Object.values(newErrors).some((error) => error !== '');
	}

	/**
	 * Handle field blur for validation
	 */
	function handleBlur(field: string) {
		touched = { ...touched, [field]: true };
		const value = field === 'item_name' ? item_name : field === 'source' ? source : expected_amount;
		errors = { ...errors, [field]: validateField(field, value) };
	}

	/**
	 * Handle form submission
	 */
	function handleSubmit(event: Event) {
		event.preventDefault();

		// Mark all fields as touched
		touched = {
			item_name: true,
			source: true,
			expected_amount: true
		};

		if (!validateForm()) {
			return;
		}

		const data: ExpectedExpenseInput = {
			item_name: item_name.trim(),
			source: source.trim(),
			expected_amount: parseFloat(expected_amount),
			expense_type
		};

		onSubmit(data);
	}

	/**
	 * Handle amount input to allow only valid currency values
	 */
	function handleAmountInput(event: Event) {
		const input = event.target as HTMLInputElement;
		// Allow only numbers and decimal point
		let value = input.value.replace(/[^0-9.]/g, '');

		// Ensure only one decimal point
		const parts = value.split('.');
		if (parts.length > 2) {
			value = parts[0] + '.' + parts.slice(1).join('');
		}

		// Limit to 2 decimal places
		const decimalPart = parts[1];
		if (parts.length === 2 && decimalPart && decimalPart.length > 2) {
			value = parts[0] + '.' + decimalPart.slice(0, 2);
		}

		expected_amount = value;
	}

	/**
	 * Get input class based on error state
	 */
	function getInputClass(hasError: boolean): string {
		const baseClass =
			'mt-1 block w-full px-3 py-2 border rounded-md shadow-sm sm:text-sm bg-surface-light text-text-primary border-border';
		if (hasError) {
			return `${baseClass} border-border-error bg-danger-dark/20 focus:ring-danger focus:border-border-error`;
		}
		return `${baseClass} focus:ring-primary focus:border-border-focus`;
	}
</script>

<form onsubmit={handleSubmit} class="space-y-4">
	<!-- Item Name -->
	<div>
		<label for="item_name" class="text-text-secondary block text-sm font-medium">
			{m.expense_form_item_name()} <span class="text-danger">*</span>
		</label>
		<input
			type="text"
			id="item_name"
			bind:value={item_name}
			onblur={() => handleBlur('item_name')}
			placeholder={m.expense_form_item_name_placeholder()}
			class={getInputClass(!!touched.item_name && !!errors.item_name)}
		/>
		{#if touched.item_name && errors.item_name}
			<p class="text-danger mt-1 flex items-center text-sm">
				<AlertCircleIcon class="mr-1 h-4 w-4 shrink-0" />
				{errors.item_name}
			</p>
		{/if}
	</div>

	<!-- Source -->
	<div>
		<label for="source" class="text-text-secondary block text-sm font-medium">
			{m.expense_form_source()} <span class="text-danger">*</span>
		</label>
		<input
			type="text"
			id="source"
			bind:value={source}
			onblur={() => handleBlur('source')}
			placeholder={m.expense_form_source_placeholder()}
			class={getInputClass(!!touched.source && !!errors.source)}
		/>
		{#if touched.source && errors.source}
			<p class="text-danger mt-1 flex items-center text-sm">
				<AlertCircleIcon class="mr-1 h-4 w-4 shrink-0" />
				{errors.source}
			</p>
		{/if}
	</div>

	<!-- Expected Amount -->
	<div>
		<label for="expected_amount" class="text-text-secondary block text-sm font-medium">
			{m.expense_form_amount()} <span class="text-danger">*</span>
		</label>
		<div class="relative mt-1 rounded-md shadow-sm">
			<div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
				<span class="text-text-secondary sm:text-sm">$</span>
			</div>
			<input
				type="text"
				inputmode="decimal"
				id="expected_amount"
				value={expected_amount}
				oninput={handleAmountInput}
				onblur={() => handleBlur('expected_amount')}
				placeholder="0.00"
				class="bg-surface-light text-text-primary border-border block w-full rounded-md border py-2 pr-3 pl-7 sm:text-sm
					{touched.expected_amount && errors.expected_amount
					? 'border-border-error bg-danger-dark/20 focus:ring-danger focus:border-border-error'
					: 'focus:ring-primary focus:border-border-focus'}"
			/>
		</div>
		{#if touched.expected_amount && errors.expected_amount}
			<p class="text-danger mt-1 flex items-center text-sm">
				<AlertCircleIcon class="mr-1 h-4 w-4 shrink-0" />
				{errors.expected_amount}
			</p>
		{/if}
	</div>

	<!-- Expense Type -->
	<div>
		<label for="expense_type" class="text-text-secondary block text-sm font-medium">
			{m.expense_form_expense_type()} <span class="text-danger">*</span>
		</label>
		<select
			id="expense_type"
			bind:value={expense_type}
			class="border-border focus:ring-primary focus:border-border-focus bg-surface-light text-text-primary mt-1 block w-full rounded-md border px-3 py-2 shadow-sm sm:text-sm"
		>
			<option value={ExpenseTypeEnum.WEEKLY}>{m.expenses_weekly()}</option>
			<option value={ExpenseTypeEnum.MONTHLY}>{m.expenses_monthly()}</option>
		</select>
	</div>

	<!-- Form Actions -->
	<div class="border-border flex flex-col-reverse gap-3 border-t pt-4 sm:flex-row sm:justify-end">
		<Button type="button" variant="secondary" onclick={onCancel} class="w-full sm:w-auto">
			{m.common_cancel()}
		</Button>
		<Button type="submit" variant="primary" disabled={!isFormValid()} class="w-full sm:w-auto">
			{m.common_save()}
		</Button>
	</div>
</form>
