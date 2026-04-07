<script>
    import Menu from '$lib/components/Menu.svelte';

	let { children } = $props();
</script>

<Menu />

<main>
	{@render children()}
</main>

<style>
	main {
		display: grid;
		overflow: hidden;
	}
</style>
