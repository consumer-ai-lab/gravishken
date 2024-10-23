import React from 'react';
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from "@/components/ui/accordion";
import { Download, LogIn, Shield, Zap, LayoutDashboard } from "lucide-react";
import { useNavigate } from "react-router-dom";

export default function LandingPage() {
    const navigate = useNavigate();

    return (
        <div className="min-h-screen flex flex-col bg-gradient-to-b from-gray-100 to-gray-200 dark:from-gray-900 dark:to-gray-800">
            <header className="bg-blue-700 text-white shadow-md">
                <div className="container mx-auto px-4 py-6 flex justify-between items-center">
                    <h1 className="text-2xl font-bold">Gravishken</h1>
                    <Button variant="outline" size="lg" className="px-4 hidden font-bold sm:flex text-blue-600  border-white" onClick={() => navigate("/login")}>
                        <LogIn className="mr-2 h-4 w-4" />
                        Login as Test Admin
                    </Button>
                </div>
            </header>

            <main className="flex-grow container mx-auto px-4 py-8">
                <section className="text-center min-h-[400px] flex flex-col justify-center items-center">
                    <h2 className="text-6xl md:text-7xl font-bold mb-6 bg-clip-text text-black">Cheating-Proof Testing Platform</h2>
                    <p className="text-xl md:text-2xl mb-10 text-muted-foreground max-w-2xl">Secure, reliable, and easy to use. Elevate your testing experience with Gravishken.</p>
                    <a href={`${import.meta.env.SERVER_URL}/release/latest/windows`} target="_blank" rel="noopener noreferrer">
                        <Button size="lg" className="text-lg bg-blue-600 hover:bg-blue-700 hover:text-white transition-colors duration-300">
                            <Download className="mr-2 h-5 w-5" />
                            Download Test App
                        </Button>
                    </a>
                </section>

                <section className="grid md:grid-cols-3 gap-8 mb-16 px-12">
                    <Card className="bg-white dark:bg-gray-800 shadow-lg hover:shadow-xl transition-shadow duration-300">
                        <CardHeader>
                            <CardTitle className="flex items-center text-blue-600">
                                <Zap className="mr-2 h-6 w-6" />
                                Single Executable Setup
                            </CardTitle>
                        </CardHeader>
                        <CardContent>
                            <p className="text-muted-foreground">
                                Our application installs quickly with a single .exe file, making deployment a breeze on any Windows device.
                            </p>
                        </CardContent>
                    </Card>
                    <Card className="bg-white dark:bg-gray-800 shadow-lg hover:shadow-xl transition-shadow duration-300">
                        <CardHeader>
                            <CardTitle className="flex items-center text-green-600">
                                <Shield className="mr-2 h-6 w-6" />
                                Cheating-Proof Environment
                            </CardTitle>
                        </CardHeader>
                        <CardContent>
                            <p className="text-muted-foreground">
                                Once installed, Gravishken disables background processes, ensuring a secure and fair testing environment for all participants.
                            </p>
                        </CardContent>
                    </Card>
                    <Card className="bg-white dark:bg-gray-800 shadow-lg hover:shadow-xl transition-shadow duration-300">
                        <CardHeader>
                            <CardTitle className="flex items-center text-purple-600">
                                <LayoutDashboard className="mr-2 h-6 w-6" />
                                Powerful Admin Panel
                            </CardTitle>
                        </CardHeader>
                        <CardContent>
                            <p className="text-muted-foreground">
                                Test administrators can easily monitor and control the testing environment through our intuitive admin panel, accessible with a simple login.
                            </p>
                        </CardContent>
                    </Card>
                </section>

                <section className="mb-16 px-12">
                    <Card className="bg-white dark:bg-gray-800 shadow-lg">
                        <CardHeader>
                            <CardTitle className="text-2xl font-bold text-blue-600">Why Choose Gravishken?</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <Accordion type="single" collapsible className="w-full">
                                <AccordionItem value="item-1">
                                    <AccordionTrigger>Secure and Tamper-Proof</AccordionTrigger>
                                    <AccordionContent>
                                        Our testing environment is designed to prevent cheating and ensure the integrity of your exams. With advanced security measures, you can trust that your tests are conducted fairly.
                                    </AccordionContent>
                                </AccordionItem>
                                <AccordionItem value="item-2">
                                    <AccordionTrigger>Easy Deployment</AccordionTrigger>
                                    <AccordionContent>
                                        With our single executable setup, deploying Gravishken is quick and hassle-free. You'll have your testing environment up and running in no time.
                                    </AccordionContent>
                                </AccordionItem>
                                <AccordionItem value="item-3">
                                    <AccordionTrigger>Real-Time Monitoring</AccordionTrigger>
                                    <AccordionContent>
                                        Administrators have full control with our real-time monitoring capabilities. Keep an eye on test progress and ensure everything runs smoothly.
                                    </AccordionContent>
                                </AccordionItem>
                                <AccordionItem value="item-4">
                                    <AccordionTrigger>Customizable Settings</AccordionTrigger>
                                    <AccordionContent>
                                        Tailor the testing experience to your needs with our customizable test settings and parameters. Gravishken adapts to your specific requirements.
                                    </AccordionContent>
                                </AccordionItem>
                                <AccordionItem value="item-5">
                                    <AccordionTrigger>Comprehensive Reporting</AccordionTrigger>
                                    <AccordionContent>
                                        Get valuable insights with our comprehensive reporting and analytics. Make data-driven decisions to improve your testing processes.
                                    </AccordionContent>
                                </AccordionItem>
                            </Accordion>
                        </CardContent>
                    </Card>
                </section>
            </main>

            <footer className="bg-blue-700 text-white mt-auto">
                <div className="container mx-auto px-4 py-6 text-center">
                    <p>&copy; 2024 Gravishken. All rights reserved.</p>
                </div>
            </footer>
        </div>
    )
}
