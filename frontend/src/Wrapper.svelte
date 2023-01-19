<script>
  import { QueryClient, QueryClientProvider } from "@sveltestack/svelte-query";
  import App from "./App.svelte";

  import "@fontsource/material-icons";
  import "@fontsource/roboto";
  import "@fontsource/roboto-mono";

  const queryClient = new QueryClient();
</script>

<QueryClientProvider client={queryClient}>
  <App />
</QueryClientProvider>
