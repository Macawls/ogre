// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

export default defineConfig({
	site: 'https://ogre.macawls.dev',
	integrations: [
		starlight({
			title: 'Ogre',
			favicon: '/favicon.png',
			social: [{ icon: 'github', label: 'GitHub', href: 'https://github.com/macawls/ogre' }],
			components: {
				Header: './src/components/Header.astro',
			},
			expressiveCode: {
				themes: ['one-dark-pro', 'one-light'],
			},
			customCss: ['./src/styles/one-dark.css'],
			sidebar: [
				{
					label: 'Overview',
					items: [
						{ label: 'Introduction', slug: '' },
						{ label: 'Installation', slug: 'getting-started/installation' },
						{ label: 'Quick Start', slug: 'getting-started/quick-start' },
						{ label: 'Examples', slug: 'getting-started/examples' },
						{ label: 'Playground', slug: 'getting-started/playground' },
					],
				},
				{
					label: 'Usage',
					items: [
						{ label: 'Go Library', slug: 'guides/library' },
						{ label: 'CLI', slug: 'guides/cli' },
						{ label: 'HTTP Server', slug: 'guides/server' },
						{ label: 'JSX Builder', slug: 'guides/jsx' },
						{ label: 'Docker', slug: 'guides/docker' },
					],
				},
				{
					label: 'Features',
					items: [
						{ label: 'Tailwind CSS', slug: 'guides/tailwind' },
						{ label: 'Custom Fonts', slug: 'guides/fonts' },
						{ label: 'Emoji & RTL', slug: 'getting-started/emoji-and-rtl' },
					],
				},
				{
					label: 'Reference',
					items: [
						{ label: 'Go API', slug: 'reference/api' },
						{ label: 'HTTP Endpoints', slug: 'reference/http' },
						{ label: 'CSS Properties', slug: 'reference/css' },
						{ label: 'Tailwind Classes', slug: 'reference/tailwind' },
						{ label: 'Architecture', slug: 'advanced/architecture' },
						{ label: 'Satori Comparison', slug: 'advanced/satori-comparison' },
					],
				},
			],
		}),
	],
});
