package catalog

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	elastic "gopkg.in/olivere/elastic.v5"
)

var (
	ErrNotFound = errors.New("product not found")
)

type Repository interface {
	Close()
	PutProduct(ctx context.Context, product *Product) error
	GetProductById(ctx context.Context, id string) (*Product, error)
	ListsProducts(ctx context.Context, skip int, take int) ([]*Product, error)
	ListsProductsWithIDs(ctx context.Context, ids []string) ([]*Product, error)
	SearchProducts(ctx context.Context, query string, skip int, take int) ([]*Product, error)
}

type elasticRepository struct {
	client *elastic.Client
}

type ProductDocument struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func NewElasticRepository(url string) (Repository, error) {
	client, err := elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(false))
	if err != nil {
		return nil, err
	}
	return &elasticRepository{client: client}, nil
}

func (r *elasticRepository) Close() {
	// No explicit close method for elastic.Client
}

func (r *elasticRepository) PutProduct(ctx context.Context, product *Product) error {
	doc := ProductDocument{
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}
	_, err := r.client.Index().
		Index("catalog").
		Type("product").
		Id(product.ID).
		BodyJson(doc).
		Do(ctx)
	return err
}

func (r *elasticRepository) GetProductById(ctx context.Context, id string) (*Product, error) {
	res, err := r.client.Get().
		Index("catalog").
		Type("product").
		Id(id).
		Do(ctx)
	if err != nil {
		if elastic.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if !res.Found {
		return nil, ErrNotFound
	}
	var doc ProductDocument
	err = json.Unmarshal(*res.Source, &doc)
	if err != nil {
		return nil, err
	}
	return &Product{
		ID:          res.Id,
		Name:        doc.Name,
		Description: doc.Description,
		Price:       doc.Price,
	}, nil
}

func (r *elasticRepository) ListsProducts(ctx context.Context, skip int, take int) ([]*Product, error) {
	query := elastic.NewMatchAllQuery()
	searchResult, err := r.client.Search().
		Index("catalog").
		Query(query).
		From(skip).Size(take).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	return r.convertSearchResults(searchResult), nil
}

func (r *elasticRepository) ListsProductsWithIDs(ctx context.Context, ids []string) ([]*Product, error) {
	query := elastic.NewIdsQuery().Ids(ids...)
	searchResult, err := r.client.Search().
		Index("catalog").
		Query(query).
		Size(len(ids)).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	return r.convertSearchResults(searchResult), nil
}

func (r *elasticRepository) SearchProducts(ctx context.Context, query string, skip int, take int) ([]*Product, error) {
	matchQuery := elastic.NewMultiMatchQuery(query, "name", "description")
	searchResult, err := r.client.Search().
		Index("catalog").
		Type("product").
		Query(matchQuery).
		From(skip).Size(take).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	return r.convertSearchResults(searchResult), nil
}

func (r *elasticRepository) convertSearchResults(searchResult *elastic.SearchResult) []*Product {
	var products []*Product
	for _, hit := range searchResult.Hits.Hits {
		var doc ProductDocument
		err := json.Unmarshal(*hit.Source, &doc)
		if err != nil {
			continue
		}
		product := &Product{
			ID:          hit.Id,
			Name:        doc.Name,
			Description: doc.Description,
			Price:       doc.Price,
		}
		products = append(products, product)
	}
	return products
}
	