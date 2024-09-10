"use client"
import { StrictMode, useEffect, useState } from 'react'
import { createRoot } from 'react-dom/client';
import {
  createBrowserRouter,
  Outlet,
  RouterProvider,
  useNavigate,
  useLocation,
} from "react-router-dom";
import './index.css';
import InstructionsPage from './pages/instructions';
import LoginPage from './pages/login-page';
import TestsPage from './pages/tests';
import EndPage from './pages/end';
import * as server from "@common/server.ts";
import * as types from "@common/types.ts";
import { Alert, AlertDescription, AlertTitle } from './components/ui/alert';
import { TestProvider } from '@/components/TestContext';

function WebSocketHandler() {
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const navigate = useNavigate();

  // TODO: some kinda progress bar of these
  let timeout = 0;
  const setErrorMessageDeffered = (msg: string, tout = 6000) => {
    clearTimeout(timeout);

    setErrorMessage(msg);

    // @ts-ignore
    timeout = setTimeout(() => {
      setErrorMessage(null);
    }, tout);
  };

  useEffect(() => {
    let disable: (() => PromiseLike<void>)[] = [];

    server.server.add_callback(types.Varient.LoadRoute, async (res) => {
      console.log(res);
      navigate(res.Route)
    }).then(d => {
      disable.push(d);
    });

    server.server.add_callback(types.Varient.Err, async (res) => {
      console.error('Error from server:', res.Message);
      setErrorMessageDeffered(res.Message);  
    }).then(d => {
      disable.push(d);
    });

    return () => {
      for (let fn of disable) {
        fn();
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
        path: "/tests",
        element: <TestProvider><TestsPage /></TestProvider>
      },
      {
        path: "/tests/:testId",
        element: <TestProvider><TestsPage /></TestProvider>
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
