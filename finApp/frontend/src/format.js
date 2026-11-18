export const rupiah = (v) => {
    const n = number(v)
    return (n < 0 ? '-' : '') + 'Rp ' + Math.abs(n).toLocaleString('id-ID')
}

export const tgl = (s) => new Date(s).toLocaleDateString('id-ID', {day:'2-digit', month:'short', year:'numeric'})
