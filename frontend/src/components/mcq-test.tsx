import { useState } from 'react';
import { ChevronLeft, ChevronRight } from 'lucide-react';
import { Button } from './ui/button';
import { Card, CardContent, CardHeader, CardTitle } from './ui/card';
import test from 'node:test';
import * as server from '@common/server';
import * as types from '@common/types';

interface MCQTestProps {
    Test: types.Test,
    testData: {
        question: string;
        options: string[];
    }[];
    handleFinishTest: (result: any) => void;
}

export default function MCQTest({
    Test,
    testData,
    handleFinishTest,
}: MCQTestProps) {
    const [currentQuestion, setCurrentQuestion] = useState(0);
    const [answers, setAnswers] = useState(new Array(testData.length).fill(null));



    const handleAnswerSelect = (index: number) => {
        const newAnswers = [...answers];
        newAnswers[currentQuestion] = index;
        setAnswers(newAnswers);
    };

    const handleNavigation = (direction: number) => {
        // Ensure the current question is within the bounds of the test data
        setCurrentQuestion(prev => Math.max(0, Math.min(testData.length, prev + direction)));
    };

    const handleSubmit = async () => {
        let resp = await fetch(server.base_url + "/get-user");
        let user: types.User = await resp.json()
        let submission: types.TestSubmission = {
          TestId: Test.Id,
          UserId: user.Id,
          TestInfo: {
            Type: Test.Type,
            McqTestInfo: {
                Answers: answers,
            }
          },
        };
        await fetch(server.base_url + "/submit-test", {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(submission),
        })
        handleFinishTest(answers);
    };

    return (
        <Card className="w-full max-w-8xl rounded-lg overflow-hidden mx-auto">
            <CardHeader >
                <CardTitle>MCQ Test</CardTitle>
            </CardHeader>
            <CardContent>
                <div className="mb-6">
                    <p className="text-lg text-primary font-medium mb-4">{testData[currentQuestion].question}</p>
                    <div className="space-y-2">
                        {testData[currentQuestion].options.map((option, index) => (
                            <Button
                                key={index}
                                onClick={() => handleAnswerSelect(index)}
                                className={`w-full text-left p-3 rounded-md transition-colors ${answers[currentQuestion] === index
                                    ? 'bg-blue-100 text-blue-800'
                                    : 'bg-gray-100 text-primary hover:bg-gray-200'
                                    }`}
                            >
                                {option}
                            </Button>
                        ))}
                    </div>
                </div>
                <div className="flex justify-between items-center">
                    <button
                        onClick={() => handleNavigation(-1)}
                        disabled={currentQuestion === 0}
                        className="p-2 rounded-full bg-gray-200 hover:bg-gray-300 disabled:opacity-50"
                    >
                        <ChevronLeft size={24} />
                    </button>
                    <span className="text-lg font-medium">Question {currentQuestion + 1} of {testData.length}</span>
                    {currentQuestion < testData.length-1 ? (
                        <button
                            onClick={() => handleNavigation(1)}
                            className="p-2 rounded-full bg-gray-200 hover:bg-gray-300"
                        >
                            <ChevronRight size={24} />
                        </button>
                    ) : (
                        <button
                            onClick={handleSubmit}
                            className="px-4 py-2 bg-green-500 text-white rounded-md hover:bg-green-600 transition-colors"
                        >
                            Submit
                        </button>
                    )}
                </div>
                <div className="mt-4 flex justify-center space-x-2">
                    {answers.map((_, index) => (
                        <div
                            key={index}
                            className={`w-3 h-3 rounded-full ${answers[index] !== null ? 'bg-blue-500' : 'bg-gray-300'
                                }`}
                        />
                    ))}
                </div>
            </CardContent>
        </Card>
    );
};

