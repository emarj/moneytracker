<script>
    import { useQuery, useQueryClient } from "@sveltestack/svelte-query";
    import { getTransactionsByAccount } from "../data";
    import AccountBalance from "./AccountBalance.svelte";
    import AccountTransactions from "./AccountTransactions.svelte";

    export let account;

    const queryClient = useQueryClient();

    const transactionsQuery = useQuery(["transactions", account.id], () =>
        getTransactionsByAccount(account.id)
    );
</script>

<div>
    <h3>{account.display_name} (id = {account.id})</h3>

    <AccountBalance id={account.id} />
    <!--{#if $transactionsQuery.isLoading}
        Loading...
    {:else if $transactionsQuery.error}
        Error: {$transactionsQuery.error.message}
    {:else}
        <AccountTransactions transactions={$transactionsQuery.data} />
    {/if}-->
</div>

<style>
    div {
        padding: 1em;
        background: rgb(34, 193, 195);
        background: linear-gradient(
            24deg,
            rgba(34, 193, 195, 1) 0%,
            rgba(253, 187, 45, 1) 100%
        );
        height: auto;
        width: 300px;
        border-radius: 15px;
    }
</style>
