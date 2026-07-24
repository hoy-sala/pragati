declare module 'tabulator-tables' {
	export class TabulatorFull {
		constructor(selector: string | HTMLElement, options: Record<string, any>);
		setData(data: any[]): void;
		getData(): any[];
		clearData(): void;
		destroy(): void;
		on(event: string, callback: (...args: any[]) => void): void;
	}
}
