function logoutAndReload() {
    localStorage.removeItem('token');
    sessionStorage.removeItem('justLoggedIn');
    window.location.reload();
}

/**
 * AuthBase Shared Logic
 */


async function handleAuthSubmit(e, endpoint) {
    e.preventDefault();
    const form = e.target;
    const formData = new FormData(form);
    const data = Object.fromEntries(formData.entries());

    const submitBtn = form.querySelector('button[type="submit"]');
    const originalContent = submitBtn.innerHTML;
    submitBtn.disabled = true;
    submitBtn.innerHTML = '<div class="spinner"></div><span>処理中...</span>';

    // エラーバッジをクリア
    const existingError = document.getElementById('error-message');
    if (existingError) existingError.remove();

    try {
        const response = await fetch(endpoint, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });

        const result = await response.json();

        if (response.ok) {
            localStorage.setItem('token', result.token);
            sessionStorage.setItem('justLoggedIn', 'true');
            window.location.reload();
        } else {
            showError(result.error || 'エラーが発生しました');
            submitBtn.disabled = false;
            submitBtn.innerHTML = originalContent;
        }
    } catch (err) {
        showError('通信エラーが発生しました');
        submitBtn.disabled = false;
        submitBtn.innerHTML = originalContent;
    }
}

function showError(message) {
    const errorDiv = document.createElement('div');
    errorDiv.className = 'error-badge';
    errorDiv.id = 'error-message';

    const icon = document.createElementNS("http://www.w3.org/2000/svg", "svg");
    icon.setAttribute("width", "18");
    icon.setAttribute("height", "18");
    icon.setAttribute("viewBox", "0 0 24 24");
    icon.setAttribute("fill", "none");
    icon.setAttribute("stroke", "currentColor");
    icon.setAttribute("stroke-width", "2");
    icon.setAttribute("stroke-linecap", "round");
    icon.setAttribute("stroke-linejoin", "round");
    icon.innerHTML = '<circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line>';

    const textSpan = document.createElement('span');
    textSpan.textContent = `エラーが発生しました: ${message}`;

    errorDiv.appendChild(icon);
    errorDiv.appendChild(textSpan);

    const container = document.getElementById('auth-card-main');
    const header = container.querySelector('.auth-header');
    header.after(errorDiv);
}

document.addEventListener('DOMContentLoaded', () => {
    // フォーム送信時のローディング表示
    const authForms = document.querySelectorAll('form');
    authForms.forEach(form => {
        if (!form.getAttribute('onsubmit')) {
            form.addEventListener('submit', (e) => {
                const submitBtn = form.querySelector('button[type="submit"]');
                if (submitBtn) {
                    const originalContent = submitBtn.innerHTML;
                    submitBtn.disabled = true;
                    submitBtn.innerHTML = '<div class="spinner"></div><span>処理中...</span>';
                }
            });
        }
    });

    // エラーメッセージのフェードイン効果
    const errorBadge = document.querySelector('.error-badge');
    if (errorBadge) {
        errorBadge.style.opacity = '0';
        errorBadge.style.transform = 'translateY(-10px)';
        setTimeout(() => {
            errorBadge.style.transition = 'all 0.3s ease-out';
            errorBadge.style.opacity = '1';
            errorBadge.style.transform = 'translateY(0)';
        }, 50);
    }
});


// ユーティリティ: クエリパラメータ取得
function getQueryParam(name) {
    const params = new URLSearchParams(window.location.search);
    return params.get(name);
}

// ユーティリティ: 相対ベースパスの解決 (必要に応じて)
function resolveStaticPath(path) {
    // 実行時のパス深度を計算することも可能だが、テンプレート側で制御するほうが確実
    return path;
}
