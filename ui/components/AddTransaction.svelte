<script>
    import SharesForm from './SharesForm.svelte';
    import {users,addTransaction, login, accounts} from '../src/stores'
    import {dateToLocalISOLikeString,dateToRFC3339} from '../src/date.js'

    let description;
    let date = dateToLocalISOLikeString(new Date());
    let notes;
    let amount;
    let ownerID = $login.userID;

    let fromID;
    let toID;

    let shares = [];
    $: shared = (shares && shares.length > 0);
    
    let paymentMethod;


    function submit() {
        let t = {
            owner: {id : 'marco', name: 'Marco'},
            date: new Date(date),
            description: description,
            notes: notes,
            amount: amount,
            fromID   : fromID,
            toID     : toID,
            //Related []Transaction
            shared: shared,
            shares: shares,
            paymentMethod: paymentMethod,
        };

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

    <div>
        <select bind:value={fromID}>
            {#each $accounts as a}
                <option value={a.id}>
                {a.displayName}
                </option>
            {/each}
        </select>
        ->
        <select bind:value={toID}>
            {#each $accounts as a}
                <option value={a.id}>
                {a.displayName}
                </option>
            {/each}
        </select>
    </div>

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
