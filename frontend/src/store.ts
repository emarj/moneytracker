import { writable } from 'svelte/store';
import type { Writable } from 'svelte/store';
import { type Entity, type Account, Expense, type Transaction, type Operation } from './model'
import { emptyOperation } from './model'

export const entityID = writable(1);


export const newOp: Writable<Operation> = writable(emptyOperation);


export const newExpense: Writable<Expense> = writable(new Expense());


export const messageStore = writable(null);

