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
    import EntitySwitcher from "./EntitySwitcher.svelte";
    import Logout from "../Logout.svelte";

    export let open = false;
    let active = "Inbox";

    function setActive(value: string) {
        active = value;
        open = false;
    }
</script>

<div class="drawer-container">
    <!-- Don't include fixed={false} if this is a page wide drawer.
          It adds a style for absolute positioning. -->
    <Drawer variant="modal" fixed={false} bind:open>
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
                <Item>
                    <EntitySwitcher style="simple" />
                </Item>
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

    .main-content {
        overflow: auto;
        padding: 16px;
        height: 100%;
        box-sizing: border-box;
    }
</style>
