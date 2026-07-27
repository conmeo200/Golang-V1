/**
 * Global Validator & Ajax Form Sender
 * Quản lý validate form và gửi Ajax chung cho toàn bộ dự án
 */

class FormValidator {
    /**
     * @param {Object} options 
     * options = {
     *    formId: 'id-cua-form',
     *    rules: { fieldName: ['required', 'email'] },
     *    messages: { fieldName: { required: 'Không được để trống' } },
     *    useRecaptchaV3: 'SITE_KEY', // Truyền site key nếu dùng reCAPTCHA v3 (Invisible)
     *    onSuccess: function(response) {}
     * }
     */
    constructor(options) {
        this.form = document.getElementById(options.formId);
        if (!this.form) return;

        this.rules = options.rules || {};
        this.messages = options.messages || {};
        this.useRecaptchaV3 = options.useRecaptchaV3 || null;
        this.onSuccess = options.onSuccess || function() {};
        this.onError = options.onError || function(err) { console.error(err); };
        
        // Cấu hình URL và Method mặc định từ HTML form
        this.action = this.form.getAttribute('action') || '';
        this.method = (this.form.getAttribute('method') || 'POST').toUpperCase();

        this.initEvents();
    }

    initEvents() {
        // Lắng nghe sự kiện submit form
        this.form.addEventListener('submit', (e) => {
            e.preventDefault();
            this.handleSubmit();
        });

        // Xóa thông báo lỗi khi user bắt đầu nhập lại
        const inputs = this.form.querySelectorAll('input, textarea, select');
        inputs.forEach(input => {
            input.addEventListener('input', () => {
                this.clearError(input.name);
            });
            input.addEventListener('change', () => {
                this.clearError(input.name);
            });
        });
    }

    validate() {
        let isValid = true;
        const formData = new FormData(this.form);

        for (const [field, fieldRules] of Object.entries(this.rules)) {
            const value = formData.get(field) ? formData.get(field).toString().trim() : '';
            
            for (const rule of fieldRules) {
                let errorMessage = '';

                // Phân tích rule (vd: min:6)
                const ruleParts = rule.split(':');
                const ruleName = ruleParts[0];
                const ruleValue = ruleParts[1];

                switch(ruleName) {
                    case 'required':
                        if (!value) errorMessage = 'Trường này không được để trống.';
                        break;
                    case 'email':
                        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
                        if (value && !emailRegex.test(value)) errorMessage = 'Email không hợp lệ.';
                        break;
                    case 'min':
                        if (value && value.length < parseInt(ruleValue)) errorMessage = `Tối thiểu ${ruleValue} ký tự.`;
                        break;
                    case 'max':
                        if (value && value.length > parseInt(ruleValue)) errorMessage = `Tối đa ${ruleValue} ký tự.`;
                        break;
                    case 'match':
                        const matchValue = formData.get(ruleValue);
                        if (value !== matchValue) errorMessage = 'Giá trị không khớp.';
                        break;
                    case 'recaptcha':
                        if (!value) errorMessage = 'Vui lòng xác nhận bạn không phải người máy.';
                        break;
                }

                if (errorMessage) {
                    // Ưu tiên custom message nếu có
                    const customMsg = this.messages[field] && this.messages[field][ruleName];
                    this.showError(field, customMsg || errorMessage);
                    isValid = false;
                    break; // Ngừng check các rule tiếp theo của field này nếu đã lỗi
                }
            }
        }

        return isValid;
    }

    showError(fieldName, message) {
        // Đối với reCAPTCHA v2 (hidden textarea), ta báo lỗi ở container của nó
        if (fieldName === 'g-recaptcha-response') {
            const recaptchaBox = this.form.querySelector('.g-recaptcha');
            if (recaptchaBox) {
                let errorEl = recaptchaBox.parentElement.querySelector('.invalid-feedback');
                if (!errorEl) {
                    errorEl = document.createElement('div');
                    errorEl.className = 'invalid-feedback';
                    errorEl.style.color = '#ef4444';
                    errorEl.style.fontSize = '0.875rem';
                    errorEl.style.marginTop = '4px';
                    recaptchaBox.parentElement.appendChild(errorEl);
                }
                errorEl.textContent = message;
            }
            return;
        }

        const input = this.form.querySelector(`[name="${fieldName}"]`);
        if (!input) return;

        input.classList.add('is-invalid'); // Add class CSS lỗi (cần có trong css)
        input.style.borderColor = '#ef4444'; // Đỏ

        // Tìm hoặc tạo thẻ span hiển thị lỗi
        let errorEl = input.parentElement.querySelector('.invalid-feedback');
        if (!errorEl) {
            errorEl = document.createElement('div');
            errorEl.className = 'invalid-feedback';
            errorEl.style.color = '#ef4444';
            errorEl.style.fontSize = '0.875rem';
            errorEl.style.marginTop = '4px';
            input.parentElement.appendChild(errorEl);
        }
        errorEl.textContent = message;
    }

