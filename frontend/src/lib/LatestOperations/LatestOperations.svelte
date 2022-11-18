<script>
    import { isExpense, isIncome, isInternal } from "../../transactions";
    import { useQuery, useQueryClient } from "@sveltestack/svelte-query";
    import { DateFMT } from "../../util/utils";
    import { getOperationsByEntity } from "../../data";
    import Amount from "../Amount.svelte";
    import OperationTransactions from "./OperationTransactions.svelte";

    const eID = 1;

    const queryClient = useQueryClient();
    const operationsQuery = useQuery(["operations", "entity", eID], () =>
        getOperationsByEntity(eID)
    );

    const computeTotal = (op) => {
        if (op.transactions) {
            return op.transactions.reduce((sum, t) => {
                if (isExpense(t, eID)) {
                    return sum - t.amount;
                } else if (isIncome(t, eID)) {
                    return sum + t.amount;
                } else {
                    return sum;
                }
            }, 0);
        }
    };
</script>

<div>
    <h2>Last transactions</h2>
    <!--<button
        on:click={() => {
            console.log(`Invalidating Transactions of Entity ${eID}`);
            queryClient.invalidateQueries(["transactions", "entity", eID]);
        }}>Refresh</button
    >-->
    {#if $operationsQuery.isLoading}
        Loading...
    {:else if $operationsQuery.error}
        Error: {$operationsQuery.error?.message}
    {:else if $operationsQuery.data}
        <ol>
            {#each $operationsQuery.data as op}
                {@const total = computeTotal(op)}
                <li>
                    <span class="date">{DateFMT(op.timestamp)}</span>
                    <span class="desc">
                        {op.description}:
                    </span>
                    {#if op.transactions && total !== 0}
                        <span>
                            <Amount
                                value={Math.abs(total)}
                                negative={total < 0}
                                hide_plus={false}
                            />
                        </span>
                    {/if}
                    <OperationTransactions
                        {eID}
                        transactions={op.transactions}
                    />
                </li>
            {/each}
        </ol>
    {/if}
</div>

<style lang="scss">
    ol {
        list-style: none;
        width: 100%;
        padding: 0;

        & > li {
            width: 100%;
            padding: 1rem 1rem 0.5rem;
            border-radius: 10px;
            margin-bottom: 1em;
            background-color: rgb(238, 238, 238);
            border: 3px solid transparent;
            box-shadow: 0 10px 10px rgba(192, 192, 192, 0.4);
            position: relative;

            font-size: 1.1rem;

            .date {
            }

            :global(.amount) {
                font-weight: bold;
            }

            ul {
                padding-left: 2rem;
                font-size: 1rem;
            }
        }
    }
</style>
