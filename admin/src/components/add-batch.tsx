import React, { useEffect, useState } from 'react'
import { Card, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { Checkbox } from '@/components/ui/checkbox'
import { useToast } from '@/hooks/use-toast'
import { Test } from '@common/types'
import api from '@/lib/api'

export default function AddBatch() {

  const [batchName, setBatchName] = useState('')
  const [selectedTests, setSelectedTests] = useState<string[]>([]);
  const [availableTests, setAvailableTests] = useState<Test[]>([]);
  const { toast } = useToast()
  console.log("Available tests: ", availableTests)

  useEffect(() => {
    async function fetchTests() {
      try {
        const response = await api.get(`/test/get_all_tests`);
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
    console.log("Available tests: ", availableTests)
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

      await api.post(`${import.meta.env.SERVER_URL}/admin/add_batch`, {
        batchName,
        selectedTests
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
                {availableTests.map((test:Test) => {
                  if (test.TestName==="" || test.TestName===null || test.Id===undefined) {
                    return null;
                  }
                  return (
                    <div key={test.Id} className="flex items-center space-x-2">
                      <Checkbox
                        id={`test-${test.Id}`}
                        checked={selectedTests.includes(test.Id)}
                        onCheckedChange={() => handleTestSelection(test.Id!)}
                      />
                      <Label htmlFor={`test-${test.Id}`}>{test.TestName.charAt(0).toUpperCase() + test.TestName.slice(1)}</Label>
                    </div>
                  )
                })}
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

