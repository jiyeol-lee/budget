<script lang="ts">
	import DocumentUploader from '$lib/components/DocumentUploader.svelte';
	import ReceiptComparison from '$lib/components/ReceiptComparison.svelte';
	import ProcessingSpinner from '$lib/components/ProcessingSpinner.svelte';
	import { Button, Dialog } from '$lib';
	import { getReceiptStore, type ExtractedItem } from '$lib/stores/receipt.svelte';
	import { actualExpensesStore } from '$lib/stores/actualExpenses.svelte';
	import { toastStore } from '$lib/stores/toast.svelte';
	import { ExpenseTypeEnum } from '$lib/types/enums';
	import {
		CheckCircleIcon,
		XCircleIcon,
		XIcon,
		LightbulbIcon,
		RefreshCwIcon,
		PlusIcon
	} from 'lucide-svelte';
	import * as m from '$lib/paraglide/messages';

	const receiptStore = getReceiptStore();

	// UI state
	let showConfirmModal = $state(false);
	let addingToExpenses = $state(false);
	let successMessage = $state<string | null>(null);
	let errorMessage = $state<string | null>(null);
	let receiptDate = $state(new Date().toISOString().split('T')[0]);

	// Current step derived from store state
	let currentStep = $derived.by(() => {
		if (receiptStore.extractedItems.length > 0) return 4; // Results ready
		if (receiptStore.isProcessing) return 3; // Processing
		if (receiptStore.selectedImage) return 2; // Image selected, ready to process
		return 1; // Initial - upload
	});

	/**
	 * Handle image selection from uploader
	 */
	function handleImageSelected(file: File): void {
		receiptStore.setImage(file);
		clearMessages();
	}

	/**
	 * Trigger receipt processing
	 */
	async function handleProcessReceipt(): Promise<void> {
		clearMessages();
		const success = await receiptStore.processReceipt();
		if (!success && receiptStore.error) {
			toastStore.error(receiptStore.error);
		} else if (success) {
			toastStore.success(m.receipt_success());
		}
	}

	/**
	 * Update an item in the extracted list
	 */
	function handleUpdateItem(
		index: number,
		field: keyof ExtractedItem,
		value: string | number | boolean
	): void {
		receiptStore.updateItem(index, field, value);
	}

	/**
	 * Remove an item from the extracted list
	 */
	function handleRemoveItem(index: number): void {
		receiptStore.removeItem(index);
	}

	/**
	 * Add a new blank row
	 */
	function handleAddItem(): void {
		receiptStore.addItem();
	}

	/**
	 * Toggle item selection
	 */
	function handleToggleSelection(index: number): void {
		receiptStore.toggleItemSelection(index);
	}

	/**
	 * Toggle select all items
	 */
	function handleToggleSelectAll(selectAll: boolean): void {
		receiptStore.toggleSelectAll(selectAll);
	}

	/**
	 * Open confirmation modal for adding to expenses
	 */
	function handleAddToExpensesClick(): void {
		const selected = receiptStore.getSelectedItems();
		if (selected.length === 0) {
			errorMessage = 'Please select at least one item to add to expenses.';
			return;
		}
		showConfirmModal = true;
	}

	/**
	 * Confirm adding selected items to expenses
	 */
	async function handleConfirmAddToExpenses(): Promise<void> {
		addingToExpenses = true;
		clearMessages();

		try {
			// Fetch the next receipt number before saving
			const receiptNumber = await actualExpensesStore.fetchNextReceiptNumber();
			receiptStore.setReceiptNumber(receiptNumber);

			const selectedItems = receiptStore.getSelectedItems().filter((item) => item.item_price !== 0);

			const inputs = selectedItems.map((item) => ({
				item_name: item.item_name,
				source: item.source,
				actual_amount: item.item_price,
				expense_type: ExpenseTypeEnum[item.type.toUpperCase() as keyof typeof ExpenseTypeEnum],
				item_code: item.item_code || undefined,
				expected_expense_id: item.expected_expense_id,
				receipt_date: new Date(receiptDate || new Date()).toISOString(),
				receipt_number: receiptNumber
			}));

			const created = await actualExpensesStore.createBatch(inputs);
			const addedCount = created.length;
			const failedCount = inputs.length - addedCount;

			if (addedCount > 0) {
				toastStore.success(
					`Successfully added ${addedCount} item${addedCount > 1 ? 's' : ''} to Actual Expenses!`
				);
				receiptStore.clearAll();
			}

			if (failedCount > 0) {
				toastStore.error(`Failed to add ${failedCount} item${failedCount > 1 ? 's' : ''}.`);
			}
		} catch {
			errorMessage = 'An error occurred while adding expenses.';
		} finally {
			addingToExpenses = false;
			showConfirmModal = false;
		}
	}

	/**
	 * Cancel the confirmation modal
	 */
	function handleCancelModal(): void {
		showConfirmModal = false;
	}

	/**
	 * Clear and start over
	 */
	function handleClearAndStartOver(): void {
		receiptStore.clearAll();
		clearMessages();
	}

	/**
	 * Clear all messages
	 */
	function clearMessages(): void {
		successMessage = null;
		errorMessage = null;
	}

	/**
	 * Get selected count
	 */
	function getSelectedCount(): number {
		return receiptStore.getSelectedItems().length;
	}
