package model

import (
	"context"
	"fmt"

	"github.com/dgraph-io/dgo/v200/protos/api"
)

func SetupSchema(ctx context.Context, schema string) error {
	return client.Alter(ctx, &api.Operation{
		Schema: schema,
	})
}

// func DropData() error {
// 	return client.Alter(context.Background(), &api.Operation{DropOp: api.Operation_DATA})
// }

// 更新field_name称，先把旧字段数据全部copy到新字段然后删除旧字段，注意：数据量大的时候代价很大
// todo: 更新数据名称后，与这个数据名关联的curd代码都要改变
func UpdatePredicateName(ctx context.Context, oldFieldName string, newFieldName string, theType string) error {
	txn := client.NewTxn()
	defer txn.Discard(ctx)
	query := fmt.Sprintf(`
	query {
		u as var(func: type(%s)){
			n as %s
		}
	}
	`, theType, oldFieldName)

	s := fmt.Sprintf(`uid(u) <%s>  val(n) .`, newFieldName)
	d := fmt.Sprintf(`uid(u) <%s> * .`, oldFieldName)

	mu := &api.Mutation{
		SetNquads: []byte(s),
		DelNquads: []byte(d),
	}

	req := &api.Request{
		Query:     query,
		Mutations: []*api.Mutation{mu},
		CommitNow: true,
	}
	if _, err := txn.Do(ctx, req); err != nil {
		return err
	}

	return nil
}

func DropAll() error {
	return client.Alter(context.Background(), &api.Operation{DropOp: api.Operation_ALL})
}

func DropData() error {
	return client.Alter(context.Background(), &api.Operation{DropOp: api.Operation_DATA})
}
