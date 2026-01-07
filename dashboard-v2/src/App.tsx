import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import DashboardLayout from "./layouts/DashboardLayout";
import UsersPage from "./pages/UsersPage";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/dashboard" element={<DashboardLayout />}>
          <Route index element={<Navigate to="/dashboard/users" replace />} />
          <Route path="users" element={<UsersPage />} />
          <Route path="labels" element={<div>Labels Page (Coming Soon)</div>} />
          <Route path="providers" element={<div>Providers Page (Coming Soon)</div>} />
          <Route path="sessions" element={<div>Sessions Page (Coming Soon)</div>} />
        </Route>
        <Route path="/" element={<Navigate to="/dashboard/users" replace />} />
        <Route path="*" element={<div>404 Not Found</div>} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;