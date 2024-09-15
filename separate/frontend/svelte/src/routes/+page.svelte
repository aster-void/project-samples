<script lang="ts">
	import Component from './comp.svelte';
	let bool: boolean;

	// this is read-only timer using svelte/store
	import { increment } from '$lib';
	let reactive: number;
	$: reactive = $increment;
</script>

<h1>Welcome to SvelteKit</h1>
<div>val={bool}</div>
<div>reactive={reactive}</div>
<Component bind:variable={bool}></Component>

<p>Visit <a href="https://kit.svelte.dev">kit.svelte.dev</a> to read the documentation</p>
