<script>
	let { children } = $props();
</script>

<main class="auth">
	<div>
		{@render children()}
	</div>
</main>

<style>
	div {
		max-width: 500px;
		height: 100dvh;
		padding: 0.5rem;
		margin: auto;
		align-content: center;
	}
</style>
