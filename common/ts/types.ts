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
    UserLogin = 3,
    LoadRoute = 4,
    ReloadUi = 5,
    GetTest = 6,
    OpenApp = 7,
    QuitApp = 8,
    Unknown = 9,
}
export enum FileType {
    PPTX = "pptx",
    DOCX = "docx",
    XLSX = "xlsx",
}
export enum TestType {
    FileTest = "file",
    TypingTest = "typing",
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
export interface TUserLogin {
    Username: string;
    Password: string;
    TestCode: string;
}
export interface TLoadRoute {
    Route: string;
}
export interface TReloadUi {

}
export interface TGetTest {
    TestPassword: string;
}
export interface TOpenApp {
    Typ: AppType;
}
export interface TQuitApp {

}
export interface Time {

}
export interface UserSubmission {
    test: number[];
    startTime: Time;
    endTime: Time;
    elapsedTime: number;
    submissionReceived: boolean;
    readingElapsedTime: number;
    readingSubmissionReceived: boolean;
    submissionFolderId: string;
    mergedFileId: string;
    wpm: number;
    wpmNormal: number;
    resultDownloaded: boolean;
}
export interface User {
    id?: number[];
    name: string;
    username: string;
    password: string;
    testPassword: string;
    batch: string;
    tests?: UserSubmission;
}

export interface UserBatchRequestData {
    From: number;
    To: number;
    ResultDownloaded: boolean;
}
export interface UserLoginRequest {
    username: string;
    password: string;
    testPassword: string;
}
export interface UserUpdateRequest {
    username: string;
    property: string;
    value: string[];
}
export interface Test {
    testId: number[];
    testType: TestType;
    fileType?: FileType;
    typingTestText?: string;
}
export interface BatchTests {
    batchId: number[];
    tests: Test[];
    testDuration: number;
    password: string;
    startTime: Time;
    endTime: Time;
}

export interface Admin {
    id?: number[];
    username: string;
    password: string;
    token: string[];
}
export interface AdminChangePassword {
    username: string;
    newPassword: string;
}
export interface AdminRequest {
    username: string;
    token: string;
}
export interface Batch {
    id?: number[];
    batchName: string;
}