<script>
  import { QueryClient, QueryClientProvider } from "@sveltestack/svelte-query";
  import BottomBar from "./lib/BottomBar.svelte";
  import Home from "./Home.svelte";
  import Status from "./lib/Status.svelte";
  import Router from "svelte-spa-router";
  import AddOperation from "./AddOperation.svelte";
  import TopBar from "./lib/TopBar.svelte";

  const queryClient = new QueryClient();

  const routes = {
    // Exact path
    "/": Home,
    "/add": AddOperation,
  };
</script>

<QueryClientProvider client={queryClient}>
  <TopBar />
  <main>
    <Router {routes} />
  </main>
  <BottomBar />
</QueryClientProvider>

<style>
  main {
    margin: 80px 0;
  }
</style>
