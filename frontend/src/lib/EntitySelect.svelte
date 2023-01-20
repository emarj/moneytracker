<script>
    import Select, { Option } from "@smui/select";
    import { useQuery } from "@sveltestack/svelte-query";
    import { getAllEntities } from "../api";
    import CircularProgress from "@smui/circular-progress";
    import { entityID } from "../store";
    import { Item } from "@smui/list";

    export let label = "Entity";
    export let value = null;
    export let not = null;
    export let disabled = false;

    const entitiesQuery = useQuery("entities", () => getAllEntities());

    let entities = [];
    $: entities = $entitiesQuery?.data?.filter(
        (e) => !e.is_system && (!not || (not && e.id != not))
    );

    export let style = "material";
</script>

<div>
    {#if $entitiesQuery.isLoading}
        <CircularProgress style="height: 32px; width: 32px;" indeterminate />
    {:else if $entitiesQuery.error}
        Error: {$entitiesQuery.error?.message}
    {:else if $entitiesQuery.data}
        {#if style == "material"}
            <Select variant="outlined" bind:value {label} {disabled}>
                {#each entities as entity (entity.id)}
                    <Option value={entity.id}>{entity.display_name}</Option>
                {/each}
                <!--  <svelte:fragment slot="helperText">{helperText}</svelte:fragment> -->
            </Select>
        {:else if style == "simple"}
            <select bind:value>
                {#each entities as entity (entity.id)}
                    <option value={entity.id}>{entity.display_name}</option>
                {/each}
            </select>
        {:else if style == "menu"}
            {#each entities as entity (entity.id)}
                <Item
                    href="javascript:void(0)"
                    on:click={() => (value = entity.id)}
                    activated={$entityID === entity.id}
                >
                    {entity.display_name}
                </Item>
            {/each}
        {/if}
    {/if}
</div>
