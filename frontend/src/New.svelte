<script lang="ts">
    import Tab, { Label } from "@smui/tab";
    import TabBar from "@smui/tab-bar";
    import OperationForm from "./lib/Operation/OperationForm.svelte";
    import { push, querystring } from "svelte-spa-router";
    import ExpenseForm from "./lib/Expense/ExpenseForm.svelte";

    let params = new URLSearchParams($querystring);

    const type: string = params.get("type");
    let active;

    switch (type) {
        case "operation":
            active = "Operation";
            break;
        default:
            active = "Expense";
            break;
    }

    $: push(`/add/?type=${active.toLowerCase()}`);
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
        <ExpenseForm />
    {:else}
        <h2>New Operation</h2>
        <OperationForm />
    {/if}
</div>
