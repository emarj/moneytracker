<script>
    import { useQuery, useQueryClient } from "@sveltestack/svelte-query";
    import { getTransactionsByAccount } from "../data";
    import Amount from "./Amount.svelte";

    export let id;

    const queryClient = useQueryClient();

    const transactionsQuery = useQuery(["transactions", id], () =>
        getTransactionsByAccount(id)
    );
</script>

{#if $transactionsQuery.isLoading}
    Loading...
{:else if $transactionsQuery.error}
    Error: {$transactionsQuery.error.message}
{:else}
    <table>
        {#each $transactionsQuery.data as t}
            <tr>
                <td>{new Date(t.timestamp).toLocaleString("en-GB")}</td>
                <td
                    ><Amount
                        value={t.amount}
                        negative={t.from_id == id}
                        hide_plus={false}
                    /></td
                >
            </tr>
        {/each}
    </table>
{/if}

<style>
    table {
        font-family: monospace;
        font-size: 0.9rem;
    }
</style>
