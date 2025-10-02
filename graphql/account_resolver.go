package main

type accountResolver struct{ 
	server *Server
}



// func (r *queryResolver) Orders(ctx context.Context, pagination *PaginationInput,  id *string) ([]*Order, error) {
// 	orders, err := r.server.orderClient.List(ctx, pagination, id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return orders, nil
// }