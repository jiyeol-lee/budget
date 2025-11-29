<script lang="ts">
	import type { Budget } from '$lib/stores/budget.svelte';
	import { Button } from '$lib';
	import { DollarSignIcon } from 'lucide-svelte';
	import * as m from '$lib/paraglide/messages';

	interface Props {
		budgets: Budget[];
		onEdit?: (budget: Budget) => void;
		onDelete?: (budget: Budget) => void;
	}

	let { budgets, onEdit, onDelete }: Props = $props();

	// Month names for display
	const monthNames = [
		'January',
		'February',
		'March',
		'April',
		'May',
		'June',
		'July',
		'August',
		'September',
		'October',
		'November',
		'December'
	];

	function formatCurrency(amount: number): string {
		return new Intl.NumberFormat('en-US', {
			style: 'currency',
			currency: 'USD'
		}).format(amount);
	}

	function formatThreshold(threshold: number): string {
		return `${Math.round(threshold * 100)}%`;
	}

	function getMonthName(month: number): string {
		return monthNames[month - 1] || 'Unknown';
	}

	function handleEdit(budget: Budget) {
		onEdit?.(budget);
	}

	function handleDelete(budget: Budget) {
		if (
			confirm(
				`Are you sure you want to delete the budget for ${getMonthName(budget.month)} ${budget.year}?`
			)
		) {
			onDelete?.(budget);
		}
	}

	// Sort budgets by year (desc) and month (desc)
	let sortedBudgets = $derived(
		[...budgets].sort((a, b) => {
			if (a.year !== b.year) {
				return b.year - a.year;
			}
			return b.month - a.month;
		})
	);

	// Check if budget is current month
	function isCurrentMonth(budget: Budget): boolean {
		const now = new Date();
		return budget.month === now.getMonth() + 1 && budget.year === now.getFullYear();
	}
</script>

<div class="bg-surface overflow-hidden rounded-lg shadow">
	<div class="border-border border-b px-4 py-4 sm:px-6">
		<h2 class="text-text-primary text-lg font-semibold">Your Budgets</h2>
	</div>

	{#if sortedBudgets.length === 0}
		<!-- Empty State -->
		<div class="px-6 py-12 text-center">
			<DollarSignIcon class="text-text-secondary mx-auto h-12 w-12" />
			<h3 class="text-text-primary mt-2 text-sm font-medium">{m.budget_no_budgets()}</h3>
			<p class="text-text-secondary mt-1 text-sm">
				{m.budget_no_budgets_description()}
			</p>
		</div>
	{:else}
		<!-- Mobile Card View -->
		<div class="divide-border divide-y sm:hidden">
			{#each sortedBudgets as budget (budget.id)}
				<div class="p-4 {isCurrentMonth(budget) ? 'bg-primary-dark/20' : ''}">
					<div class="mb-2 flex items-start justify-between">
						<div>
							<h3 class="text-text-primary text-sm font-medium">
								{getMonthName(budget.month)}
								{budget.year}
							</h3>
							{#if isCurrentMonth(budget)}
								<span
									class="bg-primary-dark text-primary-light mt-1 inline-flex items-center rounded px-2 py-0.5 text-xs font-medium"
								>
									Current
								</span>
							{/if}
						</div>
						<span
							class="bg-primary-dark text-primary-light inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium"
						>
							{formatThreshold(budget.notification_threshold)} threshold
						</span>
					</div>
					<p class="text-text-primary mb-3 text-2xl font-bold">
						{formatCurrency(budget.amount)}
					</p>
					<div class="flex justify-end space-x-3">
						<Button variant="link" onclick={() => handleEdit(budget)} aria-label={m.common_edit()}>
							{m.common_edit()}
						</Button>
						<Button
							variant="link"
							onclick={() => handleDelete(budget)}
							class="!text-danger hover:!text-danger-light"
							aria-label={m.common_delete()}
						>
							{m.common_delete()}
						</Button>
					</div>
				</div>
			{/each}
		</div>

		<!-- Desktop Table View -->
		<div class="hidden overflow-x-auto sm:block">
			<table class="divide-border min-w-full divide-y">
				<thead class="bg-surface-light">
					<tr>
						<th
							scope="col"
							class="text-text-secondary px-4 py-3 text-left text-xs font-medium tracking-wider uppercase sm:px-6"
						>
							Month / Year
						</th>
						<th
							scope="col"
							class="text-text-secondary px-4 py-3 text-left text-xs font-medium tracking-wider uppercase sm:px-6"
						>
							Budget Amount
						</th>
						<th
							scope="col"
							class="text-text-secondary px-4 py-3 text-left text-xs font-medium tracking-wider uppercase sm:px-6"
						>
							Notification Threshold
						</th>
						<th
							scope="col"
							class="text-text-secondary px-4 py-3 text-right text-xs font-medium tracking-wider uppercase sm:px-6"
						>
							Actions
						</th>
					</tr>
				</thead>
				<tbody class="bg-surface divide-border divide-y">
					{#each sortedBudgets as budget (budget.id)}
						<tr
							class="hover:bg-surface-light transition-colors {isCurrentMonth(budget)
								? 'bg-primary-dark/20 hover:bg-primary-dark/30'
								: ''}"
						>
							<td class="px-4 py-4 whitespace-nowrap sm:px-6">
								<div class="flex items-center">
									<span class="text-text-primary text-sm font-medium">
										{getMonthName(budget.month)}
										{budget.year}
									</span>
									{#if isCurrentMonth(budget)}
										<span
											class="bg-primary-dark text-primary-light ml-2 inline-flex items-center rounded px-2 py-0.5 text-xs font-medium"
										>
											Current
										</span>
									{/if}
								</div>
							</td>
							<td class="px-4 py-4 whitespace-nowrap sm:px-6">
								<div class="text-text-primary text-sm font-semibold">
									{formatCurrency(budget.amount)}
								</div>
							</td>
							<td class="px-4 py-4 whitespace-nowrap sm:px-6">
								<span
									class="bg-primary-dark text-primary-light inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium"
								>
									{formatThreshold(budget.notification_threshold)}
								</span>
							</td>
							<td class="px-4 py-4 text-right text-sm font-medium whitespace-nowrap sm:px-6">
								<Button
									variant="link"
									onclick={() => handleEdit(budget)}
									class="mr-3"
									aria-label={m.common_edit()}
								>
									{m.common_edit()}
								</Button>
								<Button
									variant="link"
									onclick={() => handleDelete(budget)}
									class="!text-danger hover:!text-danger-light"
									aria-label={m.common_delete()}
								>
									{m.common_delete()}
								</Button>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</div>
