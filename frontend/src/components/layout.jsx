import { NavLink, Outlet } from 'react-router-dom'
import { useAuth } from '../auth'

export default function Layout() {
    const {role, logout} = useAuth()
    const items = [['/', 'Dashboard'], ['/ap', 'Account Payable'], ['/ar', 'Account Receivable'], ['/reports', 'Financial Reports']]

    if (role === 'admin') items.push(['/users', 'Kelola User'])
        return (
            <div className = 'shell'>
                <aside className = 'side'>
                    <div className = 'brand'>Financial Management</div>
                    <nav>
                        {items.map(([to, label]) => (
                            <NavLink key = {to} to = {to} end = {to ==='/'} className = {({ isActive }) => (isActive ? 'on' : '')}>
                                {label}
                            </NavLink>
                        ))}
                    </nav>
                    
                    <div className = 'who'>
                        Login Sebagai {role === 'admin' ? 'Admin' : 'User'} <button onClick = {logout}>Keluar</button>
                    </div>
                    </aside>

                    <main className = 'main'> <Outlet /> </main>
                    </div>
        )
}