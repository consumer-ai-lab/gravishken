/* Do not change, this code is generated from Golang structs */


export enum AppType {
    TXT = 0,
    DOCX = 1,
    XLSX = 2,
    PPTX = 3,
}
export enum Varient {
    Err = 0,
    ExeNotFound = 1,
    Quit = 2,
    UserLoginRequest = 3,
    WarnUser = 4,
    LoadRoute = 5,
    ReloadUi = 6,
    StartTest = 7,
    OpenApp = 8,
    QuitApp = 9,
    Unknown = 10,
}
export enum TestType {
    TypingTest = "typing",
    DocxTest = "docx",
    ExcelTest = "excel",
    WordTest = "word",
}
export interface TErr {
    Message: string;
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
export interface TStartTestRequest {

}
export interface TStartTest {

}
export interface TOpenApp {
    Typ: AppType;
}
export interface TQuitApp {

}
export interface User {
    Id: string;
    Username: string;
    Password: string;
    BatchName: string;
}
export interface Time {

}
export interface UserSubmission {
    UserId: string;
    TestId: string;
    StartTime: Time;
    EndTime: Time;
    ElapsedTime: number;
    WPM: number;
    WPMNormal: number;
    ReadingSubmissionReceived: boolean;
    ReadingElapsedTime: number;
    SubmissionReceived: boolean;
    ResultDownloaded: boolean;
    MergedFileID: string;
    SubmissionFolderID: string;
}
export interface UserBatchRequestData {
    From: number;
    To: number;
    ResultDownloaded: boolean;
}
export interface Test {
    Id: string;
    Type: TestType;
    Duration: number;
    File: string;
    TypingText: string;
}
export interface Admin {
    Id: string;
    Username: string;
    Password: string;
}
export interface AdminRequest {
    Username: string;
    Token: string;
}
export interface Batch {
    Id: string;
    Name: string;
    Tests: string[];
}