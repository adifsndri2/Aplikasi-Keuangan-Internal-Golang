import { useState } from "react";
import { useAuth } from "../auth";
import { useNavigate } from "react-router-dom";

export default function Login() {
    const {login} = useAuth()
    const nav = useNavigate()
    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')
    const [error, setError] = useState('')

    async function submit(e) {
        e.PreventDefault()
        setError('')
        try {
            await login(email, password)
            nav('/', {replace : true})
        } catch (err) {
            setError(err.message)
        }
    }
    return (
        <div style = {{minHeight: '100vh', display: 'grid', placeItems: 'center'}}>
            <form className = 'c' style = {{width: '300px'}} onSubmit = {submit}>
                <h1> Masuk </h1>
                <p className = 'mute'> Sistem Keuangan Internal </p>
                <div style = {{display : 'flex', flexDirection: 'column', gap: '14px'}}>
                    <div><label> Email </label> <input type="email" value={email} onChange={(e) => setEmail(e.target.value)} required /></div>

                    <div><label> Password </label> <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} required /></div>
                    {error && <div className = 'err'>{error}</div>}
                    <button className = 'p' style = {{justifyContent: 'center', height: '44px'}}> Masuk </button>
                </div>
            </form>
        </div>
    )
}