<script lang="ts">
    import { entityID } from "../../store";
    import { emptyOperation, type Operation } from "../../model";
    import Textfield from "@smui/textfield";
    import DatePicker from "../DatePicker.svelte";
    import AccountSelect from "../AccountSelect.svelte";
    import CategorySelect from "../CategorySelect.svelte";
    import Button from "@smui/button/src/Button.svelte";

    export let op: Operation = structuredClone(emptyOperation);
</script>

<DatePicker bind:timestamp={op.timestamp} />
<Textfield
    variant="outlined"
    bind:value={op.description}
    label="Description"
    style="width: 100%;"
/>
<CategorySelect bind:value={op.category} />

<Button on:click={() => (op.transactions = [...op.transactions, { amount: 0 }])}
    >Add Transaction</Button
>
<ul>
    {#each op.transactions as t}
        <li>
            From: <AccountSelect bind:value={t.from} /> To: <AccountSelect
                bind:value={t.to}
            />
            <Textfield bind:value={t.amount} />
        </li>
    {/each}
</ul>

<h3>Preview</h3>
<pre>{JSON.stringify(op, null, 4)}</pre>
