import React, { useEffect, useState } from 'react';

import OfficeAppSwitcher from "@/components/main-test";
import TypingTest from "@/components/typing-test";
import { useNavigate, useParams } from "react-router-dom"

export default function TestsPage() {
    const { testId } = useParams();
    const navigate = useNavigate();

    const [username, setUsername] = useState<string>('');
    const [userPassword, setUserPassword] = useState<string>('');
    const [testPassword, setTestPassword] = useState<string>('');
    const [rollNumber, setRollNumber] = useState<number>(0);


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

    console.log("Test Data: ", testData);


    const renderTest = () => {
        switch (testId) {
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
                return (
                    <div>
                        <OfficeAppSwitcher/>
                    </div>
                );
        }
    };


    return (
        <div className="container mx-auto px-4 py-4">
            {renderTest()}
        </div>
    )
}