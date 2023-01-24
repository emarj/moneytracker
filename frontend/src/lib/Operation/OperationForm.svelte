<script lang="ts">
    import {
        emptyOperation,
        emptyTransaction,
        type Operation,
    } from "../../model";
    import Textfield from "@smui/textfield";
    import DatePicker from "../DatePicker.svelte";
    import AccountSelect from "../AccountSelect.svelte";
    import CategorySelect from "../CategorySelect.svelte";
    import Button from "@smui/button/src/Button.svelte";
    import { useMutation } from "@sveltestack/svelte-query";
    import { addOperation } from "../../api";
    import { push } from "svelte-spa-router";
    import { JSONPretty } from "../../util/utils";
    import DecimalInput from "../DecimalInput.svelte";

    export let op: Operation = structuredClone(emptyOperation);

    const mutation = useMutation((op) => addOperation(op), {
        onSuccess: (data: number) => {
            push("/");
        },
    });

    const removeTransaction = (i: number) => {
        op.transactions.splice(i, 1);
        op.transactions = op.transactions;
    };
</script>

<Textfield
    variant="outlined"
    bind:value={op.description}
    label="Description"
    style="width: 100%;"
/>
<CategorySelect bind:value={op.category_id} />

<div class="transactions">
    <ul>
        {#each op.transactions as t, i}
            <li>
                <DatePicker bind:timestamp={t.timestamp} />
                <div>
                    <AccountSelect
                        label="From"
                        bind:account_id={t.from_id}
                    /><AccountSelect label="To" bind:account_id={t.to_id} />
                </div>
                <DecimalInput bind:value={t.amount} decimalDigits={3} />
                <Button on:click={() => removeTransaction(i)}>Delete</Button>
            </li>
        {/each}
    </ul>
    <Button
        on:click={() =>
            (op.transactions = [
                ...op.transactions,
                structuredClone(emptyTransaction),
            ])}>Add Transaction</Button
    >
</div>
<div>
    <Button
        variant="raised"
        on:click={() => {
            $mutation.mutate(op);
        }}>Submit</Button
    >
</div>

<h3>Preview</h3>
<pre>{JSONPretty(op)}</pre>

<style lang="scss">
    ul,
    .transactions {
        margin: 1rem 0;
    }

    li {
        margin: 1rem 0;
        padding: 0.7rem;
        padding: 0.5rem;
        border: 1px solid #333;
        border-radius: 0.6rem;

        & > div {
            display: flex;
        }
    }
</style>
