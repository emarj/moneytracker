<script type="ts">
    import Textfield from "@smui/textfield";
    import Slider from "@smui/slider";
    import Switch from "@smui/switch";
    import FormField from "@smui/form-field";
    import AccountSelect from "../AccountSelect.svelte";
    import EntitySelect from "../EntitySelect.svelte";
    import { entityID } from "../../store";
    import { Share } from "../../model";

    export let share: Share = new Share(400);
</script>

<EntitySelect
    not={$entityID}
    bind:value={share.with}
    helperText={"select an entity to share with"}
/>
<Textfield
    variant="outlined"
    label="Amount"
    type="number"
    min={0}
    max={share.total}
    step={0.01}
    bind:value={share.amount}
/>
<FormField style="width:100%;">
    <Slider
        max={100}
        min={5}
        step={1}
        style="width:100%;"
        discrete
        bind:value={share.quota}
    />
</FormField>

<!--<p>Internal? <input type="checkbox" bind:checked={external} /></p>-->

<FormField>
    <Switch bind:checked={share.isCredit} color="secondary" icons={false} />
    <span slot="label">Credit</span>
</FormField>
{#key share.isCredit}
    <AccountSelect
        owner_id={$entityID}
        type={share.isCredit ? 1 : 0}
        bind:value={share.credAccount}
        label="Credited Account"
    />
    {#key share.with}
        <AccountSelect
            owner_id={share.with}
            type={share.isCredit ? 1 : 0}
            bind:value={share.debAccount}
            disabled={share.with == null}
            label="Debited Account"
        />
    {/key}
{/key}
