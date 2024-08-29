import React, { useState } from 'react';
import { File, FileText, NotepadText } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { server } from '@common/server';
import * as types from '@common/types';

const apps = [
  { name: 'NotePad', icon: NotepadText, color: 'text-green-600' },
  { name: 'PowerPoint', icon: File, color: 'text-orange-600' },
  { name: 'Word', icon: FileText, color: 'text-blue-600' },
];

const OfficeAppSwitcher = () => {
  const [activeApp, setActiveApp] = useState(apps[0].name);

  const handleStartTest = (appName: string) => {
    console.log("Starting test with testPassword:");
    setActiveApp(appName);
    server.send_message({
        Typ: types.Varient.MicrosoftApps, 
        Val: {
            AppName: appName,
            Device: "linux"
        }
    });
};

  return (
    <div className="p-4">
      <div className="flex space-x-4 mb-4">
        {apps.map((app) => (
          <Button
            key={app.name}
            onClick={() => handleStartTest(app.name)}
            variant={activeApp === app.name ? 'default' : 'outline'}
            className="flex items-center space-x-2"
          >
            <app.icon className={`h-5 w-5 ${app.color}`} />
            <span>{app.name}</span>
          </Button>
        ))}
      </div>
      <Card className='border-2 border-black'>
        <CardHeader>
          <CardTitle>Active App: {activeApp}</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="h-96 flex items-center justify-center text-gray-500">
            {activeApp} content would be displayed here
          </div>
        </CardContent>
      </Card>
    </div>
  );
};

export default OfficeAppSwitcher;