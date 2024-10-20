import { FileText, NotepadText, Sheet } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { server } from '@common/server';
import * as types from '@common/types';

const appMapping = {
  'DocTest': { icon: FileText, color: 'text-blue-600', appType: types.AppType.DOCX },
  'ExcelTest': { icon: Sheet, color: 'text-green-600', appType: types.AppType.XLSX },
  'WordTest': { icon: NotepadText, color: 'text-red-600', appType: types.AppType.TXT },
};

interface DocumentTestsProps {
  testData: {
    testType: string;
    testId: string;
    imagePath:string;
  };
  handleFinishTest: (result: any) => void;
}

export default function DocumentTests({
  testData,
  handleFinishTest,
}: DocumentTestsProps){
  const handleOpenApp = (appType: types.AppType) => {
    server.send_message({
      Typ: types.Varient.OpenApp,
      Val: { Typ: appType },
    });
  };

  const handleSubmitWork = () => {
    handleFinishTest({ /*TODO:Test data here */ });
  };

  const appConfig = appMapping[testData.testType as keyof typeof appMapping];

  return (
    <Card className="h-full flex flex-col">
      <CardHeader>
        <CardTitle>
          {testData.testType.replace(/([a-z])([A-Z])/g, '$1 $2')}
        </CardTitle>
      </CardHeader>
      <CardContent className="flex-grow">
        <div className="mb-4">
          <h3 className="text-lg font-semibold mb-2">Associated Apps:</h3>
          <div className="flex space-x-2 mb-4">

                <Button
                  variant={"default"}
                  onClick={() => handleOpenApp(appConfig.appType)}
                  className="flex items-center space-x-2"
                >
                  <span>Open Associated App</span>
                </Button>

          </div>
        </div>
        <div className='w-full h-[400px] overflow-hidden rounded-lg mb-2'>
          <img src={testData.imagePath} alt={`${testData.testId} Test`} className="w-full object-cover" />
        </div>
        <Button onClick={handleSubmitWork} className="w-full">
          Submit Work
        </Button>
      </CardContent>
    </Card>
  );
};
