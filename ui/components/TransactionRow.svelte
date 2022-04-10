<script>
    import { createEventDispatcher } from 'svelte';
    import {users, accounts, getAccountByID} from '../src/stores';

    export let transaction;

    const dispatch = createEventDispatcher();

    function deleteTx(e) {
		dispatch('delete-tx', transaction);
	}

    function editTx(e) {
		dispatch('edit-tx', transaction);
	}


</script>

<tr>
    <td>{transaction.date > 0 ? transaction.date.toISOString(): transaction.date}</td>
    <td>{transaction.owner.id}</td>
    <td contenteditable="true" bind:textContent={transaction.description}></td>
    <td>{transaction.amount}</td>
    <td>{getAccountByID(transaction.fromID).displayName} -> {getAccountByID(transaction.toID).displayName}</td>
    <td>
        {#if transaction.shared}
            {transaction.shares.reduce((sum,s) => sum + s.quota,0)}
        {/if}
    </td>
    <td>{transaction.paymentMethod}</td>
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