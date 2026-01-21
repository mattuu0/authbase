/**
 * AuthBase Shared Logic
 */

document.addEventListener('DOMContentLoaded', () => {
    // フォーム送信時のローディング表示
    const authForms = document.querySelectorAll('form');
    authForms.forEach(form => {
        form.addEventListener('submit', (e) => {
            const submitBtn = form.querySelector('button[type="submit"]');
            if (submitBtn) {
                const originalContent = submitBtn.innerHTML;
                submitBtn.disabled = true;
                submitBtn.innerHTML = '<div class="spinner"></div><span>処理中...</span>';
            }
        });
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
