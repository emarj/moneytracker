<script lang="ts">
    import Drawer, {
        AppContent,
        Content,
        Header,
        Title,
        Subtitle,
        Scrim,
    } from "@smui/drawer";
    import Button, { Label } from "@smui/button";
    import List, {
        Item,
        Text,
        Graphic,
        Separator,
        Subheader,
    } from "@smui/list";
    import Switch from "@smui/switch";
    import { authStore, showBalances } from "../store";
    import Logout from "../Logout.svelte";
    import EntitySelect from "./EntitySelect.svelte";

    import { entityID } from "../store";

    export let open = false;
</script>

<div class="drawer-container">
    <Drawer variant="modal" fixed={true} bind:open>
        <Header>
            <Title>MoneyTracker</Title>
            <Subtitle />
        </Header>
        <Content>
            <List>
                <Item>
                    Logged in as {$authStore.user.name}
                    <Logout />
                </Item>
                <Item>
                    <Switch bind:checked={$showBalances} /> Show balances
                </Item>
                <Separator />
                <Subheader tag="h3">Entities</Subheader>
                <EntitySelect bind:value={$entityID} style="menu" />
                <Separator />
                <Item on:click={() => location.reload()}>Panic Reload</Item>
            </List>
        </Content>
    </Drawer>

    <!-- Don't include fixed={false} if this is a page wide drawer.
          It adds a style for absolute positioning. -->
    <Scrim fixed={false} />
    <AppContent class="app-content">
        <slot />
    </AppContent>
</div>

<style>
    /* These classes are only needed because the
      drawer is in a container on the page. */
    .drawer-container {
        position: relative;
        display: flex;
        overflow: hidden;
        z-index: 0;
        min-height: 100vh;
    }

    * :global(.app-content) {
        flex: auto;
        overflow: auto;
        position: relative;
        flex-grow: 1;
    }
</style>
