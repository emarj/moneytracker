<script lang="ts">
    import Tab, { Label } from "@smui/tab";
    import TabBar from "@smui/tab-bar";
    import OperationForm from "../lib/Operation/OperationForm.svelte";
    import { push, querystring } from "svelte-spa-router";
    import ExpenseForm from "../lib/Expense/ExpenseForm.svelte";
    import BalanceAdjustForm from "../lib/BalanceAdjust/BalanceAdjustForm.svelte";

    let params = new URLSearchParams($querystring);

    /*const type: string = params.get("type");
    let active;
    let count = 0;*/

    type TabEntry = {
        k: string;
        label: string;
    };
    const key = (tab: TabEntry) => tab.k;

    let tabs: TabEntry[] = [
        {
            k: "expense",
            label: "Expense",
        },
        {
            k: "operation",
            label: "Operation",
        },
        {
            k: "balance",
            label: "Balance Adjust",
        },
    ];
    let active = tabs[0];

    /* $: {
        if (count == 0) {
            count += 1;
        } else {
            push(`/add/?type=${active.toLowerCase()}`);
        }
    }*/
</script>

<div>
    <TabBar {tabs} let:tab {key} bind:active>
        <!-- Note: the `tab` property is required! -->
        <Tab {tab}>
            <Label>{tab.label}</Label>
        </Tab>
    </TabBar>
    {#if active.k == "expense"}
        <h2>New Expense</h2>
        <ExpenseForm />
    {:else if active.k == "operation"}
        <h2>New Operation</h2>
        <OperationForm />
    {:else}
        <h2>New Balance Adjust</h2>
        <BalanceAdjustForm />
    {/if}
</div>
