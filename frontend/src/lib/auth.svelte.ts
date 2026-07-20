export const authState = $state({
    isAuthenticated: false,
    isChecking: true
});

export function checkAuth() {
    if (typeof sessionStorage !== 'undefined') {
        const token = sessionStorage.getItem('admin_token');
        authState.isAuthenticated = !!token;
    }
    authState.isChecking = false;
}

export function login(token: string) {
    if (typeof sessionStorage !== 'undefined') {
        sessionStorage.setItem('admin_token', token.trim());
    }
    authState.isAuthenticated = true;
}

export function logout() {
    if (typeof sessionStorage !== 'undefined') {
        sessionStorage.removeItem('admin_token');
    }
    authState.isAuthenticated = false;
}
