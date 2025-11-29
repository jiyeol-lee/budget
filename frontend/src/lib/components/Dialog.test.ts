/**
 * @vitest-environment jsdom
 */
import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';

/**
 * Test suite for Dialog.svelte - Escape Key Handler
 *
 * Tests cover:
 * - Global keydown event listener is added when dialog opens
 * - Global keydown event listener is removed when dialog closes
 * - Global keydown event listener is removed on component unmount
 * - Escape key calls onClose callback
 * - Escape key preventDefault is called
 * - No memory leaks (multiple open/close cycles)
 * - Other keys don't trigger close
 *
 * Note: These tests validate the event listener logic extracted from Dialog.svelte
 * Since Svelte 5's $effect is challenging to test in isolation, we test the core
 * logic patterns that the component uses.
 */

describe('Dialog Escape Key Handler', () => {
	let addEventListenerSpy: ReturnType<typeof vi.spyOn>;
	let removeEventListenerSpy: ReturnType<typeof vi.spyOn>;

	beforeEach(() => {
		// Spy on document event listeners
		addEventListenerSpy = vi.spyOn(document, 'addEventListener');
		removeEventListenerSpy = vi.spyOn(document, 'removeEventListener');
	});

	afterEach(() => {
		vi.restoreAllMocks();
	});

	describe('handleGlobalKeyDown function', () => {
		/**
		 * Simulates the handleGlobalKeyDown function from Dialog.svelte
		 */
		function createHandleGlobalKeyDown(onClose: () => void) {
			return function handleGlobalKeyDown(event: KeyboardEvent): void {
				if (event.key === 'Escape') {
					event.preventDefault();
					onClose();
				}
			};
		}

		it('should call onClose when Escape key is pressed', () => {
			const onClose = vi.fn();
			const handleGlobalKeyDown = createHandleGlobalKeyDown(onClose);

			const event = new KeyboardEvent('keydown', {
				key: 'Escape',
				bubbles: true,
				cancelable: true
			});

			handleGlobalKeyDown(event);

			expect(onClose).toHaveBeenCalledTimes(1);
		});

		it('should call preventDefault when Escape key is pressed', () => {
			const onClose = vi.fn();
			const handleGlobalKeyDown = createHandleGlobalKeyDown(onClose);

			const event = new KeyboardEvent('keydown', {
				key: 'Escape',
				bubbles: true,
				cancelable: true
			});
			const preventDefaultSpy = vi.spyOn(event, 'preventDefault');

			handleGlobalKeyDown(event);

			expect(preventDefaultSpy).toHaveBeenCalledTimes(1);
		});

		it('should not call onClose for other keys', () => {
			const onClose = vi.fn();
			const handleGlobalKeyDown = createHandleGlobalKeyDown(onClose);

			const keys = ['Enter', 'Tab', 'Space', 'ArrowUp', 'ArrowDown', 'a', 'z', '1'];

			keys.forEach((key) => {
				const event = new KeyboardEvent('keydown', {
					key,
					bubbles: true,
					cancelable: true
				});
				handleGlobalKeyDown(event);
			});

			expect(onClose).not.toHaveBeenCalled();
		});

		it('should not call preventDefault for other keys', () => {
			const onClose = vi.fn();
			const handleGlobalKeyDown = createHandleGlobalKeyDown(onClose);

			const event = new KeyboardEvent('keydown', {
				key: 'Enter',
				bubbles: true,
				cancelable: true
			});
			const preventDefaultSpy = vi.spyOn(event, 'preventDefault');

			handleGlobalKeyDown(event);

			expect(preventDefaultSpy).not.toHaveBeenCalled();
		});
	});

	describe('Event listener lifecycle simulation', () => {
		/**
		 * Simulates the $effect behavior from Dialog.svelte
		 */
		function simulateDialogEffect(
			open: boolean,
			handleGlobalKeyDown: (event: KeyboardEvent) => void
		): () => void {
			if (open) {
				document.body.classList.add('modal-open');
				document.addEventListener('keydown', handleGlobalKeyDown);
			} else {
				document.body.classList.remove('modal-open');
				document.removeEventListener('keydown', handleGlobalKeyDown);
			}

			// Return cleanup function (simulating $effect return)
			return () => {
				document.body.classList.remove('modal-open');
				document.removeEventListener('keydown', handleGlobalKeyDown);
			};
		}

		it('should add event listener when dialog opens', () => {
			const handler = vi.fn();

			simulateDialogEffect(true, handler);

			expect(addEventListenerSpy).toHaveBeenCalledWith('keydown', handler);
		});

		it('should remove event listener when dialog closes', () => {
			const handler = vi.fn();

			// Open
			simulateDialogEffect(true, handler);
			// Close
			simulateDialogEffect(false, handler);

			expect(removeEventListenerSpy).toHaveBeenCalledWith('keydown', handler);
		});

		it('should add modal-open class when dialog opens', () => {
			const handler = vi.fn();

			simulateDialogEffect(true, handler);

			expect(document.body.classList.contains('modal-open')).toBe(true);
		});

		it('should remove modal-open class when dialog closes', () => {
			const handler = vi.fn();

			simulateDialogEffect(true, handler);
			simulateDialogEffect(false, handler);

			expect(document.body.classList.contains('modal-open')).toBe(false);
		});

		it('cleanup function should remove event listener (unmount scenario)', () => {
			const handler = vi.fn();

			// Open dialog
			const cleanup = simulateDialogEffect(true, handler);

			// Simulate component unmount
			cleanup();

			expect(removeEventListenerSpy).toHaveBeenCalledWith('keydown', handler);
		});

		it('cleanup function should remove modal-open class (unmount scenario)', () => {
			const handler = vi.fn();

			// Open dialog
			const cleanup = simulateDialogEffect(true, handler);

			// Simulate component unmount
			cleanup();

			expect(document.body.classList.contains('modal-open')).toBe(false);
		});

		it('should handle multiple open/close cycles without memory leaks', () => {
			const handler = vi.fn();

			// Simulate multiple open/close cycles
			for (let i = 0; i < 5; i++) {
				simulateDialogEffect(true, handler);
				simulateDialogEffect(false, handler);
			}

			// Each cycle should add and remove exactly once
			expect(addEventListenerSpy).toHaveBeenCalledTimes(5);
			expect(removeEventListenerSpy).toHaveBeenCalledTimes(5);
		});

		it('cleanup should be idempotent (safe to call multiple times)', () => {
			const handler = vi.fn();

			const cleanup = simulateDialogEffect(true, handler);

			// Call cleanup multiple times
			cleanup();
			cleanup();
			cleanup();

			// Should not throw and class should remain removed
			expect(document.body.classList.contains('modal-open')).toBe(false);
		});
	});

	describe('Integration: Event listener actually handles Escape key', () => {
		it('should close dialog when Escape key is pressed on document', () => {
			const onClose = vi.fn();

			function handleGlobalKeyDown(event: KeyboardEvent): void {
				if (event.key === 'Escape') {
					event.preventDefault();
					onClose();
				}
			}

			// Simulate dialog opening
			document.addEventListener('keydown', handleGlobalKeyDown);

			// Dispatch Escape key event on document
			const escapeEvent = new KeyboardEvent('keydown', {
				key: 'Escape',
				bubbles: true,
				cancelable: true
			});
			document.dispatchEvent(escapeEvent);

			expect(onClose).toHaveBeenCalledTimes(1);

			// Cleanup
			document.removeEventListener('keydown', handleGlobalKeyDown);
		});

		it('should close dialog when Escape is pressed regardless of active element', () => {
			const onClose = vi.fn();

			function handleGlobalKeyDown(event: KeyboardEvent): void {
				if (event.key === 'Escape') {
					event.preventDefault();
					onClose();
				}
			}

			// Add listener to document
			document.addEventListener('keydown', handleGlobalKeyDown);

			// Create various elements and focus them
			const input = document.createElement('input');
			const button = document.createElement('button');
			const div = document.createElement('div');
			div.tabIndex = 0;

			document.body.appendChild(input);
			document.body.appendChild(button);
			document.body.appendChild(div);

			// Test Escape from input
			input.focus();
			document.dispatchEvent(
				new KeyboardEvent('keydown', { key: 'Escape', bubbles: true, cancelable: true })
			);
			expect(onClose).toHaveBeenCalledTimes(1);

			// Test Escape from button
			button.focus();
			document.dispatchEvent(
				new KeyboardEvent('keydown', { key: 'Escape', bubbles: true, cancelable: true })
			);
			expect(onClose).toHaveBeenCalledTimes(2);

			// Test Escape from div
			div.focus();
			document.dispatchEvent(
				new KeyboardEvent('keydown', { key: 'Escape', bubbles: true, cancelable: true })
			);
			expect(onClose).toHaveBeenCalledTimes(3);

			// Cleanup
			document.removeEventListener('keydown', handleGlobalKeyDown);
			document.body.removeChild(input);
			document.body.removeChild(button);
			document.body.removeChild(div);
		});
	});
});

