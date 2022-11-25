<script>
    import OperationForm from "./lib/Operation/OperationForm.svelte";
    import { querystring } from "svelte-spa-router";
    import { entityID } from "./entity";
    import ExpenseForm from "./lib/Expense/ExpenseForm.svelte";

    /*let params = new URLSearchParams($querystring);

    const from = params.get("from");
    const to = params.get("to");*/

    let op = {
        description: null,
        owner: { id: $entityID },
        transactions: [{}],
    };

    let advancedMode = false;
</script>

<div>
    <label>
        Advanced Mode
        <input type="checkbox" bind:checked={advancedMode} />
    </label>
    {#if !advancedMode}
        <h2>New Expense</h2>
        <ExpenseForm />
    {:else}
        <h2>New Operation</h2>
        <OperationForm {op} />
    {/if}
</div>
