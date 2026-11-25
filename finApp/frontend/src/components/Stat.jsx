export default function Stat({ n, t, labelFirst }) {
    const num = <div style = {{ fontSize: 22, fontWeight: 600 }}>{n}</div>;
    const lab = <div className = "mute">{t}</div>;
    return (
        <div style = {{ flex: 1, border : '1px solid #E5E7EB', borderRadius: 8, padding: '12 14px',  }}>
            {labelFirst ? <>{lab} {num}</> : <>{num} {lab}</>}
        </div>
    )
}