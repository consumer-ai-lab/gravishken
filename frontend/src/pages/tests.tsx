import React, { useEffect, useState, useRef, useCallback } from 'react';
import { useNavigate, useParams } from "react-router-dom";
import DocumentTests from "@/components/document-tests";
import TypingTest from "@/components/typing-test";
import MCQTest from '@/components/mcq-test';
import { Button } from "@/components/ui/button";
import * as types from "@common/types";
import { CheckCircle, UserIcon } from 'lucide-react';
import * as server from "@common/server.ts";
import { useStateContext } from '@/context/app-context';
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from '@/components/ui/alert-dialog';


interface TestResult {
    testId: string;
    testType: string;
    result: any; 
}


export default function TestsPage() {
    const [testData, setTestData] = useState<types.Test[]>([]);
    const [selectedTestIndex, setSelectedTestIndex] = useState<number | null>(0);
    const [completedTests, setCompletedTests] = useState<string[]>([]);
    const [testResults, setTestResults] = useState<types.UserSubmission[]>([]);

    const [showConfirmDialog, setShowConfirmDialog] = useState(false);
    const [isTestActive, setIsTestActive] = useState(false);
    const { username } = useStateContext();
    const navigate = useNavigate();
    const [leftWidth, setLeftWidth] = useState(250); // Initial width of left panel
    const [isResizing, setIsResizing] = useState(false);
    const containerRef = useRef<HTMLDivElement>(null);
    const startXRef = useRef<number>(0);
    const startWidthRef = useRef<number>(0);

    const [timeLeft, setTimeLeft] = useState(3600); // 1 hour in seconds

    useEffect(() => {
        const timer = setInterval(() => {
            setTimeLeft((prevTime) => (prevTime > 0 ? prevTime - 1 : 0));
        }, 1000);

        return () => clearInterval(timer);
    }, []);

    const formatTime = (seconds: number) => {
        const minutes = Math.floor(seconds / 60);
        const remainingSeconds = seconds % 60;
        return `${minutes.toString().padStart(2, '0')}:${remainingSeconds.toString().padStart(2, '0')}`;
    };

    useEffect(() => {
        fetch(server.base_url + "/get-tests").then(r => r.json()).then(json => {
            console.log(json);
            setTestData(json);
        });
    },[])

    const handleFinishTest = (result: any) => {
        if (selectedTestIndex !== null) {
            const currentTest = testData[selectedTestIndex];
            const newTestResult: types.UserSubmission = {
                TestId: currentTest.Id,
                result: result
            };

            setTestResults([...testResults, newTestResult]);
            setCompletedTests([...completedTests, currentTest.Id]);
            setSelectedTestIndex(null);
            setIsTestActive(false);
            
            if (completedTests.length + 1 === testData.length) {
                setShowConfirmDialog(true);
            }
        }
    };


    const handleConfirmSubmit = () => {
        console.log("Submitting all tests to server");
        console.log("Test results:", testResults);
        // @thrombe: Here you would send testResults to your server
        // For example:
        // sendResultsToServer(testResults).then(() => {
        //     navigate('/test-results');
        // });
        navigate('/end');
        setShowConfirmDialog(false);
    };


    const handleMouseDown = useCallback((e: React.MouseEvent) => {
        e.preventDefault();
        startXRef.current = e.clientX;
        startWidthRef.current = leftWidth;
        setIsResizing(true);
    }, [leftWidth]);

    const handleMouseMove = useCallback((e: MouseEvent) => {
        if (!isResizing) return;
        const dx = e.clientX - startXRef.current;
        const newWidth = Math.max(200, startWidthRef.current + dx);
        setLeftWidth(newWidth);
    }, [isResizing]);

    const handleMouseUp = useCallback(() => {
        setIsResizing(false);
    }, []);

    useEffect(() => {
        if (isResizing) {
            window.addEventListener('mousemove', handleMouseMove);
            window.addEventListener('mouseup', handleMouseUp);
        } else {
            window.removeEventListener('mousemove', handleMouseMove);
            window.removeEventListener('mouseup', handleMouseUp);
        }
        return () => {
            window.removeEventListener('mousemove', handleMouseMove);
            window.removeEventListener('mouseup', handleMouseUp);
        };
    }, [isResizing, handleMouseMove, handleMouseUp]);

    const renderTestContent = () => {
        if (selectedTestIndex === null) {
            return <div className="text-center text-xl mt-10">Select a test from the sidebar to begin.</div>;
        }

        const currentTest = testData[selectedTestIndex];
        if (!currentTest) {
            return <div></div>;
        }
        switch (currentTest.Type) {
            case 'typing':
                return (
                    <TypingTest
                        typingText={currentTest.TypingText!}
                        handleFinishTest={handleFinishTest}
                        isTestActive={isTestActive}
                        setIsTestActive={setIsTestActive}
                    />
                );
            case 'docx':
            case 'xlsx':
            case 'pptx':
                return <DocumentTests testData={currentTest} handleFinishTest={handleFinishTest} />;
            case 'mcq':
                return <MCQTest testData={JSON.parse(currentTest.McqJson!)} handleFinishTest={handleFinishTest} />;
            default:
                return <div>Unknown test type</div>;
        }
    };


    return (
        <div className="flex flex-col h-screen">
            <div className="bg-blue-600 text-white py-2 px-4 flex justify-between items-center">
                <div className="flex items-center space-x-4">
                    <UserIcon size={20} />
                    <span className="font-semibold">{username}</span>
                </div>
                <div className="font-bold">Time Left: {formatTime(timeLeft)}</div>
                <Button onClick={handleFinishTest} variant="destructive">Finish Test</Button>
            </div>

            <div ref={containerRef} className="flex flex-grow relative">
                {/* Left sidebar */}
                <div 
                    style={{ width: `${leftWidth}px`, minWidth: `${leftWidth}px` }} 
                    className="bg-gray-100 p-4 overflow-y-auto"
                >
                    <h2 className="text-xl font-bold mb-4">Tests</h2>
                    {testData.map((test: types.Test, index:number) => (
                        <Button
                            key={test.Id}
                            onClick={() => !isTestActive && setSelectedTestIndex(index)}
                            variant={selectedTestIndex === index ? 'default' : 'outline'}
                            className={`w-full mb-2 justify-start text-left whitespace-normal ${isTestActive ? 'opacity-50 cursor-not-allowed' : ''}`}
                            disabled={isTestActive}
                        >
                            <span className="truncate flex-grow">{test.Type.replace(/([a-z])([A-Z])/g, '$1 $2')}</span>
                            {completedTests.includes(test.Id) && (
                                <CheckCircle className="ml-2 text-green-500" size={16} />
                            )}
                        </Button>
                    ))}
                </div>

                {/* Resize handle */}
                <div
                    className="w-2 bg-gray-300 cursor-col-resize hover:bg-gray-400 transition-colors flex items-center justify-center"
                    onMouseDown={handleMouseDown}
                    style={{ flexShrink: 0 }}
                >
                    <div className="w-0.5 h-8 bg-gray-500" />
                </div>

                {/* Right content area */}
                <div className="flex-grow p-4 overflow-y-auto">
                    {renderTestContent()}
                </div>

                {/* Overlay to prevent interaction while resizing */}
                {isResizing && (
                    <div className="absolute inset-0 bg-transparent cursor-col-resize" />
                )}
            </div>
            <AlertDialog open={showConfirmDialog} onOpenChange={setShowConfirmDialog}>
                <AlertDialogContent>
                    <AlertDialogHeader>
                        <AlertDialogTitle>Are you sure you want to submit all tests?</AlertDialogTitle>
                        <AlertDialogDescription>
                            This action cannot be undone. All your test results will be submitted.
                        </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                        <AlertDialogCancel>Cancel</AlertDialogCancel>
                        <AlertDialogAction onClick={handleConfirmSubmit}>Submit All Tests</AlertDialogAction>
                    </AlertDialogFooter>
                </AlertDialogContent>
            </AlertDialog>
        </div>
    );
}
