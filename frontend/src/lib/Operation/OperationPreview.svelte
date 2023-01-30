<script lang="ts">
    import type { Operation } from "../../model";
    import Amount from "../Amount.svelte";
    import OperationTransactions from "./OperationTransactions.svelte";
    import { push } from "svelte-spa-router";
    import AccountTag from "../AccountTag.svelte";
    import IconButton from "@smui/icon-button";

    export let op: Operation;

    const computeTotal = (op: Operation): number => {
        if (op.transactions) {
            return op.transactions.reduce(
                (sum: number, t) => sum + t.amount * t.sign,
                0
            );
        }
    };

    let total: number;
    $: total = computeTotal(op);
</script>

<div class:expense={total < 0} class:income={total > 0}>
    <!-- <span class="date">{DateFMT(op.modified_on)}</span> -->
    <span class="desc">
        {#if op.description === "" || op.description.startsWith("@")}
            <strong>{op.category.full_name}</strong>
        {/if}
        {op.description}
    </span>

    {#if op.transactions && total !== 0}
        : <span>
            <Amount value={total} hide_plus={false} />
        </span>
    {:else if op.balances && op.balances.length == 1}
        {@const bal = op.balances[0]}
        <AccountTag account={bal.account} />
        <Amount value={bal.value} hide_plus={true} />
    {/if}

    <IconButton
        on:click={() => {
            push(`/operation/${op.id}`);
        }}
        class="more material-icons">more_horiz</IconButton
    >
    <div class="transactions">
        <OperationTransactions transactions={op.transactions} />
    </div>
    {#if op.balances && op.balances.length > 1}
        <div class="balances">
            {#each op.balances as bal}
                {bal.value}
            {/each}
        </div>
    {/if}
</div>

<style lang="scss">
    div {
        position: relative;

        :global(button.more) {
            position: absolute;
            top: -1rem;
            right: 1rem;
        }
    }
</style>
