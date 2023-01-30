<script lang="ts">
    import { DateFMT, JSONPretty } from "../../util/utils";
    import type { Operation } from "../../model";
    import Amount from "../Amount.svelte";
    import OperationTransactions from "./OperationTransactions.svelte";
    import { updateOperation, deleteOperation } from "../../api";
    import { useMutation, useQueryClient } from "@sveltestack/svelte-query";
    import { messageStore } from "../../store";
    import { pop, push } from "svelte-spa-router";
    import AccountTag from "../AccountTag.svelte";
    import IconButton from "@smui/icon-button";
    import Button from "@smui/button/src/Button.svelte";

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

    const deleteMutation = useMutation(
        (opID: number) => deleteOperation(opID),
        {
            onMutate: () => {
                pop();
            },
            onSuccess: () => {
                $messageStore = { text: `Operation delete successfully!` };
                queryClient.invalidateQueries();
            },
        }
    );

    const updateMutation = useMutation((op: Operation) => updateOperation(op), {
        onSuccess: () => {
            $messageStore = { text: `Operation edited successfully!` };
            queryClient.invalidateQueries();
            edit = false;
        },
    });

    let edit = false;
</script>

<div class:expense={total < 0} class:income={total > 0}>
    <span class="date">{DateFMT(op.modified_on)}</span>
    {#if edit}
        <input type="text" bind:value={op.description} />
        <Button on:click={() => $updateMutation.mutate(op)}>Save</Button>
    {:else}
        <span class="desc">
            {op.description}
        </span>
        <Button on:click={() => (edit = true)}>Edit</Button>
    {/if}

    {#if op.transactions && total !== 0}
        : <span>
            <Amount value={total} hide_plus={false} />
        </span>
    {:else if op.balances && op.balances.length == 1}
        {@const bal = op.balances[0]}
        <AccountTag account={bal.account} />
        <Amount value={bal.value} hide_plus={true} />
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
            $deleteMutation.mutate(op.id);
        }}
        disabled={$deleteMutation.isLoading}>delete</IconButton
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
