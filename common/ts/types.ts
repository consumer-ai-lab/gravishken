/* Do not change, this code is generated from Golang structs */


export enum Varient {
    ExeNotFound = 0,
    UserLogin = 1,
    Err = 2,
    Unknown = 3,
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
export interface TErr {
    Message: string;
}