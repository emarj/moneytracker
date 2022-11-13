<script>
    export let transactions;
</script>

{#each transactions as t}
    {Object.entries(t)}
{/each}
