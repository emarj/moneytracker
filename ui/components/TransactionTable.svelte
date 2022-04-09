<script>
import TransactionRow from "./TransactionRow.svelte";

export let transactions;
</script>




<table>
    {#each transactions as t}
        <TransactionRow transaction={t} on:delete-tx />
    {/each}
</table>

<style>
    table :global(td) {
        padding: 0.5em;
    }
</style>