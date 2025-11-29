/**
 * Receipt Processing Store
 * Manages temporary state for receipt processing (NOT persisted)
 */

import { createContext } from 'svelte';
import { api } from '$lib/utils/api';
import { ExpenseTypeEnum } from '$lib/types/enums';

/**
 * Extracted item from receipt OCR
 */
export interface ExtractedItem {
	source: string;
	type: ExpenseTypeEnum;
	item_code: string;
	item_price: number;
	item_name: string;
	expected_expense_id?: number; // Link to matched expected expense
	selected?: boolean;
}

/**
 * API response for receipt processing
 */
interface ProcessReceiptResponse {
	success: boolean;
	items: Omit<ExtractedItem, 'selected'>[];
	processing_time_ms: number;
}

/**
 * Receipt processing state
 */
export interface ReceiptState {
	selectedImage: File | null;
	imagePreviewUrl: string | null;
	extractedItems: ExtractedItem[];
	isProcessing: boolean;
	error: string | null;
	processingTimeMs: number | null;
	receiptNumber: number | null;
}

/**
 * Create a reactive receipt processing store
 */
export function createReceiptStore() {
	let selectedImage = $state<File | null>(null);
	let imagePreviewUrl = $state<string | null>(null);
	let extractedItems = $state<ExtractedItem[]>([]);
	let isProcessing = $state(false);
	let error = $state<string | null>(null);
	let processingTimeMs = $state<number | null>(null);
	let receiptNumber = $state<number | null>(null);

	/**
	 * Set the selected image file
	 */
	function setImage(file: File): void {
		// Clean up previous preview URL
		if (imagePreviewUrl) {
			URL.revokeObjectURL(imagePreviewUrl);
		}

		selectedImage = file;
		imagePreviewUrl = URL.createObjectURL(file);
		error = null;
		extractedItems = [];
		processingTimeMs = null;
	}

	/**
	 * Process the receipt image via API
	 */
	async function processReceipt(): Promise<boolean> {
		if (!selectedImage) {
			error = 'No image selected';
			return false;
		}

		isProcessing = true;
		error = null;

		try {
			const formData = new FormData();
			formData.append('document', selectedImage);

			const response = await api.postFormData<ProcessReceiptResponse>(
				'/receipts/process',
				formData
			);

			if (response.success) {
				// Add selected: false to each item
				extractedItems = response.items.map((item) => ({
					...item,
					selected: false
				}));
				processingTimeMs = response.processing_time_ms;
				return true;
			} else {
				error = 'Failed to process receipt';
				return false;
			}
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to process receipt';
			return false;
		} finally {
			isProcessing = false;
		}
	}

	/**
	 * Update a specific item field
	 */
	function updateItem(
		index: number,
		field: keyof ExtractedItem,
		value: string | number | boolean
	): void {
		if (index >= 0 && index < extractedItems.length) {
			extractedItems = extractedItems.map((item, i) => {
				if (i === index) {
					return { ...item, [field]: value };
				}
				return item;
			});
		}
	}

	/**
	 * Remove an item from the list
	 */
	function removeItem(index: number): void {
		if (index >= 0 && index < extractedItems.length) {
			extractedItems = extractedItems.filter((_, i) => i !== index);
		}
	}

	/**
	 * Add a blank item row
	 */
	function addItem(): void {
		const newItem: ExtractedItem = {
			source: '',
			type: ExpenseTypeEnum.WEEKLY,
			item_code: '',
			item_price: 0,
			item_name: '',
			selected: false
		};
		extractedItems = [...extractedItems, newItem];
	}

	/**
	 * Toggle selection of an item
	 */
	function toggleItemSelection(index: number): void {
		if (index >= 0 && index < extractedItems.length) {
			extractedItems = extractedItems.map((item, i) => {
				if (i === index) {
					return { ...item, selected: !item.selected };
				}
				return item;
			});
		}
	}

	/**
	 * Select or deselect all items
	 */
	function toggleSelectAll(selectAll: boolean): void {
		extractedItems = extractedItems.map((item) => ({
			...item,
			selected: selectAll
		}));
	}

	/**
	 * Get selected items
	 */
	function getSelectedItems(): ExtractedItem[] {
		return extractedItems.filter((item) => item.selected);
	}

	/**
	 * Clear all state and reset
	 */
	function clearAll(): void {
		if (imagePreviewUrl) {
			URL.revokeObjectURL(imagePreviewUrl);
		}
		selectedImage = null;
		imagePreviewUrl = null;
		extractedItems = [];
		isProcessing = false;
		error = null;
		processingTimeMs = null;
		receiptNumber = null;
	}

	/**
	 * Clear error state
	 */
	function clearError(): void {
		error = null;
	}

	/**
	 * Set the receipt number for the current batch
	 */
	function setReceiptNumber(num: number): void {
		receiptNumber = num;
	}

	return {
		// State getters
		get selectedImage() {
			return selectedImage;
		},
		get imagePreviewUrl() {
			return imagePreviewUrl;
		},
		get extractedItems() {
			return extractedItems;
		},
		get isProcessing() {
			return isProcessing;
		},
		get error() {
			return error;
		},
		get processingTimeMs() {
			return processingTimeMs;
		},
		get receiptNumber() {
			return receiptNumber;
		},
		// Actions
		setImage,
		setReceiptNumber,
		processReceipt,
		updateItem,
		removeItem,
		addItem,
		toggleItemSelection,
		toggleSelectAll,
		getSelectedItems,
		clearAll,
		clearError
	};
}

/**
 * Type for the receipt store instance
 */
export type ReceiptStore = ReturnType<typeof createReceiptStore>;

/**
 * Context getter and setter for route-level scoping
 * Key: 'receiptStore' (implicit via createContext)
 */
export const [getReceiptStore, setReceiptStore] = createContext<ReceiptStore>();

// Export a singleton store instance for backward compatibility
export const receiptStore = createReceiptStore();
