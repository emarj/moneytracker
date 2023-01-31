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
    import DecimalInput from "../DecimalInput.svelte";

    const mutation = useMutation((b: any) => adjustBalance(b), {
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

    <div>
        <AccountSelect bind:account_id={balance.account_id} />
        {#key balance.account_id}
            {#if balance.account_id}
                <AccountBalance
                    label="Current Balance"
                    id={balance.account_id}
                />
            {/if}
        {/key}
    </div>
    <DecimalInput
        bind:value={balance.value}
        decimalDigits={2}
        label="New balance"
    />

    <Textfield
        variant="outlined"
        bind:value={balance.comment}
        label="Comment"
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
        margin: 0.5rem 0.3rem;
    }
</style>
