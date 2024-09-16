import React, { useState } from 'react'
import { Card, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { Checkbox } from '@/components/ui/checkbox'
import { useToast } from '@/hooks/use-toast'
import axios from 'axios'

export default function AddBatch(){

  const [batchName, setBatchName] = useState('')
  const [selectedTests, setSelectedTests] = useState<string[]>([])
  const { toast } = useToast()

  const availableTests = [
    { id: '1', name: 'Typing Test 1' },
    { id: '2', name: 'Docx Test 1' },
    { id: '3', name: 'Excel Test 1' },
    { id: '4', name: 'Word Test 1' },
  ]

  const handleTestSelection = (testId: string) => {
    setSelectedTests(prev => 
      prev.includes(testId) 
        ? prev.filter(id => id !== testId)
        : [...prev, testId]
    )
  }

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    
    try {
      
      
      toast({
        title: "Batch added",
        description: "Successfully added the batch!",
      })
      
      setBatchName('')
      setSelectedTests([])
    } catch (error) {
      console.error('Error adding batch:', error)
      toast({
        variant: "destructive",
        title: "Failed to add",
        description: "Batch was not added due to an error on the server. Please try again later.",
      })
    }
  }
 

  return (
    <div className="w-full mx-auto p-4 space-y-6">
      <h1 className="text-3xl font-bold mb-8">Add New Batch</h1>
      <Card>
        <CardContent className="p-6 w-full">
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="batchName">Batch Name</Label>
              <Input
                type="text"
                id="batchName"
                value={batchName}
                onChange={(e) => setBatchName(e.target.value)}
                placeholder="Enter batch name"
                required
              />
            </div>

            <div className="space-y-2">
              <Label>Select Tests</Label>
              <div className="space-y-2">
                {availableTests.map((test) => (
                  <div key={test.id} className="flex items-center space-x-2">
                    <Checkbox
                      id={`test-${test.id}`}
                      checked={selectedTests.includes(test.id)}
                      onCheckedChange={() => handleTestSelection(test.id)}
                    />
                    <Label htmlFor={`test-${test.id}`}>{test.name}</Label>
                  </div>
                ))}
              </div>
            </div>

            <Button type="submit" className="w-full mt-4">
              Create Batch
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  )
};

