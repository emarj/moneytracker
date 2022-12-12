<script lang="ts">
    import Kitchen from "@smui/snackbar/kitchen";
    import BottomAppBar, { Section } from "@smui-extra/bottom-app-bar";
    import IconButton from "@smui/icon-button";
    import Fab, { Icon } from "@smui/fab";
    import { push } from "svelte-spa-router";
    import { messageStore } from "../store";

    let kitchen;

    $: {
        if ($messageStore) {
            kitchen?.push({ label: $messageStore?.text, dismissButton: true });
            $messageStore = null;
        }
    }
</script>

<Kitchen bind:this={kitchen} dismiss$class="material-icons" />
<nav>
    <BottomAppBar variant="static" color="secondary">
        <Section>
            <IconButton class="material-icons">menu</IconButton>
        </Section>
        <Section>
            <Fab
                aria-label="New"
                color="secondary"
                on:click={() => push("/add")}
            >
                <Icon class="material-icons">add</Icon>
            </Fab>
        </Section>
        <Section>
            <IconButton class="material-icons" aria-label="Search"
                >search</IconButton
            >
            <IconButton class="material-icons" aria-label="More"
                >more_vert</IconButton
            >
        </Section>
    </BottomAppBar>
</nav>

<style>
    nav {
        position: fixed;
        bottom: 0;
        left: 0;
        width: 100%;
    }
</style>
