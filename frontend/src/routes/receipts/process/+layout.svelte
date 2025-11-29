<script lang="ts">
	import type { Snippet } from 'svelte';
	import { createReceiptStore, setReceiptStore } from '$lib/stores/receipt.svelte';

	interface Props {
		children: Snippet;
	}

	let { children }: Props = $props();

	// Create store instance scoped to this layout
	const receiptStore = createReceiptStore();

	// Provide via context - will be garbage collected when layout unmounts
	setReceiptStore(receiptStore);
</script>

{@render children()}
