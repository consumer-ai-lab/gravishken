import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client';
import {
  createBrowserRouter,
  RouterProvider,
} from "react-router-dom";
import './index.css';
import InstructionsPage from './pages/instructions';
import LoginPage from './pages/login-page';
import TestsPage from './pages/tests';
import EndPage from './pages/end';

const router = createBrowserRouter([
  {
    path: "/",
    element: <InstructionsPage/>,
  },
  {
    path:"/login",
    element:<LoginPage/>,
  },
  {
    path:"/tests/:testId",
    element:<TestsPage/>
  },
  {
    path:"/end",
    element:<EndPage/>
  }
]);

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <RouterProvider router={router} />
  </StrictMode>,
)
