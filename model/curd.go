package model

import (
	"context"
	"log"

	"github.com/dgraph-io/dgo/v200/protos/api"
)

/*
param: mutationString是每一个mutation，每个mutation类似这样：
`
	_:n1 <name> "user" .
	_:n1 <email> "user@dgraphO.io" .
`
*/
func MutationSet(ctx context.Context, mutation string) (*api.Response, error) {
	txn := client.NewTxn()
	//fmt.Println("mu: ", mutation)
	// Discard可以当事务出错的时候回滚
	defer txn.Discard(ctx)
	assigned, err := txn.Mutate(ctx, &api.Mutation{
		SetNquads: []byte(mutation),
		CommitNow: true,
	})
	if err != nil {
		log.Println("mutation error: " + err.Error())
		return nil, err
	}
	return assigned, nil
}

func MutationDelete(ctx context.Context, mutationString string) (resp *api.Response, err error) {
	// 要指定uid，就先外面组装好mutationString再传进来
	txn := client.NewTxn()
	defer txn.Discard(ctx)
	mu := &api.Mutation{
		CommitNow: true,
		DelNquads: []byte(mutationString),
	}
	resp, err = txn.Mutate(ctx, mu)
	return
}

/*
upsert demo:
query := `{
	keyUids as q(func: eq(value, "` + someValue + `"))
}`

var mutations []*api.Mutation
mutations = append(mutations, &api.Mutation{
	Cond:    `@if(eq(len(keyUids), 0))`, //key does not already exist
	SetJson: []byte(`{"uid":"` + entityUid + `","hardKeys":{"uid":"_:H", "dgraph.type": "HardKey", "value": "` + someValue  + `}}`),
})
mutations = append(mutations, &api.Mutation{
	Cond:    `@if(NOT eq(len(keyUids), 0))`,  //key already exists
	SetNquads: []byte(`<` + entityUid + `> <hardKeys> uid(keyUids) .`), // Link Key to entity
})
request := &api.Request{
	Query:     query,
	Mutations: mutations,
	CommitNow: true,
}

results, err := dgraphClient.NewTxn().Do(context.Background(), request)
if err != nil {
	log.Fatal(err)
}
fmt.Println(results.Uids)
*/
func MutationSetWithConditionUpsert(ctx context.Context, mutationStrings []map[string]string, query string) (*api.Response, error) {
	txn := client.NewTxn()
	defer txn.Discard(ctx)
	var mutations []*api.Mutation
	for _, v := range mutationStrings {
		mutations = append(mutations, &api.Mutation{
			Cond:      v["cond"],
			SetNquads: []byte(v["mutation"]),
		})
	}

	request := &api.Request{
		Query:     query,
		Mutations: mutations,
		CommitNow: true,
	}

	results, err := txn.Do(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return results, nil

}

func MutationSetWithUpsert(ctx context.Context, mutationStrings []string, query string) (*api.Response, error) {
	txn := client.NewTxn()
	defer txn.Discard(ctx)
	var mutations []*api.Mutation
	for _, v := range mutationStrings {
		mutations = append(mutations, &api.Mutation{
			SetNquads: []byte(v),
		})
	}

	request := &api.Request{
		Query:     query,
		Mutations: mutations,
		CommitNow: true,
	}

	results, err := txn.Do(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func MutationDeleteWithUpsert(ctx context.Context, mutationStrings []string, query string) (*api.Response, error) {
	txn := client.NewTxn()
	defer txn.Discard(ctx)
	var mutations []*api.Mutation
	for _, v := range mutationStrings {
		mutations = append(mutations, &api.Mutation{
			DelNquads: []byte(v),
		})
	}

	request := &api.Request{
		Query:     query,
		Mutations: mutations,
		CommitNow: true,
	}

	results, err := txn.Do(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return results, nil
}

/*
params:
	setStrings: set的语句，形如：
	`
		_:n1 <name> "user" .
		_:n1 <email> "user@dgraphO.io" .
	`
	delStrings: delete的语句，形如：
	`
		<0x222> <name> * .
		<0x333> * * .
	`
returns:
	response
	error
*/
func MutationSetAndDeleteWithUpsert(ctx context.Context, setStrings []string, delStrings []string, query string) (*api.Response, error) {
	// 这是先set后delete
	txn := client.NewTxn()
	defer txn.Discard(ctx)
	var mutations []*api.Mutation

	if len(setStrings) != 0 {
		for _, v := range setStrings {
			mutations = append(mutations, &api.Mutation{
				SetNquads: []byte(v),
			})
		}
	}

	if len(delStrings) != 0 {
		for _, v := range delStrings {
			mutations = append(mutations, &api.Mutation{
				DelNquads: []byte(v),
			})
		}
	}

	request := &api.Request{
		Query:     query,
		Mutations: mutations,
		CommitNow: true,
	}

	results, err := txn.Do(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func MutationDeleteAndSetWithUpsert(ctx context.Context, delStrings []string, setStrings []string, query string) (*api.Response, error) {
	// 这是先delete后set
	txn := client.NewTxn()
	defer txn.Discard(ctx)
	var mutations []*api.Mutation

	if len(delStrings) != 0 {
		for _, v := range delStrings {
			mutations = append(mutations, &api.Mutation{
				DelNquads: []byte(v),
			})
		}
	}

	if len(setStrings) != 0 {
		for _, v := range setStrings {
			mutations = append(mutations, &api.Mutation{
				SetNquads: []byte(v),
			})
		}
	}

	request := &api.Request{
		Query:     query,
		Mutations: mutations,
		CommitNow: true,
	}

	results, err := txn.Do(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func Query(ctx context.Context, query string) (resp *api.Response, err error) {
	txn := client.NewReadOnlyTxn().BestEffort()
	resp, err = txn.Query(context.Background(), query)
	//fmt.Println("resp:", resp)
	return
}
