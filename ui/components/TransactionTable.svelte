<script>
import TransactionRow from "./TransactionRow.svelte";
import {transactions,deleteTransaction} from '../src/stores'
</script>




<table>
    {#each $transactions as t}
        <TransactionRow transaction={t} on:delete-tx={e =>  deleteTransaction(e.detail.id)} />
    {/each}
</table>

<style>
    table :global(td) {
        padding: 0.5em;
    }
</style>