import * as m from '$lib/paraglide/messages';

/**
 * Get localized month names using i18n
 * Returns array of objects with value (1-12) and localized label
 */
export function getMonths(): Array<{ value: number; label: string }> {
	return [
		{ value: 1, label: m.month_january() },
		{ value: 2, label: m.month_february() },
		{ value: 3, label: m.month_march() },
		{ value: 4, label: m.month_april() },
		{ value: 5, label: m.month_may() },
		{ value: 6, label: m.month_june() },
		{ value: 7, label: m.month_july() },
		{ value: 8, label: m.month_august() },
		{ value: 9, label: m.month_september() },
		{ value: 10, label: m.month_october() },
		{ value: 11, label: m.month_november() },
		{ value: 12, label: m.month_december() }
	];
}

/**
 * Get month name by number (1-12)
 */
export function getMonthName(month: number): string {
	const months = getMonths();
	return months.find((mo) => mo.value === month)?.label ?? '';
}

/**
 * Get array of month names (for simple dropdown displays)
 */
export function getMonthNames(): string[] {
	return getMonths().map((mo) => mo.label);
}
