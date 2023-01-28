<script>
    import { userEntitiesID } from "../store";
    export let account;

    const isExternal = !$userEntitiesID?.includes(account.owner_id);
    const isWorld = account.owner_id === 0;
</script>

<span class:account={!isExternal} class:entity={isExternal && !isWorld}>
    {account.owner?.name}:{account.name.toLowerCase()}
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
