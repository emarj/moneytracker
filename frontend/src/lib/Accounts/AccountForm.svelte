<script>
    import Select, { Option } from "@smui/select";
    import Button from "@smui/button";
    import Textfield from "@smui/textfield";
    import { useMutation, useQuery } from "@sveltestack/svelte-query";
    import { addAccount, getTypes } from "../../api";
    import { pop } from "svelte-spa-router";
    import { entityID } from "../../store";
    import { capitalize } from "../../util/utils";

    const typesQuery = useQuery(["types"], () => getTypes());

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

    {#if $typesQuery.isLoading}
        ...
    {:else if $typesQuery.error}
        error
    {:else}
        <Select variant="outlined" bind:value={account.type_id} label="Type">
            {#each $typesQuery.data.account as t}
                {#if !t.system}
                    <Option value={t.id}>{capitalize(t.name)}</Option>
                {/if}
            {/each}
        </Select>
    {/if}

    <div>
        <Button type="reset">Cancel</Button>
        <Button type="submit" on:click={handler}>Create</Button>
    </div>
</form>
