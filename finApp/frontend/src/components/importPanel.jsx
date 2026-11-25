import { useState } from "react";
import { json, download } from "../api";
import { useAuth } from "../auth";
import { rupiah } from "../format";
import Badge from "./badge";
import Stat from "./Stat";


const Stat = ({n, t}) => (
    <div style = {{flex: 1, border : '1px solid  #E5E7EB', borderRadius : 8, padding : '12px 14px' }}>
        <div style = {{fontSize : 22, fontWeight : 600}}>{n}</div>
        <div className = 'mute'> {t} </div>
    </div>
)

export default function ImportPanel( {module, title, uploadLabel, flow, fileColumn, onCommited}) {
    const {role} = useAuth()
    const [prevew, setPreview] = useState(null)
    const [fileName, setFileName] = useState('')
    const {mode, setMode} = useState('append')
    const [done, setDone] = useState(null)
    const [error, setError] = useState('')
    const [bussy, setBussy] = useState(false)

    async function pick (e) {
        const file = e.target.files[0]
        e.target.value = ''
        if (!file) return
        setBussy (true); setError(''); setDone(null); setMode('append')
        try {
            const form = new FormData()
            form.append('file', file)
            setPreview(await json(`/${module}/import/preview`, {method:'POST',body : form }))

            setFileName(file.name)
        } catch (err) {
            setPrevew(null)
            setError (err.message)
        } finally {setBusy(false) }
    }

    async function commit() {
        if (mode === 'replace' && !window.confirm('Apakah anda yakin ingin mengganti semua data yang ada di database dengan data baru?')) return setBusy(true); setError('')
        try {
            setDone (await json(`/${module}/import/commit`, {method: 'POST', body: {batch_id: preview.batch_id, mode},})) 
        setPreview(null)
        onCommited ()
        } catch (err) { setError(err.massage) } finally { setBusy(false) }
    }

    const t = preview?.totals
    const rows = preview?.rows || []
    const errors = prevew?.errors || []
    const warnings = preview?.warnings || []
    const step = preview ? 2 : done ? 3 : 1
    const st = (n, label) => <span key = {n} className = {step === n ? '' : 'mute'} style = {{ fontWeight : step === n ? 600 : 400 }}> {n}. {label} </span>

    return (
        <>
            <div className = 'row'>
                <h1> {title} </h1>
                <div className='row'>
                    <button className = 'o' onClick = {() => download(`template-${module}.xlsx`).catch((e) => setError(e.message))}> Download Template </button>
                    <label className = 'btn p'>
                        {busy ? 'memproses...' : uploadLabel}
                        <input type = 'file' accept = '.xlsx' hidden onChange={pick} disabled = {busy}/>
                    </label>
                </div>
            </div>

            <div className = 'c row' style = {{ justifyContent : 'flex-start', gap : 24, fontSize : 13 }}>
                {[st(1, 'Upload'), st(2, 'Preview dan validasi'), st(3, 'Konfirmasi import')]}
            </div>
            {error && <div className = 'err'> {error} </div>}
            {done && (
                <div className = 'c'> Import Selesai : {done.inserted} baris masuk, {done.updated > 0 && `${done.updated} baris diperbarui, `} {done.skipped_duplicates} duplikat dilewati. 
                {done.mode === 'replace' && `, ${done.deleted} baris lama dihapus`}.
                </div>
                )}

            {preview && (
                <div className = 'c'>
                    <div className = 'row'> <h2>Preview Import</h2> <span className = 'mute'>Batch #{preview.batch_id} · {fileName}</span> </div>

                    <div className = 'row' style = {{ gap: 12}}>
                        <Stat n = {preview.valid} t ='Baris Valid' />
                        <Stat n = {preview.invalid} t ='Baris Invalid' />
                        <Stat n = {preview.already_in_db} t ='Sudah ada di database' />
                    </div>

                    <h2 style = {{ marginTop : 16 }}>Pencocokan total {flow}</h2>
                    <div className = 'row' style = {{ gap: 12}}>
                        <Stat labelFirst n = {rupiah(t.file)} t = {`Total ${flow} di file`}
                        />
                        <Stat labelFirst n = {rupiah(t.valid)} t = {`Total baris valid`} />
                        <Stat labelFirst n = {rupiah(t.rejected)} t = {`Total baris invalid`} 
                        />
                    </div>

                    <p className = {t.matched ? 'mute' : error}>
                        {t.matched ? 'Cocok : total di file sama dengan total baris valid ditambah baris invalid' : 'Tidak cocok : ada baris yang terlewat atau tidak terbaca. Periksa file sebelum di import '}
                    </p>

                    <h2 style = {{  marginTop : 16 }}> Isi File </h2>
                    <div className = 'scroll'>
                        <table>
                            <thead>
                                <tr>
                                    <th>Baris</th>
                                    {fileColumn.map((c) => <th key = {c.key} className = {c.right ? 'r' : ''}>{c.label}</th>)}
                                    <th> Status </th> 
                                    <th> Keterangan </th>
                                </tr>
                            </thead>
                            <tbody>
                                {rows.map((r) => (
                                    <tr key = {r.line}>
                                        <td> {r.line} </td>
                                        {fileColumn.map((c) => (
                                            <td key = {c.key} className = {c.right ? 'r' : ''}> {c.render ? c.render(r[c.key]) : r[c.key]} </td>
                                        ))}

                                        <td><Badge text = {r.status}/></td>
                                        <td> {r.note} </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                    {preview.rows_truncated && <p className = 'mute'> Hanya 500 baris pertama yang ditampilkan</p>}

                    <h2 style = {{ marginTop : 16}}> Error ({errors.length}) dan Warning {(warnings.length)}</h2>
                    <ul>
                        {[...errors, ...warnings].map((m,i) => <li key = {i}>{m}</li>)}
                    </ul>
                    <p className = 'mute'> Daftar error dibatasi 20 baris</p>

                    <h2 style = {{ marginTop : 16 }}> Mode Import </h2>
                    <label className = 'opt'>
                        <input type = 'radio' checked = {mode === 'append'} onChange = {() => setMode('append')} /> Tambahkan (data baru ditambah, duplikat dilewati)
                    </label>
                    <label className = 'opt'>
                        <input type = 'radio' checked = {mode === 'upsert'} onChange = {() => setMode('upsert')}/> Perbarui (referensi yang sama diperbarui, yang baru ditambahkan)
                    </label>

                    <label className = 'opt'>
                        <input type = 'radio' disabled = {role !== 'ADMIN'} checked = {mode === 'replace'} onChange = {() => setMode ('replace')} /> Ganti seluruh data  (hapus data lama, lalu isi dengan file ini) 
                        {role !== 'ADMIN' && 'Khusus admin'}
                    </label>
                    
                    <div className = 'row' style = {{ justifyContent : 'flex-start', marginTop : 16}}>
                        <button className = 'p' onClick = {commit} disabled = {busy || preview.valid === 0}> {busy ? 'memproses...' : 'Konfirmasi Import'} </button>

                        <button className = 'o' onClick = {() => {setPreview(null)}}> Batal </button>

                    </div>
                </div>
            )}
        </>
    )
}