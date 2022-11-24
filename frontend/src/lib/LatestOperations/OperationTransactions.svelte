<script>
    import { isExpense, isIncome, isInternal } from "../../transactions";
    import AccountTag from "../AccountTag.svelte";
    import Amount from "../Amount.svelte";
    import { entityID } from "../../entity";

    export let transactions = [];
</script>

<div>
    {#if transactions}
        <ul>
            {#each transactions as t}
                <li>
                    <span class="fromto">
                        <AccountTag account={t.from} eID={$entityID} /> →
                        <AccountTag account={t.to} eID={$entityID} /></span
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
