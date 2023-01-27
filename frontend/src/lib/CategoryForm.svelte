<script>
    import Button from "@smui/button";
    import Textfield from "@smui/textfield";
    import { useMutation, useQueryClient } from "@sveltestack/svelte-query";

    import { createEventDispatcher } from "svelte";
    import { addCategory } from "../api";

    const dispatch = createEventDispatcher();

    const queryClient = useQueryClient();

    let emptyCategory = {
        full_name: "",
    };

    let category;

    export let reset = () => {
        category = structuredClone(emptyCategory);
    };

    reset();

    let mutation = useMutation((cat) => addCategory(cat), {
        onSuccess: () => {
            queryClient.invalidateQueries({
                queryKey: ["categories"],
            });
            dispatch("submit");
            reset();
            //TODO: Add toast
        },
    });

    const submitHandler = (e) => {
        e.preventDefault();
        $mutation.mutate(category);
    };

    const resetHandler = (e) => {
        e.preventDefault();
        reset();
    };
</script>

<form>
    <div>
        <Textfield label="Full Name" bind:value={category.full_name} />
    </div>

    <div>
        <Button type="reset" on:click={resetHandler}>Reset</Button>
        <Button type="submit" on:click={submitHandler}>Create</Button>
    </div>

    <!-- <pre>
        {JSONPretty(account)}
    </pre> -->
</form>
