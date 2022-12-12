<script lang="ts">
    import { entityID } from "../../store";
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

    export let op: Operation = structuredClone(emptyOperation);

    const mutation = useMutation((op) => addOperation(op), {
        onSuccess: (data: number) => {
            push("/");
        },
    });
</script>

<DatePicker bind:timestamp={op.timestamp} />
<Textfield
    variant="outlined"
    bind:value={op.description}
    label="Description"
    style="width: 100%;"
/>
<CategorySelect bind:value={op.category} />

<div class="transactions">
    <ul>
        {#each op.transactions as t}
            <li>
                <div>
                    <AccountSelect
                        label="From"
                        bind:value={t.from.id}
                    /><AccountSelect label="To" bind:value={t.to.id} />
                </div>
                <Textfield
                    label="Amount"
                    bind:value={t.amount}
                    suffix="€"
                    input$pattern={"\\d+(\\.\\d{2})?"}
                />
            </li>
        {/each}
    </ul>
    <Button
        on:click={() =>
            (op.transactions = op.transactions.slice(
                0,
                op.transactions.length - 1
            ))}>Remove Last</Button
    >
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
<pre>{JSON.stringify(op, null, 4)}</pre>

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
