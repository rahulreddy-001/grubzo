import "./App.css";
import Home from "./pages/Home";
import { Routes, Route, Outlet, Navigate, useLocation } from "react-router";
import type { JSX } from "@emotion/react/jsx-runtime";

import Header from "./components/Header";
import Login from "./pages/Authentication/Login";
import Signup from "./pages/Authentication/Signup";
import { useAuth } from "./context/AuthProvider";
import NotFound from "./pages/NotFound";

const ProtectedRoute = ({ children }: { children: JSX.Element }) => {
  const { isAuthenticated, loading } = useAuth();
  const location = useLocation();

  if (loading) return <></>;

  if (!isAuthenticated) {
    return <Navigate to="/login" replace state={{ from: location }} />;
  }

  return children;
};

function PublicRoute({ children }: { children: JSX.Element }) {
  const { isAuthenticated, loading } = useAuth();
  const location = useLocation();
  if (loading) return <></>;
  if (isAuthenticated) {
    const redirectTo = (location.state as any)?.from?.pathname || "/";
    return <Navigate to={redirectTo} replace />;
  }
  return children;
}

function MainLayout() {
  return (
    <>
      <Header />
      <Outlet />
    </>
  );
}

function App() {
  const { loading } = useAuth();
  if (loading) return <></>;
  return (
    <div>
      <Routes>
        <Route element={<MainLayout />}>
          <Route path="/" element={<Home />} />
        </Route>
        <Route
          path="/login"
          element={
            <PublicRoute>
              <Login />
            </PublicRoute>
          }
        />
        <Route
          path="/signup"
          element={
            <PublicRoute>
              <Signup />
            </PublicRoute>
          }
        />
        <Route path="/404" element={<NotFound />} />
        <Route path="*" element={<NotFound />} />
      </Routes>
    </div>
  );
}

export default App;
