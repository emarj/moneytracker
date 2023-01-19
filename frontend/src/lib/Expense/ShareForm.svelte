<script type="ts">
    import Textfield from "@smui/textfield";
    import Slider from "@smui/slider";
    import Switch from "@smui/switch";
    import FormField from "@smui/form-field";
    import AccountSelect from "../AccountSelect.svelte";
    import EntitySelect from "../EntitySelect.svelte";
    import { entityID } from "../../store";
    import { Share } from "../../model";

    export let share: Share = new Share();

    let external = false;

    $: if (external) share.with_id = 0; //Share with system
</script>

<FormField>
    <Switch bind:checked={external} icons={false} />
    <span slot="label">External</span>
</FormField>

{#if !external}
    <EntitySelect
        not={$entityID}
        bind:value={share.with_id}
        helperText={"select an entity to share with"}
    />
{/if}
<div>
    <Textfield
        variant="outlined"
        label="Amount"
        type="number"
        min={0}
        max={share.total}
        suffix="€"
        input$pattern={"\\d+(\\.\\d{2})?"}
        disabled={share.total === null}
        bind:value={share.amount}
    />
</div>
<FormField style="width:100%;">
    <Slider
        max={100}
        min={5}
        step={1}
        style="width:100%;"
        discrete
        disabled={share.total === null}
        bind:value={share.quota}
    />
</FormField>

<FormField>
    <Switch bind:checked={share.is_credit} color="secondary" icons={false} />
    <span slot="label">Credit</span>
</FormField>
{#key share.is_credit}
    <AccountSelect
        type_id={share.is_credit ? 1 : 0}
        bind:value={share.cred_account_id}
        label="Credited Account"
    />
    {#key share.with_id}
        <AccountSelect
            type_id={share.is_credit ? 1 : 0}
            bind:value={share.deb_account_id}
            disabled={share.with_id == null}
            label="Debited Account"
        />
    {/key}
{/key}
