<script>
    import Select, { Option } from "@smui/select";
    import { useQuery } from "@sveltestack/svelte-query";
    import { getEntities } from "../api";
    import CircularProgress from "@smui/circular-progress";

    export let label = "Entity";
    export let helperText = "";
    export let value = null;
    export let not = null;
    const entitiesQuery = useQuery("entities", () => getEntities());

    let entities = [];
    $: entities = $entitiesQuery?.data?.filter(
        (e) => !not || (not && e.id != not)
    );
</script>

<div>
    {#if $entitiesQuery.isLoading}
        <CircularProgress style="height: 32px; width: 32px;" indeterminate />
    {:else if $entitiesQuery.error}
        Error: {$entitiesQuery.error?.message}
    {:else if $entitiesQuery.data}
        <Select variant="outlined" bind:value {label}>
            {#each entities as entity (entity.id)}
                <Option value={entity.id}>{entity.name}</Option>
            {/each}
            <svelte:fragment slot="helperText">{helperText}</svelte:fragment>
        </Select>
    {/if}
</div>
