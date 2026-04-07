<script>
	let { children } = $props();
</script>

<main>
	<div>
		{@render children()}
	</div>
</main>

<style>
	main {
		width: 100dvw;
	}

	div {
		max-width: 500px;
		height: 100dvh;
		padding: 0.5rem;
		margin: auto;
		align-content: center;
	}
</style>
