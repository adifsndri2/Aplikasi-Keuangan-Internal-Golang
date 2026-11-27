export const number = (v) => {
    const n = typeof v === 'number' ? v : Number(v)
    return Number.isFinite(n) ? n : 0
}

export const rupiah = (v) => {
    const n = number(v)
    return (n < 0 ? '-' : '') + 'Rp ' + Math.abs(n).toLocaleString('id-ID')
}

export const tgl = (s) => new Date(s).toLocaleDateString('id-ID', {day:'2-digit', month:'short', year:'numeric'})
