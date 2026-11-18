import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom"
import { AuthProvider, Protected } from "./auth"
import Layout from "./components/layout"
import Login from "./pages/login"
import Dashboard from "./pages/dashboard"
import AccountPayable from "./pages/ap"
import AccountReceivable from "./pages/ar"
import Reports from "./pages/report"
import Users from "./pages/user"

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Routes>
          <Route path = '/login' element = {<Login />} />
          <Route element = {<Protected><Layout /></Protected>}> 
            <Route path = '/' element = {<Dashboard />} />
            <Route path = '/ap' element = {<AccountPayable />} />
            <Route path = '/ar' element = {<AccountReceivable />} />
            <Route path = '/reports' element = {<Navigate to = '/reports/forecasting' replace />} />
            <Route path = '/reports/:kind' element = {<Reports />} />
            <Route path = '/users' element = {<Protected admin><Users /> </Protected>} />
          </Route>
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  )
}


