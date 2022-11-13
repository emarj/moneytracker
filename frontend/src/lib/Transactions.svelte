<script>
    import { useQuery } from "@sveltestack/svelte-query";
    import { validate_component } from "svelte/internal";
    import { getTransactionsByEntity } from "../data";
    import Amount from "./Amount.svelte";

    const eID = 1;

    const transactionsQuery = useQuery(["transactions", eID], () =>
        getTransactionsByEntity(eID)
    );

    const computeClass = (feID, teID) => {
        if (feID == teID) return "tx-internal";
        else if (feID == eID) return "tx-expense";
        else return "tx-income";
    };
</script>

<h2>Last transactions</h2>
{#if $transactionsQuery.isLoading}
    Loading...
{:else if $transactionsQuery.error}
    Error: {$transactionsQuery.error.message}
{:else}
    <ol>
        {#each $transactionsQuery.data as t}
            <li class={computeClass(t.from.entity_id, t.to.entity_id)}>
                {t.timestamp}: {t.from_id} -> {t.to_id}
                <Amount value={t.amount} />
            </li>
        {/each}
    </ol>
{/if}

<style lang="scss">
    ol {
        list-style: none;

        li {
            padding: 1em;
            border-radius: 20px;
            margin-bottom: 1em;
            border: 3px solid transparent;

            &.tx-internal {
                border-color: darkgrey;
            }

            &.tx-expense {
                border-color: firebrick;
            }

            &.tx-income {
                border-color: rgb(175, 195, 34);
            }
        }
    }
</style>
