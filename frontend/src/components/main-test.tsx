import React from 'react';
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

interface TestSelectorProps {
  test: {
    id: string;
    name: string;
    description: string;
    apps: types.AppType[];
  };
}

const TestSelector: React.FC<TestSelectorProps> = ({ test }) => {
  const handleOpenApp = (appType: types.AppType) => {
    server.send_message({
      Typ: types.Varient.OpenApp,
      Val: { Typ: appType },
    });
  };

  const handleSubmitWork = () => {
    console.log('Submitting work for test:', test.name);
    // Implement submission logic here
  };

  return (
    <Card className="h-full flex flex-col">
      <CardHeader>
        <CardTitle>{test.name}</CardTitle>
      </CardHeader>
      <CardContent className="flex-grow">
        <p className="mb-4">{test.description}</p>
        <div className="mb-4">
          <h3 className="text-lg font-semibold mb-2">Associated Apps:</h3>
          <div className="flex space-x-2 mb-4">
            {test.apps.map((appType) => {
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
        <Button onClick={handleSubmitWork} className="w-full">
          Submit Work
        </Button>
      </CardContent>
    </Card>
  );
};

export default TestSelector;
