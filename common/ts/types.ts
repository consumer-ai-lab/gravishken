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
    UserLogin = 2,
    LoadRoute = 3,
    ReloadUi = 4,
    GetTest = 5,
    OpenApp = 6,
    Unknown = 7,
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