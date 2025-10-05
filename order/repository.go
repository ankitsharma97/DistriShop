package order

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Repository interface {
	Close() error
	PutOrder(ctx context.Context, order *Order) error
	GetOrdersForAccount(ctx context.Context, accountID string) ([]*Order, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(databaseURL string) (Repository, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &postgresRepository{db: db}, nil
}

func (r *postgresRepository) Close() error {
	return r.db.Close()
}

func (r *postgresRepository) PutOrder(ctx context.Context, o *Order) (err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()
	_, err = tx.ExecContext(ctx, "INSERT INTO orders (id, account_id, created_at, total_price) VALUES ($1, $2, $3, $4)", o.ID, o.AccountID, o.CreatedAt, o.Total)
	if err != nil {
		return err
	}
	stmt, err := tx.PrepareContext(ctx, pq.CopyIn("order_products", "order_id", "product_id", "quantity"))
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, p := range o.Products {
		_, err = stmt.ExecContext(ctx, o.ID, p.ProductID, p.Quantity)
		if err != nil {
			return err
		}
	}
	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *postgresRepository) GetOrdersForAccount(ctx context.Context, accountID string) ([]*Order, error) {
	rows, err := r.db.QueryContext(ctx,
		`
		SELECT 
			o.id, o.account_id, o.created_at, o.total_price, 
			op.product_id, op.quantity 
		FROM orders o 
		JOIN order_products op ON o.id = op.order_id
		WHERE o.account_id = $1
		ORDER BY o.created_at DESC
		`, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*Order
	var lastOrder *Order
	var products []*OrderedProduct

	for rows.Next() {
		var orderID, dbAccountID, productID string
		var createdAt time.Time
		var total float64
		var quantity int

		err := rows.Scan(&orderID, &dbAccountID, &createdAt, &total, &productID, &quantity)
		if err != nil {
			return nil, err
		}

		// If this row belongs to a new order, push the previous one to orders
		if lastOrder == nil || lastOrder.ID != orderID {
			if lastOrder != nil {
				lastOrder.Products = products
				orders = append(orders, lastOrder)
			}
			lastOrder = &Order{
				ID:        orderID,
				AccountID: dbAccountID,
				CreatedAt: createdAt,
				Total:     total,
			}
			products = []*OrderedProduct{}
		}

		products = append(products, &OrderedProduct{
			ProductID: productID,
			Quantity:  quantity,
		})
	}

	// Append the last order if exists
	if lastOrder != nil {
		lastOrder.Products = products
		orders = append(orders, lastOrder)
	}

	return orders, nil
}