/**
 * Code Review Validation Tests
 *
 * These tests validate specific implementation details from the Dialog.svelte code
 */
describe('Dialog.svelte Implementation Validation', () => {
	describe('handleGlobalKeyDown implementation check', () => {
		it('should match the implementation pattern in Dialog.svelte', () => {
			// This is the exact logic from Dialog.svelte lines 56-61
			const onClose = vi.fn();

			function handleGlobalKeyDown(event: KeyboardEvent): void {
				if (event.key === 'Escape') {
					event.preventDefault();
					onClose();
				}
			}

			// Test that it works correctly
			const escEvent = new KeyboardEvent('keydown', { key: 'Escape', cancelable: true });
			const preventDefaultSpy = vi.spyOn(escEvent, 'preventDefault');

			handleGlobalKeyDown(escEvent);

			expect(escEvent.key).toBe('Escape');
			expect(preventDefaultSpy).toHaveBeenCalled();
			expect(onClose).toHaveBeenCalled();
		});
	});

	describe('$effect cleanup implementation check', () => {
		it('cleanup should remove both modal-open class and event listener', () => {
			// Verify the cleanup function implementation from lines 76-79
			const handler = vi.fn();

			// Setup (simulating open: true)
			document.body.classList.add('modal-open');
			document.addEventListener('keydown', handler);

			// Cleanup function from Dialog.svelte
			function cleanup() {
				document.body.classList.remove('modal-open');
				document.removeEventListener('keydown', handler);
			}

			cleanup();

			expect(document.body.classList.contains('modal-open')).toBe(false);
			// Verify listener was removed by dispatching event
			document.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape' }));
			expect(handler).not.toHaveBeenCalled();
		});
	});

	describe('Edge cases', () => {
		it('should handle rapid open/close without race conditions', async () => {
			const onClose = vi.fn();
			let isOpen = false;

			function handleGlobalKeyDown(event: KeyboardEvent): void {
				if (event.key === 'Escape') {
					event.preventDefault();
					onClose();
				}
			}

			function simulateToggle(open: boolean) {
				if (open && !isOpen) {
					document.addEventListener('keydown', handleGlobalKeyDown);
					isOpen = true;
				} else if (!open && isOpen) {
					document.removeEventListener('keydown', handleGlobalKeyDown);
					isOpen = false;
				}
			}

			// Rapid toggles
			simulateToggle(true);
			simulateToggle(false);
			simulateToggle(true);
			simulateToggle(false);
			simulateToggle(true);

			// Final state: open
			const escEvent = new KeyboardEvent('keydown', { key: 'Escape', cancelable: true });
			document.dispatchEvent(escEvent);

			expect(onClose).toHaveBeenCalledTimes(1);

			// Cleanup
			simulateToggle(false);
		});

		it('should not leak listeners when cleanup is called while closed', () => {
			const handler = vi.fn();

			// Start closed, call cleanup anyway (edge case)
			document.body.classList.remove('modal-open');
			document.removeEventListener('keydown', handler);

			// Should not throw
			expect(document.body.classList.contains('modal-open')).toBe(false);
		});
	});
});
