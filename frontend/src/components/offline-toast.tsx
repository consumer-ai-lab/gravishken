import  { useEffect } from 'react'
import { useToast } from '@/hooks/use-toast'

export default function OfflineToast()  {
  const { toast } = useToast()

  useEffect(() => {
    const handleOnline = () => {
      toast({
        title: "You're back online!",
        description: "Your internet connection has been restored.",
      })
    }

    const handleOffline = () => {
      toast({
        title: "You're offline",
        description: "Please check your internet connection.",
        duration: Infinity,
      })
    }

    window.addEventListener('online', handleOnline)
    window.addEventListener('offline', handleOffline)

    return () => {
      window.removeEventListener('online', handleOnline)
      window.removeEventListener('offline', handleOffline)
    }
  }, [toast])

  return null
}
