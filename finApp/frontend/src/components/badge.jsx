const cls = {
    Paid : 'g', 
    posted : 'g', 
    'Not Paid' : 'y', 
    RFQ : 's', 
    'Purchase Order' : 's',
    ADMIN : 'g',
    USER : 's',
    Aktif : 'g',
    Nonaktif : 'y', 
}

export default function Badge({text}) {
    return <span className = {`b ${cls[text] || 's'}`}>{text || '-'}</span>
}