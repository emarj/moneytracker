<script lang="ts">
    import Button from "@smui/button";
    import Tab, { Label } from "@smui/tab";
    import TabBar from "@smui/tab-bar";
    import Textfield from "@smui/textfield";

    export let timestamp;
    let dt;

    const toLocale = (dt: Date): string =>
        new Date(dt.getTime() + new Date().getTimezoneOffset() * -60000)
            .toISOString()
            .slice(0, 16);

    const toDate = (str: string) => new Date(str);

    const setTimeNow = () => {
        timestamp = new Date();
        updateDT();
    };
    const setTimeYesterday = () => {
        timestamp = ((d) => new Date(d.setDate(d.getDate() - 1)))(new Date());
        updateDT();
    };

    const updateDT = () => {
        dt = toLocale(timestamp);
    };

    setTimeNow();
</script>

<input
    value={dt}
    on:change={(event) => {
        timestamp = toDate(event.target.value);
        updateDT();
    }}
    type="datetime-local"
/>
<Button variant="outlined" on:click={() => setTimeNow()}>Now</Button>
<Button variant="outlined" on:click={() => setTimeYesterday()}>Yesteday</Button>

<style>
    input {
        padding: 0.8rem 0.8rem;
        font-size: 1.2rem;
        font-family: inherit;
        border-radius: 4px;
        outline: none;
        border: 1px solid #333;
    }
</style>
