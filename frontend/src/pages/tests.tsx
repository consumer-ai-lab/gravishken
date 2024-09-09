import OfficeAppSwitcher from "@/components/main-test";
import TypingTest from "@/components/typing-test";
import { useNavigate, useParams } from "react-router-dom"

export default function TestsPage() {
    const { testId } = useParams();
    const navigate = useNavigate();

    const testData = {
        rollNumber: 12345,
        candidateName: 'Yash Thombre',
        testPassword: "securepassword123"
    };


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