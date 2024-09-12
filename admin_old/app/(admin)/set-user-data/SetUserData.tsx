"use client";

import React, { useState } from 'react';
import { AlertCircle, Loader2 } from 'lucide-react';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group';
import { Switch } from '@/components/ui/switch';

const SetUserData = () => {
  const [username, setUsername] = useState('');
  const [param, setParam] = useState('download');
  const [from, setFrom] = useState('');
  const [to, setTo] = useState('');
  const [resultDownloaded, setResultDownloaded] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [result, setResult] = useState<{ status: string, message: string } | null>(null);

  const handleSubmit = async (e: any) => {
    e.preventDefault();
    setIsLoading(true);
    setResult(null);

    const payload = {
      username,
      param,
      from: parseInt(from, 10),
      to: parseInt(to, 10),
      resultDownloaded
    };

    try {
      const response = await fetch('http://localhost:8081/admin/set_user_data', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
      });

      if (response.ok) {
        setResult({ status: 'success', message: 'User data set successfully!' });
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
        <CardTitle className="text-2xl font-bold">Set User Data</CardTitle>
        <CardDescription>Manage user data for download or reset</CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit}>
          <div className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="username">Username</Label>
              <Input
                id="username"
                placeholder="Enter username"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                required
              />
            </div>

            <RadioGroup defaultValue="download" onValueChange={setParam} className="flex space-x-4">
              <div className="flex items-center space-x-2">
                <RadioGroupItem value="download" id="download" />
                <Label htmlFor="download">Download</Label>
              </div>
              <div className="flex items-center space-x-2">
                <RadioGroupItem value="reset" id="reset" />
                <Label htmlFor="reset">Reset</Label>
              </div>
            </RadioGroup>

            <div className="space-y-2">
              <Label htmlFor="from">From</Label>
              <Input
                id="from"
                type="number"
                placeholder="Enter 'from' value"
                value={from}
                onChange={(e) => setFrom(e.target.value)}
                required
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="to">To</Label>
              <Input
                id="to"
                type="number"
                placeholder="Enter 'to' value"
                value={to}
                onChange={(e) => setTo(e.target.value)}
                required
              />
            </div>

            <div className="flex items-center space-x-2">
              <Switch
                id="result-downloaded"
                checked={resultDownloaded}
                onCheckedChange={setResultDownloaded}
              />
              <Label htmlFor="result-downloaded">Result Downloaded</Label>
            </div>
          </div>
        </form>
      </CardContent>
      <CardFooter className="flex justify-between items-center">
        <Button type="submit" onClick={handleSubmit} disabled={isLoading}>
          {isLoading ? (
            <>
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              Processing...
            </>
          ) : (
            'Set User Data'
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

export default SetUserData;