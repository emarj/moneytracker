<script>
    import { createEventDispatcher } from 'svelte';

    export let transaction

    const dispatch = createEventDispatcher();

    function deleteTx(e) {
		dispatch('delete-tx', transaction);
	}

    function editTx(e) {
		dispatch('edit-tx', transaction);
	}
</script>

<tr>
    <td>{transaction.date.toISOString()}</td>
    <td contenteditable="true" bind:textContent={transaction.description}></td>
    <td>{transaction.amount}</td>
    <td>{#if transaction.shared}
        {transaction.shares.reduce((sum,s) => sum + s.quota,0)}
    {/if}
       </td>
    <td>
        <button on:click={editTx}>Edit</button>
        <button on:click={deleteTx}>Delete</button>
    </td>
</tr>

<style>
    button {
        cursor: pointer;
    }
</style>