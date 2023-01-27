<script>
  import BottomBar from "./lib/BottomBar.svelte";
  import Home from "./pages/Home.svelte";
  import Router from "svelte-spa-router";
  import New from "./pages/New.svelte";
  import TopBar from "./lib/TopBar.svelte";
  import { authStore } from "./store";
  import Login from "./Login.svelte";
  import AllOperations from "./pages/AllOperations.svelte";
  import AccountForm from "./lib/Accounts/AccountForm.svelte";
  import MainMenu from "./lib/MainMenu.svelte";
  import Account from "./pages/Account.svelte";
  import OperationPage from "./pages/OperationPage.svelte";
  import Blank from "./pages/Blank.svelte";

  const routes = {
    "/": Home,
    "/add": New,
    "/operations": AllOperations,
    "/account/:id": Account,
    "/operation/:id": OperationPage,
    "/newaccount": AccountForm,
    "/blank": Blank,
  };

  let menuOpen = false;
</script>

{#if $authStore}
  <MainMenu open={menuOpen}>
    <TopBar bind:menuOpen />

    <main>
      <Router
        {routes}
        on:routeLoaded={() => {
          window.scrollTo(0, 0);
        }}
      />
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
