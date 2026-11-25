const cls = {
    Paid : 'g', 
    posted : 'g', 
    'Not Paid' : 'y', 
    RFQ : 's', 
    'Purchase Order' : 's',
    ADMIN : 'g',
    USER : 's',
    admin : 'g',
    user : 's',
    Aktif : 'g',
    Nonaktif : 'y', 
    Paid : 'g', 
    posted : 'g', 
    Valid : 'g',
    Warning : 'y',
    Error : 'e',
}

export default function Badge({text}) {
    return <span className = {`b ${cls[text] || 's'}`}>{text || '-'}</span>
}

