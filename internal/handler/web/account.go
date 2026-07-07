package web

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

type AccountHandler struct{}

func NewAccountHandler() *AccountHandler {
	return &AccountHandler{}
}

func (h *AccountHandler) User(c *gin.Context) {
	data := map[string]interface{}{
		"Title":      "Quản lý tài khoản",
		"ActiveMenu": "profile",
		"User": map[string]interface{}{
			"Username": "nguyenvana123",
			"Name":     "Nguyễn Văn A",
			"Email":    "nguyenvana***@gmail.com",
			"Phone":    "********89",
		},
		"RecentOrders": []map[string]interface{}{
			{"ID": "ORD-202606-001", "Date": "25/06/2026", "Total": "450.000đ", "Status": "Đang giao"},
			{"ID": "ORD-202606-002", "Date": "20/06/2026", "Total": "1.250.000đ", "Status": "Hoàn thành"},
			{"ID": "ORD-202606-003", "Date": "15/06/2026", "Total": "199.000đ", "Status": "Hoàn thành"},
		},
	}

	tmpl, err := template.ParseFiles("web/templates/layouts/base.html", "web/templates/components/sidebar_user.html", "web/templates/pages/account/user.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error parsing template")
		return
	}
	tmpl.ExecuteTemplate(c.Writer, "base", data)
}

func (h *AccountHandler) Seller(c *gin.Context) {
	data := map[string]interface{}{
		"Title":      "Kênh Người Bán",
		"ActiveMenu": "overview",
		"Seller": map[string]interface{}{
			"ShopName": "A Store Official",
		},
		"Stats": map[string]interface{}{
			"TotalProducts": 145,
			"Revenue":       "25.4Tr",
			"NewOrders":     12,
			"Followers":     "3.2k",
		},
		"Products": []map[string]interface{}{
			{"Name": "Áo Thun Nam Cổ Tròn Premium", "Category": "Thời trang Nam", "Price": "150.000đ", "Stock": 450, "Image": "https://placehold.co/400x400/E2E8F0/64748B?text=Product+Image"},
			{"Name": "Giày Thể Thao Nam Nữ Siêu Nhẹ", "Category": "Giày dép", "Price": "250.000đ", "Stock": 120, "Image": "https://placehold.co/400x400/E2E8F0/64748B?text=Product+Image"},
			{"Name": "Balo Thời Trang Đi Học", "Category": "Túi xách", "Price": "199.000đ", "Stock": 85, "Image": "https://placehold.co/400x400/E2E8F0/64748B?text=Product+Image"},
			{"Name": "Ốp Lưng Điện Thoại Trong Suốt", "Category": "Phụ kiện", "Price": "25.000đ", "Stock": 1500, "Image": "https://placehold.co/400x400/E2E8F0/64748B?text=Product+Image"},
		},
	}

	tmpl, err := template.ParseFiles("web/templates/layouts/base.html", "web/templates/components/sidebar_seller.html", "web/templates/pages/account/seller.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error parsing template")
		return
	}
	tmpl.ExecuteTemplate(c.Writer, "base", data)
}

func (h *AccountHandler) UserOrders(c *gin.Context) {
	data := map[string]interface{}{
		"Title":      "Đơn mua của tôi",
		"ActiveMenu": "orders",
		"User":       map[string]interface{}{"Name": "Nguyễn Văn A"},
		"Orders": []map[string]interface{}{
			{"ShopName": "A Store Official", "Status": "Hoàn thành", "Image": "https://placehold.co/400x400/E2E8F0/64748B?text=Product+Image", "ProductName": "Áo Thun Nam Cổ Tròn Premium", "Quantity": 2, "Price": "150.000đ", "Total": "300.000đ"},
			{"ShopName": "Giày Chính Hãng", "Status": "Đang giao", "Image": "https://placehold.co/400x400/E2E8F0/64748B?text=Product+Image", "ProductName": "Giày Thể Thao Nam Nữ", "Quantity": 1, "Price": "250.000đ", "Total": "250.000đ"},
		},
	}
	tmpl, err := template.ParseFiles("web/templates/layouts/base.html", "web/templates/components/sidebar_user.html", "web/templates/pages/account/user_orders.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error parsing template")
		return
	}
	tmpl.ExecuteTemplate(c.Writer, "base", data)
}

