import { BrowserRouter as Router, Routes, Route, Navigate } from "react-router-dom"
import { ThemeProvider } from "./components/theme-provider"
import { Toaster } from "./components/ui/toaster"
import LoginPage from "./pages/LoginPage"
import SignupPage from "./pages/SignupPage"
import DashboardLayout from "./layouts/DashboardLayout"
import UsersPage from "./pages/dashboard/UsersPage"
import LabelsPage from "./pages/dashboard/LabelsPage"
import ProvidersPage from "./pages/dashboard/ProvidersPage"
import SessionsPage from "./pages/dashboard/SessionsPage"
import { AuthProvider } from "./contexts/AuthContext"
import ProtectedRoute from "./components/auth/ProtectedRoute"

function App() {
  return (
    <ThemeProvider attribute="class" defaultTheme="light" enableSystem disableTransitionOnChange>
      <AuthProvider>
        <Router>
          <Routes>
            <Route path="/" element={<Navigate to="/login" replace />} />
            <Route path="/login" element={<LoginPage />} />
            <Route path="/signup" element={<SignupPage />} />
            <Route
              path="/dashboard"
              element={
                <ProtectedRoute>
                  <DashboardLayout />
                </ProtectedRoute>
              }
            >
              <Route index element={<UsersPage />} />
              <Route path="users" element={<UsersPage />} />
              <Route path="labels" element={<LabelsPage />} />
              <Route path="providers" element={<ProvidersPage />} />
              <Route path="sessions" element={<SessionsPage />} />
            </Route>
          </Routes>
        </Router>
        <Toaster />
      </AuthProvider>
    </ThemeProvider>
  )
}

export default App
