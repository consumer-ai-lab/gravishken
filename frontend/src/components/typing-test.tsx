import React, { useState, useEffect, useRef } from 'react';
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Progress } from "@/components/ui/progress";
import { PlayCircle, StopCircle, Send } from 'lucide-react';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";

interface TypingTestProps {
    typingText: string;
    handleFinishTest: (result: any) => void;
    isTestActive: boolean;
    setIsTestActive: (isActive: boolean) => void;
}

export default function TypingTest({
    typingText,
    handleFinishTest,
    setIsTestActive
}: TypingTestProps) {
    const testime = 300;
    const [totalCharsTyped, setTotalCharsTyped] = useState(0);
    const [totalCorrectCharacters, setTotalCorrectCharacters] = useState(0);
    const [isStarted, setIsStarted] = useState(false);
    const [inputText, setInputText] = useState('');
    const [timeLeft, setTimeLeft] = useState(testime); // 5 minutes in seconds
    const [rawWPM, setrawWPM] = useState(0);
    const [wpm, setWpm] = useState(0);
    const [feedback, setFeedback] = useState<'correct' | 'incorrect' | null>(null);
    const [traversal, setTraversal] = useState<number>(0);
    const [matched, setMatched] = useState<number[]>([]);
    const timerRef = useRef<NodeJS.Timeout | null>(null);
    const textareaRef = useRef<HTMLTextAreaElement | null>(null);
    const [showConfirmDialog, setShowConfirmDialog] = useState(false);
    const [isSubmitted, setIsSubmitted] = useState(false);

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
        if (traversal == typingText.length - 1) {
            handleSubmit();
        }
        if (isStarted) {
            if (inputText.length === 0) {
                return;
            }

            if (inputText.length - 1 < traversal) {
                if (matched.includes(traversal)) {
                    setMatched((prev) => prev.filter(item => item !== traversal));
                }
                setTraversal(traversal - 1);
            } else {
                if (typingText[traversal] === inputText[traversal]) {
                    setTotalCorrectCharacters((prev) => prev + 1);
                    setMatched((prev) => [...prev, traversal]);
                }
                setTotalCharsTyped((prev) => prev + 1);
                setTraversal(traversal + 1);
            }

            const minutesPassed = (testime - timeLeft) / 60;
            if (minutesPassed <= 0) return;

            const rawWPM = (totalCharsTyped / 5) / minutesPassed;
            const WPM = rawWPM * (totalCorrectCharacters / totalCharsTyped);

            setrawWPM(Math.round(rawWPM) || 0);
            setWpm(Math.round(WPM) || 0);

            const isCorrect = typingText.startsWith(inputText);
            setFeedback(isCorrect ? 'correct' : 'incorrect');
        }
    }, [inputText, timeLeft, isStarted]);

    const handleStart = () => {
        setIsStarted(true);
        setIsTestActive(true);
        setTimeout(() => {
            if (textareaRef.current) {
                textareaRef.current.focus();
            }
        }, 0);
    };
    const handleSubmit = () => {
        setShowConfirmDialog(true);
    };

    const confirmSubmit = () => {
        if (timerRef.current) clearInterval(timerRef.current);
        setIsStarted(false);
        setIsTestActive(false);
        setIsSubmitted(true);
        setShowConfirmDialog(false);
        console.log('TotalCharsTyped:', totalCharsTyped);
        console.log('TotalCorrectCharacters:', totalCorrectCharacters);

        const result = {
            timeTaken: testime - timeLeft,
            wpm,
            rawWPM,
            accuracy: calculateAccuracy(inputText, typingText)
        }

        console.log('Submitting results:', result);
        handleFinishTest(result);

    };

    const calculateAccuracy = (input: string, original: string) => {
        const inputWords = input.trim().split(/\s+/);
        const originalWords = original.trim().split(/\s+/);
        const correctWords = inputWords.filter((word, index) => word === originalWords[index]);
        return Math.round((correctWords.length / originalWords.length) * 100);
    };

    const getHighlightedText = () => {
        return (
            <pre className="bg-gray-800 text-white p-4 rounded-lg whitespace-pre-wrap break-words">
                {typingText.split('').map((char, index) => {
                    const isCorrect = inputText[index] === char;
                    const isSpace = char === ' ';
                    return (
                        <span
                            key={index}
                            className={`${inputText[index] !== undefined
                                    ? isCorrect
                                        ? "text-green-500"
                                        : isSpace
                                            ? "text-red-500 bg-red-900"
                                            : "text-red-500"
                                    : ""
                                }`}
                        >
                            {char}
                        </span>
                    );
                })}
            </pre>
        );
    };


    return (
        <Card className="w-full max-w-8xl rounded-lg overflow-hidden mx-auto">
            <CardHeader >
                <CardTitle>Typing Test</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4 pt-6">
                <div className="bg-gray-100 p-4 rounded-md">
                    <h3 className="font-semibold mb-2">Instructions:</h3>
                    <p>Type the following text as accurately and quickly as you can. Your time starts when you click "Start".</p>
                </div>
                <div className="bg-white border border-gray-300 p-4 rounded-md select-none">
                    {getHighlightedText()} {/* Highlighted Original Text */}
                </div>
                <Textarea
                    ref={textareaRef}
                    value={inputText}
                    onChange={(e) => setInputText(e.target.value)}
                    placeholder="Start typing here..."
                    disabled={!isStarted}
                    className={`h-40 resize-none px-4 py-3 rounded-lg shadow-sm transition-colors duration-300 focus:outline-none focus:ring-4 focus:ring-blue-500 focus:border-transparent 
                        ${isStarted ? 'bg-white' : 'bg-gray-100 border-black cursor-not-allowed'} 
                        ${feedback === 'correct' ? 'border-green-500' :
                            feedback === 'incorrect' ? 'border-red-500' : 'border-black'}
                    `}
                    onPasteCapture={(e) => {
                        console.log(e);
                        e.preventDefault();
                    }}
                />

                <div className="flex justify-between items-center">
                    <div className="text-lg font-semibold">
                        Time Left: {Math.floor(timeLeft / 60)}:{(timeLeft % 60).toString().padStart(2, '0')}
                    </div>
                    <div className="text-lg font-semibold">
                        WPM: {wpm}
                    </div>
                    <div className={`text-lg font-semibold ${feedback === 'correct' ? 'text-green-500' :
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
                    disabled={isStarted || isSubmitted}
                    className="bg-green-600 hover:bg-green-700"
                >
                    {isStarted ? <StopCircle className="mr-2" /> : <PlayCircle className="mr-2" />}
                    {isStarted ? 'In Progress' : 'Start'}
                </Button>
                <Button
                    onClick={handleSubmit}
                    disabled={!isStarted || isSubmitted}
                    className="bg-blue-600 hover:bg-blue-700"
                >
                    <Send className="mr-2" /> Submit
                </Button>
                <Dialog open={showConfirmDialog} onOpenChange={setShowConfirmDialog}>
                    <DialogContent>
                        <DialogHeader>
                            <DialogTitle>Confirm Submission</DialogTitle>
                            <DialogDescription>
                                Are you sure you want to submit your typing test?
                            </DialogDescription>
                        </DialogHeader>
                        <DialogFooter>
                            <Button onClick={() => setShowConfirmDialog(false)}>Cancel</Button>
                            <Button onClick={confirmSubmit} className="bg-blue-600 hover:bg-blue-700">
                                Confirm
                            </Button>
                        </DialogFooter>
                    </DialogContent>
                </Dialog>
            </CardFooter>
        </Card>
    );
};

