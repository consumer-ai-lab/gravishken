/* Do not change, this code is generated from Golang structs */


export enum Varient {
    Var1 = 0,
    Var2 = 1,
    Err = 2,
    Unknown = 3,
}
export interface Message {
    Type: Varient;
    Val: string;
}
export interface TVar1 {
    Field1: number;
    Field2: boolean;
}
export interface TVar2 {
    Field1: boolean;
    Field3: string;
}
export interface TErr {
    Message: string;
}