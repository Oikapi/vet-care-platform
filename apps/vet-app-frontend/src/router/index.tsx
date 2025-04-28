// src/router.jsx
import { createBrowserRouter } from 'react-router-dom';
// import App from './App';
// import AuthLayout from './layouts/AuthLayout';
// import MainLayout from './layouts/MainLayout';
// import RegisterPage from './pages/RegisterPage';
// import { App } from 'antd';
import RegisterPage from '../pages/RegisterPage';
import App from '../App';
import AuthLayout from '../layouts/AuthLayout';
import LoginPage from '../pages/LoginPage';
// import LoginPage from './pages/LoginPage';
// import HomePage from './pages/HomePage';
// import ErrorPage from './pages/ErrorPage';

export const router = createBrowserRouter([
  {
    path: '/',
    element: <App />,
    // errorElement: <ErrorPage />,
    children: [
      {
        path: 'auth',
        element: <AuthLayout />,
        children: [
          {
            path: 'register',
            element: <RegisterPage />,
          },
          {
            path: 'login',
            element: <LoginPage />,
          },
        ],
      },
    ],
  },
]);
