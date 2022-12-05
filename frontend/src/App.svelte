<script>
  import BottomBar from "./lib/BottomBar.svelte";
  import Home from "./pages/Home.svelte";
  import Router from "svelte-spa-router";
  import New from "./pages/New.svelte";
  import TopBar from "./lib/TopBar.svelte";
  import EntitySwitcher from "./lib/EntitySwitcher.svelte";
  import { entityID, authStore, historyStore, isFirstPage } from "./store";
  import Blank from "./pages/Blank.svelte";
  import Login from "./Login.svelte";
  import Logout from "./Logout.svelte";
  import AllOperations from "./AllOperations.svelte";

  const routes = {
    "/": Home,
    "/add": New,
    "/operations": AllOperations,
    "/blank": Blank,
  };

  function routeLoaded(event) {
    historyStore.push(event.detail.route);
  }
</script>

{#if $authStore}
  <TopBar showBackButton={!$isFirstPage} />

  <main>
    {#key $entityID}
      <Logout />
      <EntitySwitcher />
      <Router {routes} on:routeLoaded={routeLoaded} />
    {/key}
  </main>

  <BottomBar />
{:else}
  <Login />
{/if}

<style>
  main {
    margin: 80px 0;
  }
</style>
