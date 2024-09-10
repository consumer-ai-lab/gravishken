import React, { useState } from 'react';
import { File, FileText, NotepadText, Sheet } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { server } from '@common/server';
import * as types from '@common/types';

const appIcons = {
  [types.AppType.TXT]: { icon: NotepadText, color: 'text-green-600' },
  [types.AppType.PPTX]: { icon: File, color: 'text-orange-600' },
  [types.AppType.DOCX]: { icon: FileText, color: 'text-blue-600' },
  [types.AppType.XLSX]: { icon: Sheet, color: 'text-green-600' },
};

interface Test {
  id: string;
  name: string;
  description: string;
  apps: types.AppType[];
}

const mockTests: Test[] = [
  {
    id: '1',
    name: 'Basic Office Skills',
    description: 'Test your skills in Word, Excel, and PowerPoint.',
    apps: [types.AppType.DOCX, types.AppType.XLSX, types.AppType.PPTX],
  },
  {
    id: '2',
    name: 'Advanced Word Processing',
    description: 'Demonstrate your advanced Microsoft Word skills.',
    apps: [types.AppType.DOCX],
  },
  // Add more mock tests as needed
];

const TestSelector = () => {
  const [selectedTest, setSelectedTest] = useState<Test | null>(null);

  const handleOpenApp = (appType: types.AppType) => {
    server.send_message({
      Typ: types.Varient.OpenApp,
      Val: { Typ: appType },
    });
  };

  const handleSubmitWork = () => {
    console.log('Submitting work for test:', selectedTest?.name);
    // Implement submission logic here
  };

  return (
    <div className="flex p-4 h-full">
      <div className="w-1/3 pr-4 overflow-y-auto">
        <h2 className="text-xl font-bold mb-4">Available Tests</h2>
        {mockTests.map((test) => (
          <Button
            key={test.id}
            onClick={() => setSelectedTest(test)}
            variant={selectedTest?.id === test.id ? 'default' : 'outline'}
            className="w-full mb-2 justify-start"
          >
            {test.name}
          </Button>
        ))}
      </div>
      <div className="w-2/3 pl-4">
        {selectedTest ? (
          <Card className="h-full flex flex-col">
            <CardHeader>
              <CardTitle>{selectedTest.name}</CardTitle>
            </CardHeader>
            <CardContent className="flex-grow">
              <p className="mb-4">{selectedTest.description}</p>
              <div className="mb-4">
                <h3 className="text-lg font-semibold mb-2">Associated Apps:</h3>
                <div className="flex space-x-2">
                  {selectedTest.apps.map((appType) => {
                    const { icon: Icon, color } = appIcons[appType];
                    return (
                      <Button
                        key={appType}
                        onClick={() => handleOpenApp(appType)}
                        className="flex items-center space-x-2"
                      >
                        <Icon className={`h-5 w-5 ${color}`} />
                        <span>{types.AppType[appType]}</span>
                      </Button>
                    );
                  })}
                </div>
              </div>
            </CardContent>
            <div className="p-4 border-t">
              <Button onClick={handleSubmitWork} className="w-full">
                Submit Work
              </Button>
            </div>
          </Card>
        ) : (
          <div className="h-full flex items-center justify-center text-gray-500">
            Select a test to view details
          </div>
        )}
      </div>
    </div>
  );
};

export default TestSelector;
