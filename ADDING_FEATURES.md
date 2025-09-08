# Adding New Features Guide

คู่มือการเพิ่ม Feature ใหม่แบบ Step-by-Step (ตัวอย่าง: Product Management)

## 📋 Overview - ต้องทำอะไรบ้าง?

เมื่อต้องการเพิ่ม feature ใหม่ (เช่น ระบบจัดการสินค้า) ต้องสร้าง 7 ไฟล์ และแก้ไข 3 ไฟล์:

```
สร้างใหม่:
✅ Domain Entity        - โมเดลข้อมูล
✅ Repository Interface - กำหนด operations
✅ Repository Impl      - เชื่อมต่อ database
✅ DTO                  - รับ-ส่งข้อมูล
✅ Service              - business logic
✅ Handler              - HTTP endpoints
✅ Tests                - ทดสอบการทำงาน

แก้ไข:
✅ Container            - dependency injection
✅ Routes               - เพิ่ม URL paths
✅ Model Registry       - ลงทะเบียน migration
```

## 🚀 Step-by-Step Guide

### Step 1: วางแผนโครงสร้างข้อมูล

ก่อนเขียน code ให้คิดว่า feature ต้องการอะไร:

```yaml
Product Management:
  ข้อมูล:
    - ชื่อสินค้า (name)
    - รหัสสินค้า (SKU)
    - ราคา (price)
    - จำนวนคงเหลือ (stock)
    - หมวดหมู่ (category)
    - รูปภาพ (image_url)
    - สถานะ (status)
    
  ฟังก์ชัน:
    - สร้างสินค้า
    - แก้ไขสินค้า
    - ลบสินค้า
    - ดูรายการสินค้า
    - ค้นหาสินค้า
    - อัปเดตสต็อก
```

### Step 2: สร้างโครงสร้างโฟลเดอร์

```bash
# สร้างโฟลเดอร์สำหรับ domain ใหม่
mkdir -p internal/domain/product/{entity,repository}
```

### Step 3: สร้าง Domain Entity

📁 **internal/domain/product/entity/product.go**

```go
package entity

import (
    "time"
    "gorm.io/gorm"
)

// Product represents a product in the system
type Product struct {
    ID          uint           `json:"id" gorm:"primaryKey"`
    Name        string         `json:"name" gorm:"size:255;not null"`
    SKU         string         `json:"sku" gorm:"size:100;uniqueIndex"` 
    Description string         `json:"description" gorm:"type:text"`
    Price       float64        `json:"price" gorm:"type:decimal(10,2);not null"`
    Stock       int            `json:"stock" gorm:"default:0"`
    Category    string         `json:"category" gorm:"size:100;index"`
    ImageURL    string         `json:"image_url" gorm:"size:500"`
    Status      string         `json:"status" gorm:"size:20;default:'active'"`
    AdminID     uint           `json:"admin_id" gorm:"not null;index"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
```

💡 **Tips:**
- ใช้ `gorm` tags กำหนดคุณสมบัติ column
- `uniqueIndex` สำหรับค่าที่ห้ามซ้ำ
- `DeletedAt` สำหรับ soft delete

### Step 4: สร้าง Repository Interface

📁 **internal/domain/product/repository/product.go**

```go
package repository

import "dashboard-starter/internal/domain/product/entity"

// ProductRepository defines product data operations
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
    GetLowStockProducts(threshold int) ([]*entity.Product, error)
}
```

💡 **Tips:**
- Interface ช่วยให้ test ง่าย (mock ได้)
- เพิ่ม method ตามความต้องการ business

### Step 5: สร้าง DTO (Data Transfer Object)

📁 **internal/application/dto/product.go**

```go
package dto

// ProductInput for creating/updating product
type ProductInput struct {
    Name        string  `json:"name" binding:"required" validate:"required,min=3,max=255"`
    SKU         string  `json:"sku" binding:"required" validate:"required,alphanum,min=3,max=100"`
    Description string  `json:"description" validate:"max=1000"`
    Price       float64 `json:"price" binding:"required" validate:"required,min=0"`
    Stock       int     `json:"stock" validate:"min=0"`
    Category    string  `json:"category" validate:"max=100"`
    ImageURL    string  `json:"image_url" validate:"omitempty,url"`
    Status      string  `json:"status" validate:"omitempty,oneof=active inactive discontinued"`
}

