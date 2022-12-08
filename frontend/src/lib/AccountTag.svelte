<script>
    import AccountName from "./AccountName.svelte";
    import AccountBalance from "./Accounts/AccountBalance.svelte";

    export let account;
    export let eID;
</script>

<span>
    <span class="name" class:world={account.id === 0}
        ><AccountName {account} /></span
    >
</span>

<style lang="scss">
    span {
        position: relative;

        span.name {
            border-radius: 0.5rem;
            padding: 0.1rem 0.3rem;
            background-color: rgb(10, 141, 202);

            &.world {
                background-color: #ddd;
            }
        }

        div {
            display: none;
            position: absolute;
            top: 2rem;
            left: 2rem;
            z-index: 4;
            width: max-content;
            padding: 0.8rem;
            background-color: bisque;

            border-radius: 10px;
        }

        &:hover div {
            display: block;
        }
    }
</style>
