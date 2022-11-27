<script lang="ts">
    import { useQuery } from "@sveltestack/svelte-query";
    import { getAccountsByEntity } from "../../api";
    import { entityID } from "../../store";
    import CircularProgress from "@smui/circular-progress";

    const accountsQuery = useQuery(["accounts", $entityID], () =>
        getAccountsByEntity($entityID)
    );
    export let op;
</script>

<input type="text" placeholder="Description" value={op.description} />

<ul>
    {#each op.transactions as t}
        <li>{t.from.id} -> {t.to.id} : {t.amount}</li>
    {/each}
</ul>

<button
    on:click|preventDefault={() => (op.transactions = [...op.transactions, {}])}
    >Add Transactions</button
>

<style>
    form {
        padding: 1rem;
        border-radius: 10px;
        border: 1px solid orange;
    }
</style>
