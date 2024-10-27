import { writable, type Writable } from 'svelte/store';

export interface Message {
	type: string;
	data: string;
}

class WebSocketStore {
	private socket: WebSocket | null = null;
	public messages: Writable<Message[]> = writable([]);
	public status: Writable<string> = writable('disconnected');

	/**
	 * Establishes a WebSocket connection.
	 * @param url The WebSocket server URL
	 */
	connect(url: string): void {
		if (this.socket) {
			this.disconnect();
		}

		this.socket = new WebSocket(url);

		this.socket.onopen = () => {
			this.status.set('connected');
			console.log('WebSocket connected');
		};

		this.socket.onmessage = (event: MessageEvent) => {
			try {
				const msg: Message = JSON.parse(event.data);
				this.messages.update((current) => [...current, msg]);
			} catch (error) {
				console.error('Failed to parse message:', error);
			}
		};

		this.socket.onclose = (event: CloseEvent) => {
			this.status.set('disconnected');
			console.log(`WebSocket disconnected: ${event.reason}`);
			this.socket = null;
			// Optionally implement reconnection logic here
		};

		this.socket.onerror = (event: Event) => {
			console.error('WebSocket error observed:', event);
			this.socket?.close();
		};
	}

	/**
	 * Sends a message through the WebSocket.
	 * @param message The message to send
	 */
	send(message: Message): void {
		if (this.socket && this.socket.readyState === WebSocket.OPEN) {
			this.socket.send(JSON.stringify(message));
		} else {
			console.warn('WebSocket is not open. Unable to send message.');
		}
	}

	/**
	 * Closes the WebSocket connection.
	 */
	disconnect(): void {
		if (this.socket) {
			this.socket.close();
			this.socket = null;
		}
	}
}

export const wsStore = new WebSocketStore();
