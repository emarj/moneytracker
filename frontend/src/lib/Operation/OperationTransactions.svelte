<script>
    import { isExpense, isIncome, isInternal } from "../../transactions";
    import { TimestampFMT } from "../../util/utils";
    import AccountTag from "../AccountTag.svelte";
    import Amount from "../Amount.svelte";

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
                        <AccountTag account={t.from} /> →
                        <AccountTag account={t.to} /></span
                    >
                    <Amount
                        value={t.amount}
                        negative={isExpense(t)}
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
