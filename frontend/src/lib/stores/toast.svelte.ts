/**
 * Toast Store for Budget Tracker
 * Manages toast notifications with Svelte 5 runes
 */

import { ToastTypeEnum } from '$lib/types/enums';

export type ToastType = `${ToastTypeEnum}`;

export interface Toast {
	id: string;
	message: string;
	type: ToastType;
	duration: number;
}

interface ToastOptions {
	duration?: number; // Duration in ms, default 5000
}

/**
 * Create a reactive toast store
 */
function createToastStore() {
	let toasts = $state<Toast[]>([]);

	/**
	 * Generate a unique ID for each toast
	 */
	function generateId(): string {
		return `toast-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
	}

	/**
	 * Show a toast notification
	 */
	function showToast(
		message: string,
		type: ToastType = ToastTypeEnum.INFO,
		options: ToastOptions = {}
	): string {
		const id = generateId();
		const duration = options.duration ?? 5000;

		const toast: Toast = {
			id,
			message,
			type,
			duration
		};

		toasts = [...toasts, toast];

		// Auto-dismiss after duration
		if (duration > 0) {
			setTimeout(() => {
				dismissToast(id);
			}, duration);
		}

		return id;
	}

	/**
	 * Dismiss a specific toast by ID
	 */
	function dismissToast(id: string): void {
		toasts = toasts.filter((t) => t.id !== id);
	}

	/**
	 * Dismiss all toasts
	 */
	function dismissAll(): void {
		toasts = [];
	}

	/**
	 * Convenience methods for different toast types
	 */
	function success(message: string, options?: ToastOptions): string {
		return showToast(message, ToastTypeEnum.SUCCESS, options);
	}

	function error(message: string, options?: ToastOptions): string {
		return showToast(message, ToastTypeEnum.ERROR, { duration: 7000, ...options }); // Errors stay longer
	}

	function warning(message: string, options?: ToastOptions): string {
		return showToast(message, ToastTypeEnum.WARNING, options);
	}

	function info(message: string, options?: ToastOptions): string {
		return showToast(message, ToastTypeEnum.INFO, options);
	}

	return {
		get toasts() {
			return toasts;
		},
		showToast,
		dismissToast,
		dismissAll,
		success,
		error,
		warning,
		info
	};
}

// Export singleton instance
export const toastStore = createToastStore();
