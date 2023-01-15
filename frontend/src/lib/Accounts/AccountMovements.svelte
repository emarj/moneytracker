<script>
    import CircularProgress from "@smui/circular-progress";
    import { useQuery, useQueryClient } from "@sveltestack/svelte-query";
    import { getAccountBalances, getTransactionsByAccount } from "../../api";
    import { TimestampFMT } from "../../util/utils";
    import Amount from "../Amount.svelte";

    export let id;

    const queryClient = useQueryClient();

    const transactionsQuery = useQuery(["transactions", "account", id], () =>
        getTransactionsByAccount(id)
    );
    const balancesQuery = useQuery(["balances", "account", id], () =>
        getAccountBalances(id)
    );

    export const refresh = () => {
        queryClient.invalidateQueries({
            queryKey: ["transactions", "account", id],
        });
        queryClient.invalidateQueries({
            queryKey: ["balances", "account", id],
        });
    };
</script>

<div class:fetching={$transactionsQuery.isFetching}>
    {#if $transactionsQuery.isLoading}
        <CircularProgress style="height: 32px; width: 32px;" indeterminate />
    {:else if $transactionsQuery.error}
        Error: {$transactionsQuery.error.message}
    {:else}
        <table>
            {#each $transactionsQuery.data as t (t.id)}
                <tr>
                    <td>{TimestampFMT(t.timestamp).substring(0, 5)}</td>
                    <td><strong>{t.operation.description}</strong></td>
                    <td>
                        <Amount
                            value={t.amount}
                            negative={t.from.id == id}
                            hide_plus={false}
                        />
                    </td>
                </tr>
            {/each}
        </table>
    {/if}
    {#if $balancesQuery.isLoading}
        <CircularProgress style="height: 32px; width: 32px;" indeterminate />
    {:else if $balancesQuery.error}
        Error: {$balancesQuery.error.message}
    {:else}
        <table>
            {#each $balancesQuery.data as b (b.id + b.timestamp)}
                <tr class="balance">
                    <td>{TimestampFMT(b.timestamp).substring(0, 5)}</td>
                    <td>
                        <strong>
                            {#if b.operation}
                                {b.operation.description}
                            {:else if b.comment}
                                {b.comment}
                            {:else}
                                Balance adjust
                            {/if}
                        </strong>
                    </td>
                    <td>
                        <strong>{b.value}</strong>
                        <!--   {#if b.delta}
                            <Amount
                                value={Math.abs(b.delta)}
                                negative={b.delta < 0}
                                hide_plus={false}
                            />
                            ({b.value})
                        {:else}
                            <strong>{b.value}</strong>
                        {/if} -->
                    </td>
                </tr>
            {/each}
        </table>
    {/if}
</div>

<style lang="scss">
    div.fetching {
        opacity: 0.5;
    }
    table {
        font-family: monospace;
        font-size: 0.9rem;
        border-collapse: collapse;

        tr.balance {
            border: 1px solid black;
            border-left: none;
            border-right: none;
        }
    }

    div :global(.amount) {
        font-weight: bold;
    }

    div :global(.amount-pos) {
        color: green;
    }
    div :global(.amount-neg) {
        color: red;
    }
</style>
