/**
 * Enum for expense types (weekly, monthly, misc)
 */
export enum ExpenseTypeEnum {
	WEEKLY = 'weekly',
	MONTHLY = 'monthly',
	MISC = 'misc',
	TAX = 'tax'
}

/**
 * Enum for toast notification types
 */
export enum ToastTypeEnum {
	SUCCESS = 'success',
	ERROR = 'error',
	WARNING = 'warning',
	INFO = 'info'
}

/**
 * Enum for budget status levels
 */
export enum BudgetStatusEnum {
	SAFE = 'safe',
	WARNING = 'warning',
	DANGER = 'danger',
	OVER = 'over'
}

/**
 * Enum for notification status levels
 */
export enum NotificationStatusEnum {
	GOOD = 'good',
	CAUTION = 'caution',
	WARNING = 'warning',
	CRITICAL = 'critical'
}

/**
 * Enum for skeleton loading variants
 */
export enum SkeletonVariantEnum {
	TEXT = 'text',
	CARD = 'card',
	TABLE_ROW = 'table_row',
	CIRCLE = 'circle',
	RECTANGLE = 'rectangle'
}

/**
 * Enum for expense filter types (includes ALL for filtering)
 */
export enum ExpenseFilterTypeEnum {
	ALL = 'all',
	WEEKLY = 'weekly',
	MONTHLY = 'monthly',
	MISC = 'misc',
	TAX = 'tax'
}
