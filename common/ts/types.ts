/* Do not change, this code is generated from Golang structs */


export enum AppType {
    TXT = 0,
    DOCX = 1,
    XLSX = 2,
    PPTX = 3,
}
export enum Varient {
    Err = 0,
    Notification = 1,
    ExeNotFound = 2,
    Quit = 3,
    UserLoginRequest = 4,
    WarnUser = 5,
    LoadRoute = 6,
    ReloadUi = 7,
    StartTest = 8,
    TestFinished = 9,
    CheckSystem = 10,
    OpenApp = 11,
    QuitApp = 12,
    Unknown = 13,
}
export enum TestType {
    TypingTest = "typing",
    DocxTest = "docx",
    ExcelTest = "xlsx",
    PptTest = "pptx",
    MCQTest = "mcq",
}
export interface TErr {
    Message: string;
}
export interface TNotification {
    Message: string;
    Typ: string;
}
export interface Message {
    Typ: Varient;
    Val: string;
}
export interface TExeNotFound {
    Name: string;
    ErrMsg: string;
}
export interface TQuit {

}
export interface TUserLoginRequest {
    Username: string;
    Password: string;
}
export interface TWarnUser {
    Message: string;
}
export interface TLoadRoute {
    Route: string;
}
export interface TReloadUi {

}
export interface TStartTest {

}
export interface TTestFinished {

}
export interface TCheckSystem {

}
export interface TOpenApp {
    Typ: AppType;
    TestId: string;
}
export interface TQuitApp {

}
export interface User {
    Id: string;
    Username: string;
    Password: string;
    Batch: string;
}
export interface AppTestInfo {
    FileData: string;
}
export interface McqTestInfo {
    Answers: number[];
}
export interface TypingTestInfo {
    TimeTaken: number;
    WPM: number;
    RawWPM: number;
    Accuracy: number;
}
export interface TestInfo {
    Type: TestType;
    TypingTestInfo?: TypingTestInfo;
    McqTestInfo?: McqTestInfo;
    DocxTestInfo?: AppTestInfo;
    ExcelTestInfo?: AppTestInfo;
    PptTestInfo?: AppTestInfo;
}
export interface TestSubmission {
    UserId: string;
    TestId: string;
    TestInfo: TestInfo;
}



export interface Test {
    Id: string;
    TestName: string;
    Duration: number;
    Type: TestType;
    FilePath?: string;
    TypingText?: string;
    McqJson?: string;
}
export interface Admin {
    Id: string;
    Username: string;
    Password: string;
}
export interface Batch {
    Id: string;
    Name: string;
    Tests: string[];
}