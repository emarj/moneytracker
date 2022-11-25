<script>
  import { QueryClient, QueryClientProvider } from "@sveltestack/svelte-query";
  import BottomBar from "./lib/BottomBar.svelte";
  import Home from "./Home.svelte";
  import Router from "svelte-spa-router";
  import AddOperation from "./AddOperation.svelte";
  import TopBar from "./lib/TopBar.svelte";
  import EntitySwitcher from "./lib/EntitySwitcher.svelte";
  import { entityID } from "./entity";

  const queryClient = new QueryClient();

  const routes = {
    "/": Home,
    "/add": AddOperation,
  };
</script>

<QueryClientProvider client={queryClient}>
  <TopBar />

  <main>
    {#key $entityID}
      <EntitySwitcher />
      <Router {routes} />
    {/key}
  </main>

  <BottomBar />
</QueryClientProvider>

<style>
  main {
    margin: 80px 0;
  }
</style>
