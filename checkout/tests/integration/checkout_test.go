package integration

import (
	"context"
	"log"
	"route256/checkout/internal/api"
	cliloms "route256/checkout/internal/clients/loms"
	cliproduct "route256/checkout/internal/clients/product"
	"route256/checkout/internal/config"
	"route256/checkout/internal/domain"
	mocks "route256/checkout/internal/domain/mocks"
	pgcartitems "route256/checkout/internal/repository/postgres/cartitems"
	"route256/checkout/pkg/checkout"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProductInfo struct {
	SKU   uint32
	Count uint16
	Name  string
	Price uint32
}

type TestSuiteCheckout struct {
	suite.Suite

	ctx     context.Context
	cliLoms *cliloms.LomsService
	model   *domain.Model
	service *api.Service

	productHub map[uint32]ProductInfo
}

func (s *TestSuiteCheckout) SetupSuite() {
	cfg := config.Config{}
	cfg.Services.Loms = config.ConfigService{
		Netloc: "localhost:8081",
	}
	cfg.Postgres = config.ConfigPostgres{
		Host:     "localhost",
		Port:     5433,
		User:     "postgres",
		Password: "password",
		Database: "checkout",
	}

	connLoms, err := grpc.Dial(
		cfg.Services.Loms.Netloc,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	s.Require().NoError(err)

	pool, err := pgxpool.Connect(
		context.Background(),
		cfg.Postgres.URL(),
	)
	s.Require().NoError(err)

	s.productHub = map[uint32]ProductInfo{
		773587830: {
			SKU:   773587830,
			Name:  "product_773587830",
			Price: 10,
		},
		773596051: {
			SKU:   773596051,
			Name:  "product_773596051",
			Price: 100,
		},
	}

	// NOTE: I don't have direct access to ProductService, so I will use mock instead
	productService := mocks.NewProductClient(s.T())
	for _, item := range s.productHub {
		productService.On("GetProduct", mock.Anything, item.SKU).Return(cliproduct.ResponseGetProduct{
			Name:  item.Name,
			Price: item.Price,
		}, nil)
	}

	s.cliLoms = cliloms.NewLomsClient(connLoms)
	s.model = domain.New(
		s.cliLoms,
		productService,
		pgcartitems.NewCartItemsRepository(pool),
	)

	s.service = api.NewService(s.model)
	s.ctx = context.Background()
}

func (s *TestSuiteCheckout) TestCase1() {
	var (
		userId int64 = 1

		sku1 uint32 = 773587830
		sku2 uint32 = 773596051

		resp *checkout.ResponseListCart
		err  error
	)

	// success: delete cart before test
	err = s.model.DeleteCart(s.ctx, userId)
	s.Require().NoError(err)

	// success: add known item to cart
	_, err = s.service.AddToCart(s.ctx, &checkout.RequestAddToCart{
		User:  userId,
		Sku:   sku1,
		Count: 5,
	})
	s.Require().NoError(err)

	// success: list cart
	resp, err = s.service.ListCart(s.ctx, &checkout.RequestListCart{
		User: userId,
	})
	log.Printf("ListCart: %+v", resp)

	s.Require().NoError(err)
	s.compareCarts(resp, &checkout.ResponseListCart{
		Items: []*checkout.ResponseListCart_CartItem{
			s.createCartItem(sku1, 5),
		},
		TotalPrice: s.productHub[sku1].Price * 5,
	})

	// fail: not enough stocks
	_, err = s.service.AddToCart(s.ctx, &checkout.RequestAddToCart{
		User:  userId,
		Sku:   sku1,
		Count: 2,
	})
	s.Require().Error(err, api.ErrProductInsufficient)

	// success: delete item from cart
	_, err = s.service.DeleteFromCart(s.ctx, &checkout.RequestDeleteFromCart{
		User:  userId,
		Sku:   sku1,
		Count: 2,
	})
	s.Require().NoError(err)

	// success: list cart
	resp, err = s.service.ListCart(s.ctx, &checkout.RequestListCart{
		User: userId,
	})
	log.Printf("ListCart: %+v", resp)

	s.Require().NoError(err)
	s.compareCarts(resp, &checkout.ResponseListCart{
		Items: []*checkout.ResponseListCart_CartItem{
			s.createCartItem(sku1, 3),
		},
		TotalPrice: s.productHub[sku1].Price * 3,
	})

	// success: add one more item
	_, err = s.service.AddToCart(s.ctx, &checkout.RequestAddToCart{
		User:  userId,
		Sku:   sku2,
		Count: 3,
	})
	s.Require().NoError(err)

	// success: list cart
	resp, err = s.service.ListCart(s.ctx, &checkout.RequestListCart{
		User: userId,
	})
	log.Printf("ListCart: %+v\n", resp)

	s.Require().NoError(err)
	s.compareCarts(resp, &checkout.ResponseListCart{
		Items: []*checkout.ResponseListCart_CartItem{
			s.createCartItem(sku1, 3),
			s.createCartItem(sku2, 3),
		},
		TotalPrice: s.productHub[sku1].Price*3 + s.productHub[sku2].Price*3,
	})

	// success: purchase
	order, err := s.service.Purchase(s.ctx, &checkout.RequestPurchase{
		User: userId,
	})
	defer func() {
		err := s.cliLoms.CancelOrder(s.ctx, order.OrderID)
		if err != nil {
			log.Fatalf("failed to cancel order: %v", err)
		}
	}()
	log.Printf("Purchase: %+v\n", order)

	s.Require().NoError(err)

	// success: list cart
	resp, err = s.service.ListCart(s.ctx, &checkout.RequestListCart{
		User: userId,
	})
	log.Printf("ListCart: %+v\n", resp)

	s.Require().NoError(err)
	s.compareCarts(resp, &checkout.ResponseListCart{
		Items:      make([]*checkout.ResponseListCart_CartItem, 0),
		TotalPrice: 0,
	})
}

func (s *TestSuiteCheckout) TestCase2() {
	var (
		userId int64 = 1

		sku1 uint32 = 773587830
		sku2 uint32 = 773596051
		sku3 uint32 = 773596067

		resp *checkout.ResponseListCart
		err  error
	)

	// success: delete cart before test
	err = s.model.DeleteCart(s.ctx, userId)
	s.Require().NoError(err)

	// success: add known item to cart
	_, err = s.service.AddToCart(s.ctx, &checkout.RequestAddToCart{
		User:  userId,
		Sku:   sku1,
		Count: 5,
	})
	s.Require().NoError(err)

	// success: add known item to cart
	_, err = s.service.AddToCart(s.ctx, &checkout.RequestAddToCart{
		User:  userId,
		Sku:   sku2,
		Count: 2,
	})
	s.Require().NoError(err)

	// success: list cart
	resp, err = s.service.ListCart(s.ctx, &checkout.RequestListCart{
		User: userId,
	})
	log.Printf("ListCart: %+v", resp)

	s.Require().NoError(err)
	s.compareCarts(resp, &checkout.ResponseListCart{
		Items: []*checkout.ResponseListCart_CartItem{
			s.createCartItem(sku1, 5),
			s.createCartItem(sku2, 2),
		},
		TotalPrice: s.productHub[sku1].Price*5 + s.productHub[sku2].Price*2,
	})

	// success: delete item from cart
	_, err = s.service.DeleteFromCart(s.ctx, &checkout.RequestDeleteFromCart{
		User:  userId,
		Sku:   sku2,
		Count: 10,
	})
	s.Require().NoError(err)

	// success: delete not present item from cart
	_, err = s.service.DeleteFromCart(s.ctx, &checkout.RequestDeleteFromCart{
		User:  userId,
		Sku:   sku3,
		Count: 10,
	})
	s.Require().NoError(err)

	// success: list cart
	resp, err = s.service.ListCart(s.ctx, &checkout.RequestListCart{
		User: userId,
	})
	log.Printf("ListCart: %+v", resp)

	s.Require().NoError(err)
	s.compareCarts(resp, &checkout.ResponseListCart{
		Items: []*checkout.ResponseListCart_CartItem{
			s.createCartItem(sku1, 5),
		},
		TotalPrice: s.productHub[sku1].Price * 5,
	})

	// success: purchase
	order, err := s.service.Purchase(s.ctx, &checkout.RequestPurchase{
		User: userId,
	})
	defer func() {
		err := s.cliLoms.CancelOrder(s.ctx, order.OrderID)
		if err != nil {
			log.Fatalf("failed to cancel order: %v", err)
		}
	}()
	log.Printf("Purchase: %+v\n", order)

	s.Require().NoError(err)

	// success: list cart
	resp, err = s.service.ListCart(s.ctx, &checkout.RequestListCart{
		User: userId,
	})
	log.Printf("ListCart: %+v\n", resp)

	s.Require().NoError(err)
	s.compareCarts(resp, &checkout.ResponseListCart{
		Items:      make([]*checkout.ResponseListCart_CartItem, 0),
		TotalPrice: 0,
	})

	// fail: not enough stocks
	_, err = s.service.AddToCart(s.ctx, &checkout.RequestAddToCart{
		User:  userId,
		Sku:   sku1,
		Count: 2,
	})
	s.Require().Error(err, api.ErrProductInsufficient)
}

func (s *TestSuiteCheckout) createCartItem(sku uint32, count uint32) *checkout.ResponseListCart_CartItem {
	item := checkout.ResponseListCart_CartItem{
		Sku:   sku,
		Count: count,
		Name:  s.productHub[sku].Name,
		Price: s.productHub[sku].Price,
	}
	return &item
}

func (s *TestSuiteCheckout) compareCarts(cartA, cartB *checkout.ResponseListCart) {
	// compare number items in carts
	s.Require().Equal(len(cartA.Items), len(cartB.Items))

	cartMapA := make(map[uint32]*checkout.ResponseListCart_CartItem)
	for _, item := range cartA.Items {
		cartMapA[item.Sku] = item
	}

	cartMapB := make(map[uint32]*checkout.ResponseListCart_CartItem)
	for _, item := range cartB.Items {
		cartMapB[item.Sku] = item
	}

	// compare number items in carts (mapped)
	s.Require().Equal(len(cartMapA), len(cartMapB))

	// compare items in carts
	for sku, itemA := range cartMapA {
		itemB, isOk := cartMapB[sku]
		s.Require().True(isOk)

		s.Require().Equal(itemA.Sku, itemB.Sku)
		s.Require().Equal(itemA.Count, itemB.Count)

		if itemB.Price > 0 {
			s.Require().Equal(itemA.Price, itemB.Price)
		}

		if len(itemB.Name) > 0 {
			s.Require().Equal(itemA.Name, itemB.Name)
		}
	}

	// compare total price
	s.Require().Equal(cartA.TotalPrice, cartB.TotalPrice)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(TestSuiteCheckout))
}
