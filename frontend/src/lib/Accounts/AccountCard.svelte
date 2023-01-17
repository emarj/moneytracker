<script lang="ts">
    import AccountBalance from "./AccountBalance.svelte";
    import AccountMovements from "./AccountMovements.svelte";

    import { push } from "svelte-spa-router";
    import { getTypes } from "../../api";

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
    <h3>{account.display_name}</h3>
    <button
        class="more"
        on:click={() => {
            push(`/account/${account.id}`);
        }}
    >
        <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            width="24"
            height="24"
            ><path
                fill-rule="evenodd"
                d="M6 12a2 2 0 11-4 0 2 2 0 014 0zm8 0a2 2 0 11-4 0 2 2 0 014 0zm6 2a2 2 0 100-4 2 2 0 000 4z"
            /></svg
        >
    </button>
    <button
        class="refresh"
        on:click={() => {
            balance.refresh();
            transactions.refresh();
        }}
    >
        <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            height="24"
            width="24"
        >
            <path
                xmlns="http://www.w3.org/2000/svg"
                d="M12.7929 2.29289C13.1834 1.90237 13.8166 1.90237 14.2071 2.29289L17.2071 5.29289C17.5976 5.68342 17.5976 6.31658 17.2071 6.70711L14.2071 9.70711C13.8166 10.0976 13.1834 10.0976 12.7929 9.70711C12.4024 9.31658 12.4024 8.68342 12.7929 8.29289L14.0858 7H12.5C8.95228 7 6 9.95228 6 13.5C6 17.0477 8.95228 20 12.5 20C16.0477 20 19 17.0477 19 13.5C19 12.9477 19.4477 12.5 20 12.5C20.5523 12.5 21 12.9477 21 13.5C21 18.1523 17.1523 22 12.5 22C7.84772 22 4 18.1523 4 13.5C4 8.84772 7.84772 5 12.5 5H14.0858L12.7929 3.70711C12.4024 3.31658 12.4024 2.68342 12.7929 2.29289Z"
                fill="#0D0D0D"
            />
        </svg>
    </button>
    <AccountBalance id={account.id} bind:this={balance} />
    <div class="movements">
        <AccountMovements id={account.id} bind:this={transactions} />
    </div>
    <button
        class="add"
        on:click={() => {
            push(`/add?from=${account.name}`);
        }}
        ><svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            width="24"
            height="24"
            ><path
                d="M12.75 7.75a.75.75 0 00-1.5 0v3.5h-3.5a.75.75 0 000 1.5h3.5v3.5a.75.75 0 001.5 0v-3.5h3.5a.75.75 0 000-1.5h-3.5v-3.5z"
            /><path
                fill-rule="evenodd"
                d="M12 1C5.925 1 1 5.925 1 12s4.925 11 11 11 11-4.925 11-11S18.075 1 12 1zM2.5 12a9.5 9.5 0 1119 0 9.5 9.5 0 01-19 0z"
            /></svg
        ></button
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

        button {
            padding: 0;
            width: 24px;
            height: 24px;
            background-color: transparent;
            display: block;
        }

        button.more {
            position: absolute;
            top: 1rem;
            right: 1rem;
        }

        button.refresh {
            position: absolute;
            top: 1rem;
            left: 1rem;
        }

        button.add {
            display: none;
            margin-left: auto;
        }
    }
</style>
