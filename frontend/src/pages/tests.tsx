import { useParams } from "react-router-dom"

export default function TestsPage(){
    const {testId} = useParams();

    return (
        <div>
            Tests Page.
        </div>
    )
}