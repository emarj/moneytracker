<script>
    import AccountName from "./AccountName.svelte";
    import AccountBalance from "./Accounts/AccountBalance.svelte";

    export let account;
    export let eID;

    export let hideExternal = true;

    const isExternal = account.owner_id !== eID;
    const isWorld = account.owner_id === 0;
</script>

<span class:account={!isExternal} class:entity={isExternal && !isWorld}>
    {#if isWorld}
        🌐
    {:else}
        {account.owner?.name}:{account.display_name.toLowerCase()}

        <!-- {:else if !isExternal}
        {account.display_name.toLowerCase()}
    {:else}
        {#if hideExternal}@{/if}{account.owner
            .name}{#if !hideExternal}:{account.name}{/if} -->
    {/if}
</span>

<style lang="scss">
    span {
        font-weight: bold;

        &.account {
            border-radius: 0.5rem;
            padding: 0.1rem 0.3rem;
            background-color: rgb(10, 141, 202);
        }

        &.entity {
            border-radius: 0.5rem;
            padding: 0.1rem 0.3rem;
            background-color: rgb(36, 202, 10);
        }
    }
</style>
