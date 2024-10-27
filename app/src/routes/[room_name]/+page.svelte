<!-- src/routes/[room-name]/+page.svelte -->
<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { wsStore } from '$lib/stores/websocket.svelte';
	import { derived } from 'svelte/store';

	// Reactive declarations using Svelte 5's optimized syntax
	let roomName: string = $state('');
	let pointOptions: string[] = $state([]);
	let selectedPoint: string | null = $state(null);
	let votes: Map<string, string> = $state(new Map());
	let connected: boolean = $state(false);

	// Extract roomName from the route parameters
	roomName = $page.params['room_name'];

	// Extract point options from query parameters or default
	const query = new URLSearchParams($page.url.search);
	// const points = query.get('points') || 'fibonacci';

	// Define point options based on selection
	switch (query.get('points')) {
		case 'fibonacci':
			pointOptions = ['1', '2', '3', '5', '8', '13', '21'];
			break;
		case 't-shirt':
			pointOptions = ['XS', 'S', 'M', 'L', 'XL'];
			break;
		default:
			pointOptions = ['1', '2', '3', '5', '8', '13', '21']; // Default
	}

	// Derived store to monitor WebSocket messages
	const voteMessages = derived(wsStore.messages, ($messages) =>
		$messages.filter((msg) => msg.type === 'vote')
	);

	// Handle component mount
	onMount(() => {
		const wsProtocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
		const wsUrl = `${wsProtocol}://${window.location.host}/ws/${encodeURIComponent(roomName)}`;
		wsStore.connect(wsUrl);
		connected = true;

		const unsubscribe = voteMessages.subscribe((messages) => {
			messages.forEach((msg) => {
				const { user, vote } = JSON.parse(msg.data);
				votes.set(user, vote);
			});
		});

		return () => {
			wsStore.disconnect();
			unsubscribe();
		};
	});

	/**
	 * Casts a vote by sending a message through the WebSocket.
	 */
	function castVote(): void {
		if (selectedPoint) {
			const message = {
				type: 'vote',
				data: JSON.stringify({ user: 'currentUser', vote: selectedPoint })
			};
			wsStore.send(message);
		}
	}

	/**
	 * Resets all votes by sending a reset message and clearing local state.
	 */
	function resetVotes(): void {
		const message = { type: 'reset', data: '' };
		wsStore.send(message);
		votes.clear();
	}

	/**
	 * Leaves the room and navigates back to the home page.
	 */
	function leaveRoom(): void {
		wsStore.disconnect();
		goto('/');
	}

	onDestroy(() => {
		wsStore.disconnect();
	});
</script>

<div class="flex flex-col items-center">
	<h2 class="mb-4 text-2xl">Room: {roomName}</h2>

	<div class="mb-4 w-full max-w-md">
		<label for="voteSelect" class="mb-2 block text-gray-700">Select your vote:</label>
		<select
			id="voteSelect"
			bind:value={selectedPoint}
			class="w-full rounded border px-3 py-2 focus:outline-none focus:ring-2 focus:ring-green-500"
		>
			<option value="" disabled selected>Select a point</option>
			{#each pointOptions as point}
				<option value={point}>{point}</option>
			{/each}
		</select>
		<button
			onclick={castVote}
			class="mt-2 w-full rounded bg-green-500 py-2 text-white transition-colors duration-200 hover:bg-green-600"
			disabled={!selectedPoint}
		>
			Cast Vote
		</button>
	</div>

	<div class="w-full max-w-md">
		<h3 class="mb-2 text-xl">Votes:</h3>
		<ul class="list-disc pl-5">
			{#each Array.from(votes.entries()) as [user, vote]}
				<li>{user}: {vote}</li>
			{/each}
		</ul>
		<button
			onclick={resetVotes}
			class="mt-4 w-full rounded bg-red-500 py-2 text-white transition-colors duration-200 hover:bg-red-600"
		>
			Reset Votes
		</button>
	</div>

	<button
		onclick={leaveRoom}
		class="mt-6 rounded bg-gray-500 px-4 py-2 text-white transition-colors duration-200 hover:bg-gray-600"
	>
		Leave Room
	</button>
</div>

<style>
	/* Scoped styles for the Room page */
</style>
