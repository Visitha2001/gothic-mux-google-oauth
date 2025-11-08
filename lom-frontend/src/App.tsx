import React, { useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { useThemeStore } from './hooks/useTheme';
import { useAuth } from './hooks/useAuth';
import { Header } from './components/Header';
import { LoginForm } from './components/LoginForm';
import { RegisterForm } from './components/RegisterForm';
import { Home } from './components/Home';
import { ProtectedRoute } from './components/ProtectedRoute';

function App() {
  const { theme } = useThemeStore();
  const { isLoading } = useAuth();

  useEffect(() => {
    document.documentElement.className = theme;
  }, [theme]);

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-white dark:bg-gray-900">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  return (
    <Router>
      <div className="min-h-screen bg-white dark:bg-gray-900 transition-colors">
        <Header />
        <main>
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/login" element={<LoginForm />} />
            <Route path="/register" element={<RegisterForm />} />
            <Route
              path="/protected"
              element={
                <ProtectedRoute>
                  <div className="container mx-auto px-4 py-8">
                    <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
                      Protected Page
                    </h1>
                    <p className="text-gray-600 dark:text-gray-300 mt-4">
                      This page is only accessible to authenticated users.
                    </p>
                  </div>
                </ProtectedRoute>
              }
            />
          </Routes>
        </main>
      </div>
    </Router>
  );
}

export default App;