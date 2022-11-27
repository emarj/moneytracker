


export type Entity = {
    id?: number;
    name?: string;
}

export type Account = {
    id: number;
    owner?: Entity;
    is_credit?: boolean;
}

export type Transaction = {
    id?: number;
    amount?: number;
    from?: Account;
    to?: Account;
    operation?: Operation;
};

export type Operation = {
    id?: number;
    description: string;
    timestamp?: string;
    category?: string;
    transactions: Transaction[];
};



export type Expense = {
    timestamp: string;
    amount: number;
    description: string;
    shared: boolean;
    sharedAmount: number;
    account: number;
    credAccount: number;
    debAccount: number;
    sharedWith: number;
    category: string;
};


export const ExpenseToOperation = function (e: Expense): Operation {
    let op = {
        timestamp: e.timestamp,
        description: e.description,
        transactions: [
            {
                amount: e.amount,
                from: { id: e.account },
                to: { id: 0 },
            },

        ],
    };

    if (e.shared) {
        op.transactions.push({
            amount: e.sharedAmount,
            to: { id: e.credAccount },
            from: { id: e.debAccount },
        })
    }

    return op;
}