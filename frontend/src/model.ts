


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
    timestamp?: Date;
    amount?: number;
    from_id: number;
    from?: Account;
    to_id: number;
    to?: Account;
    operation_id?: number;
    operation?: Operation;
    sign?: number;
};

export type Operation = {
    id?: number;
    date_modified?: Date;
    description: string;
    category_id?: number;
    transactions?: Transaction[];
};

export type Tag = {
    id?: number;
    name: string;
};

export class Expense {
    timestamp: Date = new Date();
    private _amount: number = 0;
    description: string = "";
    account: number;
    shares: Share[] = [];
    category_id: number;
    tags: Tag[] = [];

    set isShared(shared: boolean) {
        this.shares = (shared) ? [new Share(this._amount)] : [];
    }

    get isShared(): boolean {
        return (this.shares.length > 0)
    }

    set amount(amount: number) {
        this._amount = amount;
        if (this.isShared) {
            for (const s of this.shares) {
                s.total = this.amount
            }
        }
    }

    get amount(): number { return this._amount }

    toOperation(): Operation {
        let op = {
            description: this.description,
            category_id: this.category_id,
            transactions: [
                {
                    timestamp: this.timestamp,
                    amount: this.amount,
                    from_id: this.account,
                    to_id: 0,
                },

            ],
        };

        if (this.isShared) {
            for (const s of this.shares) {
                op.transactions.push({
                    timestamp: this.timestamp,
                    amount: s.amount,
                    to_id: s.credAccount,
                    from_id: s.debAccount,
                })
            }
        }
        return op;
    }

};

export class Share {
    private _amount: number = 0;
    private _quota: number = 50;
    private _total: number = 0;
    with: number;
    credAccount: number;
    debAccount: number;
    isCredit: boolean = true;

    constructor(total?: number) {
        if (total) {
            this.total = total;
            this.computeAmount();
        }
    }

    set quota(quota: number) {
        this._quota = quota;
        this.computeAmount();
    }

    get quota(): number { return this._quota }

    set total(total: number) {
        this._total = total;
        this.computeAmount();
    }

    get total(): number { return this._total }

    set amount(amount: number) {
        this._amount = amount;
        this.computeQuota();
    }

    get amount(): number { return this._amount }

    private computeAmount() {

        this._amount = (this._total * this._quota) / 100;
    }

    private computeQuota() {
        if (this._amount !== 0) {
            this._quota = (this._amount / this._total) * 100;
        }
    }
};

export const emptyTransaction: Transaction = {
    amount: 0, to_id: 0, from_id: 0
};

export const emptyOperation: Operation = {
    description: "",
    category_id: 0,
    transactions: [emptyTransaction],
};





