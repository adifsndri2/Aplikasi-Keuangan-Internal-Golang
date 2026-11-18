import { createContext, useCallback, useContext, useEffect ,useState } from "react";
import {Navigate} from "react-router-dom";
import {json, setToken, setOnUnauthorized} from "./api";

const AuthCtx = createContext(null)
export const useAuth = () => useContext(AuthCtx)

export function AuthProvider({children}) {
    const [role, setRole] = useState(null)
    const logout = useCallback(() => {setToken(null); setRole(null)}, [])
    useEffect(() => {setOnUnauthorized(logout)}, [logout])

    async function login(email, password) {
        const d = await json('/auth/login', {method:'POST', body:{email, password}})
        setToken(d.token)
        setRole(d.role)
    }
    return <AuthCtx.Provider value = {{role, login, logout}}>{children}
    </AuthCtx.Provider>
}

export function Protected({children, admin}) {
    const {role} = useAuth()
    if (!role) return <Navigate to="/login" replace/>
    if (admin && role !== 'ADMIN') return <Navigate to="/" replace/>
    return children
}
