<script lang="ts" context="module">
	// Module context if needed (e.g., load functions)
</script>

<script lang="ts">
	import { goto } from '$app/navigation';
	import { createEventDispatcher } from 'svelte';

	let roomName: string = '';
	let selectedPoints: string = 'fibonacci';

	const pointOptions = [
		{ value: 'fibonacci', label: 'Fibonacci' },
		{ value: 't-shirt', label: 'T-Shirt Sizes' }
		// Add more options as needed
	];

	const dispatch = createEventDispatcher();

	/**
	 * Handles the creation of a new room.
	 */
	function createRoom(): void {
		if (roomName.trim() !== '') {
			// Navigate to the room page using SvelteKit's goto
			goto(`/${encodeURIComponent(roomName)}?points=${encodeURIComponent(selectedPoints)}`);
			dispatch('roomCreated', { roomName, selectedPoints });
		}
	}

	/**
	 * Handles form submission via Enter key.
	 * @param event Keyboard event
	 */
	function handleKeydown(event: KeyboardEvent): void {
		if (event.key === 'Enter') {
			createRoom();
		}
	}
</script>

<div class="mx-auto max-w-md rounded bg-white p-6 shadow-md">
	<h2 class="mb-4 text-2xl">Create a New Planning Poker Room</h2>

	<div class="mb-4">
		<label for="roomName" class="mb-2 block text-gray-700">Room Name</label>
		<input
			id="roomName"
			type="text"
			bind:value={roomName}
			class="w-full rounded border px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
			placeholder="Enter room name"
			on:keydown={handleKeydown}
		/>
	</div>

	<div class="mb-4">
		<label for="pointOptions" class="mb-2 block text-gray-700">Point Options</label>
		<select
			id="pointOptions"
			bind:value={selectedPoints}
			class="w-full rounded border px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
		>
			{#each pointOptions as option}
				<option value={option.value}>{option.label}</option>
			{/each}
		</select>
	</div>

	<button
		on:click={createRoom}
		class="w-full rounded bg-blue-500 py-2 text-white transition-colors duration-200 hover:bg-blue-600"
		disabled={!roomName.trim()}
	>
		Create Room
	</button>
</div>
