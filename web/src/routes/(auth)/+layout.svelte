<script>
	let { children } = $props();
</script>

<main>
	{@render children()}
</main>

<style>
	main {
		max-width: 500px;
		height: 100dvh;
		padding: 0.5rem;
		margin: auto;
		align-content: center;
	}
</style>
