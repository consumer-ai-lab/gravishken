/* Do not change, this code is generated from Golang structs */


export enum Varient {
    ExeNotFound = 0,
    Err = 1,
    Unknown = 2,
}
export interface Message {
    Type: Varient;
    Val: string;
}
export interface TExeNotFound {
    Name: string;
    ErrMsg: string;
}
export interface TErr {
    Message: string;
}