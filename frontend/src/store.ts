import { writable } from 'svelte/store';
import type { Writable } from 'svelte/store';
import type { Entity, Account, Expense, Transaction, Operation } from './model'

export const entityID = writable(1);

export const defaultOperation: Operation = {
    description: "",
    category: "",
    transactions: [],
};

export const newOp: Writable<Operation> = writable(defaultOperation);



export const defaultExpense: Expense = {
    timestamp: "2006-01-02T15:04:05.999Z",
    description: "",
    amount: null,
    account: 1,
    shared: false,
    sharedAmount: null,
    credAccount: null,
    debAccount: null,
    sharedWith: null,
    category: "",
};

export const newExpense: Writable<Expense> = writable(defaultExpense);

