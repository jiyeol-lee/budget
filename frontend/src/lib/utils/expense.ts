import { ExpenseTypeEnum } from '$lib/types/enums';

/**
 * Get badge color classes based on expense type
 */
export function getTypeBadgeClass(type: ExpenseTypeEnum | string): string {
	switch (type) {
		case ExpenseTypeEnum.WEEKLY:
			return 'bg-primary-dark text-primary-light';
		case ExpenseTypeEnum.MONTHLY:
			return 'bg-purple-900 text-purple-200';
		case ExpenseTypeEnum.MISC:
			return 'bg-surface-light text-text-secondary';
		case ExpenseTypeEnum.TAX:
			return 'bg-amber-900 text-amber-200';
		default:
			return 'bg-surface-light text-text-secondary';
	}
}
