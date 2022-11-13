<script>
  import { QueryClient, QueryClientProvider } from "@sveltestack/svelte-query";
  import Accounts from "./lib/Accounts.svelte";
  import Transactions from "./lib/Transactions.svelte";

  const queryClient = new QueryClient();
</script>

<main>
  <QueryClientProvider client={queryClient}>
    <Transactions />
    <Accounts />
  </QueryClientProvider>
</main>
