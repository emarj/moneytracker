<script lang="ts">
    import AccountBalance from "./AccountBalance.svelte";
    import AccountMovements from "./AccountMovements.svelte";

    import { push } from "svelte-spa-router";
    import { getTypes } from "../../api";
    import IconButton from "@smui/icon-button";

    export let account;

    let balance;
    let transactions;

    function getType(tID: number): string {
        let t = "";
        switch (tID) {
            case 0:
                t = "money";
                break;
            case 2:
                t = "credit";
                break;
            case 3:
                t = "investment";
                break;
        }
        return t;
    }

    function getClass(tID: number): string {
        const c = getType(tID);
        return c == "" ? "" : `type-${c}`;
    }
</script>

<div class="card {getClass(account.type_id)}">
    <h3>
        {#if account.is_default}
            [{account.display_name}]
        {:else}
            {account.display_name}
        {/if}
    </h3>

    <IconButton
        class="btn-more material-icons"
        aria-label="More"
        on:click={() => {
            push(`/account/${account.id}`);
        }}>more_horiz</IconButton
    >

    <IconButton
        class="btn-refresh material-icons"
        aria-label="Refresh"
        on:click={() => {
            balance.refresh();
            transactions.refresh();
        }}>refresh</IconButton
    >
    <AccountBalance id={account.id} bind:this={balance} />
    <div class="movements">
        <AccountMovements id={account.id} bind:this={transactions} />
    </div>
    <IconButton
        class="btn-add material-icons"
        aria-label="Add"
        on:click={() => {
            push(`/add?from=${account.name}`);
        }}>add</IconButton
    >
</div>

<style lang="scss">
    div.card {
        position: relative;
        padding: 1em;
        background: rgb(127, 127, 127);
        /*  background: linear-gradient(
            24deg,
            rgb(196, 56, 196) 0%,
            rgba(253, 187, 45, 1) 100%
        ); */
        height: auto;
        width: 300px;
        border-radius: 15px;
        transition: all 0.5s ease-in;

        &.type-credit {
            background: rgb(231, 138, 9);
        }

        &.type-money {
            background: rgb(3, 163, 81);
        }

        &.type-investment {
            background: rgb(44, 112, 207);
        }

        h3 {
            text-align: center;
            margin-top: 0;
        }

        :global(.btn-more) {
            position: absolute;
            top: 1rem;
            right: 1rem;
        }

        :global(.btn-refresh) {
            position: absolute;
            top: 1rem;
            left: 1rem;
        }

        :global(.btn-add) {
            display: none;
            margin-left: auto;
        }
    }
</style>
