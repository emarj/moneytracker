import App from './App.svelte';
import { v4 as uuidv4 } from 'uuid';

const transactions = [{
	id : uuidv4(),
	owner  : {id : 'marco', name: 'Marco'},
	date : new Date(),
	description: 'spesa 1',
	notes :       'sdada',
	amount :     23,
	from   : 'asdad',
	to     : 'asdada',
	//Related []Transaction
	shared : false,
	//Shares []Share
	paymentMethod: 'contanti',
},{
	id : uuidv4(),
	owner  : {id : 'marco', name: 'Marco'},
	date : new Date(),
	description: 'spesa 2',
	notes :       'sdada',
	amount :     23,
	from   : 'asdad',
	to     : 'asdada',
	//Related []Transaction
	shared : false,
	//Shares []Share
	paymentMethod: 'contanti',
},{
	id : uuidv4(),
	owner  : {id : 'marco', name: 'Marco'},
	date : new Date(),
	description: 'spesa 3',
	notes :       'sdada',
	amount :     23,
	from   : 'asdad',
	to     : 'asdada',
	//Related []Transaction
	shared : false,
	//Shares []Share
	paymentMethod: 'contanti',
},{
	id : uuidv4(),
	owner  : {id : 'marco', name: 'Marco'},
	date : new Date(),
	description: 'spesa 4',
	notes :       'sdada',
	amount :     23,
	from   : 'asdad',
	to     : 'asdada',
	//Related []Transaction
	shared : false,
	//Shares []Share
	paymentMethod: 'contanti',
}]

const app = new App({
	target: document.body,
	props: {
		transactions: transactions,
	}
});

export default app;