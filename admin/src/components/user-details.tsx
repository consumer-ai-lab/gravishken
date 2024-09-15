import { ScrollArea } from "./ui/scroll-area";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "./ui/table";

export default function UserDetails() {

    const users = [
        { id: 1, username: 'yash_thombre', batch: 'Batch A', testsCompleted: 2, averageWPM: 65 },
        { id: 2, username: 'adnan_husain', batch: 'Batch B', testsCompleted: 3, averageWPM: 72 },
        { id: 3, username: 'prathamesh_kurve', batch: 'Batch A', testsCompleted: 1, averageWPM: 58 },
    ]


    return (

        <div>
            <h1>
                User Details
            </h1>
            <ScrollArea className="h-[400px]">
                <Table>
                    <TableHeader>
                        <TableRow>
                            <TableHead>Username</TableHead>
                            <TableHead>Batch</TableHead>
                            <TableHead>Tests Completed</TableHead>
                            <TableHead>Average WPM</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {users.map((user) => (
                            <TableRow key={user.id}>
                                <TableCell>{user.username}</TableCell>
                                <TableCell>{user.batch}</TableCell>
                                <TableCell>{user.testsCompleted}</TableCell>
                                <TableCell>{user.averageWPM}</TableCell>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </ScrollArea>
        </div>
    )
}