func (h *AccountHandler) SellerProducts(c *gin.Context) {
	data := map[string]interface{}{
		"Title":      "Quản lý Sản phẩm",
		"ActiveMenu": "products",
		"Seller":     map[string]interface{}{"ShopName": "A Store Official"},
		"Products": []map[string]interface{}{
			{"Name": "Áo Thun Nam Cổ Tròn Premium", "Category": "Thời trang Nam", "Price": "150.000đ", "Stock": 450, "Image": "https://placehold.co/400x400/E2E8F0/64748B?text=Product+Image"},
			{"Name": "Giày Thể Thao Nam Nữ Siêu Nhẹ", "Category": "Giày dép", "Price": "250.000đ", "Stock": 0, "Image": "https://placehold.co/400x400/E2E8F0/64748B?text=Product+Image"},
			{"Name": "Balo Thời Trang Đi Học", "Category": "Túi xách", "Price": "199.000đ", "Stock": 85, "Image": "https://placehold.co/400x400/E2E8F0/64748B?text=Product+Image"},
		},
	}
	tmpl, err := template.ParseFiles("web/templates/layouts/base.html", "web/templates/components/sidebar_seller.html", "web/templates/pages/account/seller_products.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error parsing template")
		return
	}
	tmpl.ExecuteTemplate(c.Writer, "base", data)
}

func (h *AccountHandler) SellerOrders(c *gin.Context) {
	data := map[string]interface{}{
		"Title":      "Quản lý Đơn hàng",
		"ActiveMenu": "orders",
		"Seller":     map[string]interface{}{"ShopName": "A Store Official"},
		"Orders": []map[string]interface{}{
			{"ID": "ORD-12345", "Customer": "Lê Văn B", "Date": "29/06/2026", "Total": "450.000đ", "Status": "Chờ xác nhận"},
			{"ID": "ORD-12346", "Customer": "Trần Thị C", "Date": "28/06/2026", "Total": "150.000đ", "Status": "Đang giao"},
			{"ID": "ORD-12347", "Customer": "Phạm Văn D", "Date": "25/06/2026", "Total": "1.200.000đ", "Status": "Hoàn thành"},
		},
	}
	tmpl, err := template.ParseFiles("web/templates/layouts/base.html", "web/templates/components/sidebar_seller.html", "web/templates/pages/account/seller_orders.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error parsing template")
		return
	}
	tmpl.ExecuteTemplate(c.Writer, "base", data)
}

func (h *AccountHandler) SellerRevenue(c *gin.Context) {
	data := map[string]interface{}{
		"Title":      "Doanh thu & Giao dịch",
		"ActiveMenu": "revenue",
		"Seller":     map[string]interface{}{"ShopName": "A Store Official"},
		"Wallets": map[string]interface{}{
			"Available": "12.500.000đ",
			"Pending":   "4.200.000đ",
			"Withdrawn": "20.000.000đ",
		},
		"Breakdown": map[string]interface{}{
			"GrossRevenue":       "13.500.000đ",
			"PlatformFee":        "-675.000đ",
			"PaymentFee":         "-270.000đ",
			"Tax":                "-202.500đ",
			"MarketingFee":       "-150.000đ",
			"NetRevenue":         "12.202.500đ",
		},
		"ChartData": []int{120, 200, 150, 400, 300, 500, 250}, // For 7 days
		"Transactions": []map[string]interface{}{
			{"ID": "TXN-001", "Type": "Tiền vào", "Description": "Thanh toán đơn hàng #ORD-12347", "Amount": "+1.200.000đ", "Date": "28/06/2026 14:30", "Color": "green"},
			{"ID": "TXN-002", "Type": "Tiền ra", "Description": "Phí chạy Shopee Ads tuần 4", "Amount": "-150.000đ", "Date": "27/06/2026 10:00", "Color": "red"},
			{"ID": "TXN-003", "Type": "Tiền ra", "Description": "Rút tiền về Vietcombank ***123", "Amount": "-5.000.000đ", "Date": "26/06/2026 09:15", "Color": "red"},
		},
	}
	tmpl, err := template.ParseFiles("web/templates/layouts/base.html", "web/templates/components/sidebar_seller.html", "web/templates/pages/account/seller_revenue.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error parsing template")
		return
	}
	tmpl.ExecuteTemplate(c.Writer, "base", data)
}

