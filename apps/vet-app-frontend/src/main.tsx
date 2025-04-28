import React from 'react';
import ReactDOM from 'react-dom/client';
import { BrowserRouter, RouterProvider } from 'react-router-dom';
import { router } from './router';
import { App } from 'antd';

ReactDOM.createRoot(document.getElementById('root')!).render(
  <RouterProvider router={router} />
);
