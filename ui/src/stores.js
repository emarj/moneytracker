import { writable } from 'svelte/store';
import { get } from 'svelte/store';
import { v4 as uuidv4 } from 'uuid';

const getLocalStorageWritable = (key, initialValue) => {
	let value = initialValue;
	const rawValue = localStorage.getItem(key);
	if (rawValue) {
		value = JSON.parse(rawValue)
	} else {
		value = initialValue;
	}
	const store = writable(value);
	store.subscribe((value) => localStorage.setItem(key,JSON.stringify(value)));
	return store
}


export const accounts = getLocalStorageWritable('accounts',[]);

export const getAccountByID = (id)  => {return get(accounts).find(a => a.id === id)}




export const transactions = getLocalStorageWritable('transactions',[]);


export const addTransaction = (t) => {
	t.id = uuidv4()
	transactions.update(tl => [...tl, t]);
}

export const deleteTransaction = (id) => {
	transactions.update(tl => tl.filter(t => t.id !== id));
}

export const users = getLocalStorageWritable('users',[]);

export const login = getLocalStorageWritable('login',{});

//Holds content of TxForm
export const draft = writable({});

