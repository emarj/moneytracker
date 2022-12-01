<script lang="ts">
    import Tab, { Label } from "@smui/tab";
    import TabBar from "@smui/tab-bar";
    import Textfield from "@smui/textfield";

    export let timestamp;
    let dt;

    let active = "Now";

    const toLocale = (dt: Date): string =>
        new Date(dt.getTime() + new Date().getTimezoneOffset() * -60000)
            .toISOString()
            .slice(0, 16);

    const toDate = (str: string) => new Date(str);

    $: if (active == "Now") {
        timestamp = new Date();
        updateDT();
    } else if (active == "Yesterday") {
        timestamp = ((d) => new Date(d.setDate(d.getDate() - 1)))(new Date());
        updateDT();
    }

    const updateDT = () => {
        dt = toLocale(timestamp);
    };
</script>

<TabBar tabs={["Now", "Yesterday"]} let:tab bind:active>
    <Tab {tab}>
        <Label>{tab}</Label>
    </Tab>
</TabBar>
<Textfield
    variant="outlined"
    value={dt}
    on:change={(event) => {
        timestamp = toDate(event.target.value);
        updateDT();
    }}
    label="Datetime"
    type="datetime-local"
/>
