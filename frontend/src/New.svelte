<script lang="ts">
    import Tab, { Label } from "@smui/tab";
    import TabBar from "@smui/tab-bar";
    import OperationForm from "./lib/Operation/OperationForm.svelte";
    import { querystring } from "svelte-spa-router";
    import { entityID, newOp, newExpense } from "./store";
    import ExpenseForm from "./lib/Expense/ExpenseForm.svelte";
    import { ExpenseToOperation } from "./model";

    /*let params = new URLSearchParams($querystring);

    const from: number = params.get("from");*/

    let active = "Expense";

    $: $newOp = ExpenseToOperation($newExpense);
</script>

<div>
    <TabBar tabs={["Expense", "Operation"]} let:tab bind:active>
        <!-- Note: the `tab` property is required! -->
        <Tab {tab}>
            <Label>{tab}</Label>
        </Tab>
    </TabBar>
    {#if active == "Expense"}
        <h2>New Expense</h2>
        <ExpenseForm bind:e={$newExpense} />
    {:else}
        <h2>New Operation</h2>
        <OperationForm bind:op={$newOp} />
    {/if}
</div>
