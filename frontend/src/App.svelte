<script>
  import {
    QueryClient,
    QueryClientProvider,
    useQuery,
  } from "@sveltestack/svelte-query";
  import BottomBar from "./lib/BottomBar.svelte";
  import Home from "./Home.svelte";
  import Router from "svelte-spa-router";
  import New from "./New.svelte";
  import TopBar from "./lib/TopBar.svelte";
  import EntitySwitcher from "./lib/EntitySwitcher.svelte";
  import { entityID } from "./store";
  import Blank from "./Blank.svelte";
  import Login from "./Login.svelte";
  import { echo } from "./api";
  import Logout from "./Logout.svelte";

  const loginStore = useQuery(["login"], () => echo(), {
    onSuccess: () => {
      console.log("login succeded");
    },
    onError: (err) => {
      console.log("echo error", err);
    },
  });

  const routes = {
    "/": Home,
    "/add": New,
    "/blank": Blank,
  };
</script>

{#if $loginStore.isSuccess}
  <TopBar />

  <main>
    {#key $entityID}
      <Logout />
      <EntitySwitcher />
      <Router {routes} />
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
