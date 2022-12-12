<script lang="ts">
    import AccountBalance from "./AccountBalance.svelte";
    import AccountMovements from "./AccountMovements.svelte";

    import { push } from "svelte-spa-router";
    import { useMutation } from "@sveltestack/svelte-query";
    import { deleteAccount } from "../../api";
    import IconButton from "@smui/icon-button/src/IconButton.svelte";

    export let account;

    let balance;
    let transactions;
    let error;

    const mutation = useMutation((aID: number) => deleteAccount(aID), {
        onSuccess: (data: number) => {
            push("/");
        },
        onError: (data) => {
            error = data;
        },
    });
</script>

<div>
    <h1>{account.display_name}{account.type == 1 ? "*" : ""}</h1>
    type: {account.type}, name: {account.name}
    <AccountBalance id={account.id} bind:this={balance} />
    <div class="movements">
        <AccountMovements id={account.id} bind:this={transactions} />
    </div>
</div>

<IconButton
    class="material-icons"
    touch
    on:click={(event) => {
        $mutation.mutate(account.id);
    }}
    disabled={$mutation.isLoading}>delete</IconButton
>

{#if error}
    {error}
{/if}

<style lang="scss">
    div.card {
        position: relative;
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
        transition: all 0.5s ease-in;

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
