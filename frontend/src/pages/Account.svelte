<script>
    import CircularProgress from "@smui/circular-progress";
    import { useQuery } from "@sveltestack/svelte-query";
    import { getAccount } from "../api";
    import AccountDetail from "../lib/Accounts/AccountDetail.svelte";

    export let params = {};

    const accountQuery = useQuery(["account", params.id], () =>
        getAccount(params.id)
    );
</script>

{#if $accountQuery.isLoading}
    <span
        ><CircularProgress
            style="height: 32px; width: 32px;"
            indeterminate
        /></span
    >
{:else if $accountQuery.error}
    <span>An error has occurred: {$accountQuery.error.message}</span>
{:else if $accountQuery.data}
    <AccountDetail account={$accountQuery.data} />
{/if}
