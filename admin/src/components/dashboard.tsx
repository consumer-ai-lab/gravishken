import { useEffect, useState } from 'react'
import { Button } from "@/components/ui/button"
import { PlusCircle, Users, FileSpreadsheet, Database, Menu } from 'lucide-react'
import { Sheet, SheetContent, SheetTrigger } from './ui/sheet'
import AddTest from './add-test'
import UserDetails from './user-details'
import AddUser from './add-user'
import AddBatch from './add-batch'
import { useNavigate } from 'react-router-dom'
import axios from 'axios'

export default function Dashboard() {
  const [activeSection, setActiveSection] = useState('userDetails');
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    const checkAuthStatus = async () => {
      try {
        const response = await axios.get(`${import.meta.env.SERVER_URL}/admin/auth-status`, {
          withCredentials: true,
        });
        
        setIsAuthenticated(response.data.isAuthenticated);
      } catch (err:any) {
        if(err.status === 401){
          navigate('/login');
        }
      }
    }
    checkAuthStatus();
  }, [])

  

  const renderContent = () => {
    switch (activeSection) {
      case 'userDetails':
        return <UserDetails isAuthenticated={isAuthenticated}/>
      case 'addTest':
        return <AddTest />
      case 'addUsers':
        return <AddUser />
      case 'createBatch':
        return <AddBatch />
      default:
        return null
    }
  }

  return (
    <div className="min-h-screen bg-white flex flex-col">
      <header className="bg-blue-600 text-white p-4 flex justify-between items-center">
        <div className="flex items-center">
          <img src="/WCL_LOGO.png" alt="WCL Logo" className="mr-2 h-10 w-10" />
          <h1 className="text-xl font-bold">WCL Admin Panel</h1>
        </div>
        <Sheet>
          <SheetTrigger asChild>
            <Button variant="ghost" className="text-primary-foreground md:hidden">
              <Menu />
            </Button>
          </SheetTrigger>
          <SheetContent side="left" className="w-[300px] sm:w-[400px]">
            <nav className="flex flex-col space-y-4 mt-4">
              <Button
                variant="ghost"
                className="w-full justify-start mb-2 text-blue-600 hover:bg-blue-100"
                onClick={() => setActiveSection('userDetails')}
              >
                <Database className="mr-2 h-4 w-4" /> View User Details
              </Button>
              <Button
                variant="ghost"
                className="w-full justify-start mb-2 text-blue-600 hover:bg-blue-100"
                onClick={() => setActiveSection('addTest')}
              >
                <PlusCircle className="mr-2 h-4 w-4" /> Add Test
              </Button>
              <Button
                variant="ghost"
                className="w-full justify-start mb-2 text-blue-600 hover:bg-blue-100"
                onClick={() => setActiveSection('addUsers')}
              >
                <Users className="mr-2 h-4 w-4" /> Add Users from CSV
              </Button>
              <Button
                variant="ghost"
                className="w-full justify-start text-blue-600 hover:bg-blue-100"
                onClick={() => setActiveSection('createBatch')}
              >
                <FileSpreadsheet className="mr-2 h-4 w-4" /> Create Batch
              </Button>
            </nav>
          </SheetContent>
        </Sheet>
      </header>

      <div className="flex flex-1">
        <nav className={`bg-gray-100 w-64 p-4 flex-shrink-0 hidden md:block pt-10`}>
          <Button
            variant="ghost"
            className="w-full justify-start mb-2 text-blue-600 hover:bg-blue-100"
            onClick={() => setActiveSection('userDetails')}
          >
            <Database className="mr-2 h-4 w-4" /> View User Details
          </Button>
          <Button
            variant="ghost"
            className="w-full justify-start mb-2 text-blue-600 hover:bg-blue-100"
            onClick={() => setActiveSection('addTest')}
          >
            <PlusCircle className="mr-2 h-4 w-4" /> Add Test
          </Button>
          <Button
            variant="ghost"
            className="w-full justify-start mb-2 text-blue-600 hover:bg-blue-100"
            onClick={() => setActiveSection('addUsers')}
          >
            <Users className="mr-2 h-4 w-4" /> Add Users from CSV
          </Button>
          <Button
            variant="ghost"
            className="w-full justify-start text-blue-600 hover:bg-blue-100"
            onClick={() => setActiveSection('createBatch')}
          >
            <FileSpreadsheet className="mr-2 h-4 w-4" /> Create Batch
          </Button>
        </nav>

        {/* Main Content */}
        <main className="flex-1 p-6">
          <div className="max-w-6xl mx-auto">
            {renderContent()}
          </div>
        </main>
      </div>
    </div>
  )
}