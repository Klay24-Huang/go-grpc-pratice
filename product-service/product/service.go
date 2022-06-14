package product

type Service interface {
	FindAll() (*[]Product, error)
	FindById(Id int) (*Product, error)
	Create(productRequest ProductRequest) (*Product, error)
	Update(Id int, productRequest ProductRequest) (*Product, error)
	Delete(Id int) (*Product, error)
}

type service struct {
	productRepository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) FindAll() (*[]Product, error) {
	products, err := s.productRepository.FindAll()
	return products, err
}

func (s *service) FindById(Id int) (*Product, error) {
	product, err := s.productRepository.FindById(Id)
	return product, err
}

func (s *service) Create(productRequest ProductRequest) (*Product, error) {
	product := Product{
		Name:  productRequest.Name,
		Price: productRequest.Price,
	}
	newProduct, err := s.productRepository.Create(product)
	return newProduct, err
}

func (s *service) Update(Id int, productRequest ProductRequest) (*Product, error) {
	product, err := s.productRepository.FindById(Id)

	if err != nil {
		return nil, err
	}

	product.Name = productRequest.Name
	product.Price = productRequest.Price
	newProduct, err := s.productRepository.Update(*product)
	return newProduct, err
}

func (s *service) Delete(Id int) (*Product, error) {
	product, err := s.productRepository.Delete(Id)
	return product, err
}