// StockUpdateInput for updating stock
type StockUpdateInput struct {
    Quantity  int    `json:"quantity" binding:"required" validate:"required"`
    Operation string `json:"operation" binding:"required" validate:"required,oneof=add subtract set"`
    Note      string `json:"note" validate:"max=500"`
}
```

💡 **Tips:**
- ใช้ `binding` สำหรับ Gin validation
- ใช้ `validate` สำหรับ custom validation
- แยก DTO ตาม use case

### Step 6: สร้าง Application Service

📁 **internal/application/services/product_service.go**

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

// CreateProduct creates a new product
func (s *ProductApplicationService) CreateProduct(input dto.ProductInput, adminID uint) (*entity.Product, error) {
    // Check duplicate SKU
    if existing, _ := s.productRepo.GetBySKU(input.SKU); existing != nil {
        return nil, fmt.Errorf("product with SKU %s already exists", input.SKU)
    }
    
    // Set defaults
    if input.Status == "" {
        input.Status = "active"
    }
    
    // Create entity
    product := &entity.Product{
        Name:        input.Name,
        SKU:         input.SKU,
        Description: input.Description,
        Price:       input.Price,
        Stock:       input.Stock,
        Category:    input.Category,
        ImageURL:    input.ImageURL,
        Status:      input.Status,
        AdminID:     adminID,
    }
    
    // Save to database
    if err := s.productRepo.Create(product); err != nil {
        return nil, fmt.Errorf("failed to create product: %w", err)
    }
    
    return product, nil
}

// UpdateStock updates product stock
func (s *ProductApplicationService) UpdateStock(id uint, input dto.StockUpdateInput) (*entity.Product, error) {
    product, err := s.productRepo.GetByID(id)
    if err != nil {
        return nil, fmt.Errorf("product not found")
    }
    
    // Calculate new stock
    switch input.Operation {
    case "add":
        product.Stock += input.Quantity
    case "subtract":
        if product.Stock < input.Quantity {
            return nil, errors.New("insufficient stock")
        }
        product.Stock -= input.Quantity
    case "set":
        product.Stock = input.Quantity
    }
    
    // Update database
    if err := s.productRepo.Update(product); err != nil {
        return nil, fmt.Errorf("failed to update stock: %w", err)
    }
    
    return product, nil
}

// ListProducts returns paginated products
func (s *ProductApplicationService) ListProducts(page, limit int, search, category, status string) ([]*entity.Product, int64, error) {
    offset := (page - 1) * limit
    return s.productRepo.List(offset, limit, search, category, status)
}
```

💡 **Tips:**
- Service จัดการ business logic
- ตรวจสอบ business rules ที่นี่
- Return error ที่ชัดเจน

### Step 7: สร้าง Repository Implementation

📁 **internal/infrastructure/database/product_repository.go**

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
    if err != nil {
        return nil, err
    }
    return &product, nil
}

func (r *productRepository) GetBySKU(sku string) (*entity.Product, error) {
    var product entity.Product
    err := r.db.Where("sku = ?", sku).First(&product).Error
    if err != nil {
        return nil, err
    }
    return &product, nil
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
    
    // Apply filters
    if search != "" {
        searchPattern := "%" + search + "%"
        query = query.Where("name LIKE ? OR sku LIKE ? OR description LIKE ?", 
            searchPattern, searchPattern, searchPattern)
    }
    
    if category != "" {
        query = query.Where("category = ?", category)
    }
    
    if status != "" {
        query = query.Where("status = ?", status)
    }
    
    // Count total
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // Get paginated results
    err := query.
        Offset(offset).
        Limit(limit).
        Order("created_at DESC").
        Find(&products).Error
        
    return products, total, err
}

func (r *productRepository) UpdateStock(id uint, quantity int) error {
    return r.db.Model(&entity.Product{}).
        Where("id = ?", id).
        Update("stock", gorm.Expr("stock + ?", quantity)).Error
}

func (r *productRepository) GetLowStockProducts(threshold int) ([]*entity.Product, error) {
    var products []*entity.Product
    err := r.db.Where("stock <= ? AND status = ?", threshold, "active").
        Find(&products).Error
    return products, err
}
```

💡 **Tips:**
- ใช้ parameterized queries ป้องกัน SQL injection
- ใช้ `gorm.Expr` สำหรับ atomic operations
- Order by สำหรับ consistent results

### Step 8: สร้าง HTTP Handler

📁 **internal/interfaces/http/handlers/product_handler.go**

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

// CreateProduct handles POST /admin/products
func (h *ProductHandler) CreateProduct(c *gin.Context) {
    // Get admin ID from context (set by middleware)
    adminID, exists := c.Get("admin_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, Response{
            Success: false,
            Error:   "Admin ID not found",
        })
        return
    }
    
    // Bind and validate input
    var input dto.ProductInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, Response{
            Success: false,
            Error:   "Invalid input: " + err.Error(),
        })
        return
    }
    
    // Additional validation
    if err := utils.ValidateStruct(input); err != nil {
        c.JSON(http.StatusBadRequest, Response{
            Success: false,
            Error:   err.Error(),
        })
        return
    }
    
    // Create product
    product, err := h.productService.CreateProduct(input, adminID.(uint))
    if err != nil {
        c.JSON(http.StatusBadRequest, Response{
            Success: false,
            Error:   err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusCreated, Response{
        Success: true,
        Data:    product,
    })
}

// ListProducts handles GET /admin/products
func (h *ProductHandler) ListProducts(c *gin.Context) {
    // Parse query parameters
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    search := c.Query("search")
    category := c.Query("category")
    status := c.Query("status")
    
    // Validate pagination
    if page < 1 {
        page = 1
    }
    if limit < 1 || limit > 100 {
        limit = 10
    }
    
    // Get products
    products, total, err := h.productService.ListProducts(page, limit, search, category, status)
    if err != nil {
        c.JSON(http.StatusInternalServerError, Response{
            Success: false,
            Error:   "Failed to retrieve products",
        })
        return
    }
    
    // Calculate total pages
    totalPages := (total + int64(limit) - 1) / int64(limit)
    
    c.JSON(http.StatusOK, Response{
        Success: true,
        Data:    products,
        Meta: gin.H{
            "page":       page,
            "limit":      limit,
            "total":      total,
            "totalPages": totalPages,
        },
    })
}

// UpdateStock handles PUT /admin/products/:id/stock
func (h *ProductHandler) UpdateStock(c *gin.Context) {
    // Get product ID
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, Response{
            Success: false,
            Error:   "Invalid product ID",
        })
        return
    }
    
    // Bind input
    var input dto.StockUpdateInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, Response{
            Success: false,
            Error:   "Invalid input: " + err.Error(),
        })
        return
    }
    
    // Update stock
    product, err := h.productService.UpdateStock(uint(id), input)
    if err != nil {
        c.JSON(http.StatusBadRequest, Response{
            Success: false,
            Error:   err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, Response{
        Success: true,
        Data:    product,
    })
}
```

