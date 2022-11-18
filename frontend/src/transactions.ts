export const isInternal = (t) => t.from.entity_id === t.to.entity_id;
export const isExpense = (t, eID) =>
    t.from.entity_id !== t.to.entity_id && t.from.entity_id === eID;
export const isIncome = (t, eID) =>
    t.from.entity_id !== t.to.entity_id && t.to.entity_id === eID;