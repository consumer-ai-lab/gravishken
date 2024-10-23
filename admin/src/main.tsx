import React from 'react'
import ReactDOM from 'react-dom/client'
import { createHashRouter, RouterProvider } from 'react-router-dom'
import './index.css'
import { Toaster } from "@/components/ui/toaster"
import Dashboard from './components/dashboard'
import Login from './components/login'
import LandingPage from './components/landing-page'

const router = createHashRouter([
  {
    path: "/dashboard",
    element: <Dashboard/>,
  },
  {
    path:'/login',
    element:<Login/>
  },
  {
    path:'/',
    element:<LandingPage/>
  }
])

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <RouterProvider router={router} />
    <Toaster />
  </React.StrictMode>,
)
