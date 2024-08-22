import { StrictMode, useEffect } from 'react'
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

function WebSocketHandler() {
  const navigate = useNavigate();

  useEffect(() => {
    let disable: () => PromiseLike<void>;

    server.server.add_callback(types.Varient.LoadRoute, async (res) => {
      console.log(res);
      navigate(res.Route)
    }).then(d => {
      disable = d;
    });

    return () => {
      if (disable) {
        disable();
      }
    };
  }, [navigate]);

  return <Outlet />;
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
