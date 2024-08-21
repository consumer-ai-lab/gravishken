import React, { useState, useEffect, useRef } from 'react';
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Progress } from "@/components/ui/progress";
import { PlayCircle, StopCircle, Send } from 'lucide-react';

const mockText = "It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout. The point of using Lorem Ipsum is that it has a more-or-less normal distribution of letters, as opposed to using 'Content here, content here', making it look like readable English. Many desktop publishing packages and web page editors now use Lorem Ipsum as their default model text, and a search for 'lorem ipsum' will uncover many web sites still in their infancy. Various versions have evolved over the years, sometimes by accident, sometimes on purpose (injected humour and the like).";

interface TypingTestProps {
    testId: string;
    rollNumber: number;
    candidateName: string;
}

export default function TypingTest({
    testId,
    rollNumber,
    candidateName
}: TypingTestProps) {
    const [isStarted, setIsStarted] = useState(false);
    const [inputText, setInputText] = useState('');
    const [timeLeft, setTimeLeft] = useState(300); // 5 minutes in seconds
    const [wpm, setWpm] = useState(0);
    const [feedback, setFeedback] = useState<'correct' | 'incorrect' | null>(null);
    const timerRef = useRef<NodeJS.Timeout | null>(null);
    const textareaRef = useRef<HTMLTextAreaElement | null>(null);

    useEffect(() => {
        if (isStarted && timeLeft > 0) {
            timerRef.current = setInterval(() => {
                setTimeLeft((prevTime) => prevTime - 1);
            }, 1000);
        } else if (timeLeft === 0) {
            handleSubmit();
        }

        return () => {
            if (timerRef.current) clearInterval(timerRef.current);
        };
    }, [isStarted, timeLeft]);

    useEffect(() => {
        if (isStarted) {
            const wordsTyped = inputText.trim().split(/\s+/).length;
            const minutesPassed = (300 - timeLeft) / 60;
            setWpm(Math.round(wordsTyped / minutesPassed) || 0);

            const isCorrect = mockText.startsWith(inputText);
            setFeedback(isCorrect ? 'correct' : 'incorrect');
        }
    }, [inputText, timeLeft, isStarted]);

    const handleStart = () => {
        setIsStarted(true);
        if (textareaRef.current) textareaRef.current.focus();
    };


    const handleSubmit = () => {
        if (timerRef.current) clearInterval(timerRef.current);
        setIsStarted(false);

        console.log('Submitting results:', {
            testId,
            rollNumber,
            candidateName,
            timeTaken: 300 - timeLeft,
            wpm,
            accuracy: calculateAccuracy(inputText, mockText)
        });
    };

    const calculateAccuracy = (input:string, original:string) => {
        const inputWords = input.trim().split(/\s+/);
        const originalWords = original.trim().split(/\s+/);
        const correctWords = inputWords.filter((word, index) => word === originalWords[index]);
        return Math.round((correctWords.length / originalWords.length) * 100);
    };

    return (
        <Card className="w-full max-w-8xl rounded-lg overflow-hidden mx-auto">
            <CardHeader className="bg-blue-600 text-white">
                <div className="flex justify-between items-center">
                    <CardTitle>WCL Recruitment Test - Typing Speed</CardTitle>
                    <div className="text-sm">
                        <div>Roll Number: {rollNumber}</div>
                        <div>Name: {candidateName}</div>
                    </div>
                </div>
            </CardHeader>
            <CardContent className="space-y-4 pt-6">
                <div className="bg-gray-100 p-4 rounded-md">
                    <h3 className="font-semibold mb-2">Instructions:</h3>
                    <p>Type the following text as accurately and quickly as you can. Your time starts when you click "Start".</p>
                </div>
                <div className="bg-white border border-gray-300 p-4 rounded-md">
                    <p className="text-gray-700">{mockText}</p>
                </div>
                <Textarea
                    ref={textareaRef}
                    value={inputText}
                    onChange={(e) => setInputText(e.target.value)}
                    placeholder="Start typing here..."
                    disabled={!isStarted}
                    className={`h-40 resize-none ${
                        feedback === 'correct' ? 'border-green-500' : 
                        feedback === 'incorrect' ? 'border-red-500' : ''
                    }`}
                />
                <div className="flex justify-between items-center">
                    <div className="text-lg font-semibold">
                        Time Left: {Math.floor(timeLeft / 60)}:{(timeLeft % 60).toString().padStart(2, '0')}
                    </div>
                    <div className="text-lg font-semibold">
                        WPM: {wpm}
                    </div>
                    <div className={`text-lg font-semibold ${
                        feedback === 'correct' ? 'text-green-500' : 
                        feedback === 'incorrect' ? 'text-red-500' : ''
                    }`}>
                        {feedback && (feedback === 'correct' ? 'Correct' : 'Incorrect')}
                    </div>
                </div>
                <Progress value={(300 - timeLeft) / 3} className="w-full" />
            </CardContent>
            <CardFooter className="flex justify-between">
                <Button
                    onClick={handleStart}
                    disabled={isStarted}
                    className="bg-green-600 hover:bg-green-700"
                >
                    {isStarted ? <StopCircle className="mr-2" /> : <PlayCircle className="mr-2" />}
                    {isStarted ? 'In Progress' : 'Start'}
                </Button>
                <Button
                    onClick={handleSubmit}
                    disabled={!isStarted}
                    className="bg-blue-600 hover:bg-blue-700"
                >
                    <Send className="mr-2" /> Submit
                </Button>
            </CardFooter>
        </Card>
    );
};

