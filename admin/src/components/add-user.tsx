import { useState } from 'react';
import { Input } from './ui/input';
import { Button } from './ui/button';
import { Card,  CardContent, CardFooter } from '@/components/ui/card';
import { Label } from '@/components/ui/label';
import { Upload, FileText, Users } from 'lucide-react';
import api from '@/lib/api';
import axios from 'axios';
import { useToast } from '@/hooks/use-toast';

export default function AddUser() {
  const [file, setFile] = useState<File | null>(null);
  const [fileName, setFileName] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const { toast } = useToast();

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFile = event.target.files?.[0];
    if (selectedFile) {
      setFile(selectedFile);
      setFileName(selectedFile.name);
    }
  };

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    if (!file) {
      toast({
        variant:"destructive",
        description:"Please select a CSV file to upload."
      })
      return;
    }

    setIsLoading(true);

    const formData = new FormData();
    formData.append('file', file);

    try {
      const response = await api.post(`/admin/add_users_from_csv`, formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });

      
      toast({
        variant:"default",
        title:"Successfully uploaded the CSV",
        description:response.data.message 
      })
      setFile(null);
      setFileName('');
    } catch (error:any) {
      if (axios.isAxiosError(error) && error.response) {
        toast({
          variant:"destructive",
          title:"Error while uploading the student CSV",
          description: error.response.data.error || 'An error occurred while uploading the file.'
        })
      } else {
        toast({
          variant:"destructive",
          title:"Error while uploading the student CSV",
          description:'An unknown error occurred while uploading the file.' 
        })
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="container mx-auto p-4 space-y-6">
      <div className='flex gap-2 items-center mb-8'>
        <Users className="mr-2" />
        <h1 className='text-3xl font-bold'>
          Add Users from CSV
        </h1>
      </div>
      <Card className="p-6 w-full">
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="csvFile" className="text-sm font-medium">
                Upload CSV File
              </Label>
              <div className="flex items-center space-x-2">
                <div className="relative flex-grow">
                <Input
                    type="file"
                    id="csvFile"
                    accept=".csv"
                    onChange={handleFileChange}
                    className="sr-only"
                  />
                   <Label
                    htmlFor="csvFile"
                    className="flex items-center justify-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 cursor-pointer"
                  >
                    <Upload className="w-5 h-5 mr-2" />
                    Choose File
                  </Label>
                </div>
                {fileName && (
                  <div className="flex items-center text-sm text-gray-500">
                    <FileText className="w-4 h-4 mr-1" />
                    {fileName}
                  </div>
                )}
              </div>
            </div>
          </form>
        </CardContent>
        <CardFooter>
        <Button 
            type="submit" 
            className="w-full" 
            onClick={handleSubmit}
            disabled={isLoading || !file}
          >
            {isLoading ? 'Uploading...' : 'Submit'}
          </Button>
        </CardFooter>
      </Card>
    </div>
  );
}