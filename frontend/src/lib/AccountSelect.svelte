<script lang="ts">
    import { useQuery } from "@sveltestack/svelte-query";
    import { getAccounts } from "../data";
    import AccountName from "./AccountName.svelte";

    const accountsQuery = useQuery(["accounts"], () => getAccounts());

    export let neg = false;
    export let credit = false;
    export let firstSelected = true;

    export let owner_id = null;

    export let value;

    let accounts = [];

    $: accounts = $accountsQuery?.data?.filter(
        (a) =>
            owner_id &&
            ((a.owner.id == owner_id && !neg) ||
                (a.owner.id != owner_id && neg)) &&
            ((a.is_credit && credit) || (!a.is_credit && !credit))
    );

    $: if (firstSelected && accounts && accounts.length > 0)
        value = accounts[0].id;
</script>

{#if $accountsQuery.isLoading}
    <span>Loading...</span>
{:else if $accountsQuery.error}
    <span>An error has occurred: {$accountsQuery.error.message}</span>
{:else}
    <select bind:value>
        {#each accounts as account}
            <option value={account.id}><AccountName {account} /></option>
        {/each}
    </select>
{/if}
