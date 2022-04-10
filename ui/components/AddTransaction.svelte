<script>
    import SharesForm from './SharesForm.svelte';
    import {users,addTransaction, login} from '../src/stores'

    let description;
    let date = new Date().toISOString().slice(0,16);
    let notes;
    let amount;
    let ownerID = $login.userID;

    let shares = [];
    $: shared = (shares && shares.length > 0);
    
    let paymentMethod;


    function submit() {
        let newDate = new Date(date + 'Z');
        let t = {
            //owner: {id : 'marco', name: 'Marco'},
            date: newDate,
            description: description,
            notes: notes,
            amount: amount,
            //from   : 'asdad',
            //to     : 'asdada',
            //Related []Transaction
            shared: shared,
            shares: shares,
            paymentMethod: paymentMethod,
        };

        console.log(shares);

        addTransaction(t)
    }
</script>

<form method="post" on:submit|preventDefault={submit}>
    <select bind:value={ownerID}>
        {#each $users as u,i}
            <option value={u.id}>
            {u.name}
            </option>
        {/each}
    </select>
    <input type="datetime-local" bind:value={date} required>
    <input bind:value={description} placeholder="Description" required/>
    <textarea bind:value={notes} placeholder="Notes"></textarea>
    <input type="number" bind:value={amount} step=".01" required placeholder="0.00">

    <SharesForm ownerID={ownerID} bind:shares={shares} amount={amount} />

    <button type="reset">Reset</button>
    <button type="submit">Submit</button>
</form>

<style>
    form {
        width: 50%;
    }
    form > :global(*) {
        display: block;
        width: 100%;
    }
</style>
