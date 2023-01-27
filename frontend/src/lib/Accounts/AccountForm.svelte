<script>
    import Select, { Option } from "@smui/select";
    import Button from "@smui/button";
    import Textfield from "@smui/textfield";
    import {
        useMutation,
        useQuery,
        useQueryClient,
    } from "@sveltestack/svelte-query";
    import { addAccount, getTypes } from "../../api";
    import { pop } from "svelte-spa-router";
    import { capitalize, JSONPretty } from "../../util/utils";
    import Switch from "@smui/switch";
    import EntitySelect from "../EntitySelect.svelte";
    import { userEntitiesID } from "../../store";
    import { createEventDispatcher } from "svelte";

    const dispatch = createEventDispatcher();

    const typesQuery = useQuery(["types"], () => getTypes());

    const queryClient = useQueryClient();

    export let defaultEntityID = null;

    let emptyAccount = {
        name: "",
        display_name: "",
        type_id: 0,
        is_default: false,
        owner_id: defaultEntityID,
    };

    let account;

    export let reset = () => {
        account = structuredClone(emptyAccount);
    };

    reset();

    let mutation = useMutation((a) => addAccount(a), {
        onSuccess: () => {
            queryClient.invalidateQueries({
                queryKey: ["accounts", "entity", account.owner_id],
            });
            dispatch("submit");
            reset();
            //TODO: Add toast
        },
    });

    const submitHandler = (e) => {
        e.preventDefault();
        $mutation.mutate(account);
    };

    const resetHandler = (e) => {
        e.preventDefault();
        reset();
    };

    let linked = true;

    $: {
        if (linked) {
            account.name = account.display_name
                .toLowerCase()
                .replace(/ /g, "-")
                .replace(/[^\w-]+/g, "");
        }
    }
</script>

<form>
    <div>
        <Textfield label="Display Name" bind:value={account.display_name} />
        <div
            contenteditable
            bind:innerHTML={account.name}
            on:input={() => (linked = false)}
        />
        Linked <Switch bind:checked={linked} />
    </div>
    <br />
    <br />

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
    <EntitySelect entities={$userEntitiesID} bind:value={account.owner_id} />
    Default: <Switch bind:checked={account.is_default} />

    <div>
        <Button type="reset" on:click={resetHandler}>Reset</Button>
        <Button type="submit" on:click={submitHandler}>Create</Button>
    </div>

    <!-- <pre>
        {JSONPretty(account)}
    </pre> -->
</form>
