import React from 'react';
import { useAuth } from '../hooks/useAuth';
import { Link } from 'react-router-dom';

export const Home: React.FC = () => {
  const { isAuthenticated, user } = useAuth();

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="max-w-4xl mx-auto text-center">
        <h1 className="text-4xl font-bold text-gray-900 dark:text-white mb-6">
          Welcome to List of Maldives
        </h1>
        
        {isAuthenticated ? (
          <div className="space-y-4">
            <p className="text-xl text-gray-600 dark:text-gray-300">
              Hello, {user?.nickname || user?.email}! You are logged in.
            </p>
            <div className="space-x-4">
              <Link
                to="/protected"
                className="inline-block bg-blue-500 hover:bg-blue-600 text-white px-6 py-3 rounded-lg font-medium transition-colors"
              >
                Go to Protected Page
              </Link>
            </div>
          </div>
        ) : (
          <div className="space-y-4">
            <p className="text-xl text-gray-600 dark:text-gray-300 mb-8">
              Please log in or register to access all features.
            </p>
            <div className="space-x-4">
              <Link
                to="/login"
                className="inline-block bg-blue-500 hover:bg-blue-600 text-white px-6 py-3 rounded-lg font-medium transition-colors"
              >
                Login
              </Link>
              <Link
                to="/register"
                className="inline-block bg-green-500 hover:bg-green-600 text-white px-6 py-3 rounded-lg font-medium transition-colors"
              >
                Sign Up
              </Link>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};