package queries

import (
	"fmt"
	"strconv"
	"strings"

	"ronche.se/moneytracker"
)

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

func GetTransactionByAccount(aID int) string {
	return fmt.Sprintf(strings.ReplaceAll(GetTransactionsByAccountQuery, "?", "%d"), aID, aID)
}

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

func GetOperation(opID int) string {
	return fmt.Sprintf(strings.ReplaceAll(GetTransactionsByAccountQuery, "?", "%d"), opID)
}

const InsertTransactionQuery string = `INSERT INTO  transactions (
													from_id,
													to_id,
													amount,
													operation_id)
											VALUES (?,?,?,?);`

func InsertTransaction(t moneytracker.Transaction) string {
	return fmt.Sprintf(strings.ReplaceAll(InsertTransactionQuery, "?", "%s"),
		strconv.Itoa(t.From.ID),
		strconv.Itoa(t.To.ID),
		t.Amount.String(),
		strconv.Itoa(t.Operation.ID),
	)
}

const InsertOperationQuery string = `INSERT INTO  operations (
	timestamp,
	created_by_id,
	description)
VALUES (?,?,?);`

func InsertOperation(op moneytracker.Operation) string {

	return fmt.Sprintf(strings.ReplaceAll(InsertOperationQuery, "?", "%s"),
		"'"+op.Timestamp.String()+"'", //op.Timestamp,
		strconv.Itoa(op.CreatedByID),
		op.Description,
	)
}

const InsertEntityQuery string = `INSERT INTO  entities (id, name, is_system) VALUES (?,?,?);`

func InsertEntity(e moneytracker.Entity) string {
	return fmt.Sprintf(strings.ReplaceAll(InsertEntityQuery, "?", "%s"),
		strconv.Itoa(e.ID), e.Name, strconv.FormatBool(e.System),
	)
}

const InsertAccountQuery string = `INSERT INTO accounts (id,owner_id,name,display_name,is_system,is_world,is_credit)
											   VALUES (?,?,?,?,?,?,?);`

func InsertAccount(a moneytracker.Account) string {
	return fmt.Sprintf(strings.ReplaceAll(InsertAccountQuery, "?", "%s"),
		strconv.Itoa(a.ID), strconv.Itoa(a.EntityID), a.Name, a.DisplayName, strconv.FormatBool(a.IsSystem), strconv.FormatBool(a.IsWorld), strconv.FormatBool(a.IsCredit),
	)
}

const InsertBalanceQuery string = `INSERT INTO balances (timestamp,account_id,value) VALUES (?,?,?);`

func InsertBalance(b moneytracker.Balance) string {
	return fmt.Sprintf(strings.ReplaceAll(InsertBalanceQuery, "?", "%s"),
		"'"+b.Timestamp.String()+"'", strconv.Itoa(b.AccountID), b.Value,
	)
}
