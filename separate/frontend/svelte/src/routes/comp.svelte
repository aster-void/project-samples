<script lang="ts">
	export let variable=false;
</script>

{#if variable}
	<div>you passed true</div>
{:else}
	<div>you passed false</div>
{/if}
