import adapter from '@sveltejs/adapter-node';

export default {
	kit: {
		adapter: adapter({
			out: 'build',
			precompress: true,
			polyfill: false
		}),
		alias: {
			$lib: './src/lib'
		}
	}
};
