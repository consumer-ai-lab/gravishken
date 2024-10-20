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
    SubmitTest = 8,
    OpenApp = 9,
    QuitApp = 10,
    Unknown = 11,
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
export interface TSubmitTest {
    TestId: string;
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