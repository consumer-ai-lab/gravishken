import { FileText, NotepadText, Sheet } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import * as server from '@common/server';
import * as types from '@common/types';

const appMapping = {
  'docx': { icon: FileText, color: 'text-blue-600', appType: types.AppType.DOCX },
  'xlsx': { icon: Sheet, color: 'text-green-600', appType: types.AppType.XLSX },
  'pptx': { icon: NotepadText, color: 'text-red-600', appType: types.AppType.TXT },
};

interface DocumentTestsProps {
  testData: types.Test;
  handleFinishTest: (result: any) => void;
}

export default function DocumentTests({
  testData,
  handleFinishTest,
}: DocumentTestsProps) {
  const handleOpenApp = (app: types.Test) => {
    /// @ts-ignore
    let typ: types.AppType = appMapping[testData.Type].appType;
    server.server.send_message({
      Typ: types.Varient.OpenApp,
      Val: { Typ: typ, TestId: app.Id },
    });
  };

  const handleSubmitWork = async () => {
    let resp = await fetch(server.base_url + "/get-user");
    let user: types.User = await resp.json()
    let submission: types.TestSubmission = {
      TestId: testData.Id,
      UserId: user.Id,
      TestInfo: {
        Type: testData.Type,
      },
    };
    await fetch(server.base_url + "/submit-test", {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(submission),
    })
  };

  // @ts-ignore
  const appConfig = appMapping[testData.Type];

  return (
    <Card className="h-full flex flex-col">
      <CardHeader>
        <CardTitle>
          {testData.Type.replace(/([a-z])([A-Z])/g, '$1 $2')}
        </CardTitle>
      </CardHeader>
      <CardContent className="flex-grow">
        <div className="mb-4">
          <h3 className="text-lg font-semibold mb-2">Associated Apps:</h3>
          <div className="flex space-x-2 mb-4">

            <Button
              variant={"default"}
              onClick={() => handleOpenApp(testData)}
              className="flex items-center space-x-2"
            >
              <span>Open Associated App</span>
            </Button>

          </div>
        </div>
        <div className='w-full h-[400px] overflow-hidden rounded-lg mb-2'>
          <img src={testData.FilePath} alt={`${testData.Id} Test`} className="w-full object-cover" />
        </div>
        <Button onClick={handleSubmitWork} className="w-full">
          Submit Work
        </Button>
      </CardContent>
    </Card>
  );
};
