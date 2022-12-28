<script>
    import Select, { Option } from "@smui/select";
    import { useQuery } from "@sveltestack/svelte-query";
    import { getCategories } from "../api";
    import CircularProgress from "@smui/circular-progress";

    export let label = "Category";
    export let value = null;

    const entitiesQuery = useQuery("categories", () => getCategories());
</script>

<div>
    {#if $entitiesQuery.isLoading}
        <CircularProgress style="height: 32px; width: 32px;" indeterminate />
    {:else if $entitiesQuery.error}
        Error: {$entitiesQuery.error?.message}
    {:else if $entitiesQuery.data}
        <Select variant="outlined" bind:value {label}>
            {#each $entitiesQuery.data as cat (cat.id)}
                <Option value={cat.id}
                    >{#if cat.parent_id}{cat.parent
                            .name}/{/if}{cat.name}</Option
                >
            {/each}
        </Select>
    {/if}
</div>
