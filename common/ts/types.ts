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