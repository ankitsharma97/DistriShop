package main

type queryResolver struct{
	server *Server
}


// func (r *queryResolver) Accounts(ctx context.Context, pagination *PaginationInput, id string) ([]*Account, error) {
// 	accounts, err := r.server.accountClient.List(ctx, pagination, id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return accounts, nil
// }

// func (r *queryResolver) Products(ctx context.Context, pagination *PaginationInput, Query *string, id *string) ([]*Product, error) {
// 	products, err := r.server.catalogClient.List(ctx, pagination, id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return products, nil
// }

