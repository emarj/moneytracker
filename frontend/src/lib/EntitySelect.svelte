<script lang="ts">
    import Select, { Option } from "@smui/select";
    import { useQuery } from "@sveltestack/svelte-query";
    import { getAllEntities } from "../api";
    import CircularProgress from "@smui/circular-progress";
    import { Item } from "@smui/list";

    export let label = "Entity";
    export let value = null;
    export let entities = [];
    export let invert = false;
    export let disabled = false;

    const entitiesQuery = useQuery("entities", () => getAllEntities());

    const filterByID = (list) => {
        if (!entities) return list;

        return list.filter((e) => {
            const includes = entities.includes(e.id);
            return (!invert && includes) || (invert && !includes);
        });
    };

    const filter = (list: any[]) =>
        filterByID(list.filter((e) => !e.is_system));

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
                {#each filter($entitiesQuery.data) as entity (entity.id)}
                    <Option value={entity.id}>{entity.display_name}</Option>
                {/each}
                <!--  <svelte:fragment slot="helperText">{helperText}</svelte:fragment> -->
            </Select>
            <!-- {:else if style == "simple"}
            <select bind:value>
                {#each entities as entity (entity.id)}
                    <option value={entity.id}>{entity.display_name}</option>
                {/each}
            </select> -->
            <!-- {:else if style == "menu"}
            {#each entities as entity (entity.id)}
                <Item
                    href="javascript:void(0)"
                    on:click={() => (value = entity.id)}
                    activated={$entityID === entity.id}
                >
                    {entity.display_name}
                </Item>
            {/each} -->
        {/if}
    {/if}
</div>
