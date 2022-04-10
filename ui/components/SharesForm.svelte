<script>
    import { users } from "../src/stores";
    import { v4 as uuidv4 } from 'uuid';

    export let ownerID;
    export let shares = [];
    export let amount;

    //$: shares.forEach(s => { s.amount = 0  });
    
    //const sharedWith = (uID) => shares.findIndex(s => s.with && s.with.id === uID) > 0;

    $: otherUsers = $users.filter(u => u.id !== ownerID);

    const addShare = () => {
        shares = [...shares, { id: uuidv4(), quota: amount/2 }]; 
        //shares.push({ id: uuidv4(), quota: 50 });// WE NEED THIS TO UPDATE PARENT
        //shares = shares;
    };

    const resetShares = () => {shares = [];};
</script>

<div>
    <label>Shared: <input type="checkbox" checked={shares && shares.length > 0} disabled /></label>
    <button on:click|preventDefault={addShare} disabled={!amount || amount === 0 || shares.length === otherUsers.length}>Add Share</button>
    <button on:click|preventDefault={resetShares} disabled={shares.length === 0}>Reset</button>
    <ul>
        {#each shares as share}
            <li>
                <select bind:value={share.with}>
                    {#each otherUsers as u}
                        <option value={u}>
                            {u.name}
                        </option>
                    {/each}
                </select>
                <input
                    type="range"
                    min="0"
                    max={amount}
                    step="0.01"
                    bind:value={share.quota}
                />
                <input type="number"
                min="0"
                max={amount}
                step="0.01"
                bind:value={share.quota}>/{amount} ~ {share.quota * 100/amount}%
            </li>
        {/each}
    </ul>
</div>

<style>
    ul {
        list-style: none;
        padding: 0;
    }

    li {
        border: 1px solid #333;
    }
</style>
