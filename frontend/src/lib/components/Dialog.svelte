<script lang="ts">
	import type { Snippet } from 'svelte';
	import { focusTrap } from '$lib/utils/focusTrap';
	import { fade, scale } from 'svelte/transition';

	type DialogSize = 'sm' | 'md' | 'lg' | 'xl' | 'full';

	interface Props {
		open: boolean;
		onClose: () => void;
		title: string;
		size?: DialogSize;
		/** Hide the default header - useful for image viewers or custom layouts */
		hideHeader?: boolean;
		/** Remove padding from content area */
		noPadding?: boolean;
		children?: Snippet;
	}

	let {
		open,
		onClose,
		title,
		size = 'md',
		hideHeader = false,
		noPadding = false,
		children
	}: Props = $props();

	// Generate unique ID for aria-labelledby
	const titleId = `dialog-title-${Math.random().toString(36).slice(2, 11)}`;

	// Size mappings
	const sizeStyles: Record<DialogSize, string> = {
		sm: 'max-w-sm',
		md: 'max-w-md',
		lg: 'max-w-lg',
		xl: 'max-w-4xl',
		full: 'max-w-[90vw]'
	};

	/**
	 * Handle backdrop click to close dialog
	 */
	function handleBackdropClick(event: MouseEvent): void {
		// Only close if clicking directly on the backdrop, not its children
		if (event.target === event.currentTarget) {
			onClose();
		}
	}

	/**
	 * Global keydown handler for Escape key
	 * This ensures Escape works regardless of which element has focus
	 */
	function handleGlobalKeyDown(event: KeyboardEvent): void {
		if (event.key === 'Escape') {
			event.preventDefault();
			onClose();
		}
	}

	/**
	 * Manage body scroll and global Escape key handler when dialog opens/closes
	 */
	$effect(() => {
		if (open) {
			document.body.classList.add('modal-open');
			document.addEventListener('keydown', handleGlobalKeyDown);
		} else {
			document.body.classList.remove('modal-open');
			document.removeEventListener('keydown', handleGlobalKeyDown);
		}

		// Cleanup on unmount
		return () => {
			document.body.classList.remove('modal-open');
			document.removeEventListener('keydown', handleGlobalKeyDown);
		};
	});
</script>

{#if open}
	<div
		class="bg-surface-overlay fixed top-0 left-0 z-50 flex h-screen w-screen items-center justify-center p-4"
		role="dialog"
		aria-modal="true"
		aria-labelledby={hideHeader ? undefined : titleId}
		aria-label={hideHeader ? title : undefined}
		tabindex="-1"
		onclick={handleBackdropClick}
		onkeydown={() => {}}
		transition:fade={{ duration: 150 }}
	>
		<div
			class="bg-surface {sizeStyles[size]} max-h-[90vh] w-full overflow-y-auto rounded-lg shadow-lg"
			use:focusTrap
			transition:scale={{ duration: 150, start: 0.95 }}
		>
			{#if !hideHeader}
				<!-- Dialog Header -->
				<div class="border-border border-b px-6 py-4">
					<h2 id={titleId} class="text-text-primary text-lg font-semibold">
						{title}
					</h2>
				</div>
			{/if}

			<!-- Dialog Content -->
			<div class={noPadding ? '' : 'p-6'}>
				{#if children}
					{@render children()}
				{/if}
			</div>
		</div>
	</div>
{/if}