func (h *AccountHandler) UserNotifications(c *gin.Context) {
	data := map[string]interface{}{
		"Title":      "Thông báo",
		"ActiveMenu": "notifications",
		"User":       map[string]interface{}{"Name": "Nguyễn Văn A"},
		"Notifications": []map[string]interface{}{
			{"Title": "Đơn hàng đang giao", "Message": "Đơn hàng #ORD-123 đang trên đường giao đến bạn.", "Time": "2 giờ trước"},
			{"Title": "Khuyến mãi mới", "Message": "Bạn nhận được mã giảm giá 50k.", "Time": "1 ngày trước"},
		},
	}
	tmpl, err := template.ParseFiles("web/templates/layouts/base.html", "web/templates/components/sidebar_user.html", "web/templates/pages/account/user_notifications.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error parsing template")
		return
	}
	tmpl.ExecuteTemplate(c.Writer, "base", data)
}

func (h *AccountHandler) UserAddresses(c *gin.Context) {
	data := map[string]interface{}{
		"Title":      "Địa chỉ nhận hàng",
		"ActiveMenu": "addresses",
		"User":       map[string]interface{}{"Name": "Nguyễn Văn A"},
		"Addresses": []map[string]interface{}{
			{"Name": "Nguyễn Văn A", "Phone": "0912345678", "Address": "Số 1 Đại Cồ Việt, Hai Bà Trưng, Hà Nội", "IsDefault": true},
			{"Name": "Nguyễn Văn A", "Phone": "0987654321", "Address": "Tòa nhà Landmark 81, Bình Thạnh, TP.HCM", "IsDefault": false},
		},
	}
	tmpl, err := template.ParseFiles("web/templates/layouts/base.html", "web/templates/components/sidebar_user.html", "web/templates/pages/account/user_addresses.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error parsing template")
		return
	}
	tmpl.ExecuteTemplate(c.Writer, "base", data)
}

func (h *AccountHandler) UserPassword(c *gin.Context) {
	data := map[string]interface{}{
		"Title":      "Đổi mật khẩu",
		"ActiveMenu": "password",
		"User":       map[string]interface{}{"Name": "Nguyễn Văn A"},
	}
	tmpl, err := template.ParseFiles("web/templates/layouts/base.html", "web/templates/components/sidebar_user.html", "web/templates/pages/account/user_password.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error parsing template")
		return
	}
	tmpl.ExecuteTemplate(c.Writer, "base", data)
}

func (h *AccountHandler) SellerSettings(c *gin.Context) {
	data := map[string]interface{}{
		"Title":      "Thiết lập Shop",
		"ActiveMenu": "settings",
		"Seller":     map[string]interface{}{"ShopName": "A Store Official"},
	}
	tmpl, err := template.ParseFiles("web/templates/layouts/base.html", "web/templates/components/sidebar_seller.html", "web/templates/pages/account/seller_settings.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error parsing template")
		return
	}
	tmpl.ExecuteTemplate(c.Writer, "base", data)
}
func (h *AccountHandler) SellerProductAdd(c *gin.Context) {
	data := map[string]interface{}{
		"Title":      "Thêm Sản Phẩm Mới",
		"ActiveMenu": "products",
		"Seller":     map[string]interface{}{"ShopName": "A Store Official"},
	}
	tmpl, err := template.ParseFiles("web/templates/layouts/base.html", "web/templates/components/sidebar_seller.html", "web/templates/pages/account/seller_product_add.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error parsing template")
		return
	}
	tmpl.ExecuteTemplate(c.Writer, "base", data)
}

func (h *AccountHandler) SellerProductAddPost(c *gin.Context) {
	// 1. Get form data
	name := c.PostForm("name")
	price := c.PostForm("price")
	stock := c.PostForm("stock")
	category := c.PostForm("category")
	description := c.PostForm("description")

	_ = name
	_ = price
	_ = stock
	_ = category
	_ = description

	// 2. Handle file upload
	file, err := c.FormFile("image")
	if err == nil {
		// Create directory if not exists
		uploadDir := "web/static/uploads"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			c.String(http.StatusInternalServerError, "Failed to create upload directory")
			return
		}

		// Save the file
		filename := filepath.Base(file.Filename)
		uploadPath := filepath.Join(uploadDir, filename)
		if err := c.SaveUploadedFile(file, uploadPath); err != nil {
			c.String(http.StatusInternalServerError, "Failed to save file")
			return
		}
	}

	// In a real application, we would save the product to the database here.
	// For now, we just redirect back to the product list.
	c.Redirect(http.StatusFound, "/seller/products")
}
