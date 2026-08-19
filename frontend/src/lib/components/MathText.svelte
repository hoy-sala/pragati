<script lang="ts">
	import katex from 'katex';
	import 'katex/contrib/mhchem';

	let { text = '' }: { text?: string } = $props();

	function renderMath(src: string, display: boolean): string {
		try {
			return katex.renderToString(src, {
				displayMode: display,
				throwOnError: false,
				output: 'html'
			});
		} catch {
			return src;
		}
	}

	const inlineRe = /\\\((.*?)\\\)/gs;
	const displayRe = /\\\[([\s\S]*?)\\\]/g;

	let html = $derived.by(() => {
		if (!text) return '';
		let out = '';
		let last = 0;
		while (last < text.length) {
			// Find the earliest math token (display or inline) from `last`.
			let m: RegExpExecArray | null = null;
			let kind: 'display' | 'inline' | null = null;
			displayRe.lastIndex = last;
			const dm = displayRe.exec(text);
			inlineRe.lastIndex = last;
			const im = inlineRe.exec(text);
			if (dm && im) {
				if (dm.index <= im.index) {
					m = dm;
					kind = 'display';
				} else {
					m = im;
					kind = 'inline';
				}
			} else if (dm) {
				m = dm;
				kind = 'display';
			} else if (im) {
				m = im;
				kind = 'inline';
			}
			if (!m) {
				out += text.slice(last);
				break;
			}
			if (m.index > last) out += text.slice(last, m.index);
			out += renderMath(m[1], kind === 'display');
			last = m.index + m[0].length;
		}
		return out;
	});
</script>

{@html html}