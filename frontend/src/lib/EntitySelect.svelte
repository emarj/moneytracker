<script>
    import { useQuery } from "@sveltestack/svelte-query";
    import { getEntities } from "../data";

    export let value = null;
    export let not = null;
    const entitiesQuery = useQuery("entities", () => getEntities());
</script>

{#if $entitiesQuery.isLoading}
    Loading...
{:else if $entitiesQuery.error}
    Error: {$entitiesQuery.error?.message}
{:else if $entitiesQuery.data}
    Entity: <select bind:value>
        {#each $entitiesQuery.data.filter((e) => !not || (not && e.id != not)) as ent}
            <option value={ent.id}>{ent.name}</option>
        {/each}
    </select>
{/if}
