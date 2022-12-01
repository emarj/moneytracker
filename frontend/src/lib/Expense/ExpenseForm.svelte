<script lang="ts">
    import Textfield from "@smui/textfield";
    import Slider from "@smui/slider";
    import Switch from "@smui/switch";
    import FormField from "@smui/form-field";

    import { messageStore } from "../../store";

    import Button, { Label } from "@smui/button";
    import AccountSelect from "../AccountSelect.svelte";
    import EntitySelect from "../EntitySelect.svelte";
    import { entityID } from "../../store";
    import { Share, Expense } from "../../model";
    import Operation from "../Operation/Operation.svelte";
    import { addOperation } from "../../api";
    import { useMutation } from "@sveltestack/svelte-query";
    import { push } from "svelte-spa-router";
    import TagInput from "../TagInput.svelte";
    import ShareForm from "./ShareForm.svelte";
    import New from "../../New.svelte";
    import DatePicker from "../DatePicker.svelte";

    const mutation = useMutation((op) => addOperation(op), {
        onSuccess: (data: number) => {
            $messageStore = { text: `Operation added successfully!` };
            push("/");
        },
    });

    export let e: Expense;

    const init = () => {
        e = new Expense();
    };

    const reset = (event) => {
        event.preventDefault();
        init();
    };

    init();

    let op;

    $: op = e.toOperation();
</script>

<form>
    <DatePicker bind:timestamp={e.timestamp} />
    <Textfield
        variant="outlined"
        bind:value={e.description}
        label="Description"
        style="width: 100%;"
    />
    <AccountSelect owner_id={$entityID} credit={false} bind:value={e.account} />

    <Textfield
        variant="outlined"
        bind:value={e.amount}
        label="Amount"
        suffix="€"
        input$pattern={"\\d+(\\.\\d{2})?"}
    />

    <Textfield
        variant="outlined"
        bind:value={e.category}
        label="Category"
        style="width: 100%;"
    />
    <TagInput
        existing={[
            { id: 1, name: "tag1" },
            { id: 3, name: "tag3" },
        ]}
        bind:tags={e.tags}
    />

    <FormField>
        <Switch bind:checked={e.isShared} icons={false} />
        <span slot="label">Shared</span>
    </FormField>
    {#if e.isShared}
        {#each e.shares as s}
            <ShareForm bind:share={s} />
        {/each}
    {/if}
    <div class="buttons">
        <Button color="secondary" on:click={reset} variant="raised">
            <Label>Reset</Label>
        </Button>
        <Button
            color="primary"
            on:click={(event) => {
                event.preventDefault();
                $mutation.mutate(op);
            }}
            variant="outlined"
            disabled={$mutation.isLoading}
        >
            <Label>Add</Label>
        </Button>
    </div>
</form>

<h3>Preview</h3>
<pre>
    {JSON.stringify(e, null, 4)}
</pre>

<style>
    form > :global(*) {
        /*  display: block; */
        margin: 1rem 0.3rem;
    }

    .buttons > :global(*) {
        width: 48%;
    }
</style>
