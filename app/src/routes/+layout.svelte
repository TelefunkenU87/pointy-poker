<script lang="ts">
	import '../app.css';
	let { children } = $props();
</script>

<div class="flex min-h-screen flex-col">
	<header class="bg-blue-600 p-4 text-white">
		<h1 class="text-xl">Pointy Poker</h1>
	</header>

	<main class="container mx-auto flex-grow p-4">
		{@render children()}
	</main>

	<footer class="bg-gray-800 p-4 text-center text-white">&copy; 2024 Tele Funk</footer>
</div>
