import React, { useEffect, useState } from 'react'
import { Card, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { Checkbox } from '@/components/ui/checkbox'
import { useToast } from '@/hooks/use-toast'
import axios from 'axios'

export default function AddBatch() {

  const [batchName, setBatchName] = useState('')
  const [selectedTests, setSelectedTests] = useState<string[]>([]);
  const [availableTests, setAvailableTests] = useState<any[]>([]);
  const { toast } = useToast()

  useEffect(() => {
    const fetchTests = async () => {
      try {
        const response = await axios.get(`${import.meta.env.SERVER_URL}/test/get_all_tests`);
        console.log("Response: ", response)
        setAvailableTests(response.data.tests);
      } catch (error) {
        console.error('Error fetching tests:', error);
        toast({
          variant: "destructive",
          title: "Failed to fetch tests",
          description: "Unable to load available tests. Please try again later."
        });
      }
    };
    fetchTests();
    console.log("Available tests: ",availableTests)
  }, []);

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

      await axios.post(`${import.meta.env.SERVER_URL}/admin/add_batch`, {
        batchName,
        selectedTests
      }, {
        withCredentials: true,
      });

      toast({
        title: "Batch added",
        description: "Successfully added the batch!",
      });

      setBatchName('');
      setSelectedTests([]);
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
                    <Label htmlFor={`test-${test.id}`}>{test.type.charAt(0).toUpperCase() + test.type.slice(1)} Test</Label>
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

