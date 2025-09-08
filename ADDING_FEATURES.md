# Adding New Features - Step-by-Step Guide

This guide shows how to add a new feature using **Product Management** as an example.

## Overview - Files to Create & Modify

**Create 7 files:**
1. Domain Entity - Data model
2. Repository Interface - Define operations  
3. Repository Implementation - Database connection
4. DTO - Input/Output objects
5. Service - Business logic
6. Handler - HTTP endpoints
7. Tests - Unit tests

**Modify 3 files:**
1. Container - Dependency injection
2. Routes - URL paths
3. Model Registry - Migration registration

---

## Step-by-Step Implementation

### Step 1: Create Domain Entity

**📁 internal/domain/product/entity/product.go**

```go
package entity

import (
    "time"
    "gorm.io/gorm"
)

type Product struct {
    ID          uint           `json:"id" gorm:"primaryKey"`
    Name        string         `json:"name" gorm:"size:255;not null"`
    SKU         string         `json:"sku" gorm:"size:100;uniqueIndex"`
    Price       float64        `json:"price" gorm:"type:decimal(10,2);not null"`
    Stock       int            `json:"stock" gorm:"default:0"`
    Category    string         `json:"category" gorm:"size:100;index"`
    Status      string         `json:"status" gorm:"size:20;default:'active'"`
    AdminID     uint           `json:"admin_id" gorm:"not null;index"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
```

### Step 2: Define Repository Interface

**📁 internal/domain/product/repository/product.go**

```go
package repository

import "dashboard-starter/internal/domain/product/entity"

type ProductRepository interface {
    // Basic CRUD
    GetByID(id uint) (*entity.Product, error)
    GetBySKU(sku string) (*entity.Product, error)
    Create(product *entity.Product) error
    Update(product *entity.Product) error
    Delete(id uint) error
    
    // List with filters
    List(offset, limit int, search, category, status string) ([]*entity.Product, int64, error)
    
    // Business specific
    UpdateStock(id uint, quantity int) error
}
```

### Step 3: Create DTO

**📁 internal/application/dto/product.go**

```go
package dto

type ProductInput struct {
    Name     string  `json:"name" binding:"required" validate:"required,min=3,max=255"`
    SKU      string  `json:"sku" binding:"required" validate:"required,alphanum"`
    Price    float64 `json:"price" binding:"required" validate:"required,min=0"`
    Stock    int     `json:"stock" validate:"min=0"`
    Category string  `json:"category" validate:"max=100"`
    Status   string  `json:"status" validate:"omitempty,oneof=active inactive"`
}

type StockUpdateInput struct {
    Quantity  int    `json:"quantity" binding:"required"`
    Operation string `json:"operation" binding:"required" validate:"oneof=add subtract set"`
}
```

### Step 4: Create Service

**📁 internal/application/services/product_service.go**

```go
package services

import (
    "dashboard-starter/internal/application/dto"
    "dashboard-starter/internal/domain/product/entity"
    "dashboard-starter/internal/domain/product/repository"
    "errors"
    "fmt"
)

type ProductApplicationService struct {
    productRepo repository.ProductRepository
}

func NewProductApplicationService(productRepo repository.ProductRepository) *ProductApplicationService {
    return &ProductApplicationService{
        productRepo: productRepo,
    }
}

func (s *ProductApplicationService) CreateProduct(input dto.ProductInput, adminID uint) (*entity.Product, error) {
    // Check duplicate SKU
    if existing, _ := s.productRepo.GetBySKU(input.SKU); existing != nil {
        return nil, fmt.Errorf("SKU %s already exists", input.SKU)
    }
    
    product := &entity.Product{
        Name:     input.Name,
        SKU:      input.SKU,
        Price:    input.Price,
        Stock:    input.Stock,
        Category: input.Category,
        Status:   input.Status,
        AdminID:  adminID,
    }
    
    if product.Status == "" {
        product.Status = "active"
    }
    
    if err := s.productRepo.Create(product); err != nil {
        return nil, err
    }
    
    return product, nil
}

