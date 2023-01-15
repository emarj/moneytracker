<script>
    import { isExpense, isIncome, isInternal } from "../../transactions";
    import { TimestampFMT } from "../../util/utils";
    import AccountOrEntityTag from "../AccountOrEntityTag.svelte";
    import Amount from "../Amount.svelte";
    import { entityID } from "../../store";

    export let transactions = [];
</script>

<div>
    {#if transactions}
        <ul>
            {#each transactions as t}
                <li>
                    <span class="date"
                        >{TimestampFMT(t.timestamp).slice(0, 10)}</span
                    >
                    <span class="fromto">
                        <AccountOrEntityTag account={t.from} eID={$entityID} /> →
                        <AccountOrEntityTag
                            account={t.to}
                            eID={$entityID}
                        /></span
                    >
                    <Amount
                        value={t.amount}
                        negative={isExpense(t, $entityID)}
                        hide_plus={isInternal(t)}
                    />
                </li>
            {/each}
        </ul>
    {/if}
</div>

<style lang="scss">
    ul {
        padding-left: 2rem;
        font-size: 1rem;

        li {
            margin-top: 0.5rem;
        }
    }
</style>
