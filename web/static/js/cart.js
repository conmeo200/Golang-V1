class CartManager {
    constructor() {
        this.items = [];
        this.totalQuantity = 0;
        this.badgeElement = document.getElementById('cart-badge');
        
        // Cố gắng tải dữ liệu giỏ hàng từ localStorage nếu có
        this.loadCart();
        this.updateBadge();
    }

    // Tải giỏ hàng từ Local Storage
    loadCart() {
        const savedCart = localStorage.getItem('shopping_cart');
        if (savedCart) {
            try {
                this.items = JSON.parse(savedCart);
                this.calculateTotal();
            } catch (e) {
                console.error("Lỗi khi đọc giỏ hàng", e);
                this.items = [];
            }
        }
    }

    // Lưu giỏ hàng vào Local Storage
    saveCart() {
        localStorage.setItem('shopping_cart', JSON.stringify(this.items));
        this.calculateTotal();
        this.updateBadge();
    }

    // Tính tổng số lượng
    calculateTotal() {
        this.totalQuantity = this.items.reduce((sum, item) => sum + item.quantity, 0);
    }

    // Cập nhật giao diện badge
    updateBadge() {
        if (!this.badgeElement) return;

        let displayCount = this.totalQuantity;
        
        // Nếu số lượng > 9 thì hiển thị 9+
        if (this.totalQuantity > 9) {
            displayCount = '9+';
        }

        this.badgeElement.textContent = displayCount;
        
        // Thêm hiệu ứng rung/pop cho badge khi cập nhật
        this.badgeElement.style.transform = 'scale(1.5)';
        setTimeout(() => {
            this.badgeElement.style.transform = 'scale(1)';
        }, 200);
    }

    // Thêm sản phẩm vào giỏ
    addToCart(product) {
        // Kiểm tra xem sản phẩm đã có trong giỏ chưa
        const existingItem = this.items.find(item => item.id === product.id);
        
        if (existingItem) {
            existingItem.quantity += (product.quantity || 1);
        } else {
            this.items.push({
                id: product.id,
                name: product.name,
                price: product.price,
                image: product.image,
                quantity: product.quantity || 1
            });
        }

        this.saveCart();
        this.showToastEffect(product);
    }

    // Hiển thị hiệu ứng Toast thông báo
    showToastEffect(product) {
        // Tạo element toast
        const toast = document.createElement('div');
        toast.style.position = 'fixed';
        toast.style.bottom = '20px';
        toast.style.right = '20px';
        toast.style.backgroundColor = '#10B981'; // Xanh lá mạ
        toast.style.color = 'white';
        toast.style.padding = '12px 24px';
        toast.style.borderRadius = '8px';
        toast.style.boxShadow = '0 4px 12px rgba(0,0,0,0.15)';
        toast.style.display = 'flex';
        toast.style.alignItems = 'center';
        toast.style.gap = '12px';
        toast.style.zIndex = '9999';
        toast.style.transform = 'translateY(100px)';
        toast.style.opacity = '0';
        toast.style.transition = 'all 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275)';

        // Nội dung toast
        toast.innerHTML = `
            <i class="ph-fill ph-check-circle" style="font-size: 24px;"></i>
            <div>
                <div style="font-weight: 700; font-size: 14px;">Thêm thành công!</div>
                <div style="font-size: 12px; opacity: 0.9;">${product.name || 'Sản phẩm'} đã vào giỏ hàng.</div>
            </div>
        `;

        document.body.appendChild(toast);

        // Kích hoạt animation hiện
        requestAnimationFrame(() => {
            toast.style.transform = 'translateY(0)';
            toast.style.opacity = '1';
        });

        // Tự động ẩn sau 3 giây
        setTimeout(() => {
            toast.style.transform = 'translateY(100px)';
            toast.style.opacity = '0';
            setTimeout(() => {
                toast.remove();
            }, 300); // Đợi animation xong rồi xóa
        }, 3000);
    }
}

// Khởi tạo instance global để có thể gọi từ HTML
const cartManager = new CartManager();

// Sửa lại hàm addToCart global nếu trước đây bạn dùng nó trên các thẻ HTML
window.addToCart = function(productId, productName = 'Sản phẩm', productPrice = 0, productImage = '') {
    cartManager.addToCart({
        id: productId,
        name: productName,
        price: productPrice,
        image: productImage,
        quantity: 1
    });
};

// Gắn event listener cho các nút có class .add-to-cart-btn (nếu có)
document.addEventListener('DOMContentLoaded', () => {
    const addButtons = document.querySelectorAll('.add-to-cart-btn');
    addButtons.forEach(btn => {
        btn.addEventListener('click', (e) => {
            e.preventDefault(); // Tránh reload trang nếu là thẻ a
            
            // Lấy thông tin từ data attributes
            const id = btn.dataset.id;
            const name = btn.dataset.name || 'Sản phẩm';
            const price = parseFloat(btn.dataset.price) || 0;
            const image = btn.dataset.image || '';
            
            cartManager.addToCart({ id, name, price, image, quantity: 1 });
        });
    });
});
