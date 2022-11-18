<script>
    import AccountBalance from "./Accounts/AccountBalance.svelte";

    export let name;
    export let id;
</script>

<span>
    <span class="name">{name}</span>
    <div>
        {name}
        <AccountBalance {id} />
    </div>
</span>

<style lang="scss">
    span {
        position: relative;

        span.name {
            background-color: aquamarine;
            border-radius: 0.5rem;
            padding: 0.1rem 0.3rem;
        }

        div {
            display: none;
            position: absolute;
            top: 2rem;
            left: 2rem;
            z-index: 999;
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
