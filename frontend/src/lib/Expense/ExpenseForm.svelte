<script lang="ts">
    import AccountSelect from "../AccountSelect.svelte";
    import EntitySelect from "../EntitySelect.svelte";
    import { entityID } from "../../entity";

    type Expense = {
        amount: number;
        description: string;
        shared: boolean;
        sharedAmount: number;
        account: number;
        credAccount: number;
        debAccount: number;
        sharedWith: number;
    };

    const defaultExpense: Expense = {
        amount: null,
        description: "",
        shared: true,
        sharedAmount: null,
        account: 1,
        credAccount: null,
        debAccount: null,
        sharedWith: null,
    };

    const defaultQuota = 50;

    let e: Expense;
    let quota: number;
    let alreadyPaid: boolean;

    const reset = () => {
        e = structuredClone(defaultExpense); //without structuredClone e is just a reference
        quota = defaultQuota;
        alreadyPaid = false;
    };

    reset();

    $: e.sharedAmount = e.amount ? (e.amount * quota) / 100 : null;
</script>

<form>
    <input type="datetime-local" />
    <input type="description" placeholder="Description" />
    <AccountSelect owner_id={$entityID} credit={false} bind:value={e.account} />
    <input
        type="number"
        placeholder="Amount"
        step="0.01"
        bind:value={e.amount}
    />

    <input type="text" placeholder="Category" />
    <input type="text" placeholder="Tags" />

    <p>Shared?<input type="checkbox" bind:checked={e.shared} /></p>
    {#if e.shared}
        <input
            type="number"
            placeholder="Amount"
            step="0.01"
            max={e.amount}
            value={e.sharedAmount}
            on:change={(event) => {
                quota = (event.target.value / e.amount) * 100;
            }}
            required
        />
        <input
            type="range"
            placeholder="Percentage"
            min="5"
            max="100"
            bind:value={quota}
        />
        {quota}

        <!--<p>Internal? <input type="checkbox" bind:checked={external} /></p>-->
        <div>
            <label>
                Already Paid
                <input
                    type="radio"
                    name="paid"
                    bind:group={alreadyPaid}
                    value={true}
                /></label
            >
            <label>
                Credit
                <input
                    type="radio"
                    name="paid"
                    bind:group={alreadyPaid}
                    value={false}
                />
            </label>
        </div>

        <AccountSelect
            owner_id={$entityID}
            credit={!alreadyPaid}
            bind:value={e.credAccount}
        />
        <EntitySelect not={$entityID} bind:value={e.sharedWith} />
        <AccountSelect
            owner_id={e.sharedWith}
            credit={!alreadyPaid}
            bind:value={e.debAccount}
        />
    {/if}
    <button type="reset" on:click|preventDefault={reset}>Reset</button>
    <button
        on:click|preventDefault={() => {
            console.log(e);
        }}>Send</button
    >
</form>

<style>
    form > input {
        display: block;
        width: 100%;
        margin: 1rem;
        height: 2rem;
        border-radius: 6px;
        border: 1px solid black;
    }

    input:invalid {
        border-color: orange;
    }
</style>
