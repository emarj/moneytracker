<script>
  import { QueryClient, QueryClientProvider } from "@sveltestack/svelte-query";
  import App from "./App.svelte";

  const queryClient = new QueryClient();
</script>

<QueryClientProvider client={queryClient}>
  <App />
</QueryClientProvider>
