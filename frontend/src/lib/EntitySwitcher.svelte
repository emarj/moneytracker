<script>
    import { useQuery, useQueryClient } from "@sveltestack/svelte-query";
    import { getEntities } from "../data";
    import { entityID } from "../entity";

    const queryClient = useQueryClient();

    const entitiesQuery = useQuery("entities", () => getEntities());

    entityID.subscribe((value) => {
        console.log(`entity changed (${value}), invalidating queries`);
        //queryClient.invalidateQueries();
    });
</script>

{#if $entitiesQuery.isLoading}
    Loading...
{:else if $entitiesQuery.error}
    Error: {$entitiesQuery.error?.message}
{:else if $entitiesQuery.data}
    Entity: <select bind:value={$entityID}>
        {#each $entitiesQuery.data as ent}
            <option value={ent.id}>{ent.name}</option>
        {/each}
    </select>
{/if}
