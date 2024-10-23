"use client"
import { StrictMode, useEffect, useState } from 'react'
import { createRoot } from 'react-dom/client';
import {
  createHashRouter,
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
import OfflineToast from './components/offline-toast';
import { toast, useToast } from './hooks/use-toast';
import { Toaster } from './components/ui/toaster';
import { StateContextProvider } from './context/app-context';

function WebSocketHandler() {
  const navigate = useNavigate();

 const {toast} = useToast();

  useEffect(() => {
    let disable: (() => PromiseLike<void>)[] = [];

    server.server.add_callback(types.Varient.ExeNotFound, async (res) => {
      console.error('Error from server:', res.ErrMsg);
      toast({
        title: "Error",
        description: res.ErrMsg,
        variant:"destructive"
      })
    }).then(d => {
      disable.push(d);
    });

    server.server.add_callback(types.Varient.Notification, async (res) => {
      console.log('Notification from App:', res.Message);
      toast({
        title: "Notification",
        description: res.Message,
        variant:"destructive"
      })
    }).then(d => {
      disable.push(d);
    });

    server.server.add_callback(types.Varient.LoadRoute, async (res) => {
      console.log(res);
      navigate(res.Route)
    }).then(d => {
      disable.push(d);
    });

    server.server.add_callback(types.Varient.TestFinished, async (res) => {
      console.log(res);
      navigate("/end")
    }).then(d => {
      disable.push(d);
    });

    server.server.add_callback(types.Varient.Err, async (res) => {
      console.error('Error from server:', res.Message);
      toast({
        title: "Error",
        description: res.Message,
        variant:"destructive"
      })
    }).then(d => {
      disable.push(d);
    });

    server.server.add_callback(types.Varient.WarnUser, async (res) => {
      console.error('Warning to the user:', res.Message);
      toast({
        title: "Warning",
        description: res.Message,
        variant:"destructive"
      })  
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
    <StateContextProvider>
      <Outlet />
    </StateContextProvider>
  );
}

const router = createHashRouter([
  {
    path: "/",
    element: <WebSocketHandler />,
    children: [
      {
        path: "/",
        element: <LoginPage/>,
      },
      {
        path: "/instructions",
        element: <InstructionsPage />,
      },
      {
        path: "/tests",
        element: <TestsPage />
      },
      // TODO: add a quit button to this page. send common.Quit message on that button press
      {
        path: "/end",
        element: <EndPage />
      },
    ],
  },
]);


server.init().then(async () => {
  toast({
    title: "Connected",
    description: "Connected to the application server",
    variant: "default"
  })
  createRoot(document.getElementById('root')!).render(
    <StrictMode>
      <RouterProvider router={router} />
      <OfflineToast/>
      <Toaster/>
    </StrictMode>,
  );
});
