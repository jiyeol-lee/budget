<script lang="ts">
	import {
		ZoomInIcon,
		ZoomOutIcon,
		MaximizeIcon,
		ImageIcon,
		XIcon
	} from 'lucide-svelte';
	import { Button } from '$lib';

	interface Props {
		imageUrl: string | null;
		alt?: string;
	}

	let { imageUrl, alt = 'Receipt image' }: Props = $props();

	let zoom = $state(1);
	let showModal = $state(false);
	const minZoom = 0.5;
	const maxZoom = 3;
	const zoomStep = 0.25;

	/**
	 * Zoom in
	 */
	function zoomIn(): void {
		if (zoom < maxZoom) {
			zoom = Math.min(zoom + zoomStep, maxZoom);
		}
	}

	/**
	 * Zoom out
	 */
	function zoomOut(): void {
		if (zoom > minZoom) {
			zoom = Math.max(zoom - zoomStep, minZoom);
		}
	}

	/**
	 * Reset zoom to fit
	 */
	function resetZoom(): void {
		zoom = 1;
	}

	/**
	 * Open modal view
	 */
	function openModal(): void {
		showModal = true;
	}

	/**
	 * Close modal view
	 */
	function closeModal(): void {
		showModal = false;
	}

	/**
	 * Handle keyboard events in modal
	 */
	function handleModalKeydown(event: KeyboardEvent): void {
		if (event.key === 'Escape') {
			closeModal();
		}
	}
</script>

{#if imageUrl}
	<div class="flex h-full flex-col">
		<!-- Zoom Controls -->
		<div class="bg-surface border-border flex items-center justify-between border-b p-2">
			<span class="text-text-secondary text-sm">Receipt Preview</span>
			<div class="flex items-center gap-2">
				<Button
					variant="ghost"
					onclick={zoomOut}
					disabled={zoom <= minZoom}
					class="p-1.5"
					aria-label="Zoom out"
				>
					<ZoomOutIcon class="h-4 w-4" />
				</Button>
				<span class="text-text-secondary min-w-12 text-center text-sm"
					>{Math.round(zoom * 100)}%</span
				>
				<Button
					variant="ghost"
					onclick={zoomIn}
					disabled={zoom >= maxZoom}
					class="p-1.5"
					aria-label="Zoom in"
				>
					<ZoomInIcon class="h-4 w-4" />
				</Button>
				<Button
					variant="ghost"
					onclick={resetZoom}
					class="p-1.5"
					aria-label="Reset zoom"
				>
					Fit
				</Button>
				<Button
					variant="ghost"
					onclick={openModal}
					class="p-1.5"
					aria-label="Expand"
				>
					<MaximizeIcon class="h-4 w-4" />
				</Button>
			</div>
		</div>

		<!-- Image Container -->
		<div
			class="bg-surface-dark border-border flex-1 overflow-auto rounded-b-lg border"
			style="min-height: 300px;"
		>
			<div class="flex min-h-full items-center justify-center p-4">
				<button type="button" onclick={openModal} class="cursor-zoom-in">
					<img
						src={imageUrl}
						{alt}
						class="h-auto max-w-full shadow-lg transition-transform duration-200"
						style="transform: scale({zoom}); transform-origin: center center;"
					/>
				</button>
			</div>
		</div>
	</div>
{:else}
	<div
		class="bg-surface-dark border-border flex h-full min-h-64 items-center justify-center rounded-lg border"
	>
		<div class="text-text-secondary text-center">
			<ImageIcon class="text-text-tertiary mx-auto h-12 w-12" />
			<p class="mt-2">No image to preview</p>
		</div>
	</div>
{/if}

<!-- Modal View -->
{#if showModal && imageUrl}
	<div
		class="bg-surface-overlay fixed top-0 left-0 z-50 flex h-screen w-screen items-center justify-center"
		role="dialog"
		aria-modal="true"
		aria-label="Receipt image expanded view"
	>
		<!-- Backdrop button to close -->
		<button
			type="button"
			class="absolute inset-0 h-full w-full cursor-default"
			onclick={closeModal}
			onkeydown={handleModalKeydown}
			aria-label="Close modal"
		></button>
		<div class="relative z-10 max-h-[90vh] max-w-[90vw]">
			<Button
				variant="ghost"
				onclick={closeModal}
				class="absolute -top-10 right-0"
				aria-label="Close modal"
			>
				<XIcon class="h-8 w-8" />
			</Button>
			<img
				src={imageUrl}
				{alt}
				class="max-h-[85vh] max-w-full rounded-lg object-contain shadow-2xl"
			/>
		</div>
	</div>
{/if}
