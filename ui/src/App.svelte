<script>
	import AddTransaction from "../components/AddTransaction.svelte";
import Login from "../components/Login.svelte";
	import TransactionTable from "../components/TransactionTable.svelte";
	import {populate} from './data';

</script>

<main>
	<button on:click={populate}>Populate</button>
	<Login />
	<AddTransaction/>
	<TransactionTable/>
</main>

<style>
</style>
