<script lang="ts">
    import { DateFMT } from "../../util/utils";
    import type { Operation } from "../../model";
    import Amount from "../Amount.svelte";
    import OperationTransactions from "./OperationTransactions.svelte";
    import { isExpense, isIncome, isInternal } from "../../transactions";
    import { entityID } from "../../store";
    import { deleteOperation } from "../../api";
    import { useMutation, useQueryClient } from "@sveltestack/svelte-query";
    import { messageStore } from "../../store";

    export let op: Operation;

    const computeTotal = (op: Operation): number => {
        if (op.transactions) {
            return op.transactions.reduce((sum: number, t) => {
                if (isExpense(t, $entityID)) {
                    return sum - t.amount;
                } else if (isIncome(t, $entityID)) {
                    return sum + t.amount;
                } else {
                    return sum;
                }
            }, 0);
        }
    };

    let total: number;
    $: total = computeTotal(op);

    const queryClient = useQueryClient();

    const mutation = useMutation((opID: number) => deleteOperation(opID), {
        onSuccess: (data: number) => {
            $messageStore = { text: `Operation delete successfully!` };
            queryClient.invalidateQueries();
        },
    });
</script>

<div class:expense={total < 0} class:income={total > 0}>
    <span class="date">{DateFMT(op.modified_on)}</span>
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
    <button
        class="delete"
        on:click={(event) => {
            event.preventDefault();
            $mutation.mutate(op.id);
        }}
        disabled={$mutation.isLoading}
    >
        X
    </button>
    <div class="transactions">
        <OperationTransactions transactions={op.transactions} />
    </div>
</div>

<style>
    button.delete {
        position: absolute;
        top: 1rem;
        right: 1rem;
    }
</style>
