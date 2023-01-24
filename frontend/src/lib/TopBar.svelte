<script>
    import { pop } from "svelte-spa-router";
    import { useQueryClient } from "@sveltestack/svelte-query";
    import { showBalances } from "../store";
    import Switch from "@smui/switch";
    import IconButton from "@smui/icon-button";
    const queryClient = useQueryClient();

    export let menuOpen = false;
</script>

<nav>
    <button
        class="back"
        on:click={() => {
            pop();
        }}
        >← Back
    </button>

    <a href="/#/">
        <img class="logo" src="/moneytracker.png" alt="MoneyTracker" />
    </a>

    <IconButton
        class="material-icons"
        on:click={() => {
            menuOpen = !menuOpen;
        }}
        touch>menu</IconButton
    >
</nav>

<style lang="scss">
    nav {
        position: fixed;
        z-index: 4;
        top: 0;
        left: 0;
        width: 100%;
        height: 60px;
        box-shadow: 0px 1px 15px 0px rgba(0, 0, 0, 0.75);
        background-color: whitesmoke;
        display: flex;
        justify-content: space-between;

        & > a {
            margin-left: auto;
            margin-right: auto;
        }

        & > * {
            display: block;
            height: 100%;
        }

        & > a > img {
            display: block;
            height: 100%;
        }

        button.back {
            font-size: 1.1rem;
            background-color: transparent;
        }
    }
</style>
