package queries

const GetTransactionsByAccountQuery string = `SELECT  t.id,
													t.from_id,
													t.to_id,
													t.amount,
													t.operation_id,
													op.id,
													op.timestamp,
													op.created_by_id,
													op.description
											FROM transactions t
											INNER JOIN operations op
											ON t.operation_id = op.id
											WHERE from_id = ?
											OR to_id = ?
											ORDER BY op.timestamp DESC;`

const GetOperationByEntityQuery string = `SELECT  t.*,
												op.*,
												fa.name AS from_name,
												fa.display_name AS from_display_name,
												ta.name AS to_name,
												ta.display_name AS to_display_name,
												fe.id,
												fe.name,
												te.id,
												te.name
												FROM transactions t
													INNER JOIN operations op ON t.operation_id = op.id
													INNER JOIN accounts AS fa ON t.from_id = fa.id
													INNER JOIN accounts AS ta ON t.to_id = ta.id
													INNER JOIN entities AS fe ON fa.owner_id = fe.id
													INNER JOIN entities AS te ON ta.owner_id = te.id
												WHERE fa.owner_id = ?
													OR ta.owner_id = ?
												ORDER BY op.timestamp DESC,op.id,t.id;`

const GetOperationQuery string = `SELECT  			op.id,
													op.timestamp,
													op.created_by_id,
													op.description,
													t.id,
													t.from_id,
													t.to_id,
													t.amount
											FROM operations op
											INNER JOIN transactions t
											ON t.operation_id = op.id
											WHERE op.id = ?;`

const InsertTransactionQuery string = `INSERT INTO  transactions (
													from_id,
													to_id,
													amount,
													operation_id)
											VALUES (?,?,?,?);`

const InsertOperationQuery string = `INSERT INTO  operations (
	timestamp,
	created_by_id,
	description)
VALUES (?,?,?);`

const InsertEntityQuery string = `INSERT INTO  entities (id, name, is_system) VALUES (?,?,?);`

const InsertAccountQuery string = `INSERT INTO accounts (id,owner_id,name,display_name,is_system,is_world,is_credit)
											   VALUES (?,?,?,?,?,?,?);`

const InsertBalanceQuery string = `INSERT INTO balances (timestamp,account_id,value) VALUES (?,?,?);`