func (s *ProductApplicationService) UpdateStock(id uint, input dto.StockUpdateInput) error {
    product, err := s.productRepo.GetByID(id)
    if err != nil {
        return errors.New("product not found")
    }
    
    switch input.Operation {
    case "add":
        product.Stock += input.Quantity
    case "subtract":
        if product.Stock < input.Quantity {
            return errors.New("insufficient stock")
        }
        product.Stock -= input.Quantity
    case "set":
        product.Stock = input.Quantity
    }
    
    return s.productRepo.Update(product)
}

func (s *ProductApplicationService) ListProducts(page, limit int, search, category, status string) ([]*entity.Product, int64, error) {
    offset := (page - 1) * limit
    return s.productRepo.List(offset, limit, search, category, status)
}
```

### Step 5: Implement Repository

**📁 internal/infrastructure/database/product_repository.go**

```go
package database

import (
    "dashboard-starter/internal/domain/product/entity"
    "dashboard-starter/internal/domain/product/repository"
    "gorm.io/gorm"
)

type productRepository struct {
    db *gorm.DB
}

func NewProductRepository(db *gorm.DB) repository.ProductRepository {
    return &productRepository{db: db}
}

func (r *productRepository) GetByID(id uint) (*entity.Product, error) {
    var product entity.Product
    err := r.db.First(&product, id).Error
    return &product, err
}

func (r *productRepository) GetBySKU(sku string) (*entity.Product, error) {
    var product entity.Product
    err := r.db.Where("sku = ?", sku).First(&product).Error
    return &product, err
}

func (r *productRepository) Create(product *entity.Product) error {
    return r.db.Create(product).Error
}

func (r *productRepository) Update(product *entity.Product) error {
    return r.db.Save(product).Error
}

func (r *productRepository) Delete(id uint) error {
    return r.db.Delete(&entity.Product{}, id).Error
}

func (r *productRepository) List(offset, limit int, search, category, status string) ([]*entity.Product, int64, error) {
    var products []*entity.Product
    var total int64
    
    query := r.db.Model(&entity.Product{})
    
    if search != "" {
        searchPattern := "%" + search + "%"
        query = query.Where("name LIKE ? OR sku LIKE ?", searchPattern, searchPattern)
    }
    
    if category != "" {
        query = query.Where("category = ?", category)
    }
    
    if status != "" {
        query = query.Where("status = ?", status)
    }
    
    // Count total
    query.Count(&total)
    
    // Get paginated results
    err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&products).Error
    
    return products, total, err
}

func (r *productRepository) UpdateStock(id uint, quantity int) error {
    return r.db.Model(&entity.Product{}).
        Where("id = ?", id).
        Update("stock", gorm.Expr("stock + ?", quantity)).Error
}
```

### Step 6: Create HTTP Handler

**📁 internal/interfaces/http/handlers/product_handler.go**

```go
package handlers

import (
    "dashboard-starter/internal/application/dto"
    "dashboard-starter/internal/application/services"
    "dashboard-starter/utils"
    "net/http"
    "strconv"
    
    "github.com/gin-gonic/gin"
)

type ProductHandler struct {
    productService *services.ProductApplicationService
}

func NewProductHandler(productService *services.ProductApplicationService) *ProductHandler {
    return &ProductHandler{
        productService: productService,
    }
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
    adminID, exists := c.Get("admin_id")
    if !exists {
        RespondWithError(c, http.StatusUnauthorized, "Admin authentication required")
        return
    }
    
    var input dto.ProductInput
    if err := c.ShouldBindJSON(&input); err != nil {
        RespondWithError(c, http.StatusBadRequest, err.Error())
        return
    }
    
    if err := utils.ValidateStruct(input); err != nil {
        RespondWithError(c, http.StatusBadRequest, err.Error())
        return
    }
    
    product, err := h.productService.CreateProduct(input, adminID.(uint))
    if err != nil {
        RespondWithError(c, http.StatusBadRequest, err.Error())
        return
    }
    
    RespondWithSuccess(c, http.StatusCreated, product)
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    search := c.Query("search")
    category := c.Query("category")
    status := c.Query("status")
    
    if page < 1 {
        page = 1
    }
    if limit < 1 || limit > 100 {
        limit = 10
    }
    
    products, total, err := h.productService.ListProducts(page, limit, search, category, status)
    if err != nil {
        RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve products")
        return
    }
    
    totalPages := (total + int64(limit) - 1) / int64(limit)
    
    RespondWithSuccess(c, http.StatusOK, products, gin.H{
        "page":       page,
        "limit":      limit,
        "total":      total,
        "totalPages": totalPages,
    })
}

