<script lang="ts">
	import { ExpenseFilterTypeEnum } from '$lib/types/enums';
	import * as m from '$lib/paraglide/messages';

	type FilterType = `${ExpenseFilterTypeEnum}`;

	interface Props {
		activeTab: string;
		onTabChange: (tab: FilterType) => void;
		showMisc?: boolean;
	}

	let { activeTab, onTabChange, showMisc = false }: Props = $props();

	let tabs = $derived([
		{ id: ExpenseFilterTypeEnum.ALL as FilterType, label: m.expenses_all() },
		{ id: ExpenseFilterTypeEnum.WEEKLY as FilterType, label: m.expenses_weekly() },
		{ id: ExpenseFilterTypeEnum.MONTHLY as FilterType, label: m.expenses_monthly() },
		...(showMisc
			? [{ id: ExpenseFilterTypeEnum.MISC as FilterType, label: m.expenses_misc() }]
			: [])
	]);
</script>

<div class="bg-surface-light flex space-x-1 rounded-lg p-1">
	{#each tabs as tab (tab.id)}
		<button
			type="button"
			onclick={() => onTabChange(tab.id)}
			class="cursor-pointer rounded-md px-3 py-1.5 text-sm font-medium transition-colors {tab.id ===
			activeTab
				? 'bg-primary text-text-on-primary'
				: 'text-text-secondary hover:text-text-primary hover:bg-surface-lighter'}"
			aria-label={tab.id === activeTab ? `${tab.label} (active)` : tab.label}
			aria-pressed={tab.id === activeTab}
		>
			{tab.label}
		</button>
	{/each}
</div>
