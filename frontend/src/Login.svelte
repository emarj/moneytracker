<script>
    import Button from "@smui/button";
    import Textfield from "@smui/textfield";
    import { useMutation } from "@sveltestack/svelte-query";
    import { login } from "./api";
    import { authStore } from "./store";

    let l = { user: "", password: "" };

    const mutation = useMutation((l) => login(l), {
        onSuccess: (data) => {
            authStore.set(data);
        },
    });
</script>

<div>
    <h1>Login</h1>
    <form>
        <Textfield
            variant="outlined"
            bind:value={l.user}
            label="User"
            autocomplete="username"
        />
        <Textfield
            variant="outlined"
            bind:value={l.password}
            label="Password"
            type="password"
            autocomplete="current-password"
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
