<script>
  import BottomBar from "./lib/BottomBar.svelte";
  import Home from "./pages/Home.svelte";
  import Router from "svelte-spa-router";
  import New from "./pages/New.svelte";
  import TopBar from "./lib/TopBar.svelte";
  import { entityID, authStore, historyStore } from "./store";
  import Login from "./Login.svelte";
  import AllOperations from "./pages/AllOperations.svelte";
  import AccountForm from "./lib/Accounts/AccountForm.svelte";
  import MainMenu from "./lib/MainMenu.svelte";
  import Account from "./pages/Account.svelte";
  import OperationPage from "./pages/OperationPage.svelte";

  const routes = {
    "/": Home,
    "/add": New,
    "/operations": AllOperations,
    "/account/:id": Account,
    "/operation/:id": OperationPage,
    "/newaccount": AccountForm,
  };

  function routeLoaded(event) {
    historyStore.push(event.detail.route);
    window.scrollTo(0, 0);
  }

  let menuOpen = false;
</script>

{#if $authStore}
  <MainMenu open={menuOpen}>
    <TopBar bind:menuOpen />

    <main>
      {#key $entityID}
        <Router {routes} on:routeLoaded={routeLoaded} />
      {/key}
    </main>

    <BottomBar />
  </MainMenu>
{:else}
  <Login />
{/if}

<style>
  main {
    margin: 80px 0;
  }
</style>
