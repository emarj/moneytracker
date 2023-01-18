export const isInternal = (t) => t.sign === 0;
export const isExpense = (t, eID) => t.sign === -1;
//t.from.owner.id !== t.to.owner.id && t.from.owner.id === eID;
export const isIncome = (t, eID) => t.sign === 1;
    //t.from.owner.id !== t.to.owner.id && t.to.owner.id === eID;