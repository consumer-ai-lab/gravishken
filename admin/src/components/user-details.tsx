import { useState, useEffect } from 'react';
import { User } from '@common/types';
import axios from 'axios';
import { ScrollArea } from "@/components/ui/scroll-area";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from './ui/card';
import { Loader2, Trash2 } from 'lucide-react';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from './ui/tooltip';


interface UserModelUpdationRequest {
    id: string
    username: string | undefined
    password: string | undefined
    testPassword: string | undefined
    batch: string | undefined
}

interface UserDetailsProps{
    isAuthenticated: boolean;
}

export default function UserDetails(
    {isAuthenticated}: UserDetailsProps
) {
    const [users, setUsers] = useState<User[]>([]);
    const [currentPage, setCurrentPage] = useState(1);
    const [totalPages, setTotalPages] = useState(0);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [editingUserId, setEditingUserId] = useState<string | null>(null);
    const [editedUser, setEditedUser] = useState<Partial<User> | null>(null);

    const fetchUsers = async (page: number) => {
        setIsLoading(true);
        setError(null);
        if(isAuthenticated === false){  
            return;
        }
        try {
            const response = await axios.get(`${import.meta.env.SERVER_URL}/user/paginated_users`, {
                params: {
                    page: page,
                    limit: 50
                },
                withCredentials: true
            });
            setUsers(response.data.users || []);
            setTotalPages(response.data.totalPages || 0);
        } catch (error) {
            console.error('Error fetching users:', error);
            setError('Failed to fetch users. Please try again.');
            setUsers([]);
            setTotalPages(0);
        }
        setIsLoading(false);
    };

    useEffect(() => {
        fetchUsers(currentPage);
    }, [currentPage, isAuthenticated]);

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

    const handleEditUser = (user: User) => {
        setEditingUserId(user.id!);
        setEditedUser(user);
    };

    const handleCancelEdit = () => {
        setEditingUserId(null);
        setEditedUser(null);
    };

    const handleSaveUser = async () => {
        if (editedUser && editingUserId) {
            try {
                // Prepare the request body
                const updateUserData: UserModelUpdationRequest = {
                    id: editingUserId,
                    username: editedUser.username,
                    password: editedUser.password,
                    testPassword: editedUser.testPassword,
                    batch: editedUser.batch
                };

                // Update the user in the backend
                await axios.put(`${import.meta.env.SERVER_URL}/user/update_user`, updateUserData, {
                    withCredentials: true
                });

                // Update the users list in the state
                setUsers(users.map(user => user.id === editingUserId ? { ...user, ...editedUser } : user));
                setEditingUserId(null);
                setEditedUser(null);
            } catch (error) {
                console.error('Error updating user:', error);
                setError('Failed to update user. Please try again.');
            }
        }
    };

    const handleDeleteUser = async (userId: string | undefined) => {
        try {
            console.log("Userid: ", userId);
            await axios.delete(`${import.meta.env.SERVER_URL}/user/delete_user`, {
                data: { userId: userId },
                withCredentials: true
            });

            console.log("User deleted successfully!!");

            fetchUsers(currentPage);
        } catch (error) {
            console.log("Error in deleting user: ", error);
            setError("Failed to delete user!");
        }
    }

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (editedUser) {
            setEditedUser({
                ...editedUser,
                [e.target.name]: e.target.value
            });
        }
    };

    if (error) {
        return <div className="text-red-500">{error}</div>;
    }

    return (
        <div className="w-full mx-auto p-4 space-y-6">
            <h1 className="text-3xl font-bold mb-8">View All Users</h1>
            <Card>
                <CardContent >
                    <div className="relative mt-6">
                        <div className="overflow-hidden border rounded-md">
                            <Table>
                                <colgroup>
                                    <col className="w-[15%]" /><col className="w-[18%]" /><col className="w-[15%]" /><col className="w-[26%]" /><col className="w-[10%]" /><col className="w-[10%]" />
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
                                                {editingUserId === user.id ? (
                                                    <>
                                                        <TableCell>
                                                            <input
                                                                type="text"
                                                                name="id"
                                                                value={editedUser?.id || ''}
                                                                disabled
                                                                className="border p-1"
                                                            />
                                                        </TableCell>
                                                        <TableCell>
                                                            <input
                                                                type="text"
                                                                name="username"
                                                                value={editedUser?.username || ''}
                                                                onChange={handleInputChange}
                                                                className="border p-1"
                                                            />
                                                        </TableCell>
                                                        <TableCell>
                                                            <input
                                                                type="text"
                                                                name="password"
                                                                value={editedUser?.password || ''}
                                                                onChange={handleInputChange}
                                                                className="border p-1"
                                                            />
                                                        </TableCell>
                                                        <TableCell>
                                                            <input
                                                                type="text"
                                                                name="testPassword"
                                                                value={editedUser?.testPassword || ''}
                                                                onChange={handleInputChange}
                                                                className="border p-1"
                                                            />
                                                        </TableCell>
                                                        <TableCell>
                                                            <input
                                                                type="text"
                                                                name="batch"
                                                                value={editedUser?.batch || ''}
                                                                onChange={handleInputChange}
                                                                className="border p-1"
                                                            />
                                                        </TableCell>
                                                        {/* <TableCell className='flex justify-center'>
                                                            <Button variant={"ghost"} onClick={handleSaveUser}>
                                                                Save
                                                            </Button>
                                                            <Button variant={"ghost"} onClick={handleCancelEdit}>
                                                                Cancel
                                                            </Button>
                                                        </TableCell> */}
                                                    </>
                                                ) : (
                                                    <>
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
                                                            <Button variant={"ghost"} onClick={() => handleEditUser(user)}>
                                                                Edit
                                                            </Button>
                                                            <Button variant="ghost" className="p-2" onClick={() => handleDeleteUser(user.id?.toString())}>
                                                                <Trash2 className="h-4 w-4" />
                                                            </Button>
                                                        </TableCell>
                                                    </>
                                                )}
                                            </TableRow>
                                        ))}
                                    </TableBody>
                                </Table>
                            </ScrollArea>
                        </div>
                    </div>
                    <div className="flex justify-between items-center mt-4">
                        {editingUserId ? (
                            <>
                                <Button
                                    onClick={handleCancelEdit}
                                    variant="outline"
                                >
                                    Cancel
                                </Button>
                                <Button
                                    onClick={handleSaveUser}
                                    variant="outline"
                                >
                                    Save
                                </Button>
                            </>
                        ) : (
                            <>
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
                            </>
                        )}
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
