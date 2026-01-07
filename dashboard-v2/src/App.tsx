import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import DashboardLayout from "./layouts/DashboardLayout";
import UsersPage from "./pages/UsersPage";
import LabelsPage from "./pages/LabelsPage";
import ProvidersPage from "./pages/ProvidersPage";
import SessionsPage from "./pages/SessionsPage";
import LoginPage from "./pages/LoginPage";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/dashboard" element={<DashboardLayout />}>
          <Route index element={<Navigate to="/dashboard/users" replace />} />
          <Route path="users" element={<UsersPage />} />
          <Route path="labels" element={<LabelsPage />} />
          <Route path="providers" element={<ProvidersPage />} />
          <Route path="sessions" element={<SessionsPage />} />
        </Route>
        <Route path="/" element={<Navigate to="/dashboard/users" replace />} />
        <Route path="*" element={<div>404 Not Found</div>} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;