func (h *ProductHandler) UpdateStock(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        RespondWithError(c, http.StatusBadRequest, "Invalid product ID")
        return
    }
    
    var input dto.StockUpdateInput
    if err := c.ShouldBindJSON(&input); err != nil {
        RespondWithError(c, http.StatusBadRequest, err.Error())
        return
    }
    
    if err := h.productService.UpdateStock(uint(id), input); err != nil {
        RespondWithError(c, http.StatusBadRequest, err.Error())
        return
    }
    
    RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Stock updated successfully"})
}
```

### Step 7: Update Container

**📁 internal/interfaces/http/container.go**

Add to imports:
```go
import productRepo "dashboard-starter/internal/domain/product/repository"
```

Add to Container struct:
```go
ProductHandler *handlers.ProductHandler
```

Add in NewContainer function:
```go
// Create product repository
productRepository := database.NewProductRepository(db.DB)

// Create product service  
productService := services.NewProductApplicationService(productRepository)

// Create product handler
productHandler := handlers.NewProductHandler(productService)

// Add to return
ProductHandler: productHandler,
```

### Step 8: Add Routes

**📁 routes/routes.go**

```go
// Product management routes
products := admin.Group("/products")
{
    products.GET("", container.ProductHandler.ListProducts)
    products.POST("", container.ProductHandler.CreateProduct)
    products.GET("/:id", container.ProductHandler.GetProduct)
    products.PUT("/:id", container.ProductHandler.UpdateProduct)
    products.DELETE("/:id", container.ProductHandler.DeleteProduct)
    products.PUT("/:id/stock", container.ProductHandler.UpdateStock)
}
```

### Step 9: Register Model

**📁 db/model.go**

```go
import productEntity "dashboard-starter/internal/domain/product/entity"

func RegisterAllModels() {
    modelRegistry = []interface{}{
        // ... existing models
        &productEntity.Product{}, // Add this line
    }
}
```

### Step 10: Create Tests

**📁 tests/product_test.go**

```go
package tests

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/stretchr/testify/assert"
)

func TestCreateProduct(t *testing.T) {
    router := setupTestRouter()
    
    product := map[string]interface{}{
        "name":     "Test Product",
        "sku":      "TEST001",
        "price":    99.99,
        "stock":    100,
        "category": "Electronics",
    }
    
    body, _ := json.Marshal(product)
    req := httptest.NewRequest("POST", "/api/v1/admin/products", bytes.NewBuffer(body))
    req.Header.Set("Authorization", "Bearer "+getTestAdminToken())
    req.Header.Set("Content-Type", "application/json")
    
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusCreated, w.Code)
    
    var response map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.True(t, response["success"].(bool))
}
```

---

## Quick Copy Template

Create folder structure:
```bash
FEATURE=product  # Change to your feature name
mkdir -p internal/domain/$FEATURE/{entity,repository}
touch internal/domain/$FEATURE/entity/$FEATURE.go
touch internal/domain/$FEATURE/repository/$FEATURE.go
touch internal/application/dto/$FEATURE.go
touch internal/application/services/${FEATURE}_service.go
touch internal/infrastructure/database/${FEATURE}_repository.go
touch internal/interfaces/http/handlers/${FEATURE}_handler.go
touch tests/${FEATURE}_test.go
```

---

## Checklist

- [ ] Domain Entity created
- [ ] Repository Interface defined
- [ ] Repository Implementation done
- [ ] DTO with validation
- [ ] Service with business logic
- [ ] Handler for HTTP endpoints
- [ ] Container updated with dependencies
- [ ] Routes added
- [ ] Model registered for migration
- [ ] Tests written

---

## Common Mistakes to Avoid

1. **Forgot to validate input** - Always use binding tags and ValidateStruct
2. **No error handling** - Check and return clear errors
3. **Forgot AdminID** - Admin-created features need AdminID
4. **No pagination** - List endpoints must have pagination
5. **Forgot to register model** - Must add to db/model.go
6. **No transaction** - Multi-table operations need transactions