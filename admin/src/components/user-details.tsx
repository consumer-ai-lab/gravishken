import { useEffect, useState } from 'react';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Loader2, Trash2, PencilIcon, Search } from 'lucide-react';
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { Alert, AlertDescription } from "@/components/ui/alert";
import axios from 'axios';
import { useDebounce } from 'use-debounce';
import { User as ImportedUser } from '@common/types';

// @thrombe: fix this type in common, it is changing every time on compliation
type UserWithoutId = Omit<ImportedUser, 'id'>;

export interface User extends UserWithoutId {
    id: string;
}

interface UserDetailsProps {
    isAuthenticated: boolean;
}

export default function UserDetails({ isAuthenticated }: UserDetailsProps) {
    const [users, setUsers] = useState<User[]>([]);
    const [searchTerm, setSearchTerm] = useState('');
    const [currentPage, setCurrentPage] = useState(1);
    const [isLoading, setIsLoading] = useState(false);
    const [totalPages, setTotalPages] = useState(0);
    const [totalUsers, setTotalUsers] = useState(0);
    const [value] = useDebounce(searchTerm, 1000);
    const [isDeleted, setIsDeleted] = useState(false);
    const itemsPerPage = 8;



    useEffect(() => {
        async function fetchUsers() {
            if (!isAuthenticated) return;
            setIsLoading(true);
            try {
                const response = await axios.get(`${import.meta.env.SERVER_URL}/user/paginated_users`, {
                    params: {
                        page: currentPage,
                        limit: itemsPerPage,
                        search: value
                    },
                    withCredentials: true
                });
                console.log("Response", response.data);

                setUsers(response.data.users || []);
                setTotalPages(response.data.totalPages || 0);
                setTotalUsers(response.data.totalUsers || 0);
                setCurrentPage(response.data.currentPage || 1);
            } catch (error) {
                console.error('Error fetching users:', error);
                setUsers([]);
                setTotalPages(0);
                setTotalUsers(0);
            }
            setIsLoading(false);
        };

        fetchUsers();
    }, [isAuthenticated, currentPage, value, isDeleted]);


    async function handleDeleteUser(userId: string | undefined) {
        try {
            console.log("Userid: ", userId);
            await axios.delete(`${import.meta.env.SERVER_URL}/user/delete_user`, {
                data: { userId: userId },
                withCredentials: true
            });

            console.log("User deleted successfully!!");

            setIsDeleted((prev) => !prev);
        } catch (error) {
            console.log("Error in deleting user: ", error);
        }
    }

    return (
        <div className="w-full mx-auto p-6 space-y-6 max-w-7xl">
            <h1 className="text-3xl font-bold mb-8">User Management</h1>
            <Card>
                <CardHeader>

                    <div className="flex items-center space-x-4 mt-4">
                        <div className="relative flex-1">
                            <Search className="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
                            <Input
                                placeholder="Search users..."
                                value={searchTerm}
                                onChange={(e) => setSearchTerm(e.target.value)}
                                className="pl-8"
                            />
                        </div>
                        <Badge variant="secondary" className="px-3 py-1">
                            {totalUsers} Users
                        </Badge>
                    </div>
                </CardHeader>
                <CardContent>
                    <div className="rounded-md border">
                        <Table>
                            <TableHeader>
                                <TableRow>
                                    <TableHead className="w-[16.66%]">ID</TableHead>
                                    <TableHead className="w-[16.66%]">Username</TableHead>
                                    <TableHead className="w-[16.66%]">Password</TableHead>
                                    <TableHead className="w-[16.66%]">Test Password</TableHead>
                                    <TableHead className="w-[16.66%]">Batch</TableHead>
                                    <TableHead className="w-[16.66%] text-right">Actions</TableHead>
                                </TableRow>
                            </TableHeader>
                            <TableBody>
                                {isLoading ? (
                                    <TableRow>
                                        <TableCell colSpan={6} className="h-24 text-center">
                                            <div className="flex items-center justify-center">
                                                <Loader2 className="h-6 w-6 animate-spin" />
                                            </div>
                                        </TableCell>
                                    </TableRow>

                                ) : (
                                    users.map((user: User) => {
                                        if (user.id === undefined) {
                                            return null;
                                        }
                                        return (
                                            <TableRow key={user.id}>
                                                <TableCell className="font-mono text-sm">
                                                    {user.id.slice(0, 8)}...
                                                </TableCell>
                                                <TableCell>
                                                    <div className="flex items-center space-x-2">
                                                        <div className="h-8 w-8 rounded-full bg-slate-100 flex items-center justify-center">
                                                            {user.username.charAt(0).toUpperCase()}
                                                        </div>
                                                        <span>{user.username}</span>
                                                    </div>
                                                </TableCell>
                                                <TableCell>{user.password}</TableCell>
                                                <TableCell>{user.testPassword}</TableCell>
                                                <TableCell>
                                                    <Badge variant="outline">{user.batch}</Badge>
                                                </TableCell>
                                                <TableCell className="text-right">
                                                    <Button variant="ghost" size="icon" onClick={() => handleDeleteUser(user.id)}>
                                                        <Trash2 className="h-4 w-4" />
                                                    </Button>
                                                </TableCell>
                                            </TableRow>
                                        );
                                    })
                                )}
                            </TableBody>
                        </Table>
                    </div>

                    {users.length === 0 && (
                        <Alert className="mt-4">
                            <AlertDescription>
                                No users found matching your search criteria.
                            </AlertDescription>
                        </Alert>
                    )}

                    <div className="flex items-center justify-between mt-4">
                        <Button
                            variant="outline"
                            disabled={currentPage === 1 || isLoading}
                            onClick={() => setCurrentPage(currentPage - 1)}
                        >
                            Previous
                        </Button>
                        <div className="text-sm text-muted-foreground">
                            Page {currentPage} of {totalPages}
                        </div>
                        <Button
                            variant="outline"
                            disabled={currentPage === totalPages || isLoading}
                            onClick={() => setCurrentPage(currentPage + 1)}
                        >
                            Next
                        </Button>
                    </div>


                </CardContent>
            </Card>
        </div>
    );
}