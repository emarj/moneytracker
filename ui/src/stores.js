import { writable } from 'svelte/store';
import { v4 as uuidv4 } from 'uuid';

const dummyTransactions = [{
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

export const transactions = writable(dummyTransactions);

export const addTransaction = (t) => {
	t.id = uuidv4()
	transactions.update(tl => [...tl, t]);
}

export const deleteTransaction = (id) => {
	transactions.update(tl => tl.filter(t => t.id !== id));
}

export const users = writable([{id: 'marco',name:'Marco'},{id: 'arianna',name:'Arianna'},{id: 'dunno',name:'Dunno'}]);

export const login = writable({userID: 'marco'});

export const draft = writable({});