import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Clock, AlertTriangle } from 'lucide-react';
import { useState } from "react";
import { server } from '@common/server';
import * as types from '@common/types';

export default function InstructionsPage() {

    // Function to handle Start Test
    const handleStartTest = () => {

        server.send_message({
            Typ: types.Varient.StartTest,
            Val: {}
        });
    };

    function handleEnd() {
        server.send_message({
            Typ: types.Varient.Quit,
            Val: {}
        });
    }


    return (
        <div className="min-h-screen bg-gray-100 p-4 flex items-center justify-center">
            <Card className="w-full max-w-6xl rounded-lg overflow-hidden">
                <CardContent className="pt-6">
                    <CardTitle className="text-2xl mb-4">Instructions for the Test</CardTitle>
                    <ol className="list-decimal pl-6 space-y-3">
                        <li>The total duration of this test is 10 minutes, and it carries a maximum of 10 marks.</li>
                        <li>You will be given 5 minutes to read the question. Click the "Start Test" button to begin the test timer.</li>
                        <li>The test must be submitted within 10 minutes by clicking the "Submit Test" button for final submission.</li>
                        <li>If not done so, the test will be automatically submitted once the time is up.</li>
                    </ol>
                </CardContent>
                <CardFooter className="flex flex-col items-stretch gap-4">
                    <Alert variant="destructive">
                        <AlertTriangle className="h-4 w-4" />
                        <AlertTitle>Important</AlertTitle>
                        <AlertDescription>
                            Ensure you have a stable internet connection before starting the test.
                        </AlertDescription>
                    </Alert>
                    <div className="w-full flex flex-row gap-2">
                        <Button className="bg-green-600 hover:bg-green-700 text-white py-3 w-full" onClick={handleStartTest}>
                            <Clock className="mr-2 h-4 w-4" /> Start Test
                        </Button>
                        <Button className="bg-blue-600 hover:bg-blue-700 text-white py-3 mr-2 w-full" onClick={handleEnd}>
                            End Test
                        </Button>
                    </div>
                </CardFooter>
            </Card>
        </div>
    );

}
