<script>
    import { isExpense, isIncome, isInternal } from "../../transactions";
    import { useQuery, useQueryClient } from "@sveltestack/svelte-query";
    import { DateFMT } from "../../util/utils";
    import { getOperationsByEntity } from "../../data";
    import Amount from "../Amount.svelte";
    import OperationTransactions from "./OperationTransactions.svelte";
    import { entityID } from "../../entity";

    const queryClient = useQueryClient();
    const operationsQuery = useQuery(["operations", "entity", $entityID], () =>
        getOperationsByEntity($entityID)
    );

    let expanded = false;

    const computeTotal = (op) => {
        if (op.transactions) {
            return op.transactions.reduce((sum, t) => {
                const amount = new Number(t.amount);
                if (isExpense(t, $entityID)) {
                    return sum - amount;
                } else if (isIncome(t, $entityID)) {
                    return sum + amount;
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
            {#each $operationsQuery.data as op (op.id)}
                {@const total = computeTotal(op)}

                <!-- svelte-ignore a11y-click-events-have-key-events -->
                <li
                    on:click={() => (expanded = !expanded)}
                    class:expense={total < 0}
                    class:income={total > 0}
                    class:expanded
                >
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
                    <div class="transactions">
                        <OperationTransactions transactions={op.transactions} />
                    </div>
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
            --top-color: rgb(185, 185, 185);

            width: 100%;
            padding: 1rem 1rem 0.5rem;
            border-radius: 10px;
            margin-bottom: 1em;
            background: rgb(238, 238, 238);
            background: linear-gradient(
                180deg,
                var(--top-color) 5px,
                rgb(238, 238, 238) 5px
            );
            box-shadow: 10px 10px 10px rgba(192, 192, 192, 0.4);
            position: relative;

            font-size: 1.1rem;

            &.expense {
                --top-color: red;
            }

            &.income {
                --top-color: green;
            }

            :global(.amount) {
                font-weight: bold;
            }

            .transactions {
                display: none;
            }

            &.expanded .transactions {
                display: block;
            }
        }
    }
</style>
