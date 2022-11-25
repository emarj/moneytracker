<script lang="ts">
    import { useQuery } from "@sveltestack/svelte-query";
    import { getAccountsByEntity } from "../../data";
    import { entityID } from "../../entity";

    const accountsQuery = useQuery(["accounts", $entityID], () =>
        getAccountsByEntity($entityID)
    );
    export let op;
</script>

<form>
    <input type="text" placeholder="Description" value={op.description} />

    {#each op.transactions as t}
        <fieldset>
            {#if $accountsQuery.isLoading}
                <span>Loading...</span>
            {:else if $accountsQuery.error}
                <span
                    >An error has occurred: {$accountsQuery.error.message}</span
                >
            {:else}
                <select>
                    {#each $accountsQuery.data as account}
                        <option value={account.id}
                            >{account.display_name}</option
                        >
                    {/each}
                </select>
                <select>
                    {#each $accountsQuery.data as account}
                        <option value={account.id}
                            >{account.display_name}</option
                        >
                    {/each}
                </select>
            {/if}
            <input type="number" min="0" step="0.01" />
        </fieldset>
    {/each}

    <button
        on:click|preventDefault={() =>
            (op.transactions = [...op.transactions, {}])}
        >Add Transactions</button
    >
</form>

<style>
    form {
        padding: 1rem;
        border-radius: 10px;
        border: 1px solid orange;
    }
</style>
