<script>
    export let value;

    export let currency = "euro";

    let symbols = new Map([
        ["euro", "€"],
        ["dollar", "$"],
    ]);

    $: symbol = symbols.has(currency) ? symbols.get(currency) : currency;
</script>

{value}
{symbol}
