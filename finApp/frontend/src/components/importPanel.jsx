import { useState } from "react";
import { json, download } from "../api";

const Stat = ({n, t}) => (
    <div style = {{flex: 1, border : '1px solid  #E5E7EB', borderRadius : 8, padding : '12px 14px' }}>
        <div style = {{fontSize : 22, fontWeight : 600}}>{n}</div>
        <div className = 'mute'> {t} </div>
    </div>
)

export default function ImportPanel( {module, title, uploadLabel, onCommited}) {
    const [prevew, setPreview] = useState(null)
    const [fileName, setFileName] = useState('')
    const [done, setDone] = useState(null)
    const [error, setError] = useState('')
    const [bussy, setBussy] = useState(false)

    async function pick (e) {
        const file = e.target.files[0]
        e.target.value = ''
        if (!file) return
        setBussy (true); setError(''); setDone(null)
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
        setBust(true); setError('')
        try {
            setDone(await json(`/{module}/import/commit`, {method:'POST', body: {batch_id: preview.batch_id} }))
            setPreview(null)
            onCommited()
        } catch (err) { setError(err.message) } finally { setBusy(false) }
    }

    const errors = prevew?.errors || []
    const warnings = preview?.warnings || []
    const step = preview ? 2 : done ? 3 : 1
    const st = (n, t) => <span key = {n} className = {step === n ? '' : 'mute'} style = {{ fontWeight : step === n ? 600 : 400 }}> {n}. {t} </span>

    return (
        <>
            <div className = 'row'>
                <h1> {title} </h1>
                <div className='row'>
                    <button className = 'o' onClick = {() => download(`template-${module}.xlsx`).catch((e) => setError(e.message))}> Download Template </button>
                    <label className = 'btn p'>
                        {busy ? 'memproses...' : uploadLabel}
                        <input type = 'file' accept = '.xlsx' hidden onChange={pick} disabled = {busy} />
                    </label>
                </div>
            </div>

            <div className = 'c row' style = {{ justifyContent : 'flex-start', gap : 24, fontSize : 13 }}>
                {[st(1, 'Upload'), st(2, 'Preview dan validasi'), st(3, 'Konfirmasi import')]}
            </div>
            {error && <div className = 'err'> {error} </div>}
            {done && <div className = 'c'> Import Selesai : {done.inserted} baris masuk, {done.skipped_duplicates} duplikat dilewati. </div>}
            {preview && (
                <div className = 'c'>
                    <div className = 'row'> <h2>Preview Import</h2> <span className = 'mute'>Batch #{preview.batch_id} . {fileName}</span> </div>

                    <div className = 'row' style = {{ gap: 12}}>
                        <Stat n = {preview.valid} t ='Baris Valid' />
                        <Stat n = {preview.invalid} t ='Baris Invalid' />
                        <Stat n = {preview.already_in_db} t ='Sudah ada di database' />
                    </div>

                    <h2 style = {{ marginTop : 16 }}>Error ({errors.length})</h2>
                    {errors.length === 0 ? <p className = 'mute'> Tidak ada error.</p> :<ul>{errors.map((m,i) => <li key = {i}>{m}</li>)} </ul>}

                    <h2 style = {{  marginTop : 16 }}> Warning ({warnings.length}) </h2>
                    <ul>{warnings.map((m,i) => <li key = {i}>{m}</li>)} </ul>

                    <p className = 'mute'> Daftar error dibatasi 20 baris</p>

                    <div className = 'row' style = {{ justifyContent : 'flex-start', marginTop : 16}}>
                        <button className = 'p' onClick = {commit} disabled = {busy || preview.valid === 0}> {busy ? 'memproses...' : 'Konfirmasi Import'} </button>

                        <button className = 'o' onClick = {() => {setPreview(null)}}> Batal </button>

                    </div>
                </div>
            )}
        </>
    )
}