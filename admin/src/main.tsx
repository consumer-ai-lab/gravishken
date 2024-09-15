import React from 'react'
import ReactDOM from 'react-dom/client'
import { createBrowserRouter, RouterProvider } from 'react-router-dom'
import './index.css'
import { Toaster } from "@/components/ui/toaster"
import Layout from './components/Layout'
import Dashboard from './components/Dashboard'
import AddTest from './components/AddTest'
import AddBatch from './components/AddBatch'
import AddAllUsers from './components/AddAllUsers'
import UpdateUserData from './components/UpdateUserData'
import GetBatchwiseData from './components/GetBatchwiseData'
import IncreaseTestTime from './components/IncreaseTestTime'
import SetUserData from './components/SetUserData'
import Login from './components/login'

const router = createBrowserRouter([
  {
    path: "/",
    element: <Layout/>,
    children: [
      {
        path: "add-test",
        element: <AddTest />,
      },
      {
        path: "add-batch",
        element: <AddBatch />,
      },
      {
        path: "add-users",
        element: <AddAllUsers />,
      },
      {
        path: "update-user",
        element: <UpdateUserData />,
      },
      {
        path: "get-batchwise-data",
        element: <GetBatchwiseData />,
      },
      {
        path: "increase-test-time",
        element: <IncreaseTestTime />,
      },
      {
        path: "set-user-data",
        element: <SetUserData />,
      },
    ],
  },
  {
    path:'/login',
    element:<Login/>
  }
])

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <RouterProvider router={router} />
    <Toaster />
  </React.StrictMode>,
)
