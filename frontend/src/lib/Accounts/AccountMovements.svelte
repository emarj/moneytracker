<script>
    import CircularProgress from "@smui/circular-progress";
    import { useQuery, useQueryClient } from "@sveltestack/svelte-query";
    import { getTransactionsByAccount } from "../../api";
    import { DateFMT } from "../../util/utils";
    import Amount from "../Amount.svelte";

    export let id;

    const queryClient = useQueryClient();

    const transactionsQuery = useQuery(["transactions", "account", id], () =>
        getTransactionsByAccount(id)
    );

    export const refresh = () => {
        queryClient.invalidateQueries({
            queryKey: ["transactions", "account", id],
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
                    <td>{DateFMT(t.operation.timestamp).substring(0, 5)}</td>
                    <td><strong>{t.operation.description}</strong></td>
                    <td>
                        <Amount
                            value={t.amount}
                            negative={t.from.id == id}
                            hide_plus={false}
                        /></td
                    >
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
