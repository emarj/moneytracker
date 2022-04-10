import { v4 as uuidv4 } from 'uuid';
import {accounts,transactions,users,login} from './stores'

const mockAccounts = [
	{
		id: uuidv4(),
		owners:[{id: 'marco',name:'Marco'}],
		name: "main",
		displayName: 'Main',
		description: '',
		default: false,
	},
	{
		id: uuidv4(),
		owners:[{id: 'marco',name:'Marco'}],
		name: "secondary",
		displayName: 'Secondary',
		description: '',
		default: false,
	}
];

const mockTransactions = [{
	id : uuidv4(),
	owner  : {id : 'marco', name: 'Marco'},
	date : new Date(),
	description: 'spesa 1',
	notes :       'sdada',
	amount :     23,
	fromID   : mockAccounts[0].id,
	toID     : mockAccounts[1].id,
	//Related []Transaction
	shared : true,
	shares : [{id:uuidv4(),with: {id : 'arianna', name: 'Arianna'},quota: 12.89 }],
	paymentMethod: 'contanti',
},{
	id : uuidv4(),
	owner  : {id : 'marco', name: 'Marco'},
	date : new Date(),
	description: 'spesa 2',
	notes :       '',
	amount :     67.09,
	fromID   : mockAccounts[0].id,
	toID     : mockAccounts[1].id,
	//Related []Transaction
	shared : false,
	//Shares []Share
	paymentMethod: 'contanti',
},{
	id : uuidv4(),
	owner  : {id : 'marco', name: 'Marco'},
	date : new Date(),
	description: 'spesa 3',
	notes :       '',
	amount :     124.54,
	fromID   : mockAccounts[0].id,
	toID     : mockAccounts[1].id,
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
	fromID   : mockAccounts[0].id,
	toID     : mockAccounts[1].id,
	//Related []Transaction
	shared : false,
	//Shares []Share
	paymentMethod: 'contanti',
}]

const mockUsers = [
    {id: 'marco',name:'Marco'},
    {id: 'arianna',name:'Arianna'},
    {id: 'dunno',name:'Dunno'}
];

const mockLogin = {userID: 'marco'};


export const populate = () => {
    users.set(mockUsers);
    login.set(mockLogin);
    accounts.set(mockAccounts);
    transactions.set(mockTransactions);
}
