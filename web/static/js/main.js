// Main Javascript interactions

document.addEventListener('DOMContentLoaded', () => {
    // Quantity Selector Logic
    const qtyInputs = document.querySelectorAll('.qty-input');
    
    qtyInputs.forEach(input => {
        const container = input.parentElement;
        const minusBtn = container.querySelector('.qty-btn:first-child');
        const plusBtn = container.querySelector('.qty-btn:last-child');
        
        if (minusBtn && plusBtn) {
            minusBtn.addEventListener('click', () => {
                let val = parseInt(input.value);
                if (val > 1) input.value = val - 1;
            });
            
            plusBtn.addEventListener('click', () => {
                let val = parseInt(input.value);
                if (val < 99) input.value = val + 1;
            });
        }
    });
});

// Toast notification for adding to cart
function addToCart(productId) {
    // Simple custom toast implementation
    const toast = document.createElement('div');
    toast.textContent = 'Đã thêm vào giỏ hàng thành công!';
    toast.style.position = 'fixed';
    toast.style.bottom = '30px';
    toast.style.right = '30px';
    toast.style.background = 'var(--secondary)';
    toast.style.color = 'white';
    toast.style.padding = '16px 24px';
    toast.style.borderRadius = 'var(--radius-full)';
    toast.style.boxShadow = 'var(--shadow-lg)';
    toast.style.fontWeight = '500';
    toast.style.zIndex = '1000';
    toast.style.transition = 'opacity 0.3s ease, transform 0.3s ease';
    toast.style.opacity = '0';
    toast.style.transform = 'translateY(20px)';
    
    document.body.appendChild(toast);
    
    // Animate in
    setTimeout(() => {
        toast.style.opacity = '1';
        toast.style.transform = 'translateY(0)';
    }, 10);
    
    // Animate out
    setTimeout(() => {
        toast.style.opacity = '0';
        toast.style.transform = 'translateY(20px)';
        setTimeout(() => {
            document.body.removeChild(toast);
        }, 300);
    }, 3000);
}
