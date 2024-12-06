package catalog

import (
	"context"
	"encoding/json"
	"errors"

	elastic "gopkg.in/olivere/elastic.v5"
)

type Repository interface {
	Close()
	CreateProduct(ctx context.Context, product Product) (*Product, error)
	GetProductById(ctx context.Context, id string) (*Product, error)
	ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error)
	ListProductsWithIds(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error)
}

type elasticRepository struct {
	client *elastic.Client
}

func NewElasticRepository(url string) (Repository, error) {
	client, err := elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(false))
	if err != nil {
		return nil, err
	}
	return &elasticRepository{client}, nil
}

type ProductDocument struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func (e *elasticRepository) CreateProduct(ctx context.Context, p Product) (*Product, error) {
	_, err := e.client.Index().Index("catalog").Type("product").Id(p.Id).BodyJson(ProductDocument{
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}).Do(ctx)
	if err != nil {
		return nil, err
	}
	// pd := ProductDocument{}
	// if err = json.Unmarshal(*res.Source, &pd); err != nil {
	// 	return nil, err
	// }
	return &Product{
		Id:          "200",
		Name:        "name",
		Description: "pd.Description",
		Price:       0000,
	}, nil
}

func (e *elasticRepository) GetProductById(ctx context.Context, id string) (*Product, error) {
	res, err := e.client.Get().Index("catalog").Type("product").Id(id).Do(ctx)
	if err != nil {
		return nil, err
	}
	if !res.Found {
		return nil, errors.New("item not found")
	}
	return &Product{}, nil
}

func (e *elasticRepository) ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error) {
	res, err := e.client.Search().Index("catalog").Type("product").Query(elastic.NewMatchAllQuery()).From(int(skip)).Size(int(take)).Do(ctx)
	if err != nil {
		return nil, err
	}
	products := []Product{}
	for _, hit := range res.Hits.Hits {
		p := ProductDocument{}
		if err = json.Unmarshal(*hit.Source, &p); err == nil {
			products = append(products, Product{
				Id:          hit.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
		}
	}
	return products, nil
}

func (e *elasticRepository) ListProductsWithIds(ctx context.Context, ids []string) ([]Product, error) {

	items := []*elastic.MultiGetItem{}

	for _, id := range ids {
		items = append(items, elastic.NewMultiGetItem().Index("catalog").Type("product").Id(id))
	}
	res, err := e.client.MultiGet().Add(items...).Do(ctx)
	if err != nil {
		return nil, err
	}
	products := []Product{}

	for _, doc := range res.Docs {
		p := ProductDocument{}
		if err := json.Unmarshal(*doc.Source, &p); err == nil {
			products = append(products, Product{
				Id:          doc.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
		}
	}
	return products, nil
}

func (e *elasticRepository) SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error) {
	res, err := e.client.Search().Index("catalog").Type("Product").Query(elastic.NewMultiMatchQuery(query, "name", "description")).From(int(skip)).Size(int(take)).Do(ctx)
	if err != nil {
		return nil, err
	}
	products := []Product{}
	for _, hit := range res.Hits.Hits {
		p := ProductDocument{}
		if err = json.Unmarshal(*hit.Source, &p); err == nil {
			products = append(products, Product{
				Id:          hit.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
		}
	}
	return products, nil
}

func (e *elasticRepository) Close() {

}
