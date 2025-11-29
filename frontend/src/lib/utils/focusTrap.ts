import type { Action } from 'svelte/action';

/**
 * Options for the focusTrap action
 */
export interface FocusTrapOptions {
	/**
	 * Which element to focus when the trap is activated:
	 * - 'first': Focus the first focusable element (default)
	 * - 'last': Focus the last focusable element
	 * - 'none': Don't auto-focus any element
	 * - HTMLElement: Focus a specific element
	 */
	initialFocus?: 'first' | 'last' | 'none' | HTMLElement;
}

/**
 * Selector for all focusable elements
 */
const FOCUSABLE_SELECTOR = [
	'a[href]',
	'button:not([disabled])',
	'input:not([disabled])',
	'textarea:not([disabled])',
	'select:not([disabled])',
	'[tabindex]:not([tabindex="-1"])'
].join(', ');

/**
 * Gets all focusable elements within a container
 */
function getFocusableElements(container: HTMLElement): HTMLElement[] {
	const elements = Array.from(container.querySelectorAll<HTMLElement>(FOCUSABLE_SELECTOR));

	// Filter out elements that are not visible or have tabindex="-1"
	return elements.filter((el) => {
		// Check if element is visible
		if (el.offsetParent === null && getComputedStyle(el).position !== 'fixed') {
			return false;
		}
		// Check tabindex (already filtered in selector, but double-check)
		const tabindex = el.getAttribute('tabindex');
		if (tabindex !== null && parseInt(tabindex, 10) < 0) {
			return false;
		}
		return true;
	});
}

/**
 * A Svelte action that traps focus within a container element.
 * Useful for modals, dialogs, and other overlay components.
 *
 * Features:
 * - Traps Tab/Shift+Tab navigation within the container
 * - Auto-focuses the first/last/specific element when activated
 * - Restores focus to the previously focused element when destroyed
 * - Handles dynamically added/removed focusable elements
 *
 * @example
 * ```svelte
 * <div use:focusTrap>
 *   <button>First</button>
 *   <input type="text" />
 *   <button>Last</button>
 * </div>
 *
 * <div use:focusTrap={{ initialFocus: 'last' }}>
 *   ...
 * </div>
 *
 * <div use:focusTrap={{ initialFocus: 'none' }}>
 *   ...
 * </div>
 * ```
 */
export const focusTrap: Action<HTMLElement, FocusTrapOptions | undefined> = (
	node,
	options = {}
) => {
	// Store the element that had focus before the trap was activated
	const previouslyFocusedElement = document.activeElement as HTMLElement | null;

	let currentOptions = options;

	/**
	 * Handles keydown events to trap Tab navigation
	 */
	function handleKeyDown(event: KeyboardEvent): void {
		if (event.key !== 'Tab') {
			return;
		}

		const focusableElements = getFocusableElements(node);

		if (focusableElements.length === 0) {
			// No focusable elements, prevent tab from leaving the container
			event.preventDefault();
			return;
		}

		const firstElement = focusableElements[0];
		const lastElement = focusableElements[focusableElements.length - 1];
		const activeElement = document.activeElement;

		// These are guaranteed to exist since we checked length > 0 above
		if (!firstElement || !lastElement) {
			return;
		}

		if (event.shiftKey) {
			// Shift+Tab: going backwards
			if (activeElement === firstElement || !node.contains(activeElement)) {
				event.preventDefault();
				lastElement.focus();
			}
		} else {
			// Tab: going forwards
			if (activeElement === lastElement || !node.contains(activeElement)) {
				event.preventDefault();
				firstElement.focus();
			}
		}
	}

	/**
	 * Sets initial focus based on options
	 */
	function setInitialFocus(): void {
		const initialFocus = currentOptions?.initialFocus ?? 'first';

		if (initialFocus === 'none') {
			return;
		}

		if (initialFocus instanceof HTMLElement) {
			// Focus a specific element
			initialFocus.focus();
			return;
		}

		const focusableElements = getFocusableElements(node);

		if (focusableElements.length === 0) {
			// If no focusable elements, make the container itself focusable
			if (!node.hasAttribute('tabindex')) {
				node.setAttribute('tabindex', '-1');
			}
			node.focus();
			return;
		}

		const targetElement =
			initialFocus === 'last'
				? focusableElements[focusableElements.length - 1]
				: focusableElements[0];

		if (targetElement) {
			targetElement.focus();
		}
	}

	/**
	 * Restores focus to the previously focused element
	 */
	function restoreFocus(): void {
		if (previouslyFocusedElement && typeof previouslyFocusedElement.focus === 'function') {
			// Use requestAnimationFrame to ensure DOM is ready
			requestAnimationFrame(() => {
				previouslyFocusedElement.focus();
			});
		}
	}

	// Set up the focus trap
	node.addEventListener('keydown', handleKeyDown);

	// Set initial focus after a microtask to ensure the DOM is ready
	queueMicrotask(() => {
		setInitialFocus();
	});

	return {
		update(newOptions: FocusTrapOptions | undefined) {
			currentOptions = newOptions ?? {};
		},
		destroy() {
			node.removeEventListener('keydown', handleKeyDown);
			restoreFocus();
		}
	};
};

export default focusTrap;