### Step 9: เชื่อมต่อ Dependencies

📁 **internal/interfaces/http/container.go**

```go
// เพิ่มในฟังก์ชัน NewContainer()

import productRepo "dashboard-starter/internal/domain/product/repository"

// ใน struct Container เพิ่ม:
ProductHandler *handlers.ProductHandler

// ใน function NewContainer() เพิ่ม:
// Create product repository
productRepository := database.NewProductRepository(db.DB)

// Create product service
productService := services.NewProductApplicationService(productRepository)

// Create product handler
ProductHandler := handlers.NewProductHandler(productService)

// Return container พร้อม ProductHandler
```

### Step 10: เพิ่ม Routes

📁 **routes/routes.go**

```go
// เพิ่มใน admin group

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

### Step 11: ลงทะเบียน Model

📁 **db/model.go**

```go
import productEntity "dashboard-starter/internal/domain/product/entity"

func RegisterAllModels() {
    modelRegistry = []interface{}{
        // ... existing models
        &productEntity.Product{}, // เพิ่มบรรทัดนี้
    }
}
```

### Step 12: เขียน Tests (Optional แต่ควรทำ)

📁 **tests/product_test.go**

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
    
    // Prepare request
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
    
    // Execute request
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    // Assert response
    assert.Equal(t, http.StatusCreated, w.Code)
    
    var response map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.True(t, response["success"].(bool))
}
```

## ✅ Checklist - ตรวจสอบความสมบูรณ์

```
□ Domain Entity สร้างแล้ว
□ Repository Interface กำหนดแล้ว
□ Repository Implementation ทำงานได้
□ DTO validation ครบถ้วน
□ Service logic ถูกต้อง
□ Handler จัดการ error ครบ
□ Container เชื่อม dependencies
□ Routes เพิ่มแล้ว
□ Model ลงทะเบียนสำหรับ migration
□ Test cases ครอบคลุม
```

## 🎯 Quick Copy Templates

### สร้างโครงสร้างโฟลเดอร์
```bash
FEATURE=product  # เปลี่ยนเป็นชื่อ feature
mkdir -p internal/domain/$FEATURE/{entity,repository}
touch internal/domain/$FEATURE/entity/$FEATURE.go
touch internal/domain/$FEATURE/repository/$FEATURE.go
touch internal/application/dto/$FEATURE.go
touch internal/application/services/${FEATURE}_service.go
touch internal/infrastructure/database/${FEATURE}_repository.go
touch internal/interfaces/http/handlers/${FEATURE}_handler.go
touch tests/${FEATURE}_test.go
```

### Template สำหรับ Entity
```go
package entity

import (
    "time"
    "gorm.io/gorm"
)

type YourEntity struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    Name      string         `json:"name" gorm:"size:255;not null"`
    AdminID   uint           `json:"admin_id" gorm:"not null"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
```

## 🚨 Common Mistakes - ข้อผิดพลาดที่พบบ่อย

1. **ลืม validate input** - ใช้ binding tags และ ValidateStruct
2. **ไม่ handle error** - ตรวจสอบและ return error ที่ชัดเจน
3. **ลืมเพิ่ม AdminID** - features ที่ admin สร้างต้องมี AdminID
4. **ไม่ทำ pagination** - API ที่ return list ต้องมี pagination
5. **ลืม register model** - ต้องเพิ่มใน db/model.go
6. **ไม่ใช้ transaction** - operations ที่เกี่ยวข้องหลายตารางควรใช้ transaction

## 📚 Next Steps

- อ่าน [Development Guide](./DEVELOPMENT.md) สำหรับรายละเอียดเพิ่มเติม
- ดูตัวอย่างจาก Article, User, Device ที่มีอยู่
- เขียน test cases ให้ครอบคลุม
- ทำ API documentation ใน Swagger