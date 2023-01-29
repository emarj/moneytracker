import { removeUnderscore } from "./util/utils";



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
    modified_on?: Date;
    description: string;
    category_id?: number;
    transactions?: Transaction[];
    balances?: any[];
};

export type Tag = {
    id?: number;
    name: string;
};

export class Expense {
    timestamp: Date = new Date();
    private _amount: number = null;
    description: string = "";
    account_id: number;
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

    toJSON() {
        return removeUnderscore(this)
    };

};

export class Share {
    private _amount: number = null;
    private _quota: number = 50;
    private _total: number = null;
    with_id: number = null;
    cred_account_id: number;
    deb_account_id: number;
    is_credit: boolean = true;

    constructor(total?: number) {
        if (total) {
            this.total = total;
            this.computeAmount();
        }
    }

    set quota(quota: number) {
        //console.log("set quota", quota)
        this._quota = quota;
        this.computeAmount();
    }

    get quota(): number { return this._quota }

    set total(total: number) {
        //console.log("set total", total)
        this._total = total;
        this.computeAmount();
    }

    get total(): number { return this._total }

    set amount(amount: number) {
        //console.log("set amount", amount)
        this._amount = amount;
        this.computeQuota();
    }

    get amount(): number { return this._amount }

    private computeAmount() {
        if (this._total === null) {
            this._amount = null
        } else {
            this._amount = (this._total * this._quota) / 100;
        }
    }

    private computeQuota() {
        if (this._amount !== 0) {
            this._quota = (this._amount / this._total) * 100;
        }
    }

    toJSON() {
        return removeUnderscore(this)
    };
};



export const emptyTransaction: Transaction = {
    amount: null, to_id: 0, from_id: 0
};

export const emptyOperation: Operation = {
    description: "",
    category_id: 0,
    transactions: [emptyTransaction],
};





