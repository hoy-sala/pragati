export type ToastType = 'success' | 'error' | 'info';

export type ToastItem = {
	id: number;
	message: string;
	type: ToastType;
};

class ToastStore {
	toasts = $state<ToastItem[]>([]);
	#next = 0;

	show(message: string, type: ToastType = 'info', duration = 3200) {
		const id = ++this.#next;
		this.toasts.push({ id, message, type });
		setTimeout(() => this.dismiss(id), duration);
	}

	dismiss(id: number) {
		this.toasts = this.toasts.filter(t => t.id !== id);
	}
}

export const toastStore = new ToastStore();

export function toast(message: string, type: ToastType = 'info') {
	toastStore.show(message, type);
}
