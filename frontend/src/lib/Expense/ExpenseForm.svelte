<script lang="ts">
    import Textfield from "@smui/textfield";
    import Slider from "@smui/slider";
    import Switch from "@smui/switch";
    import FormField from "@smui/form-field";

    import { messageStore } from "../../store";

    import Button, { Label } from "@smui/button";
    import AccountSelect from "../AccountSelect.svelte";
    import EntitySelect from "../EntitySelect.svelte";
    import { entityID, defaultExpense } from "../../store";
    import type { Expense } from "../../model";
    import { ExpenseToOperation } from "../../model";
    import Operation from "../Operation/Operation.svelte";
    import { addOperation } from "../../api";
    import { useMutation } from "@sveltestack/svelte-query";
    import { DateFMT } from "../../util/utils";
    import { push } from "svelte-spa-router";

    const mutation = useMutation((op) => addOperation(op), {
        onSuccess: (data: number) => {
            $messageStore = { text: `Operation added successfully!` };
            push("/");
        },
    });

    const defaultQuota = 50;

    export let e: Expense = structuredClone(defaultExpense);
    let quota: number = defaultQuota;
    let alreadyPaid: boolean = false;

    const reset = (event) => {
        event?.preventDefault();
        e = structuredClone(defaultExpense); //without structuredClone e is just a reference
        quota = defaultQuota;
        alreadyPaid = false;
    };

    $: e.sharedAmount = e.amount ? (e.amount * quota) / 100 : null;

    let op;
    let submitted = false;

    e.timestamp = new Date().toISOString();

    $: op = ExpenseToOperation(e);
</script>

<form>
    <Textfield
        variant="outlined"
        bind:value={e.timestamp}
        label="Datetime"
        style="width: 100%;"
    />
    <Textfield
        variant="outlined"
        bind:value={e.description}
        label="Description"
        style="width: 100%;"
    />
    <AccountSelect owner_id={$entityID} credit={false} bind:value={e.account} />

    <Textfield
        variant="outlined"
        bind:value={e.amount}
        label="Amount"
        suffix="€"
        input$pattern={"\\d+(\\.\\d{2})?"}
    />

    <Textfield
        variant="outlined"
        bind:value={e.category}
        label="Category"
        style="width: 100%;"
    />
    <Textfield
        variant="outlined"
        label="Tags"
        style="width: 100%;"
        value={"asasa"}
    />

    <FormField>
        <Switch bind:checked={e.shared} icons={false} />
        <span slot="label">Shared</span>
    </FormField>
    {#if e.shared}
        <EntitySelect
            not={$entityID}
            bind:value={e.sharedWith}
            helperText="select an entity to share with"
        />
        <Textfield
            variant="outlined"
            label="Amount"
            type="number"
            min={0}
            max={e.amount}
            value={e.sharedAmount}
            on:change={(event) => {
                quota = (event.target.value / e.amount) * 100;
            }}
        />
        <FormField style="width:100%;">
            <Slider
                max={100}
                min={5}
                step={1}
                style="width:100%;"
                discrete
                bind:value={quota}
            />
            <!-- <span slot="label"> Percentage </span> -->
        </FormField>

        <!--<p>Internal? <input type="checkbox" bind:checked={external} /></p>-->

        <FormField>
            <Switch
                bind:checked={alreadyPaid}
                color="secondary"
                icons={false}
                on:change={() => {
                    // e.credAccount = null;
                    //e.debAccount = null;
                }}
            />
            <span slot="label">Already Paid</span>
        </FormField>
        {#key alreadyPaid}
            <AccountSelect
                owner_id={$entityID}
                credit={!alreadyPaid}
                bind:value={e.credAccount}
                label="Credited Account"
                helperText="Select where to receive credit"
            />
            {#key e.sharedWith}
                <AccountSelect
                    owner_id={e.sharedWith}
                    credit={!alreadyPaid}
                    bind:value={e.debAccount}
                    disabled={e.sharedWith === null ||
                        e.sharedWith === undefined}
                    label="Debited Account"
                    helperText="Select where to get credit from"
                />
            {/key}
        {/key}
    {/if}
    <div>
        <Button color="secondary" on:click={reset} variant="raised">
            <Label>Reset</Label>
        </Button>
        <Button
            color="primary"
            on:click={(event) => {
                event.preventDefault();
                $mutation.mutate(op);
            }}
            variant="outlined"
            disabled={$mutation.isLoading}
        >
            <Label>Add</Label>
        </Button>
    </div>
</form>

<style>
    form > :global(*) {
        /*  display: block; */
        margin: 1rem 0.3rem;
    }
</style>
