<script lang="ts">
	// import { sleep } from '$lib';
	import Component from './comp.svelte';
	let val: boolean;
</script>

<h1>Welcome to SvelteKit</h1>
<div>val={val}</div>
<Component bind:variable={val}></Component>

<p>Visit <a href="https://kit.svelte.dev">kit.svelte.dev</a> to read the documentation</p>
