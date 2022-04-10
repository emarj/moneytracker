<script>
    import {users,login} from '../src/stores'

</script>

Logged in as <strong>{$login.userID}</strong>. Change 
<select bind:value={$login.userID}>
    {#each $users as u}
        <option value={u.id}>
        {u.name}
        </option>
    {/each}
</select>
