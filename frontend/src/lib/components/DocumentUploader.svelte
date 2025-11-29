<script lang="ts">
	import { UploadIcon, FileIcon, XIcon } from 'lucide-svelte';
	import * as m from '$lib/paraglide/messages';
	import { Button } from '$lib';

	interface Props {
		onFileSelected: (file: File) => void;
		maxSizeMB?: number;
		accept?: string;
	}

	let { onFileSelected, maxSizeMB = 10, accept = 'application/pdf' }: Props = $props();

	let isDragging = $state(false);
	let previewUrl = $state<string | null>(null);
	let selectedFile = $state<File | null>(null);
	let errorMessage = $state<string | null>(null);
	let fileInputRef = $state<HTMLInputElement | null>(null);

	const maxSizeBytes = maxSizeMB * 1024 * 1024;

	/**
	 * Validate the selected file
	 */
	function validateFile(file: File): string | null {
		// Check file type
		const validTypes = accept.split(',').map((t) => t.trim());
		if (!validTypes.includes(file.type)) {
			return `Invalid file type. Please upload ${validTypes.join(' or ')}.`;
		}

		// Check file size
		if (file.size > maxSizeBytes) {
			return `File too large. Maximum size is ${maxSizeMB}MB.`;
		}

		return null;
	}

	/**
	 * Handle file selection
	 */
	function handleFile(file: File): void {
		errorMessage = null;

		const validationError = validateFile(file);
		if (validationError) {
			errorMessage = validationError;
			return;
		}

		// Clean up previous preview
		if (previewUrl) {
			URL.revokeObjectURL(previewUrl);
		}

		previewUrl = URL.createObjectURL(file);
		selectedFile = file;
		onFileSelected(file);
	}

	/**
	 * Handle drag over event
	 */
	function handleDragOver(event: DragEvent): void {
		event.preventDefault();
		isDragging = true;
	}

	/**
	 * Handle drag leave event
	 */
	function handleDragLeave(event: DragEvent): void {
		event.preventDefault();
		isDragging = false;
	}

	/**
	 * Handle drop event
	 */
	function handleDrop(event: DragEvent): void {
		event.preventDefault();
		isDragging = false;

		const files = event.dataTransfer?.files;
		if (files && files.length > 0) {
			const file = files[0];
			if (file) handleFile(file);
		}
	}

	/**
	 * Handle file input change
	 */
	function handleFileInputChange(event: Event): void {
		const input = event.target as HTMLInputElement;
		if (input.files && input.files.length > 0) {
			const file = input.files[0];
			if (file) handleFile(file);
		}
	}

	/**
	 * Open file dialog
	 */
	function openFileDialog(): void {
		fileInputRef?.click();
	}

	/**
	 * Clear preview and reset
	 */
	function clearPreview(): void {
		if (previewUrl) {
			URL.revokeObjectURL(previewUrl);
			previewUrl = null;
		}
		selectedFile = null;
		errorMessage = null;
		if (fileInputRef) {
			fileInputRef.value = '';
		}
	}
</script>

<div class="space-y-4">
	{#if previewUrl}
		<!-- Preview Mode -->
		<div class="relative">
			<div class="border-border bg-surface-dark overflow-hidden rounded-lg border-2">
				{#if previewUrl.endsWith('.pdf') || selectedFile?.type === 'application/pdf'}
					<div
						class="flex h-64 w-full flex-col items-center justify-center bg-gray-100 dark:bg-gray-800"
					>
						<svg class="h-16 w-16 text-red-500" fill="currentColor" viewBox="0 0 24 24">
							<path
								d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8l-6-6zM6 20V4h7v5h5v11H6z"
							/>
							<text x="8" y="16" font-size="6" fill="currentColor">PDF</text>
						</svg>
						<p class="text-text-secondary mt-2 text-sm">PDF Document</p>
					</div>
				{:else}
					<img src={previewUrl} alt="Receipt preview" class="h-64 w-full object-contain" />
				{/if}
			</div>
			<Button
				variant="danger"
				onclick={clearPreview}
				class="absolute top-2 right-2 !rounded-full !p-2"
			>
				{#snippet leading()}<XIcon class="h-4 w-4" />{/snippet}
			</Button>
		</div>
	{:else}
		<!-- Drop Zone -->
		<div
			role="button"
			tabindex="0"
			class="cursor-pointer rounded-lg border-2 border-dashed p-8 text-center transition-colors {isDragging
				? 'border-primary bg-primary-dark/20'
				: 'border-border bg-surface hover:border-primary hover:bg-surface-light'}"
			ondragover={handleDragOver}
			ondragleave={handleDragLeave}
			ondrop={handleDrop}
			onclick={openFileDialog}
			onkeydown={(e) => e.key === 'Enter' && openFileDialog()}
		>
			<UploadIcon class="text-text-secondary mx-auto h-12 w-12" />
			<p class="text-text-secondary mt-4">
				{m.receipt_upload_description()}
			</p>
			<p class="text-text-tertiary mt-2 text-sm">{m.receipt_upload_formats()}</p>
		</div>

		<!-- Hidden File Input -->
		<input
			bind:this={fileInputRef}
			type="file"
			{accept}
			class="hidden"
			onchange={handleFileInputChange}
		/>

		<!-- Action Buttons -->
		<div class="flex justify-center gap-4">
			<Button onclick={openFileDialog}>
				{#snippet leading()}<FileIcon class="h-5 w-5" />{/snippet}
				{m.file_choose()}
			</Button>
		</div>
	{/if}

	<!-- Error Message -->
	{#if errorMessage}
		<div class="bg-danger-dark/20 border-danger text-danger-light rounded-lg border px-4 py-3">
			<p class="text-sm">{errorMessage}</p>
		</div>
	{/if}
</div>
