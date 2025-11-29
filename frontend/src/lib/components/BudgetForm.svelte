<script lang="ts">
	import { budgetStore, type Budget, type CreateBudgetRequest } from '$lib/stores/budget.svelte';
	import { toastStore } from '$lib/stores/toast.svelte';
	import { getMonths, Button } from '$lib';
	import { AlertCircleIcon } from 'lucide-svelte';
	import * as m from '$lib/paraglide/messages';

	interface Props {
		budget?: Budget | null;
		onSave?: () => void;
		onCancel?: () => void;
	}

	let { budget = null, onSave, onCancel }: Props = $props();

	// Form state
	let month = $state(budget?.month ?? new Date().getMonth() + 1);
	let year = $state(budget?.year ?? new Date().getFullYear());
	let amount = $state(budget?.amount ?? 0);
	let notificationThreshold = $state(
		budget?.notification_threshold ? budget.notification_threshold * 100 : 80
	);
	let submitting = $state(false);

	// Track which fields have been touched (for on-blur validation)
	let touched = $state<{
		month: boolean;
		year: boolean;
		amount: boolean;
		threshold: boolean;
	}>({
		month: false,
		year: false,
		amount: false,
		threshold: false
	});

	// Validation errors
	let errors = $state<{
		month?: string;
		year?: string;
		amount?: string;
		threshold?: string;
	}>({});

	// Computed values
	let isEditMode = $derived(budget !== null);

	// Check if form is valid for submit button
	let isFormValid = $derived(() => {
		const monthValid = month >= 1 && month <= 12;
		const yearValid = year >= 2000 && year <= 2100;
		const amountValid = amount > 0;
		const thresholdValid = notificationThreshold >= 0 && notificationThreshold <= 100;
		return monthValid && yearValid && amountValid && thresholdValid;
	});

	/**
	 * Validate a single field
	 */
	function validateField(field: keyof typeof errors): string | undefined {
		switch (field) {
			case 'month':
				if (!month || month < 1 || month > 12) {
					return 'Please select a valid month';
				}
				break;
			case 'year':
				if (!year || year < 2000 || year > 2100) {
					return 'Please enter a valid year (2000-2100)';
				}
				break;
			case 'amount':
				if (!amount || amount <= 0) {
					return 'Amount must be greater than 0';
				}
				break;
			case 'threshold':
				if (notificationThreshold < 0 || notificationThreshold > 100) {
					return 'Threshold must be between 0 and 100';
				}
				break;
		}
		return undefined;
	}

	/**
	 * Handle field blur - validate on blur
	 */
	function handleBlur(field: keyof typeof touched) {
		touched = { ...touched, [field]: true };
		errors = { ...errors, [field]: validateField(field) };
	}

	/**
	 * Validate all fields
	 */
	function validate(): boolean {
		// Mark all fields as touched
		touched = {
			month: true,
			year: true,
			amount: true,
			threshold: true
		};

		errors = {
			month: validateField('month'),
			year: validateField('year'),
			amount: validateField('amount'),
			threshold: validateField('threshold')
		};

		return !Object.values(errors).some((error) => error !== undefined);
	}

	async function handleSubmit(event: Event) {
		event.preventDefault();

		if (!validate()) {
			return;
		}

		submitting = true;

		const data: CreateBudgetRequest = {
			month,
			year,
			amount,
			notification_threshold: notificationThreshold / 100 // Convert percentage to decimal
		};

		try {
			let result;
			if (isEditMode && budget) {
				result = await budgetStore.updateBudget(budget.id, data);
				if (result) {
					toastStore.success('Budget updated successfully');
				}
			} else {
				result = await budgetStore.createBudget(data);
				if (result) {
					toastStore.success('Budget created successfully');
				}
			}

			if (result) {
				onSave?.();
			} else if (budgetStore.error) {
				toastStore.error(budgetStore.error);
			}
		} catch {
			toastStore.error('An unexpected error occurred');
		} finally {
			submitting = false;
		}
	}

	function handleCancel() {
		onCancel?.();
	}

	/**
	 * Get input class based on error state
	 */
	function getInputClass(field: keyof typeof errors, baseClass: string): string {
		const hasError = touched[field] && errors[field];
		if (hasError) {
			return `${baseClass} border-border-error bg-danger-dark/20 focus:ring-danger focus:border-border-error`;
		}
		return `${baseClass} border-border bg-surface-light text-text-primary focus:ring-primary focus:border-border-focus`;
	}
