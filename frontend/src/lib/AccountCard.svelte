<script>
    import AccountBalance from "./AccountBalance.svelte";

    export let account;
</script>

<div>
    <h3>{account.display_name}</h3>
    {account.name}

    <AccountBalance id={account.id} />
</div>

<style>
    div {
        padding: 1em;
        background: rgb(34, 193, 195);
        background: linear-gradient(
            24deg,
            rgba(34, 193, 195, 1) 0%,
            rgba(253, 187, 45, 1) 100%
        );
        height: auto;
        width: 300px;
        border-radius: 15px;
    }
</style>
