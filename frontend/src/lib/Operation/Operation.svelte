<script lang="ts">
    import { DateFMT, JSONPretty } from "../../util/utils";
    import type { Operation } from "../../model";
    import Amount from "../Amount.svelte";
    import OperationTransactions from "./OperationTransactions.svelte";
    import { entityID } from "../../store";
    import { deleteOperation } from "../../api";
    import { useMutation, useQueryClient } from "@sveltestack/svelte-query";
    import { messageStore } from "../../store";
    import { pop, push } from "svelte-spa-router";
    import AccountOrEntityTag from "../AccountOrEntityTag.svelte";
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

    const queryClient = useQueryClient();

    const mutation = useMutation((opID: number) => deleteOperation(opID), {
        onMutate: () => {
            pop();
        },
        onSuccess: (data: number) => {
            $messageStore = { text: `Operation delete successfully!` };
            queryClient.invalidateQueries();
        },
    });
</script>

<div class:expense={total < 0} class:income={total > 0}>
    <span class="date">{DateFMT(op.modified_on)}</span>
    <span class="desc">
        {op.description}
    </span>

    {#if op.transactions && total !== 0}
        : <span>
            <Amount
                value={Math.abs(total)}
                negative={total < 0}
                hide_plus={false}
            />
        </span>
    {:else if op.balances && op.balances.length == 1}
        {@const bal = op.balances[0]}
        <AccountOrEntityTag account={bal.account} eID={$entityID} />
        <Amount
            value={Math.abs(bal.value)}
            negative={bal.value < 0}
            hide_plus={true}
        />
    {/if}
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
    <IconButton
        class="material-icons"
        on:click={(event) => {
            event.preventDefault();
            $mutation.mutate(op.id);
        }}
        disabled={$mutation.isLoading}>delete</IconButton
    >
</div>
<pre>
    {JSONPretty(op)}
</pre>

<style lang="scss">
    div {
        position: relative;
        button.more {
            position: absolute;
            top: 1rem;
            right: 1rem;
        }
    }
</style>
