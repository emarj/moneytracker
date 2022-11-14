<script>
    import AccountBalance from "./AccountBalance.svelte";
    import AccountTransactions from "./AccountTransactions.svelte";

    export let account;
</script>

<div>
    <h3>{account.display_name} (id = {account.id})</h3>

    <AccountBalance id={account.id} />
    <hr />
    <AccountTransactions id={account.id} />
    <hr />
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
