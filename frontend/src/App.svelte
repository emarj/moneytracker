<script>
  import { QueryClient, QueryClientProvider } from "@sveltestack/svelte-query";
  import Accounts from "./lib/Accounts.svelte";

  const queryClient = new QueryClient();
</script>

<main>
  <QueryClientProvider client={queryClient}>
    <Accounts />
  </QueryClientProvider>
</main>
