<script lang="ts">
	import { FileTextIcon } from 'lucide-svelte';

	interface Props {
		message?: string;
		fullScreen?: boolean;
	}

	let { message = 'Processing receipt with AI...', fullScreen = true }: Props = $props();
</script>

{#if fullScreen}
	<!-- Full-screen overlay -->
	<div
		class="bg-surface-dark/90 fixed top-0 left-0 z-50 flex h-screen w-screen items-center justify-center backdrop-blur-sm"
	>
		<div class="text-center">
			<!-- Animated spinner -->
			<div class="relative mx-auto h-20 w-20">
				<!-- Outer ring -->
				<div class="border-primary-dark absolute inset-0 animate-pulse rounded-full border-4"></div>
				<!-- Spinning ring -->
				<div
					class="border-t-primary absolute inset-0 animate-spin rounded-full border-4 border-transparent"
				></div>
				<!-- Inner icon -->
				<div class="absolute inset-0 flex items-center justify-center">
					<FileTextIcon class="text-primary h-8 w-8" />
				</div>
			</div>

			<!-- Message -->
			<p class="text-text-primary mt-6 text-lg font-medium">{message}</p>
			<p class="text-text-secondary mt-2 text-sm">This may take a few seconds...</p>

			<!-- Animated dots -->
			<div class="mt-4 flex justify-center gap-1">
				<span class="bg-primary h-2 w-2 animate-bounce rounded-full" style="animation-delay: 0ms;"
				></span>
				<span class="bg-primary h-2 w-2 animate-bounce rounded-full" style="animation-delay: 150ms;"
				></span>
				<span class="bg-primary h-2 w-2 animate-bounce rounded-full" style="animation-delay: 300ms;"
				></span>
			</div>
		</div>
	</div>
{:else}
	<!-- Inline spinner -->
	<div class="flex items-center justify-center py-8">
		<div class="text-center">
			<!-- Spinning ring -->
			<div
				class="border-primary-dark border-t-primary mx-auto h-12 w-12 animate-spin rounded-full border-4"
			></div>

			<!-- Message -->
			<p class="text-text-secondary mt-4 text-sm font-medium">{message}</p>
		</div>
	</div>
{/if}

<style>
	@keyframes bounce {
		0%,
		100% {
			transform: translateY(0);
		}
		50% {
			transform: translateY(-8px);
		}
	}
</style>
