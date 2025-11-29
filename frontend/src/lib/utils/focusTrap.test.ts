/**
 * @vitest-environment jsdom
 */
import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { focusTrap, type FocusTrapOptions } from './focusTrap';

/**
 * Test suite for focusTrap Svelte action
 *
 * Tests cover:
 * - Tab cycling through focusable elements
 * - Shift+Tab backwards cycling
 * - Initial focus options (first, last, none, specific element)
 * - Focus restoration on destroy
 * - Edge cases (empty containers, single element, dynamic elements)
 */

describe('focusTrap', () => {
	let container: HTMLElement;

	// Helper to create a container with focusable elements
	function createTestContainer(html: string): HTMLElement {
		const div = document.createElement('div');
		div.innerHTML = html;
		document.body.appendChild(div);
		return div;
	}

	// Helper to simulate Tab key press
	function pressTab(element: HTMLElement, shiftKey = false): void {
		const event = new KeyboardEvent('keydown', {
			key: 'Tab',
			shiftKey,
			bubbles: true,
			cancelable: true
		});
		element.dispatchEvent(event);
	}

	// Helper to simulate other key press
	function pressKey(element: HTMLElement, key: string): void {
		const event = new KeyboardEvent('keydown', {
			key,
			bubbles: true,
			cancelable: true
		});
		element.dispatchEvent(event);
	}

	// Wait for microtask queue to flush
	async function flushMicrotasks(): Promise<void> {
		await new Promise((resolve) => queueMicrotask(resolve));
	}

	beforeEach(() => {
		document.body.innerHTML = '';
		// Create a button to serve as "previously focused" element
		const triggerButton = document.createElement('button');
		triggerButton.id = 'trigger';
		triggerButton.textContent = 'Open Dialog';
		document.body.appendChild(triggerButton);
		triggerButton.focus();
	});

	afterEach(() => {
		// Clean up DOM
		document.body.innerHTML = '';
		// Reset any mocked timers
		vi.restoreAllMocks();
	});

	describe('Tab cycling', () => {
		it('should cycle forward through focusable elements with Tab', async () => {
			container = createTestContainer(`
				<button id="btn1">Button 1</button>
				<button id="btn2">Button 2</button>
				<button id="btn3">Button 3</button>
			`);

			const action = focusTrap(container);
			await flushMicrotasks();

			expect(document.activeElement?.id).toBe('btn1');

			// Tab from last element should wrap to first
			const btn3 = container.querySelector('#btn3') as HTMLElement;
			btn3.focus();
			expect(document.activeElement?.id).toBe('btn3');

			// Tab from last should wrap to first
			pressTab(container);
			expect(document.activeElement?.id).toBe('btn1');

			action.destroy?.();
		});

		it('should cycle backward through focusable elements with Shift+Tab', async () => {
			container = createTestContainer(`
				<button id="btn1">Button 1</button>
				<button id="btn2">Button 2</button>
				<button id="btn3">Button 3</button>
			`);

			const action = focusTrap(container);
			await flushMicrotasks();

			expect(document.activeElement?.id).toBe('btn1');

			// Shift+Tab from first should wrap to last
			pressTab(container, true);
			expect(document.activeElement?.id).toBe('btn3');

			action.destroy?.();
		});

		it('should handle Tab when focus is outside the container', async () => {
			container = createTestContainer(`
				<button id="btn1">Button 1</button>
				<button id="btn2">Button 2</button>
			`);

			const action = focusTrap(container);
			await flushMicrotasks();

			expect(document.activeElement?.id).toBe('btn1');

			// Simulate focus being outside (e.g., on body)
			document.body.focus();

			// Tab should focus first element
			pressTab(container);
			expect(document.activeElement?.id).toBe('btn1');

			action.destroy?.();
		});

		it('should handle Shift+Tab when focus is outside the container', async () => {
			container = createTestContainer(`
				<button id="btn1">Button 1</button>
				<button id="btn2">Button 2</button>
			`);

			const action = focusTrap(container);
			await flushMicrotasks();

			expect(document.activeElement?.id).toBe('btn1');

			// Simulate focus being outside
			document.body.focus();

			// Shift+Tab should focus last element
			pressTab(container, true);
			expect(document.activeElement?.id).toBe('btn2');

			action.destroy?.();
		});
	});

	describe('Initial focus options', () => {
		it('should focus first element by default', async () => {
			container = createTestContainer(`
				<button id="btn1">Button 1</button>
				<button id="btn2">Button 2</button>
			`);

			const action = focusTrap(container);
			await flushMicrotasks();

			expect(document.activeElement?.id).toBe('btn1');

			action.destroy?.();
		});

		it('should focus first element when initialFocus is "first"', async () => {
			container = createTestContainer(`
				<button id="btn1">Button 1</button>
				<button id="btn2">Button 2</button>
			`);

			const action = focusTrap(container, { initialFocus: 'first' });
			await flushMicrotasks();

			expect(document.activeElement?.id).toBe('btn1');

			action.destroy?.();
		});

		it('should focus last element when initialFocus is "last"', async () => {
			container = createTestContainer(`
				<button id="btn1">Button 1</button>
				<button id="btn2">Button 2</button>
				<button id="btn3">Button 3</button>
			`);

			const action = focusTrap(container, { initialFocus: 'last' });
			await flushMicrotasks();

			expect(document.activeElement?.id).toBe('btn3');

			action.destroy?.();
		});

		it('should not auto-focus when initialFocus is "none"', async () => {
			const triggerButton = document.getElementById('trigger') as HTMLElement;
			triggerButton.focus();

			container = createTestContainer(`
				<button id="btn1">Button 1</button>
				<button id="btn2">Button 2</button>
			`);

			const action = focusTrap(container, { initialFocus: 'none' });
			await flushMicrotasks();

			// Focus should remain on trigger button (not moved to dialog content)
			expect(document.activeElement?.id).not.toBe('btn1');
			expect(document.activeElement?.id).not.toBe('btn2');

			action.destroy?.();
		});

		it('should focus specific element when initialFocus is an HTMLElement', async () => {
			container = createTestContainer(`
				<button id="btn1">Button 1</button>
				<input id="input1" type="text" />
				<button id="btn2">Button 2</button>
			`);

			const specificElement = container.querySelector('#input1') as HTMLElement;
			const action = focusTrap(container, { initialFocus: specificElement });
			await flushMicrotasks();

			expect(document.activeElement?.id).toBe('input1');

			action.destroy?.();
		});
	});

	describe('Focus restoration', () => {
		it('should restore focus to previously focused element on destroy', async () => {
			const triggerButton = document.getElementById('trigger') as HTMLElement;
			triggerButton.focus();
			expect(document.activeElement?.id).toBe('trigger');

			container = createTestContainer(`
				<button id="btn1">Button 1</button>
				<button id="btn2">Button 2</button>
			`);

			const action = focusTrap(container);
			await flushMicrotasks();

			expect(document.activeElement?.id).toBe('btn1');

			action.destroy?.();

			// Wait for requestAnimationFrame
			await new Promise((resolve) => requestAnimationFrame(resolve));

			expect(document.activeElement?.id).toBe('trigger');
		});
	});

	describe('Edge cases', () => {
		it('should handle container with no focusable elements', async () => {
			container = createTestContainer(`
				<div>Non-focusable content</div>
				<span>More non-focusable content</span>
			`);

			const action = focusTrap(container);
			await flushMicrotasks();

			// Container should have been given tabindex and focused
			expect(container.hasAttribute('tabindex')).toBe(true);

			// Tab should be prevented (no elements to cycle through)
			pressTab(container);
			// Should not throw

			action.destroy?.();
		});

		it('should handle container with single focusable element', async () => {
			container = createTestContainer(`
				<button id="onlyBtn">Only Button</button>
			`);

			const action = focusTrap(container);
			await flushMicrotasks();

			expect(document.activeElement?.id).toBe('onlyBtn');

			// Tab from the only element should wrap back to itself
			// (since first === last, Tab on last wraps to first)
			pressTab(container);
			expect(document.activeElement?.id).toBe('onlyBtn');

			// Shift+Tab should also keep focus on the only element
			pressTab(container, true);
			expect(document.activeElement?.id).toBe('onlyBtn');

			action.destroy?.();
		});

		it('should ignore non-Tab key events', async () => {
			container = createTestContainer(`
				<button id="btn1">Button 1</button>
				<button id="btn2">Button 2</button>
			`);

			const action = focusTrap(container);
			await flushMicrotasks();

			expect(document.activeElement?.id).toBe('btn1');

			// Press Enter - should not affect focus
			pressKey(container, 'Enter');
			expect(document.activeElement?.id).toBe('btn1');

			// Press Escape - should not affect focus (that's Dialog's job)
			pressKey(container, 'Escape');
			expect(document.activeElement?.id).toBe('btn1');

			action.destroy?.();
		});

		it('should exclude disabled buttons from focusable elements', async () => {
			container = createTestContainer(`
				<button id="btn1">Button 1</button>
				<button id="btn2" disabled>Disabled Button</button>
				<button id="btn3">Button 3</button>
			`);

			const action = focusTrap(container);
			await flushMicrotasks();

			expect(document.activeElement?.id).toBe('btn1');

			// Focus last, then Tab should go to first (skipping disabled)
			const btn3 = container.querySelector('#btn3') as HTMLElement;
			btn3.focus();

			pressTab(container);
			expect(document.activeElement?.id).toBe('btn1');

			// Shift+Tab from first should go to btn3, skipping disabled btn2
			pressTab(container, true);
			expect(document.activeElement?.id).toBe('btn3');

			action.destroy?.();
		});

		it('should exclude elements with tabindex="-1"', async () => {
			container = createTestContainer(`
				<button id="btn1">Button 1</button>
				<button id="btn2" tabindex="-1">Hidden from Tab</button>
				<button id="btn3">Button 3</button>
			`);

			const action = focusTrap(container);
			await flushMicrotasks();

			expect(document.activeElement?.id).toBe('btn1');

			const btn3 = container.querySelector('#btn3') as HTMLElement;
			btn3.focus();

			pressTab(container);
			expect(document.activeElement?.id).toBe('btn1');

			action.destroy?.();
		});

		it('should include various focusable element types', async () => {
			container = createTestContainer(`
				<a href="#" id="link1">Link</a>
				<button id="btn1">Button</button>
				<input id="input1" type="text" />
				<textarea id="textarea1"></textarea>
				<select id="select1"><option>Option</option></select>
				<div id="div1" tabindex="0">Focusable div</div>
			`);

			const action = focusTrap(container);
			await flushMicrotasks();

			expect(document.activeElement?.id).toBe('link1');

			// Focus last element and Tab should wrap to first
			const div1 = container.querySelector('#div1') as HTMLElement;
			div1.focus();
			expect(document.activeElement?.id).toBe('div1');

			pressTab(container);
			expect(document.activeElement?.id).toBe('link1');

			action.destroy?.();
		});
	});

	describe('Update functionality', () => {
		it('should update options via update method', async () => {
			container = createTestContainer(`
				<button id="btn1">Button 1</button>
				<button id="btn2">Button 2</button>
			`);

			const action = focusTrap(container, { initialFocus: 'first' });
			await flushMicrotasks();

			expect(document.activeElement?.id).toBe('btn1');

			// Update options (in real usage, this would affect new initialFocus calls)
			action.update?.({ initialFocus: 'last' });

			// The update itself doesn't re-trigger initial focus,
			// just updates the stored options
			expect(document.activeElement?.id).toBe('btn1');

			action.destroy?.();
		});
	});

	describe('Event listener cleanup', () => {
		it('should remove event listener on destroy', async () => {
			container = createTestContainer(`
				<button id="btn1">Button 1</button>
				<button id="btn2">Button 2</button>
			`);

			const removeEventListenerSpy = vi.spyOn(container, 'removeEventListener');

			const action = focusTrap(container);
			await flushMicrotasks();

			action.destroy?.();

			expect(removeEventListenerSpy).toHaveBeenCalledWith('keydown', expect.any(Function));
		});
	});
});
