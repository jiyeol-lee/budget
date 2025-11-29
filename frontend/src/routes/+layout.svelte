<script lang="ts">
	import '../app.css';
	import { Button } from '$lib';
	import Toast from '$lib/components/Toast.svelte';
	import { MenuIcon, XIcon } from 'lucide-svelte';
	import type { Snippet } from 'svelte';
	import * as m from '$lib/paraglide/messages';
	import { page } from '$app/state';

	interface Props {
		children: Snippet;
	}

	let { children }: Props = $props();

	// Mobile menu state
	let mobileMenuOpen = $state(false);

	function toggleMobileMenu() {
		mobileMenuOpen = !mobileMenuOpen;
	}

	function closeMobileMenu() {
		mobileMenuOpen = false;
	}

	// Navigation items configuration
	const navItems = [
		{ href: '/', label: () => m.nav_dashboard(), exact: true },
		{ href: '/budget', label: () => m.nav_budget(), exact: true },
		{ href: '/expected-expenses', label: () => m.nav_expected_expenses(), exact: true },
		{ href: '/actual-expenses', label: () => m.nav_actual_expenses(), exact: true },
		{
			href: '/receipts/process',
			label: () => m.nav_receipts(),
			exact: false,
			matchPrefix: '/receipts'
		}
	];

	// Helper function to determine if a nav item is active
	function isActive(item: (typeof navItems)[number]): boolean {
		const pathname = page.url.pathname;
		if (item.exact) {
			return pathname === item.href;
		}
		const prefix = item.matchPrefix ?? item.href;
		// Match exact prefix path OR any child paths under the prefix
		return pathname === item.href || pathname === prefix || pathname.startsWith(prefix + '/');
	}

	// CSS class strings for navigation items
	const baseClasses = 'rounded-md px-3 py-2 text-sm font-medium transition-colors';
	const inactiveClasses = 'text-nav-item hover:text-nav-item-hover hover:bg-surface-light';
	const activeClasses = 'text-nav-item-active bg-nav-item-active-bg';
</script>

<div class="bg-background text-text-primary flex min-h-screen flex-col">
	<!-- Header -->
	<header class="bg-nav-bg border-border sticky top-0 z-40 border-b shadow-sm">
		<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
			<div class="flex h-16 items-center justify-between">
				<div class="flex items-center">
					<a href="/" class="text-primary text-xl font-bold" onclick={closeMobileMenu}>
						{m.app_title()}
					</a>
				</div>

				<!-- Desktop Navigation -->
				<nav class="hidden space-x-1 md:flex">
					{#each navItems as item (item.href)}
						{@const active = isActive(item)}
						<a
							href={item.href}
							class="{baseClasses} {active ? activeClasses : inactiveClasses}"
							aria-current={active ? 'page' : undefined}
						>
							{item.label()}
						</a>
					{/each}
				</nav>

				<!-- Mobile menu button -->
				<Button
					variant="ghost"
					onclick={toggleMobileMenu}
					aria-label={m.nav_toggle_menu()}
					class="p-2 md:hidden"
				>
					{#if mobileMenuOpen}
						<XIcon class="h-6 w-6" />
					{:else}
						<MenuIcon class="h-6 w-6" />
					{/if}
				</Button>
			</div>
		</div>

		<!-- Mobile Navigation Menu -->
		{#if mobileMenuOpen}
			<div class="border-border bg-nav-bg border-t md:hidden">
				<div class="space-y-1 px-2 pt-2 pb-3">
					{#each navItems as item (item.href)}
						{@const active = isActive(item)}
						<a
							href={item.href}
							onclick={closeMobileMenu}
							class="block {baseClasses} text-base {active ? activeClasses : inactiveClasses}"
							aria-current={active ? 'page' : undefined}
						>
							{item.label()}
						</a>
					{/each}
				</div>
			</div>
		{/if}
	</header>

	<!-- Main Content -->
	<main class="mx-auto w-full max-w-7xl flex-1 px-4 py-4 sm:px-6 sm:py-6 lg:px-8 lg:py-8">
		{@render children()}
	</main>

	<!-- Footer -->
	<footer class="bg-nav-bg border-border mt-auto border-t">
		<div class="mx-auto max-w-7xl px-4 py-4 sm:px-6 lg:px-8">
			<p class="text-text-secondary text-center text-sm">
				{m.app_title()} &copy; {new Date().getFullYear()}
			</p>
		</div>
	</footer>
</div>

<!-- Toast Notifications -->
<Toast />
