<script lang="ts">
    import { emptyOperation, type Operation } from "../../model";
    import Textfield from "@smui/textfield";
    import DatePicker from "../DatePicker.svelte";
    import AccountSelect from "../AccountSelect.svelte";
    import Button from "@smui/button/src/Button.svelte";
    import AccountBalance from "../Accounts/AccountBalance.svelte";
    import { push } from "svelte-spa-router";
    import { useMutation } from "@sveltestack/svelte-query";
    import { adjustBalance } from "../../api";
    import { JSONPretty } from "../../util/utils";

    const mutation = useMutation((b) => adjustBalance(b), {
        onSuccess: (data: number) => {
            push("/");
        },
    });

    let balance = {
        timestamp: new Date(),
        comment: "",
        account_id: null,
        value: null,
    };
</script>

<div>
    <DatePicker bind:timestamp={balance.timestamp} />
    <Textfield
        variant="outlined"
        bind:value={balance.comment}
        label="Description"
        style="width: 100%;"
    />
    <AccountSelect bind:account_id={balance.account_id} />
    {#key balance.account_id}
        {#if balance.account_id}
            <AccountBalance id={balance.account_id} />
        {/if}
    {/key}

    <Textfield
        variant="outlined"
        bind:value={balance.value}
        placeholder="New balance"
        suffix="€"
        input$pattern={"\\d+(\\.\\d{2})?"}
    />

    <Button
        on:click={() => {
            $mutation.mutate(balance);
        }}>Submit</Button
    >
</div>

<h3>Preview</h3>
<pre>{JSONPretty(balance)}</pre>

<style>
    div > :global(*) {
        /*  display: block; */
        margin: 1rem 0.3rem;
    }

    .buttons > :global(*) {
        margin: 0;
        width: 48%;
    }
</style>
