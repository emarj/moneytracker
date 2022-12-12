<script>
    import AccountTag from "../AccountTag.svelte";
    import EntityTag from "../EntityTag.svelte";

    export let eID;
    export let account;
</script>

{#if eID == account.owner.id}
    <AccountTag name={account.name} id={account.id} />
{:else}
    <EntityTag eID={account.owner.id} name={account.owner.name} />
{/if}
