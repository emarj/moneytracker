<script lang="ts">
    import CircularProgress from "@smui/circular-progress";
    import Select, { Option } from "@smui/select";

    import { useQuery } from "@sveltestack/svelte-query";
    import { getAccounts } from "../api";
    import AccountName from "./AccountName.svelte";

    const accountsQuery = useQuery(["accounts"], () => getAccounts());

    export let label = "Account";

    export let type_id = null;
    export let firstSelected = true;
    export let disabled = false;
    export let entity_ids = null;
    export let invert = false;

    let value;
    export let account_id;
    export let entity_id = null;

    $: account_id = value?.id;
    $: entity_id = value?.owner_id;

    const filterByOwner = (list) =>
        list.filter((a) => {
            const res = entity_ids.includes(a.owner_id);
            return (res && !invert) || (res && invert);
        });
    const filterByType = (list) =>
        list.filter((a) => a.type_id === type_id || a.type_id === 1);

    const filter = (list) => {
        if (entity_ids && entity_ids.length > 0) {
            list = filterByOwner(list);
        }

        if (type_id !== null) {
            list = filterByType(list);
        }

        return list;
    };
</script>

<div>
    {#if $accountsQuery.isLoading}
        <span
            ><CircularProgress
                style="height: 32px; width: 32px;"
                indeterminate
            /></span
        >
    {:else if $accountsQuery.error}
        <span>An error has occurred: {$accountsQuery.error.message}</span>
    {:else}
        <Select variant="outlined" bind:value {label} {disabled}>
            <!--     {#if !firstSelected}
                <Option value={null} />
            {/if} -->
            {#each filter($accountsQuery.data) as account (account.id)}
                <Option value={account}><AccountName {account} /></Option>
            {/each}
        </Select>
    {/if}
</div>
