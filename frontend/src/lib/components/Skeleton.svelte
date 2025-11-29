<script lang="ts">
	import { SkeletonVariantEnum } from '$lib/types/enums';

	type SkeletonVariant = `${SkeletonVariantEnum}`;

	interface Props {
		variant?: SkeletonVariant;
		width?: string;
		height?: string;
		count?: number;
		className?: string;
	}

	let {
		variant = SkeletonVariantEnum.TEXT,
		width = '100%',
		height,
		count = 1,
		className = ''
	}: Props = $props();

	// Compute default heights based on variant
	let computedHeight = $derived(() => {
		if (height) return height;
		switch (variant) {
			case SkeletonVariantEnum.TEXT:
				return '1rem';
			case SkeletonVariantEnum.CARD:
				return '8rem';
			case SkeletonVariantEnum.TABLE_ROW:
				return '3rem';
			case SkeletonVariantEnum.CIRCLE:
				return width; // Circles should be square
			case SkeletonVariantEnum.RECTANGLE:
				return '2rem';
			default:
				return '1rem';
		}
	});

	// Generate array for count
	let items = $derived(Array.from({ length: count }, (_, i) => i));
</script>

{#if variant === SkeletonVariantEnum.CARD}
	<!-- Card Skeleton -->
	{#each items as i (i)}
		<div class="bg-surface animate-pulse rounded-lg p-6 shadow {className}" style="width: {width};">
			<div class="bg-skeleton mb-4 h-4 w-1/3 rounded"></div>
			<div class="bg-skeleton mb-2 h-8 w-2/3 rounded"></div>
			<div class="bg-skeleton h-3 w-1/2 rounded"></div>
		</div>
	{/each}
{:else if variant === SkeletonVariantEnum.TABLE_ROW}
	<!-- Table Row Skeleton -->
	{#each items as i (i)}
		<tr class="animate-pulse {className}">
			<td class="px-6 py-4">
				<div class="bg-skeleton h-4 w-3/4 rounded"></div>
			</td>
			<td class="px-6 py-4">
				<div class="bg-skeleton h-4 w-1/2 rounded"></div>
			</td>
			<td class="px-6 py-4">
				<div class="bg-skeleton h-4 w-1/4 rounded"></div>
			</td>
			<td class="px-6 py-4">
				<div class="bg-skeleton h-6 w-16 rounded-full"></div>
			</td>
			<td class="px-6 py-4">
				<div class="bg-skeleton h-4 w-12 rounded"></div>
			</td>
			<td class="px-6 py-4 text-right">
				<div class="flex justify-end gap-2">
					<div class="bg-skeleton h-4 w-10 rounded"></div>
					<div class="bg-skeleton h-4 w-12 rounded"></div>
				</div>
			</td>
		</tr>
	{/each}
{:else if variant === SkeletonVariantEnum.CIRCLE}
	<!-- Circle Skeleton -->
	{#each items as i (i)}
		<div
			class="bg-skeleton animate-pulse rounded-full {className}"
			style="width: {width}; height: {computedHeight()};"
		></div>
	{/each}
{:else}
	<!-- Default/Text/Rectangle Skeleton -->
	{#each items as i (i)}
		<div
			class="bg-skeleton animate-pulse {variant === SkeletonVariantEnum.TEXT
				? 'rounded'
				: 'rounded-md'} {className}"
			style="width: {width}; height: {computedHeight()};"
		></div>
	{/each}
{/if}
