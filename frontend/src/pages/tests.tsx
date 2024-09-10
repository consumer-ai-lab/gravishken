import React, { useEffect, useState, useRef, useCallback } from 'react';
import { useNavigate, useParams } from "react-router-dom";
import TestSelector from "@/components/main-test";
import TypingTest from "@/components/typing-test";
import { Button } from "@/components/ui/button";
import * as types from "@common/types";
import { UserIcon } from 'lucide-react';

const testList = [
    { id: '1', name: 'Typing Test' },
    { id: '2', name: 'Basic Office Skills', description: 'Test your skills in Word, Excel, and PowerPoint.', apps: [types.AppType.DOCX, types.AppType.XLSX, types.AppType.PPTX] },
    { id: '3', name: 'Advanced Word Processing', description: 'Demonstrate your advanced Microsoft Word skills.', apps: [types.AppType.DOCX] },
    // Add more tests as needed
];

export default function TestsPage() {
    const { testId } = useParams();
    const navigate = useNavigate();
    const [leftWidth, setLeftWidth] = useState(250); // Initial width of left panel
    const [isResizing, setIsResizing] = useState(false);
    const containerRef = useRef<HTMLDivElement>(null);
    const startXRef = useRef<number>(0);
    const startWidthRef = useRef<number>(0);

    const [username, setUsername] = useState<string>('');
    const [userPassword, setUserPassword] = useState<string>('');
    const [testPassword, setTestPassword] = useState<string>('');
    const [rollNumber, setRollNumber] = useState<number>(0);

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

    const handleFinishTest = () => {
        // Implement finish test logic here
        console.log('Finishing test');
        // You might want to navigate to a results page or show a confirmation dialog
    };

    // Load stored data from localStorage when the component mounts
    useEffect(() => {
        const storedUsername = localStorage.getItem('username');
        const storedUserPassword = localStorage.getItem('userPassword');
        const storedTestPassword = localStorage.getItem('testPassword');

        console.log("Username: ", storedUsername);
        console.log("UserPassword: ", storedUserPassword);
        console.log("TestPassword: ", storedTestPassword);

        if (storedUsername) setUsername(storedUsername);
        if (storedUserPassword) {
            setUserPassword(storedUserPassword);
            const digits = storedUserPassword.split('_')[1];
            setRollNumber(parseInt(digits));
            console.log("Digits: ", digits);
            console.log("Roll Number: ", rollNumber);
        }
        if (storedTestPassword) setTestPassword(storedTestPassword);
    }, []);

    const testData = {
        rollNumber: rollNumber,
        candidateName: username,
        testPassword: testPassword
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
        switch (testId) {
            case undefined:
                return (<>
                    TODO: all tests
                    </>)
            case '1':
                return (
                    <TypingTest
                        testId={testId}
                        rollNumber={testData.rollNumber}
                        candidateName={testData.candidateName}
                        testPassword={testData.testPassword}
                    />
                );
            default:
                const selectedTest = testList.find(test => test.id === testId);
                return selectedTest ? (
                    <TestSelector test={selectedTest} />
                ) : (
                    <div>Test not found</div>
                );
        }
    };

    return (
        <div className="flex flex-col h-screen">
            {/* Updated top bar */}
            <div className="bg-gray-200 py-2 px-4 flex justify-between items-center">
                <div className="flex items-center space-x-4">
                    <UserIcon size={20} />
                    <span className="font-semibold">{username}</span>
                    <span>Roll: {rollNumber}</span>
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
                    {testList.map((test) => (
                        <Button
                            key={test.id}
                            onClick={() => navigate(`/tests/${test.id}`)}
                            variant={testId === test.id ? 'default' : 'outline'}
                            className="w-full mb-2 justify-start text-left whitespace-normal"
                        >
                            <span className="truncate">{test.name}</span>
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
        </div>
    );
}
