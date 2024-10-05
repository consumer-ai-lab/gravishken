import { useState, useEffect } from 'react';
import { User } from '@common/types';
import axios from 'axios';
import { ScrollArea } from "@/components/ui/scroll-area";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from './ui/card';
import { Loader2 } from 'lucide-react';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from './ui/tooltip';

export default function UserDetails() {
    const [users, setUsers] = useState([]);
    const [currentPage, setCurrentPage] = useState(1);
    const [totalPages, setTotalPages] = useState(0);
    const [isLoading, setIsLoading] = useState(false);

    const fetchUsers = async (page: number) => {
        setIsLoading(true);
        try {
            const response = await axios.get(`${import.meta.env.SERVER_URL}/user/paginated_users`, {
                params: {
                    page: page,
                    limit: 50
                },
                withCredentials: true
            });
            setUsers(response.data.users);
            setTotalPages(response.data.totalPages);
        } catch (error) {
            console.error('Error fetching users:', error);
        }
        setIsLoading(false);
    };

    useEffect(() => {
        fetchUsers(currentPage);
    }, [currentPage]);

    const handlePrevPage = () => {
        if (currentPage > 1) {
            setCurrentPage(currentPage - 1);
        }
    };

    const handleNextPage = () => {
        if (currentPage < totalPages) {
            setCurrentPage(currentPage + 1);
        }
    };

    return (
        <div className="w-full mx-auto p-4 space-y-6">
            <h1 className="text-3xl font-bold mb-8">View All Users</h1>
            <Card>
                <CardContent >
                    <div className="relative mt-6">
                        <div className="overflow-hidden border rounded-md">
                            <Table>
                                <colgroup>
                                    <col className="w-[15%]" /> {/* Adjust the width as needed */}
                                    <col className="w-[18%]" />
                                    <col className="w-[15%]" />
                                    <col className="w-[26%]" />
                                    <col className="w-[10%]" />
                                    <col className="w-[10%]" />
                                </colgroup>
                                <TableHeader className="sticky top-0 bg-background z-10">
                                    <TableRow>
                                        <TableHead className="font-semibold mr-60">ID</TableHead>
                                        <TableHead className="font-semibold">Username</TableHead>
                                        <TableHead className="font-semibold">Password</TableHead>
                                        <TableHead className="font-semibold">Test Password</TableHead>
                                        <TableHead className="font-semibold">Batch</TableHead>
                                    </TableRow>
                                </TableHeader>
                            </Table>
                            <ScrollArea className="h-[360px]">
                                <Table>
                                    <TableBody>
                                        {users.map((user: User) => (
                                            <TableRow key={user.id}>
                                                <TableCell>
                                                    {user.id && <TooltipProvider>
                                                        <Tooltip>
                                                            <TooltipTrigger>{user.id.slice(0, 4)}...{user.id.slice(-4)}</TooltipTrigger>
                                                            <TooltipContent>
                                                                <div>{user.id}</div>
                                                            </TooltipContent>
                                                        </Tooltip>
                                                    </TooltipProvider>}
                                                </TableCell>
                                                <TableCell>{user.username}</TableCell>
                                                <TableCell>{user.password}</TableCell>
                                                <TableCell>{user.testPassword}</TableCell>
                                                <TableCell>{user.batch}</TableCell>
                                                <TableCell className='flex justify-center'>
                                                    <Button variant={"ghost"}>
                                                        Edit
                                                    </Button>
                                                </TableCell>
                                            </TableRow>
                                        ))}
                                    </TableBody>
                                </Table>
                            </ScrollArea>
                        </div>
                    </div>
                    <div className="flex justify-between items-center mt-4">
                        <Button
                            onClick={handlePrevPage}
                            disabled={currentPage === 1 || isLoading}
                            variant="outline"
                        >
                            Previous
                        </Button>
                        <span className="text-sm text-muted-foreground">
                            Page {currentPage} of {totalPages}
                        </span>
                        <Button
                            onClick={handleNextPage}
                            disabled={currentPage === totalPages || isLoading}
                            variant="outline"
                        >
                            Next
                        </Button>
                    </div>
                    {isLoading && (
                        <div className="flex justify-center items-center mt-4">
                            <Loader2 className="h-6 w-6 animate-spin" />
                        </div>
                    )}
                </CardContent>
            </Card>
        </div>
    )
}