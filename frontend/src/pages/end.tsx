import { useNavigate } from 'react-router-dom';
import { Button } from "@/components/ui/button";
import { CheckCircle } from 'lucide-react';
import { server } from '@common/server';
import * as types from '@common/types';


export default function ResultPage() {
    
    function handleEnd() {
        server.send_message({
            Typ:types.Varient.Quit,
            Val: {}
        });
    }

    return (
        <div className="min-h-screen bg-gray-100 flex flex-col items-center justify-center p-4">
            <div className="bg-white rounded-lg shadow-xl p-8 max-w-md w-full text-center">
                <CheckCircle className="mx-auto mb-4 text-green-500" size={64} />
                <h1 className="text-2xl font-bold mb-4">Submission Recorded</h1>
                <p className="text-gray-600 mb-6">
                    Your test results have been successfully submitted. Thank you for your participation.
                </p>
                <div className="bg-blue-50 border border-blue-200 rounded-md p-4 mb-6">
                    <p className="text-blue-800 font-medium">
                        Thank you for taking the test!
                    </p>
                    <p className="text-blue-600 mt-2">
                        - Team Gravishken
                    </p>
                </div>
                <Button
                    onClick={handleEnd}
                    className="w-full"
                >
                    Return to Login
                </Button>
            </div>
        </div>
    );
}