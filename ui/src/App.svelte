<script>
	import { v4 as uuidv4 } from 'uuid';
	import AddTransaction from "../components/AddTransaction.svelte";

	import TransactionTable from "../components/TransactionTable.svelte";

	export let transactions;

	$: txs = transactions.filter(t => false)

	function addTx(e) {
		let t = e.detail.transaction;
		t.id = uuidv4()
		transactions = [...transactions, t];
	}

	function deleteTx(e) {
		transactions = transactions.filter(t => t.id !== e.detail.id);
	}
</script>

<main>
	<AddTransaction on:add-tx={addTx}/>
	<TransactionTable {transactions} on:delete-tx={deleteTx} />
</main>

<style>
	main {
		text-align: center;
		padding: 1em;
		max-width: 240px;
		margin: 0 auto;
	}

	@media (min-width: 640px) {
		main {
			max-width: none;
		}
	}
</style>
