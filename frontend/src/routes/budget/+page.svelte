<script lang="ts">
	import { onMount } from 'svelte';
	import { budgetStore, type Budget } from '$lib/stores/budget.svelte';
	import { toastStore } from '$lib/stores/toast.svelte';
	import { Button, Dialog } from '$lib';
	import BudgetForm from '$lib/components/BudgetForm.svelte';
	import BudgetList from '$lib/components/BudgetList.svelte';
	import Skeleton from '$lib/components/Skeleton.svelte';
	import { PlusIcon, LightbulbIcon } from 'lucide-svelte';
	import * as m from '$lib/paraglide/messages';

	// UI state
	let showForm = $state(false);
	let editingBudget = $state<Budget | null>(null);

	// Load data on mount
	onMount(() => {
		budgetStore.fetchBudgets();
	});

	function handleAddBudget() {
		editingBudget = null;
		showForm = true;
	}

	function handleEditBudget(budget: Budget) {
		editingBudget = budget;
		showForm = true;
	}

	async function handleDeleteBudget(budget: Budget) {
		const success = await budgetStore.deleteBudget(budget.id);
		if (success) {
			toastStore.success('Budget deleted successfully');
		} else if (budgetStore.error) {
			toastStore.error(budgetStore.error);
		}
	}

	function handleFormSave() {
		showForm = false;
		editingBudget = null;
	}

	function handleFormCancel() {
		showForm = false;
		editingBudget = null;
	}

	// Computed dialog title based on editing state
	let dialogTitle = $derived(editingBudget ? m.budget_edit_title() : m.budget_create_title());
</script>

<svelte:head>
	<title>Budget Settings | Budget Tracker</title>
</svelte:head>

<div class="space-y-4 sm:space-y-6">
	<!-- Page Header -->
	<div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
		<div>
			<h1 class="text-text-primary text-xl font-bold sm:text-2xl">{m.budget_title()}</h1>
			<p class="text-text-secondary mt-1 text-sm">
				{m.budget_description()}
			</p>
		</div>
		{#if !showForm}
			<Button onclick={handleAddBudget}>
				{#snippet leading()}<PlusIcon class="h-4 w-4" />{/snippet}
				{m.budget_add()}
			</Button>
		{/if}
	</div>

	<!-- Budget Form Modal -->
	<Dialog open={showForm} onClose={handleFormCancel} title={dialogTitle}>
		<BudgetForm budget={editingBudget} onSave={handleFormSave} onCancel={handleFormCancel} />
	</Dialog>

	<!-- Loading State -->
	{#if budgetStore.loading && budgetStore.budgets.length === 0}
		<div class="bg-surface overflow-hidden rounded-lg shadow">
			<div class="border-border border-b px-6 py-4">
				<div class="bg-skeleton h-6 w-32 animate-pulse rounded"></div>
			</div>
			<div class="space-y-4 p-6">
				<Skeleton variant="rectangle" height="4rem" />
				<Skeleton variant="rectangle" height="4rem" />
				<Skeleton variant="rectangle" height="4rem" />
			</div>
		</div>
	{:else}
		<!-- Budget List -->
		<BudgetList
			budgets={budgetStore.budgets}
			onEdit={handleEditBudget}
			onDelete={handleDeleteBudget}
		/>
	{/if}

	<!-- Quick Tips Section -->
	<div class="bg-info-bg border-info rounded-lg border p-4">
		<h3 class="text-info mb-2 flex items-center text-sm font-medium">
			<LightbulbIcon class="mr-2 h-5 w-5" />
			{m.budget_tips_title()}
		</h3>
		<ul class="text-info ml-7 space-y-1 text-sm">
			<li>{m.budget_tips_set_threshold()}</li>
		</ul>
	</div>
</div>
