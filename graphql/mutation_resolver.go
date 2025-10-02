package main


type mutationResolver struct{ 
	server *Server
}


// func (r *mutationResolver) CreateAccount(ctx context.Context, input AccountInput) (*Account, error) {
// 	account, err := r.server.accountClient.Create(ctx, input)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return account, nil
// }


// func (r *mutationResolver) CreateProduct(ctx context.Context, input ProductInput) (*Product, error) {
// 	product, err := r.server.catalogClient.Create(ctx, input)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return product, nil
// }


// func (r *mutationResolver) CreateOrder(ctx context.Context, input OrderInput) (*Order, error) {
// 	order, err := r.server.orderClient.Create(ctx, input)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return order, nil
// }