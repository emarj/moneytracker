<script>
    import Button from "@smui/button/src/Button.svelte";
    import { QueryClient, useMutation } from "@sveltestack/svelte-query";
    import { logout } from "./api";
    import { authStore } from "./store";

    const queryClient = new QueryClient();

    const mutation = useMutation(() => logout(), {
        onSuccess: (data) => {
            console.log("logged out successfully");
            authStore.set(null);
            window.location.replace("/");
        },
    });
</script>

<Button on:click={() => $mutation.mutate()}>Logout</Button>
