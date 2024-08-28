"use client"
import { StrictMode, useEffect, useState } from 'react'
import { createRoot } from 'react-dom/client';
import {
  createBrowserRouter,
  Outlet,
  RouterProvider,
  useNavigate,
} from "react-router-dom";
import './index.css';
import InstructionsPage from './pages/instructions';
import LoginPage from './pages/login-page';
import TestsPage from './pages/tests';
import EndPage from './pages/end';
import * as server from "@common/server.ts";
import * as types from "@common/types.ts";
import { Alert, AlertDescription, AlertTitle } from './components/ui/alert';

function WebSocketHandler() {
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    let disable: () => PromiseLike<void>;

    server.server.add_callback(types.Varient.LoadRoute, async (res) => {
      console.log(res);
      navigate(res.Route)
    }).then(d => {
      disable = d;
    });

    server.server.add_callback(types.Varient.Err, async (res) => {
      console.error('Error from server:', res.Message);
      setErrorMessage(res.Message);  
    }).then(d => {
      disable = d;
    });

    return () => {
      if (disable) {
        disable();
      }
    };
  }, [navigate]);

  return (
    <div>
      {errorMessage && (
        <Alert variant="destructive">
          <AlertTitle>Error</AlertTitle>
          <AlertDescription>{errorMessage}</AlertDescription>
        </Alert>
      )}
      <Outlet />
    </div>
  );
}

const router = createBrowserRouter([
  {
    path: "/",
    element: <WebSocketHandler />,
    children: [
      {
        path: "/",
        element: <LoginPage />,
      },
      {
        path: "/instructions",
        element: <InstructionsPage />,
      },
      {
        path: "/tests/:testId",
        element: <TestsPage />
      },
      {
        path: "/end",
        element: <EndPage />
      }
    ],
  },
]);


server.init().then(async () => {
  createRoot(document.getElementById('root')!).render(
    <StrictMode>
      <RouterProvider router={router} />
    </StrictMode>,
  );
});
