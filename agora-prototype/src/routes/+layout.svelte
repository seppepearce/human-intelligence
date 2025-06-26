<script>
	import '../app.css';
</script>

<svelte:head>
	<title>Digital Agora - Ancient Wisdom, Modern Tools</title>
	<meta name="description" content="Experience discourse as the ancient Greeks intended - with natural constraints that enhance rather than limit conversation." />
</svelte:head>

<main class="min-h-screen">
	<slot />
</main>

<style>
	:global(body) {
		margin: 0;
		padding: 0;
		min-height: 100vh;
	}

	:global(#svelte) {
		min-height: 100vh;
		display: flex;
		flex-direction: column;
	}

	main {
		flex: 1;
		display: flex;
		flex-direction: column;
	}
</style>
