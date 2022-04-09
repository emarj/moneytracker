<script>
        import { createEventDispatcher } from 'svelte';

    let description;
    let date = new Date().toISOString().slice(0,16);
    let notes;
    let amount;
    let shared;
    let paymentMethod;

    const dispatch = createEventDispatcher();


    function submit() {
        let newDate = new Date(date + 'Z');
        let t = {
            //owner  : {id : 'marco', name: 'Marco'},
            date: newDate,
            description: description,
            notes: notes,
            amount: amount,
            //from   : 'asdad',
            //to     : 'asdada',
            //Related []Transaction
            shared: shared,
            //Shares []Share
            paymentMethod: paymentMethod,
        };

        dispatch('add-tx', {
			transaction: t,
		});
    }
</script>

<form method="post" on:submit|preventDefault={submit}>
    <input type="datetime-local" bind:value={date} required>
    <input bind:value={description} placeholder="Description" required/>
    <textarea bind:value={notes} placeholder="Notes"></textarea>
    <input type="number" bind:value={amount} step=".01" required placeholder="0.00">
    <label>Shared: <input type="checkbox" bind:checked="{shared}"></label>
    <button type="reset">Reset</button>
    <button type="submit">Submit</button>
</form>

<style>
    form {
        width: 50%;
    }
    form > * {
        display: block;
        width: 100%;
    }
</style>
