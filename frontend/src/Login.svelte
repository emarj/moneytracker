<script>
    import Button from "@smui/button";
    import Textfield from "@smui/textfield";
    import { QueryClient, useMutation } from "@sveltestack/svelte-query";
    import { login } from "./api";

    let l = { user: "", password: "" };

    const queryClient = new QueryClient();

    const mutation = useMutation((l) => login(l), {
        onSuccess: (data) => {
            //$messageStore = { text: `Operation added successfully!` };
            //push("/");
            setTimeout(() => queryClient.invalidateQueries(["login"]), 2000);
        },
    });
</script>

<div>
    <h1>Login</h1>
    <form>
        <Textfield variant="outlined" bind:value={l.user} label="User" />
        <Textfield
            variant="outlined"
            bind:value={l.password}
            label="Password"
            type="password"
        />
        <div>
            <Button type="reset" variant="outlined">Cancel</Button>
            <Button
                type="submit"
                variant="outlined"
                on:click={(e) => {
                    e.preventDefault();
                    $mutation.mutate(l);
                }}>Login</Button
            >
        </div>
    </form>
    {#if $mutation.isLoading}
        Loading
    {:else if $mutation.isError}
        Error
    {/if}
</div>

<style>
    div > :global(*) {
        display: block;
    }
</style>
