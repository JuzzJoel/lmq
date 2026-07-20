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

export async function login(token: string) {
    if (typeof sessionStorage !== 'undefined') {
        const encoder = new TextEncoder();
        const data = encoder.encode(token.trim());
        const hashBuffer = await crypto.subtle.digest('SHA-256', data);
        const hashArray = Array.from(new Uint8Array(hashBuffer));
        const hashHex = hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
        sessionStorage.setItem('admin_token', hashHex);
    }
    authState.isAuthenticated = true;
}

export function logout() {
    if (typeof sessionStorage !== 'undefined') {
        sessionStorage.removeItem('admin_token');
    }
    authState.isAuthenticated = false;
}
