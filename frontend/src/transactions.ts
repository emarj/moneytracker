export const isInternal = (t) => t.from.owner.id === t.to.owner.id;
export const isExpense = (t, eID) =>
    t.from.owner.id !== t.to.owner.id && t.from.owner.id === eID;
export const isIncome = (t, eID) =>
    t.from.owner.id !== t.to.owner.id && t.to.owner.id === eID;