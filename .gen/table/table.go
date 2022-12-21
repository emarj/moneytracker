//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

// UseSchema changes all global tables/views with the value returned
// returned by calling FromSchema on them. Passing an empty string to this function
// will cause queries to be generated without any table/view alias.
func UseSchema(schema string) {
	Account = Account.FromSchema(schema)
	AccountType = AccountType.FromSchema(schema)
	Balance = Balance.FromSchema(schema)
	Category = Category.FromSchema(schema)
	Entity = Entity.FromSchema(schema)
	Info = Info.FromSchema(schema)
	Operation = Operation.FromSchema(schema)
	OperationType = OperationType.FromSchema(schema)
	Transaction = Transaction.FromSchema(schema)
	User = User.FromSchema(schema)
}
