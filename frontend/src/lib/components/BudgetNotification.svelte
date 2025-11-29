<script lang="ts">
	import { formatCurrency, formatPercentageValue, safeNumber } from '$lib/utils/format';
	import { NotificationStatusEnum } from '$lib/types/enums';
	import { XCircleIcon, AlertTriangleIcon, CheckCircleIcon, XIcon } from 'lucide-svelte';
	import * as m from '$lib/paraglide/messages';
	import { Button } from '$lib';

	interface Props {
		currentSpending: number | null | undefined;
		budgetAmount: number;
		threshold: number;
	}

	let { currentSpending, budgetAmount, threshold }: Props = $props();

	let dismissed = $state(false);

	// Safe spending value (default to 0 if null/undefined)
	let spending = $derived(safeNumber(currentSpending, 0));

	// Computed values
	let percentage = $derived(budgetAmount > 0 ? (spending / budgetAmount) * 100 : 0);
	let thresholdPercentage = $derived(threshold * 100);

	// Determine status based on spending percentage
	let status = $derived.by(() => {
		if (percentage >= thresholdPercentage) {
			return NotificationStatusEnum.CRITICAL;
		} else if (percentage >= 80) {
			return NotificationStatusEnum.WARNING;
		} else if (percentage >= 50) {
			return NotificationStatusEnum.CAUTION;
		} else {
			return NotificationStatusEnum.GOOD;
		}
	});

	// Status-based styling
	let statusConfig = $derived.by(() => {
		switch (status) {
			case NotificationStatusEnum.CRITICAL:
				return {
					bgColor: 'bg-budget-danger-bg',
					borderColor: 'border-danger-dark',
					textColor: 'text-danger-light',
					progressColor: 'bg-danger',
					iconColor: 'text-danger',
					message: m.notification_critical_message({
						percentage: formatPercentageValue(percentage)
					})
				};
			case NotificationStatusEnum.WARNING:
				return {
					bgColor: 'bg-warning-dark/20',
					borderColor: 'border-warning-dark',
					textColor: 'text-warning-light',
					progressColor: 'bg-warning',
					iconColor: 'text-warning',
					message: m.notification_warning_message({ percentage: formatPercentageValue(percentage) })
				};
			case NotificationStatusEnum.CAUTION:
				return {
					bgColor: 'bg-budget-warning-bg',
					borderColor: 'border-warning-dark',
					textColor: 'text-warning-light',
					progressColor: 'bg-warning',
					iconColor: 'text-warning',
					message: m.notification_caution_message({ percentage: formatPercentageValue(percentage) })
				};
			default:
				return {
					bgColor: 'bg-budget-safe-bg',
					borderColor: 'border-success-dark',
					textColor: 'text-success-light',
					progressColor: 'bg-success',
					iconColor: 'text-success',
					message: m.notification_good_message({ percentage: formatPercentageValue(percentage) })
				};
		}
	});

	let remainingAmount = $derived(Math.max(0, budgetAmount - spending));
	let overBudgetAmount = $derived(Math.max(0, spending - budgetAmount));

	function handleDismiss() {
		dismissed = true;
	}
</script>

{#if !dismissed && budgetAmount > 0}
	<div class="rounded-lg border {statusConfig.bgColor} {statusConfig.borderColor} p-4">
		<div class="flex items-start">
			<!-- Icon -->
			<div class="shrink-0 {statusConfig.iconColor}">
				{#if status === NotificationStatusEnum.CRITICAL}
					<XCircleIcon class="h-5 w-5" />
				{:else if status === NotificationStatusEnum.WARNING || status === NotificationStatusEnum.CAUTION}
					<AlertTriangleIcon class="h-5 w-5" />
				{:else}
					<CheckCircleIcon class="h-5 w-5" />
				{/if}
			</div>

			<!-- Content -->
			<div class="ml-3 flex-1">
				<h3 class="text-sm font-medium {statusConfig.textColor}">
					{m.dashboard_budget_status()}
				</h3>
				<p class="mt-1 text-sm {statusConfig.textColor}">
					{statusConfig.message}
				</p>

				<!-- Progress Bar -->
				<div class="mt-3">
					<div class="flex justify-between text-xs {statusConfig.textColor} mb-1">
						<span>Spent: {formatCurrency(spending)}</span>
						<span>Budget: {formatCurrency(budgetAmount)}</span>
					</div>
					<div class="bg-surface-light h-2.5 w-full overflow-hidden rounded-full">
						<div
							class="{statusConfig.progressColor} h-2.5 rounded-full transition-all duration-500"
							style="width: {Math.min(percentage, 100)}%"
						></div>
					</div>
					<div class="mt-1 flex justify-between text-xs {statusConfig.textColor}">
						{#if percentage > 100}
							<span class="font-medium">Over budget by {formatCurrency(overBudgetAmount)}</span>
						{:else}
							<span>Remaining: {formatCurrency(remainingAmount)}</span>
						{/if}
						<span>Threshold: {thresholdPercentage}%</span>
					</div>
				</div>
			</div>

			<!-- Dismiss Button -->
			<div class="ml-4 shrink-0">
				<Button
					variant="ghost"
					onclick={handleDismiss}
					aria-label={m.common_close()}
					class="{statusConfig.textColor} !p-1"
				>
					<XIcon class="h-5 w-5" />
				</Button>
			</div>
		</div>
	</div>
{/if}
