<script>
    import Select, { Option } from "@smui/select";
    import { useQuery } from "@sveltestack/svelte-query";
    import { getCategories } from "../api";
    import CircularProgress from "@smui/circular-progress";
    import Button from "@smui/button/src/Button.svelte";
    import Dialog, { Content, Title } from "@smui/dialog";
    import IconButton from "@smui/icon-button";
    import CategoryForm from "./CategoryForm.svelte";

    export let label = "Category";
    export let value = null;

    const entitiesQuery = useQuery("categories", () => getCategories());

    let openDialog = false;
    let resetForm;
</script>

<div>
    {#if $entitiesQuery.isLoading}
        <CircularProgress style="height: 32px; width: 32px;" indeterminate />
    {:else if $entitiesQuery.error}
        Error: {$entitiesQuery.error?.message}
    {:else if $entitiesQuery.data}
        <Select variant="outlined" bind:value {label}>
            {#each $entitiesQuery.data as cat (cat.id)}
                <Option value={cat.id}
                    >{#if cat.parent_id}{cat.parent.name}/{/if}<strong
                        >{cat.name}</strong
                    ></Option
                >
            {/each}
        </Select>
    {/if}
    <Button on:click={() => (openDialog = true)}>Add New</Button>
    <Dialog
        bind:open={openDialog}
        on:SMUIDialog:closing={resetForm}
        sheet
        aria-labelledby="event-title"
        aria-describedby="event-content"
    >
        <Title id="event-title">New Category</Title>
        <Content id="event-content">
            <IconButton action="close" class="material-icons">close</IconButton>
            <CategoryForm
                on:submit={() => (openDialog = false)}
                bind:reset={resetForm}
            />
        </Content>
    </Dialog>
</div>
