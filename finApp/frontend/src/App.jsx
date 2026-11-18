import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom"
import { AuthProvider, Protected } from "./auth"
import Layout from "./components/layout"
import Login from "./pages/login"
import Dashboard from "./pages/dashboard"
import AccountPayable from "./pages/ap"
import AccountReceivable from "./pages/ar"
import Reports from "./pages/reports"
import Users from "./pages/users"
import { AuthProvider, Protected } from "./auth"

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Routes>
          <Route path = '/login' element = {<Login />} />
          <Route element = {<Protected><Layout /></Protected>}> 
            <Route path = '/ap' element = {<AccountPayable />} />
            <Route path = '/ar' element = {<AccountReceivable />} />
            <Route path = '/reports' element = {<Navigate to = '/reports/forecasting' replace />} />
            <Route path = 'reports/:kind' element = {<Reports />} />
            <Route path = '/users' element = {<Protected admin><Users /> </Protected>} />
          </Route>
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  )
}

function App() {
  return (
    <main className="app">
      <section className="card">
        <h1>FinApp</h1>
        <p>Aplikasi Keuangan Internal</p>
      </section>
    </main>
  )
}


