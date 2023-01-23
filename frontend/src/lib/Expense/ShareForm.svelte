<script type="ts">
    import Textfield from "@smui/textfield";
    import Slider from "@smui/slider";
    import Switch from "@smui/switch";
    import FormField from "@smui/form-field";
    import AccountSelect from "../AccountSelect.svelte";
    import EntitySelect from "../EntitySelect.svelte";
    import { entityID } from "../../store";
    import { Share } from "../../model";
    import DecimalInput from "../DecimalInput.svelte";
    import Button from "@smui/button/src/Button.svelte";

    export let entity_id = null;

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
        disabled={entity_id == null}
        not={entity_id}
        bind:value={share.with_id}
    />
{/if}
<div>
    <DecimalInput
        decimalDigits={3}
        disabled={share.total == null}
        bind:value={share.amount}
    />

    <Button
        disabled={share.total == null}
        variant="outlined"
        on:click={() => (share.quota = 50)}>50%</Button
    >
    <Button
        disabled={share.total == null}
        variant="outlined"
        on:click={() => (share.quota = 100)}>100%</Button
    >
    <div class="level" style={`--level: ${share.quota}%`} />
</div>

<FormField>
    <Switch bind:checked={share.is_credit} color="secondary" icons={false} />
    <span slot="label">Credit</span>
</FormField>
{#key share.is_credit}
    <AccountSelect
        entity_ids={[entity_id]}
        type_id={share.is_credit ? 2 : 0}
        bind:account_id={share.cred_account_id}
        disabled={share.with_id == null}
        label="Credited Account"
    />
    {#key share.with_id}
        <AccountSelect
            entity_ids={[share.with_id]}
            type_id={share.is_credit ? 2 : 0}
            bind:account_id={share.deb_account_id}
            disabled={share.with_id == null}
            label="Debited Account"
        />
    {/key}
{/key}

<style>
    .level {
        --color: rgb(73, 36, 223);
        margin: 1rem;
        width: 90%;
        height: 10px;
        border-radius: 5px;
        border: 2px solid rgb(205, 205, 205);
        background: linear-gradient(
            to right,
            var(--color) var(--level),
            transparent var(--level)
        );
    }
</style>
