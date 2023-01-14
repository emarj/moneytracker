<script>
    import Select, { Option } from "@smui/select";
    import Button from "@smui/button";
    import Textfield from "@smui/textfield";
    import { useMutation } from "@sveltestack/svelte-query";
    import { addAccount } from "../../api";
    import { pop } from "svelte-spa-router";
    import { entityID } from "../../store";

    let account = {
        name: "",
        display_name: "",
        type_id: 0,
        owner_id: $entityID,
    };

    let mutation = useMutation((a) => addAccount(a), {
        onSuccess: () => {
            pop();
        },
    });

    const handler = (e) => {
        e.preventDefault();
        $mutation.mutate(account);
    };
</script>

<form>
    <Textfield label="Name" bind:value={account.name} />
    <Textfield label="Display Name" bind:value={account.display_name} />
    <Select variant="outlined" bind:value={account.type_id} label="Type">
        <Option value={0}>Money</Option>
        <Option value={1}>Credit</Option>
    </Select>
    <div>
        <Button type="reset">Cancel</Button>
        <Button type="submit" on:click={handler}>Create</Button>
    </div>
</form>
