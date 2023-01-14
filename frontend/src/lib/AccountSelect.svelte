<script lang="ts">
    import CircularProgress from "@smui/circular-progress";
    import Select, { Option } from "@smui/select";

    import { useQuery } from "@sveltestack/svelte-query";
    import { getAccounts } from "../api";
    import AccountName from "./AccountName.svelte";

    const accountsQuery = useQuery(["accounts"], () => getAccounts());

    export let neg = false;
    export let type_id = null;
    export let firstSelected = true;
    export let disabled = false;

    export let owner_id = null;

    export let value;
    export let label = "Account";

    let accounts = [];

    $: accounts = $accountsQuery?.data?.filter(
        (a) =>
            owner_id === null ||
            (((a.owner.id == owner_id && !neg) ||
                (a.owner.id != owner_id && neg)) &&
                (type_id === null || a.type_id == type_id || a.is_world))
    );
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
            {#if !firstSelected}
                <Option value={null} />
            {/if}
            {#each accounts as account (account.id)}
                <Option value={account.id}><AccountName {account} /></Option>
            {/each}
            <!-- <svelte:fragment slot="helperText">{helperText}</svelte:fragment> -->
        </Select>
    {/if}
</div>
