<script lang="ts">
	import { Button } from '$lib';
	import { toastStore, type ToastType } from '$lib/stores/toast.svelte';
	import { ToastTypeEnum } from '$lib/types/enums';
	import { fly, fade } from 'svelte/transition';
	import { CheckCircleIcon, XCircleIcon, AlertTriangleIcon, InfoIcon, XIcon } from 'lucide-svelte';

	// Get toasts from store
	let toasts = $derived(toastStore.toasts);

	/**
	 * Get colors for toast type
	 */
	function getColors(type: ToastType): { bg: string; border: string; icon: string; text: string } {
		switch (type) {
			case ToastTypeEnum.SUCCESS:
				return {
					bg: 'bg-success-dark',
					border: 'border-success',
					icon: 'text-success',
					text: 'text-text-primary'
				};
			case ToastTypeEnum.ERROR:
				return {
					bg: 'bg-danger-dark',
					border: 'border-danger',
					icon: 'text-danger',
					text: 'text-text-primary'
				};
			case ToastTypeEnum.WARNING:
				return {
					bg: 'bg-warning-dark',
					border: 'border-warning',
					icon: 'text-warning',
					text: 'text-text-primary'
				};
			case ToastTypeEnum.INFO:
			default:
				return {
					bg: 'bg-info-dark',
					border: 'border-info',
					icon: 'text-info',
					text: 'text-text-primary'
				};
		}
	}

	function handleDismiss(id: string) {
		toastStore.dismissToast(id);
	}
</script>

<!-- Toast Container - Fixed position at top-right -->
<div
	class="pointer-events-none fixed top-4 right-4 z-[100] flex w-full max-w-sm flex-col gap-3"
	aria-live="polite"
	aria-atomic="true"
>
	{#each toasts as toast (toast.id)}
		{@const colors = getColors(toast.type)}
		<div
			class="pointer-events-auto w-full {colors.bg} {colors.border} overflow-hidden rounded-lg border shadow-lg"
			in:fly={{ x: 300, duration: 300 }}
			out:fade={{ duration: 200 }}
			role="alert"
		>
			<div class="p-4">
				<div class="flex items-start">
					<!-- Icon -->
					<div class="shrink-0 {colors.icon}">
						{#if toast.type === ToastTypeEnum.SUCCESS}
							<CheckCircleIcon class="h-5 w-5" />
						{:else if toast.type === ToastTypeEnum.ERROR}
							<XCircleIcon class="h-5 w-5" />
						{:else if toast.type === ToastTypeEnum.WARNING}
							<AlertTriangleIcon class="h-5 w-5" />
						{:else}
							<InfoIcon class="h-5 w-5" />
						{/if}
					</div>

					<!-- Message -->
					<div class="ml-3 flex-1">
						<p class="text-sm font-medium {colors.text}">
							{toast.message}
						</p>
					</div>

					<!-- Dismiss Button -->
					<div class="ml-4 shrink-0">
						<Button
							variant="ghost"
							onclick={() => handleDismiss(toast.id)}
							aria-label="Dismiss notification"
							class="!p-1 {colors.text}"
						>
							{#snippet leading()}
								<XIcon class="h-5 w-5" />
							{/snippet}
						</Button>
					</div>
				</div>
			</div>

			<!-- Progress bar for auto-dismiss -->
			{#if toast.duration > 0}
				<div class="h-1 {colors.bg}">
					<div
						class="h-full bg-current {colors.icon} opacity-30"
						style="animation: shrink {toast.duration}ms linear forwards;"
					></div>
				</div>
			{/if}
		</div>
	{/each}
</div>

<style>
	@keyframes shrink {
		from {
			width: 100%;
		}
		to {
			width: 0%;
		}
	}
</style>
