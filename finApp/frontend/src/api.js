let token = null
let onUnauthorized = () => {}
export const setToken = (t) => {
    token = t
}
export const setOnUnauthorized = (fn) => {
    onUnauthorized = fn
}

async function request(path, {method = 'GET', body, form} = {}) {
    const headers = {}

    if (token) headers.Authorization = `Bearer ${token}`
    let playload
    if (form) playload = form
    else if (body) {
        headers['Content-Type'] = 'application/json'
        playload = JSON.stringify(body)
    }

    const res = await fetch(`/api/v1${path}`, {method, headers, body: playload})
    if (res.status === 401 && path != 'auth/login') onUnauthorized()
        if (!res.ok) {
            const data = await res.json().catch(() => ({}))
            throw new Error(data.error || 'Terjadi kesalahan (${res.status})')
        } 
        return res
}

//Endpoint export dan template butuh header Authorization, jadi tidak bisa memakai link biasa.
export async function download (path, filename) {
    const res = await request(path)
    const url = URL.createObjectURL(await res.blob())
    const a = document.createElement('a')
    a.href = url
    a.download = filename
    a.click()
    URL.revokeObjectURL(url)
}

