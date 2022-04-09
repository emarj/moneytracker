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
    <td>{transaction.description}</td>
    <td>{transaction.amount}</td>
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