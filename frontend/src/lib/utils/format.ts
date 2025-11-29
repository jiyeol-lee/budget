/**
 * Utility functions for formatting values safely
 */

/**
 * Format a number as currency, returns "$0.00" for invalid values
 */
export function formatCurrency(amount: number | null | undefined): string {
	if (amount == null || isNaN(amount)) {
		return '$0.00';
	}
	return new Intl.NumberFormat('en-US', {
		style: 'currency',
		currency: 'USD'
	}).format(amount);
}

/**
 * Format a decimal as percentage, returns "0%" for invalid values
 */
export function formatPercentage(value: number | null | undefined, decimals: number = 0): string {
	if (value == null || isNaN(value)) {
		return '0%';
	}
	return `${(value * 100).toFixed(decimals)}%`;
}

/**
 * Format a percentage value (already multiplied by 100), returns "0%" for invalid values
 */
export function formatPercentageValue(
	value: number | null | undefined,
	decimals: number = 1
): string {
	if (value == null || isNaN(value)) {
		return '0%';
	}
	return `${value.toFixed(decimals)}%`;
}

/**
 * Format a date string, returns "N/A" for invalid dates
 */
export function formatDate(
	dateString: string | null | undefined,
	options?: Intl.DateTimeFormatOptions
): string {
	if (!dateString) {
		return 'N/A';
	}

	const date = new Date(dateString);
	if (isNaN(date.getTime())) {
		return 'N/A';
	}

	return date.toLocaleDateString(
		'en-US',
		options ?? {
			year: 'numeric',
			month: 'short',
			day: 'numeric'
		}
	);
}

/**
 * Format month and year into a readable string, returns "N/A" for invalid values
 */
export function formatMonthYear(
	month: number | null | undefined,
	year: number | null | undefined
): string {
	if (month == null || year == null || isNaN(month) || isNaN(year) || month < 1 || month > 12) {
		return 'N/A';
	}

	const date = new Date(year, month - 1);
	if (isNaN(date.getTime())) {
		return 'N/A';
	}

	return date.toLocaleDateString('en-US', { month: 'long', year: 'numeric' });
}

/**
 * Format a number with fallback, returns "0" for invalid values
 */
export function formatNumber(value: number | null | undefined, fallback: string = '0'): string {
	if (value == null || isNaN(value)) {
		return fallback;
	}
	return value.toLocaleString();
}

/**
 * Safely get a number value, returns 0 for invalid values
 */
export function safeNumber(value: number | null | undefined, fallback: number = 0): number {
	if (value == null || isNaN(value)) {
		return fallback;
	}
	return value;
}