</script>

<form onsubmit={handleSubmit} class="space-y-4">
	<!-- Month Select -->
	<div>
		<label for="month" class="text-text-secondary block text-sm font-medium">
			{m.budget_form_month()} <span class="text-danger">*</span>
		</label>
		<select
			id="month"
			bind:value={month}
			onblur={() => handleBlur('month')}
			class={getInputClass('month', 'mt-1 block w-full rounded-md border py-2 pr-10 pl-3')}
		>
			{#each getMonths() as mo (mo.value)}
				<option value={mo.value}>{mo.label}</option>
			{/each}
		</select>
		{#if touched.month && errors.month}
			<p class="text-danger mt-1 flex items-center text-sm">
				<AlertCircleIcon class="mr-1 h-4 w-4" />
				{errors.month}
			</p>
		{/if}
	</div>

	<!-- Year Input -->
	<div>
		<label for="year" class="text-text-secondary block text-sm font-medium">
			{m.budget_form_year()} <span class="text-danger">*</span>
		</label>
		<input
			type="number"
			id="year"
			bind:value={year}
			onblur={() => handleBlur('year')}
			min="2000"
			max="2100"
			class={getInputClass('year', 'mt-1 block w-full rounded-md border px-3 py-2')}
		/>
		{#if touched.year && errors.year}
			<p class="text-danger mt-1 flex items-center text-sm">
				<AlertCircleIcon class="mr-1 h-4 w-4" />
				{errors.year}
			</p>
		{/if}
	</div>

	<!-- Amount Input -->
	<div>
		<label for="amount" class="text-text-secondary block text-sm font-medium">
			{m.budget_form_amount()} <span class="text-danger">*</span>
		</label>
		<div class="relative mt-1 rounded-md shadow-sm">
			<div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
				<span class="text-text-secondary sm:text-sm">$</span>
			</div>
			<input
				type="number"
				id="amount"
				bind:value={amount}
				onblur={() => handleBlur('amount')}
				min="0"
				step="0.01"
				placeholder="0.00"
				class={getInputClass('amount', 'block w-full rounded-md border py-2 pr-4 pl-7')}
			/>
		</div>
		{#if touched.amount && errors.amount}
			<p class="text-danger mt-1 flex items-center text-sm">
				<AlertCircleIcon class="mr-1 h-4 w-4" />
				{errors.amount}
			</p>
		{/if}
	</div>

	<!-- Notification Threshold Input -->
	<div>
		<label for="threshold" class="text-text-secondary block text-sm font-medium">
			{m.budget_form_threshold()} <span class="text-danger">*</span>
		</label>
		<div class="relative mt-1 rounded-md shadow-sm">
			<input
				type="number"
				id="threshold"
				bind:value={notificationThreshold}
				onblur={() => handleBlur('threshold')}
				min="0"
				max="100"
				placeholder="80"
				class={getInputClass('threshold', 'block w-full rounded-md border py-2 pr-10')}
			/>
			<div class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-3">
				<span class="text-text-secondary sm:text-sm">%</span>
			</div>
		</div>
		<p class="text-text-tertiary mt-1 text-xs">
			{m.budget_form_threshold_description()}
		</p>
		{#if touched.threshold && errors.threshold}
			<p class="text-danger mt-1 flex items-center text-sm">
				<AlertCircleIcon class="mr-1 h-4 w-4" />
				{errors.threshold}
			</p>
		{/if}
	</div>

	<!-- Form Actions -->
	<div class="flex flex-col-reverse gap-3 pt-4 sm:flex-row sm:justify-end">
		<Button variant="secondary" onclick={handleCancel} class="w-full sm:w-auto">
			{m.common_cancel()}
		</Button>
		<Button type="submit" disabled={!isFormValid()} loading={submitting} class="w-full sm:w-auto">
			{submitting ? m.common_processing() : isEditMode ? m.budget_edit() : m.budget_add()}
		</Button>
	</div>
</form>
