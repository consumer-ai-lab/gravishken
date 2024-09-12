"use client";

import React, { useState } from 'react';
import { AlertCircle, Clock } from 'lucide-react';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group';
import { Textarea } from '@/components/ui/textarea';

const IncreaseTestTime = () => {
  const [param, setParam] = useState('user');
  const [username, setUsername] = useState('');
  const [timeToIncrease, setTimeToIncrease] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [result, setResult] = useState<{ status: string; message: string } | null>(null);

  const handleSubmit = async (e: any) => {
    e.preventDefault();
    setIsLoading(true);
    setResult(null);

    const payload = {
      param,
      username: username.split('\n').map(name => name.trim()).filter(name => name !== ''),
      time_to_increase: parseInt(timeToIncrease, 10)
    };

    try {
      const response = await fetch('http://localhost:8081/admin/increase_test_time', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
      });

      if (response.ok) {
        setResult({ status: 'success', message: 'Test time increased successfully!' });
      } else {
        const errorData = await response.json();
        setResult({ status: 'error', message: errorData.message || 'An error occurred' });
      }
    } catch (error: any) {
      setResult({ status: 'error', message: 'Network error: ' + error.message });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Card className="w-full max-w-2xl mx-auto">
      <CardHeader>
        <CardTitle className="text-2xl font-bold">Increase Test Time</CardTitle>
        <CardDescription>Adjust test duration for users or batches</CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit}>
          <div className="space-y-4">
            <RadioGroup defaultValue="user" onValueChange={setParam} className="flex space-x-4">
              <div className="flex items-center space-x-2">
                <RadioGroupItem value="user" id="user" />
                <Label htmlFor="user">User</Label>
              </div>
              <div className="flex items-center space-x-2">
                <RadioGroupItem value="batch" id="batch" />
                <Label htmlFor="batch">Batch</Label>
              </div>
            </RadioGroup>

            <div className="space-y-2">
              <Label htmlFor="username">
                {param === 'user' ? 'Username(s)' : 'Batch Name(s)'}
              </Label>
              <Textarea
                id="username"
                placeholder={`Enter ${param === 'user' ? 'username(s)' : 'batch name(s)'}, one per line`}
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                className="min-h-[100px]"
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="time">Time to Increase (seconds)</Label>
              <Input
                id="time"
                type="number"
                placeholder="e.g., 3600"
                value={timeToIncrease}
                onChange={(e) => setTimeToIncrease(e.target.value)}
              />
            </div>
          </div>
        </form>
      </CardContent>
      <CardFooter className="flex justify-between items-center">
        <Button type="submit" onClick={handleSubmit} disabled={isLoading}>
          {isLoading ? (
            <>
              <Clock className="mr-2 h-4 w-4 animate-spin" />
              Processing...
            </>
          ) : (
            'Increase Time'
          )}
        </Button>
        {result && (
          <Alert variant={result.status === 'success' ? 'default' : 'destructive'} className="mt-4">
            <AlertCircle className="h-4 w-4" />
            <AlertTitle>{result.status === 'success' ? 'Success' : 'Error'}</AlertTitle>
            <AlertDescription>{result.message}</AlertDescription>
          </Alert>
        )}
      </CardFooter>
    </Card>
  );
};

export default IncreaseTestTime;