</script>

<svelte:head>
	<title>Process Receipt | Budget Tracker</title>
</svelte:head>

<!-- Processing Overlay -->
{#if receiptStore.isProcessing}
	<ProcessingSpinner message={m.receipt_processing()} fullScreen={true} />
{/if}

<div class="space-y-4 sm:space-y-6">
	<!-- Page Header -->
	<div>
		<h1 class="text-text-primary text-xl font-bold sm:text-2xl">{m.receipt_title()}</h1>
		<p class="text-text-secondary mt-1 text-sm">
			{m.receipt_description()}
		</p>
	</div>

	<!-- Instructions Card -->
	<div class="bg-info-bg border-info rounded-lg border p-4">
		<h3 class="text-info font-medium">{m.receipt_instructions_title()}:</h3>
		<ol class="text-info mt-2 list-inside list-decimal space-y-1 text-sm">
			<li>{m.receipt_instructions_1()}</li>
			<li>{m.receipt_instructions_2()}</li>
			<li>{m.receipt_instructions_3()}</li>
			<li>{m.receipt_instructions_4()}</li>
		</ol>
	</div>

	<!-- Messages -->
	{#if successMessage}
		<div
			class="bg-success-dark border-success text-success-light flex items-center gap-3 rounded-lg border px-4 py-3"
		>
			<CheckCircleIcon class="h-5 w-5 shrink-0" />
			<span>{successMessage}</span>
			<Button
				variant="ghost"
				onclick={() => (successMessage = null)}
				class="text-success hover:text-success-hover ml-auto p-1"
				aria-label="Dismiss success message"
			>
				<XIcon class="h-4 w-4" />
			</Button>
		</div>
	{/if}

	{#if errorMessage}
		<div
			class="bg-danger-dark border-danger text-danger-light flex items-center gap-3 rounded-lg border px-4 py-3"
		>
			<XCircleIcon class="h-5 w-5 shrink-0" />
			<span>{errorMessage}</span>
			<Button
				variant="ghost"
				onclick={() => (errorMessage = null)}
				class="text-danger hover:text-danger-hover ml-auto p-1"
				aria-label="Dismiss error message"
			>
				<XIcon class="h-4 w-4" />
			</Button>
		</div>
	{/if}

	<!-- Step 1: Upload Section -->
	{#if currentStep === 1}
		<div class="bg-surface rounded-lg p-6 shadow">
			<h2 class="text-text-primary mb-4 text-lg font-semibold">{m.receipt_step_1()}</h2>
		<DocumentUploader
			onFileSelected={handleImageSelected}
			maxSizeMB={10}
			accept="application/pdf"
		/>
		</div>
	{/if}

	<!-- Step 2: Process Button -->
	{#if currentStep === 2}
		<div class="bg-surface rounded-lg p-6 shadow">
			<h2 class="text-text-primary mb-4 text-lg font-semibold">{m.receipt_step_2()}</h2>

			<!-- Preview of selected image -->
			<div class="mb-6">
				<div
					class="border-border bg-surface-dark mx-auto max-w-md overflow-hidden rounded-lg border"
				>
					{#if receiptStore.selectedImage}
						<div
							class="flex h-64 w-full flex-col items-center justify-center bg-gray-100 dark:bg-gray-800"
						>
							<svg class="h-16 w-16 text-red-500" fill="currentColor" viewBox="0 0 24 24">
								<path
									d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8l-6-6zM6 20V4h7v5h5v11H6z"
								/>
								<text x="8" y="16" font-size="6" fill="currentColor">PDF</text>
							</svg>
							<p class="text-text-secondary mt-2 text-sm">PDF Document</p>
						</div>
					{/if}
				</div>
				<p class="text-text-secondary mt-2 text-center text-sm">
					{receiptStore.selectedImage?.name}
				</p>
			</div>

			<div class="flex justify-center gap-4">
				<Button variant="primary" onclick={handleProcessReceipt} class="px-6 py-3">
					{#snippet leading()}<LightbulbIcon class="h-5 w-5" />{/snippet}
					{m.receipt_process()}
				</Button>
				<Button variant="secondary" onclick={handleClearAndStartOver} class="px-6 py-3">
					{m.receipt_clear_start_over()}
				</Button>
			</div>
		</div>
	{/if}

	<!-- Step 3: Processing (handled by overlay) -->

	<!-- Step 4: Results -->
	{#if currentStep === 4}
		<!-- Processing time indicator -->
		{#if receiptStore.processingTimeMs}
			<div class="text-text-secondary text-sm">
				{m.receipt_processed_time({ seconds: (receiptStore.processingTimeMs / 1000).toFixed(2) })}
			</div>
		{/if}

		<!-- Comparison View -->
		<ReceiptComparison
			imageUrl={receiptStore.imagePreviewUrl}
			items={receiptStore.extractedItems}
			onUpdateItem={handleUpdateItem}
			onRemoveItem={handleRemoveItem}
			onAddItem={handleAddItem}
			onToggleSelection={handleToggleSelection}
			onToggleSelectAll={handleToggleSelectAll}
		/>

		<!-- Action Buttons -->
		<div class="flex flex-wrap justify-end gap-4">
			<Button variant="secondary" onclick={handleClearAndStartOver}>
				{#snippet leading()}<RefreshCwIcon class="h-4 w-4" />{/snippet}
				{m.receipt_clear_start_over()}
			</Button>
			<Button
				variant="primary"
				onclick={handleAddToExpensesClick}
				disabled={getSelectedCount() === 0}
				class="bg-success hover:bg-success-hover text-text-on-success"
			>
				{#snippet leading()}<PlusIcon class="h-4 w-4" />{/snippet}
				{m.receipt_add_to_expenses()} ({getSelectedCount()})
			</Button>
		</div>
	{/if}
</div>

<!-- Confirmation Modal -->
<Dialog
	open={showConfirmModal}
	onClose={handleCancelModal}
	title={m.receipt_confirm_add()}
	size="md"
>
	<p class="text-text-tertiary">
		{m.receipt_confirm_add_description({ count: getSelectedCount().toString() })}
	</p>

	<div class="mt-4">
		<label for="receiptDate" class="text-text-tertiary block text-sm font-medium"
			>{m.actual_expenses_month()}</label
		>
		<input
			type="date"
			id="receiptDate"
			bind:value={receiptDate}
			class="border-input-border bg-input-bg text-text-primary focus:border-input-focus focus:ring-input-focus mt-1 block w-full rounded-md shadow-sm sm:text-sm"
		/>
	</div>

	<!-- Selected items summary -->
	<div class="mt-4 max-h-48 overflow-y-auto">
		<ul class="divide-border divide-y text-sm">
			{#each receiptStore.getSelectedItems() as item, i (i)}
				<li class="flex justify-between py-2">
					<span class="text-text-tertiary">{item.item_name}</span>
					<span class="text-text-primary font-medium">
						${item.item_price.toFixed(2)}
					</span>
				</li>
			{/each}
		</ul>
	</div>

	<div class="mt-6 flex justify-end gap-3">
		<Button variant="secondary" onclick={handleCancelModal} disabled={addingToExpenses}>
			{m.common_cancel()}
		</Button>
		<Button
			variant="primary"
			onclick={handleConfirmAddToExpenses}
			disabled={addingToExpenses}
			loading={addingToExpenses}
			class="bg-success hover:bg-success-hover text-text-on-success"
		>
			{addingToExpenses ? m.common_processing() : m.common_confirm()}
		</Button>
	</div>
</Dialog>