    clearError(fieldName) {
        if (fieldName === 'g-recaptcha-response') {
            const recaptchaBox = this.form.querySelector('.g-recaptcha');
            if (recaptchaBox) {
                const errorEl = recaptchaBox.parentElement.querySelector('.invalid-feedback');
                if (errorEl) errorEl.remove();
            }
            return;
        }

        const input = this.form.querySelector(`[name="${fieldName}"]`);
        if (!input) return;

        input.classList.remove('is-invalid');
        input.style.borderColor = ''; // Trả về mặc định

        const errorEl = input.parentElement.querySelector('.invalid-feedback');
        if (errorEl) {
            errorEl.textContent = '';
            errorEl.remove();
        }
    }

    async handleSubmit() {
        if (!this.validate()) {
            return; // Dừng lại nếu validate xịt
        }

        // Lấy dữ liệu dạng JSON
        const formData = new FormData(this.form);
        const data = Object.fromEntries(formData.entries());

        // Thay đổi trạng thái nút submit để tránh double-click spam
        const submitBtn = this.form.querySelector('button[type="submit"]');
        const originalBtnText = submitBtn ? submitBtn.innerHTML : '';
        if (submitBtn) {
            submitBtn.disabled = true;
            submitBtn.innerHTML = '<i class="ph ph-spinner ph-spin"></i> Đang xử lý...';
        }

        try {
            // Xử lý Google reCAPTCHA v3 (Invisible) nếu được cấu hình
            if (this.useRecaptchaV3 && typeof grecaptcha !== 'undefined') {
                const token = await new Promise((resolve) => {
                    grecaptcha.ready(() => {
                        grecaptcha.execute(this.useRecaptchaV3, {action: 'submit'}).then(resolve);
                    });
                });
                data['g-recaptcha-response'] = token;
            }

            // Gửi request
            const response = await fetchAjax(this.action, this.method, data);
            
            // Gọi callback thành công
            this.onSuccess(response);

        } catch (error) {
            // Gọi callback lỗi
            this.onError(error);
            showToast(error.message || 'Lỗi hệ thống hoặc kết nối. Vui lòng thử lại!', 'error');
        } finally {
            // Phục hồi nút submit
            if (submitBtn) {
                submitBtn.disabled = false;
                submitBtn.innerHTML = originalBtnText;
            }

            // Tự động reset reCAPTCHA v2 (nếu có trên form) để chống gửi lại mã cũ (Replay Attack spam)
            if (typeof grecaptcha !== 'undefined' && !this.useRecaptchaV3) {
                try {
                    grecaptcha.reset();
                } catch(e) {
                    console.log("reCAPTCHA v2 not found or already reset");
                }
            }
        }
    }
}

/**
 * Hàm gửi AJAX dùng chung cho toàn dự án bằng Fetch API
 * @param {string} url - Đường dẫn API
 * @param {string} method - 'GET', 'POST', 'PUT', 'DELETE'
 * @param {Object} data - Payload data (sẽ tự chuyển thành query params nếu là GET)
 * @returns Promise
 */
async function fetchAjax(url, method = 'POST', data = null) {
    const options = {
        method: method,
        headers: {
            'Content-Type': 'application/json',
            'Accept': 'application/json'
        }
    };

    if (data) {
        if (method === 'GET') {
            const params = new URLSearchParams(data).toString();
            url += `?${params}`;
        } else {
            options.body = JSON.stringify(data);
        }
    }

    const response = await fetch(url, options);

    // Bắt các lỗi HTTP thông thường
    if (!response.ok) {
        let errorData = {};
        try {
            errorData = await response.json();
        } catch(e) {}
        throw new Error(errorData.message || `Lỗi HTTP: ${response.status}`);
    }

    // Trả về JSON Data
    return await response.json();
}

/**
 * Global Toast Notification cho các phản hồi từ Server
 */
function showToast(message, type = 'success') {
    const toast = document.createElement('div');
    toast.style.position = 'fixed';
    toast.style.bottom = '20px';
    toast.style.right = '20px';
    toast.style.backgroundColor = type === 'success' ? '#10B981' : '#ef4444';
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

    const icon = type === 'success' ? 'ph-check-circle' : 'ph-warning-circle';

    toast.innerHTML = `
        <i class="ph-fill ${icon}" style="font-size: 24px;"></i>
        <div>
            <div style="font-weight: 600; font-size: 14px;">${message}</div>
        </div>
    `;

    document.body.appendChild(toast);

    requestAnimationFrame(() => {
        toast.style.transform = 'translateY(0)';
        toast.style.opacity = '1';
    });

    setTimeout(() => {
        toast.style.transform = 'translateY(100px)';
        toast.style.opacity = '0';
        setTimeout(() => toast.remove(), 300);
    }, 3000);
}
