import TypingTest from "@/components/typing-test";
import { useNavigate, useParams } from "react-router-dom"

export default function TestsPage() {
    const { testId } = useParams();
    const navigate = useNavigate();

    const testData = {
        rollNumber: 12345,
        candidateName: 'Yash Thombre'
    };


    const renderTest = () => {
        switch (testId) {
            case '1':
                return (
                    <TypingTest
                        testId={testId}
                        rollNumber={testData.rollNumber}
                        candidateName={testData.candidateName}
                    />
                );
            default:
                return (
                    <div className="text-center p-8">
                        <h2 className="text-2xl font-bold mb-4">Test Not Found</h2>
                        <p className="mb-4">The requested test does not exist or you don't have access to it.</p>
                        <button
                            onClick={() => navigate('/')}
                            className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
                        >
                            Go to Home
                        </button>
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